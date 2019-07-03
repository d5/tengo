package runtime

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/d5/tengo"
	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/compiler/source"
	"github.com/d5/tengo/compiler/token"
)

const (
	// StackSize is the maximum stack size.
	StackSize = 2048

	// GlobalsSize is the maximum number of global variables.
	GlobalsSize = 1024

	// MaxFrames is the maximum number of function frames.
	MaxFrames = 1024
)

// VM is a virtual machine that executes the bytecode compiled by Compiler.
type VM struct {
	constants   []tengo.Object
	stack       [StackSize]tengo.Object
	sp          int
	globals     []tengo.Object
	fileSet     *source.FileSet
	frames      [MaxFrames]Frame
	framesIndex int
	curFrame    *Frame
	curInsts    []byte
	ip          int
	aborting    int64
	maxAllocs   int64
	allocs      int64
	lastPanic   interface{}
}

// NewVM creates a VM.
func NewVM(bytecode *compiler.Bytecode, globals []tengo.Object, maxAllocs int64) *VM {
	if globals == nil {
		globals = make([]tengo.Object, GlobalsSize)
	}

	v := &VM{
		constants:   bytecode.Constants,
		sp:          0,
		globals:     globals,
		fileSet:     bytecode.FileSet,
		framesIndex: 1,
		ip:          -1,
		maxAllocs:   maxAllocs,
	}

	v.frames[0].fn = bytecode.MainFunction
	v.frames[0].ip = -1
	v.curFrame = &v.frames[0]
	v.curInsts = v.curFrame.fn.Instructions

	return v
}

// Abort aborts the execution.
func (v *VM) Abort() {
	atomic.StoreInt64(&v.aborting, 1)
}

// Run starts the execution.
func (v *VM) Run() error {
	defer func() {
		atomic.StoreInt64(&v.aborting, 0)
		v.lastPanic = nil
	}()

	// reset VM states
	v.sp = 0
	v.curFrame = &(v.frames[0])
	v.curInsts = v.curFrame.fn.Instructions
	v.framesIndex = 1
	v.ip = -1
	v.allocs = v.maxAllocs + 1

	if err := v.run(); err != nil {
		filePos := v.fileSet.Position(v.curFrame.fn.SourcePos(v.ip - 1))
		err = fmt.Errorf("Runtime Error: %s\n\tat %s", err.Error(), filePos)
		for v.framesIndex > 1 {
			v.framesIndex--
			v.curFrame = &v.frames[v.framesIndex-1]

			filePos = v.fileSet.Position(v.curFrame.fn.SourcePos(v.curFrame.ip - 1))
			err = fmt.Errorf("%s\n\tat %s", err.Error(), filePos)
		}
		return err
	}

	return nil
}

// InteropCall is an interop call method for Go functions.
func (v *VM) InteropCall(
	fn tengo.Object,
	args ...tengo.Object,
) (ret tengo.Object, err error) {
	numArgs := len(args)
	numLocals := numArgs
	if cf, ok := fn.(*tengo.CompiledFunction); ok {
		numLocals = cf.NumLocals
	}

	// check for stack overflow
	if v.sp+numLocals >= StackSize {
		return nil, ErrStackOverflow
	}

	// create a micro-function to handle the call
	callee := &tengo.CompiledFunction{
		Instructions: []byte{
			compiler.OpCall,
			byte(numArgs),
			compiler.OpSuspend,
		},
		SourceMap: map[int]source.Pos{}, // TODO: fix this
	}

	// enter new frame
	v.curFrame.ip = v.ip // store current ip before call
	v.curFrame = &(v.frames[v.framesIndex])
	v.curFrame.fn = callee
	v.curFrame.freeVars = nil
	v.curFrame.basePointer = v.sp
	v.curInsts = callee.Instructions
	v.ip = -1
	v.framesIndex++

	// set up the stack for the call in the micro-function
	v.stack[v.sp] = fn
	copy(v.stack[v.sp+1:], args)
	v.sp += numArgs + 1

	// run
	defer func() {
		if p := recover(); p != nil {
			v.lastPanic = p
		}
	}()
	err = v.run()
	if err == nil {
		ret = v.stack[v.sp-1]
	}

	// leave frame
	v.framesIndex--
	v.sp = v.frames[v.framesIndex].basePointer
	v.curFrame = &v.frames[v.framesIndex-1]
	v.curInsts = v.curFrame.fn.Instructions
	v.ip = v.curFrame.ip
	return
}

func (v *VM) run() error {
	for atomic.LoadInt64(&v.aborting) == 0 {
		v.ip++

		switch v.curInsts[v.ip] {
		case compiler.OpSuspend:
			return nil

		case compiler.OpConstant:
			v.ip += 2
			cidx := int(v.curInsts[v.ip]) | int(v.curInsts[v.ip-1])<<8

			v.stack[v.sp] = v.constants[cidx]
			v.sp++

		case compiler.OpNull:
			v.stack[v.sp] = tengo.UndefinedValue
			v.sp++

		case compiler.OpBinaryOp:
			v.ip++
			right := v.stack[v.sp-1]
			left := v.stack[v.sp-2]

			tok := token.Token(v.curInsts[v.ip])
			res, err := left.BinaryOp(tok, right)
			if err != nil {
				v.sp -= 2

				if err == tengo.ErrInvalidOperator {
					return fmt.Errorf("invalid operation: %s %s %s",
						left.TypeName(), tok.String(), right.TypeName())
				}
				return err
			}

			v.allocs--
			if v.allocs == 0 {
				return ErrObjectAllocLimit
			}

			v.stack[v.sp-2] = res
			v.sp--

		case compiler.OpEqual:
			right := v.stack[v.sp-1]
			left := v.stack[v.sp-2]
			v.sp -= 2

			if left.Equals(right) {
				v.stack[v.sp] = tengo.TrueValue
			} else {
				v.stack[v.sp] = tengo.FalseValue
			}
			v.sp++

		case compiler.OpNotEqual:
			right := v.stack[v.sp-1]
			left := v.stack[v.sp-2]
			v.sp -= 2

			if left.Equals(right) {
				v.stack[v.sp] = tengo.FalseValue
			} else {
				v.stack[v.sp] = tengo.TrueValue
			}
			v.sp++

		case compiler.OpPop:
			v.sp--

		case compiler.OpTrue:
			v.stack[v.sp] = tengo.TrueValue
			v.sp++

		case compiler.OpFalse:
			v.stack[v.sp] = tengo.FalseValue
			v.sp++

		case compiler.OpLNot:
			operand := v.stack[v.sp-1]
			v.sp--

			if operand.IsFalsy() {
				v.stack[v.sp] = tengo.TrueValue
			} else {
				v.stack[v.sp] = tengo.FalseValue
			}
			v.sp++

		case compiler.OpBComplement:
			operand := v.stack[v.sp-1]
			v.sp--

			switch x := operand.(type) {
			case *tengo.Int:
				var res tengo.Object = &tengo.Int{Value: ^x.Value}

				v.allocs--
				if v.allocs == 0 {
					return ErrObjectAllocLimit
				}

				v.stack[v.sp] = res
				v.sp++
			default:
				return fmt.Errorf("invalid operation: ^%s", operand.TypeName())
			}

		case compiler.OpMinus:
			operand := v.stack[v.sp-1]
			v.sp--

			switch x := operand.(type) {
			case *tengo.Int:
				var res tengo.Object = &tengo.Int{Value: -x.Value}

				v.allocs--
				if v.allocs == 0 {
					return ErrObjectAllocLimit
				}

				v.stack[v.sp] = res
				v.sp++
			case *tengo.Float:
				var res tengo.Object = &tengo.Float{Value: -x.Value}

				v.allocs--
				if v.allocs == 0 {
					return ErrObjectAllocLimit
				}

				v.stack[v.sp] = res
				v.sp++
			default:
				return fmt.Errorf("invalid operation: -%s", operand.TypeName())
			}

		case compiler.OpJumpFalsy:
			v.ip += 2
			v.sp--
			if v.stack[v.sp].IsFalsy() {
				pos := int(v.curInsts[v.ip]) | int(v.curInsts[v.ip-1])<<8
				v.ip = pos - 1
			}

		case compiler.OpAndJump:
			v.ip += 2

			if v.stack[v.sp-1].IsFalsy() {
				pos := int(v.curInsts[v.ip]) | int(v.curInsts[v.ip-1])<<8
				v.ip = pos - 1
			} else {
				v.sp--
			}

		case compiler.OpOrJump:
			v.ip += 2

			if v.stack[v.sp-1].IsFalsy() {
				v.sp--
			} else {
				pos := int(v.curInsts[v.ip]) | int(v.curInsts[v.ip-1])<<8
				v.ip = pos - 1
			}

		case compiler.OpJump:
			pos := int(v.curInsts[v.ip+2]) | int(v.curInsts[v.ip+1])<<8
			v.ip = pos - 1

		case compiler.OpSetGlobal:
			v.ip += 2
			v.sp--

			globalIndex := int(v.curInsts[v.ip]) | int(v.curInsts[v.ip-1])<<8
			v.globals[globalIndex] = v.stack[v.sp]

		case compiler.OpSetSelGlobal:
			v.ip += 3
			globalIndex := int(v.curInsts[v.ip-1]) | int(v.curInsts[v.ip-2])<<8
			numSelectors := int(v.curInsts[v.ip])

			// selectors and RHS value
			selectors := make([]tengo.Object, numSelectors)
			for i := 0; i < numSelectors; i++ {
				selectors[i] = v.stack[v.sp-numSelectors+i]
			}

			val := v.stack[v.sp-numSelectors-1]
			v.sp -= numSelectors + 1

			if err := indexAssign(v.globals[globalIndex], val, selectors); err != nil {
				return err
			}

		case compiler.OpGetGlobal:
			v.ip += 2
			globalIndex := int(v.curInsts[v.ip]) | int(v.curInsts[v.ip-1])<<8

			val := v.globals[globalIndex]

			v.stack[v.sp] = val
			v.sp++

		case compiler.OpArray:
			v.ip += 2
			numElements := int(v.curInsts[v.ip]) | int(v.curInsts[v.ip-1])<<8

			var elements []tengo.Object
			for i := v.sp - numElements; i < v.sp; i++ {
				elt := v.stack[i]
				if spread, ok := elt.(*tengo.Spread); ok {
					elements = append(elements, spread.Values...)
				} else {
					elements = append(elements, elt)
				}
			}

			v.sp -= numElements

			var arr tengo.Object = &tengo.Array{Value: elements}

			v.allocs--
			if v.allocs == 0 {
				return ErrObjectAllocLimit
			}

			v.stack[v.sp] = arr
			v.sp++

		case compiler.OpMap:
			v.ip += 2
			numElements := int(v.curInsts[v.ip]) | int(v.curInsts[v.ip-1])<<8

			kv := make(map[string]tengo.Object)
			for i := v.sp - numElements; i < v.sp; i += 2 {
				key := v.stack[i]
				value := v.stack[i+1]
				kv[key.(*tengo.String).Value] = value
			}
			v.sp -= numElements

			var m tengo.Object = &tengo.Map{Value: kv}

			v.allocs--
			if v.allocs == 0 {
				return ErrObjectAllocLimit
			}

			v.stack[v.sp] = m
			v.sp++

		case compiler.OpError:
			value := v.stack[v.sp-1]

			var e tengo.Object = &tengo.Error{
				Value: value,
			}

			v.allocs--
			if v.allocs == 0 {
				return ErrObjectAllocLimit
			}

			v.stack[v.sp-1] = e

		case compiler.OpImmutable:
			value := v.stack[v.sp-1]

			switch value := value.(type) {
			case *tengo.Array:
				var immutableArray tengo.Object = &tengo.ImmutableArray{
					Value: value.Value,
				}

				v.allocs--
				if v.allocs == 0 {
					return ErrObjectAllocLimit
				}

				v.stack[v.sp-1] = immutableArray
			case *tengo.Map:
				var immutableMap tengo.Object = &tengo.ImmutableMap{
					Value: value.Value,
				}

				v.allocs--
				if v.allocs == 0 {
					return ErrObjectAllocLimit
				}

				v.stack[v.sp-1] = immutableMap
			}

		case compiler.OpIndex:
			index := v.stack[v.sp-1]
			left := v.stack[v.sp-2]
			v.sp -= 2

			val, err := left.IndexGet(index)
			if err != nil {
				if err == tengo.ErrNotIndexable {
					return fmt.Errorf("not indexable: %s", index.TypeName())
				}
				if err == tengo.ErrInvalidIndexType {
					return fmt.Errorf("invalid index type: %s", index.TypeName())
				}
				return err
			}

			if val == nil {
				val = tengo.UndefinedValue
			}

			v.stack[v.sp] = val
			v.sp++

		case compiler.OpSliceIndex:
			high := v.stack[v.sp-1]
			low := v.stack[v.sp-2]
			left := v.stack[v.sp-3]
			v.sp -= 3

			var lowIdx int64
			if low != tengo.UndefinedValue {
				if low, ok := low.(*tengo.Int); ok {
					lowIdx = low.Value
				} else {
					return fmt.Errorf("invalid slice index type: %s", low.TypeName())
				}
			}

			switch left := left.(type) {
			case *tengo.Array:
				numElements := int64(len(left.Value))
				var highIdx int64
				if high == tengo.UndefinedValue {
					highIdx = numElements
				} else if high, ok := high.(*tengo.Int); ok {
					highIdx = high.Value
				} else {
					return fmt.Errorf("invalid slice index type: %s", high.TypeName())
				}

				if lowIdx > highIdx {
					return fmt.Errorf("invalid slice index: %d > %d", lowIdx, highIdx)
				}

				if lowIdx < 0 {
					lowIdx = 0
				} else if lowIdx > numElements {
					lowIdx = numElements
				}

				if highIdx < 0 {
					highIdx = 0
				} else if highIdx > numElements {
					highIdx = numElements
				}

				var val tengo.Object = &tengo.Array{Value: left.Value[lowIdx:highIdx]}

				v.allocs--
				if v.allocs == 0 {
					return ErrObjectAllocLimit
				}

				v.stack[v.sp] = val
				v.sp++

			case *tengo.ImmutableArray:
				numElements := int64(len(left.Value))
				var highIdx int64
				if high == tengo.UndefinedValue {
					highIdx = numElements
				} else if high, ok := high.(*tengo.Int); ok {
					highIdx = high.Value
				} else {
					return fmt.Errorf("invalid slice index type: %s", high.TypeName())
				}

				if lowIdx > highIdx {
					return fmt.Errorf("invalid slice index: %d > %d", lowIdx, highIdx)
				}

				if lowIdx < 0 {
					lowIdx = 0
				} else if lowIdx > numElements {
					lowIdx = numElements
				}

				if highIdx < 0 {
					highIdx = 0
				} else if highIdx > numElements {
					highIdx = numElements
				}

				var val tengo.Object = &tengo.Array{Value: left.Value[lowIdx:highIdx]}

				v.allocs--
				if v.allocs == 0 {
					return ErrObjectAllocLimit
				}

				v.stack[v.sp] = val
				v.sp++

			case *tengo.String:
				numElements := int64(len(left.Value))
				var highIdx int64
				if high == tengo.UndefinedValue {
					highIdx = numElements
				} else if high, ok := high.(*tengo.Int); ok {
					highIdx = high.Value
				} else {
					return fmt.Errorf("invalid slice index type: %s", high.TypeName())
				}

				if lowIdx > highIdx {
					return fmt.Errorf("invalid slice index: %d > %d", lowIdx, highIdx)
				}

				if lowIdx < 0 {
					lowIdx = 0
				} else if lowIdx > numElements {
					lowIdx = numElements
				}

				if highIdx < 0 {
					highIdx = 0
				} else if highIdx > numElements {
					highIdx = numElements
				}

				var val tengo.Object = &tengo.String{Value: left.Value[lowIdx:highIdx]}

				v.allocs--
				if v.allocs == 0 {
					return ErrObjectAllocLimit
				}

				v.stack[v.sp] = val
				v.sp++

			case *tengo.Bytes:
				numElements := int64(len(left.Value))
				var highIdx int64
				if high == tengo.UndefinedValue {
					highIdx = numElements
				} else if high, ok := high.(*tengo.Int); ok {
					highIdx = high.Value
				} else {
					return fmt.Errorf("invalid slice index type: %s", high.TypeName())
				}

				if lowIdx > highIdx {
					return fmt.Errorf("invalid slice index: %d > %d", lowIdx, highIdx)
				}

				if lowIdx < 0 {
					lowIdx = 0
				} else if lowIdx > numElements {
					lowIdx = numElements
				}

				if highIdx < 0 {
					highIdx = 0
				} else if highIdx > numElements {
					highIdx = numElements
				}

				var val tengo.Object = &tengo.Bytes{Value: left.Value[lowIdx:highIdx]}

				v.allocs--
				if v.allocs == 0 {
					return ErrObjectAllocLimit
				}

				v.stack[v.sp] = val
				v.sp++
			}

		case compiler.OpSpread:
			spreadSP := v.sp - 1
			target := v.stack[spreadSP]
			if !target.CanSpread() {
				return fmt.Errorf("cannot spread value of type %s", target.TypeName())
			}

			v.stack[spreadSP] = &tengo.Spread{
				Values: target.Spread(),
			}

		case compiler.OpCall:
			numArgs := int(v.curInsts[v.ip+1])
			v.ip++
			spBase := v.sp - 1 - numArgs
			value := v.stack[spBase]

			if !value.CanCall() {
				return fmt.Errorf("not callable: %s", value.TypeName())
			}

			if numArgs > 0 {
				i := v.sp - 1
				arg := v.stack[i]
				if spread, ok := arg.(*tengo.Spread); ok {
					list := spread.Values
					numSpreadValues := len(list)
					if v.sp+numSpreadValues >= StackSize {
						return ErrStackOverflow
					}
					ebStart, ebEnd := i, i+numSpreadValues
					rxStart, rxEnd := i+1, spBase+numArgs+1
					rmStart, rmEnd := i+numSpreadValues, spBase+numArgs+numSpreadValues
					copy(v.stack[rmStart:rmEnd], v.stack[rxStart:rxEnd])
					copy(v.stack[ebStart:ebEnd], list)
					numArgs += numSpreadValues - 1
					v.sp += numSpreadValues - 1
				}
			}

			if callee, ok := value.(*tengo.CompiledFunction); ok {
				if callee.VarArgs {
					// if the closure is variadic,
					// roll up all variadic parameters into an array
					realArgs := callee.NumParameters - 1
					varArgs := numArgs - realArgs
					if varArgs >= 0 {
						numArgs = realArgs + 1
						args := make([]tengo.Object, varArgs)
						spStart := v.sp - varArgs
						for i := spStart; i < v.sp; i++ {
							args[i-spStart] = v.stack[i]
						}
						v.stack[spStart] = &tengo.Array{Value: args}
						v.sp = spStart + 1
					}
				}

				if numArgs != callee.NumParameters {
					if callee.VarArgs {
						return fmt.Errorf("wrong number of arguments: want>=%d, got=%d",
							callee.NumParameters-1, numArgs)
					}
					return fmt.Errorf("wrong number of arguments: want=%d, got=%d",
						callee.NumParameters, numArgs)
				}

				// test if it's tail-call
				if callee == v.curFrame.fn { // recursion
					nextOp := v.curInsts[v.ip+1]
					if nextOp == compiler.OpReturn ||
						(nextOp == compiler.OpPop && compiler.OpReturn == v.curInsts[v.ip+2]) {
						for p := 0; p < numArgs; p++ {
							v.stack[v.curFrame.basePointer+p] = v.stack[v.sp-numArgs+p]
						}
						v.sp -= numArgs + 1
						v.ip = -1 // reset IP to beginning of the frame
						continue
					}
				}

				if v.framesIndex >= MaxFrames {
					return ErrStackOverflow
				}

				// update call frame
				v.curFrame.ip = v.ip // store current ip before call
				v.curFrame = &(v.frames[v.framesIndex])
				v.curFrame.fn = callee
				v.curFrame.freeVars = callee.Free
				v.curFrame.basePointer = v.sp - numArgs
				v.curInsts = callee.Instructions
				v.ip = -1
				v.framesIndex++
				v.sp = v.sp - numArgs + callee.NumLocals
			} else {
				var args []tengo.Object
				args = append(args, v.stack[v.sp-numArgs:v.sp]...)

				ret, err := value.Call(v, args...)
				if v.lastPanic != nil {
					// if the runtime had a panic while executing user function
					// it must not continue its execution even if the panic
					// was handled or ignored by the user function.
					panic(v.lastPanic)
				}

				v.sp -= numArgs + 1

				// runtime error
				if err != nil {
					if err == tengo.ErrWrongNumArguments {
						return fmt.Errorf("wrong number of arguments in call to '%s'",
							value.TypeName())
					}
					if err, ok := err.(tengo.ErrInvalidArgumentType); ok {
						return fmt.Errorf("invalid type for argument '%s' in call to '%s': expected %s, found %s",
							err.Name, value.TypeName(), err.Expected, err.Found)
					}
					return err
				}

				if ret == nil {
					// nil return -> undefined
					ret = tengo.UndefinedValue
				}

				v.allocs--
				if v.allocs == 0 {
					return ErrObjectAllocLimit
				}

				v.stack[v.sp] = ret
				v.sp++
			}

		case compiler.OpReturnIfError:
			errv, isErr := v.stack[v.sp-1].(*tengo.Error)
			if !isErr {
				v.ip++
				// no error, just move on
				continue
			}

			// special case: try expressions in the top level
			// will be treated as runtime errors and returned
			if v.framesIndex == 1 {
				errStr, _ := tengo.ToString(errv.Value)
				return errors.New(errStr)
			}

			fallthrough

		case compiler.OpReturn:
			v.ip++
			var retVal tengo.Object
			if int(v.curInsts[v.ip]) == 1 {
				retVal = v.stack[v.sp-1]
			} else {
				retVal = tengo.UndefinedValue
			}
			//v.sp--

			v.framesIndex--
			v.curFrame = &v.frames[v.framesIndex-1]
			v.curInsts = v.curFrame.fn.Instructions
			v.ip = v.curFrame.ip

			//v.sp = lastFrame.basePointer - 1
			v.sp = v.frames[v.framesIndex].basePointer

			// skip stack overflow check because (newSP) <= (oldSP)
			v.stack[v.sp-1] = retVal
			//v.sp++

		case compiler.OpDefineLocal:
			v.ip++
			localIndex := int(v.curInsts[v.ip])

			sp := v.curFrame.basePointer + localIndex

			// local variables can be mutated by other actions
			// so always store the copy of popped value
			val := v.stack[v.sp-1]
			v.sp--

			v.stack[sp] = val

		case compiler.OpSetLocal:
			localIndex := int(v.curInsts[v.ip+1])
			v.ip++

			sp := v.curFrame.basePointer + localIndex

			// update pointee of v.stack[sp] instead of replacing the pointer itself.
			// this is needed because there can be free variables referencing the same local variables.
			val := v.stack[v.sp-1]
			v.sp--

			if obj, ok := v.stack[sp].(*tengo.ObjectPtr); ok {
				*obj.Value = val
				val = obj
			}
			v.stack[sp] = val // also use a copy of popped value

		case compiler.OpSetSelLocal:
			localIndex := int(v.curInsts[v.ip+1])
			numSelectors := int(v.curInsts[v.ip+2])
			v.ip += 2

			// selectors and RHS value
			selectors := make([]tengo.Object, numSelectors)
			for i := 0; i < numSelectors; i++ {
				selectors[i] = v.stack[v.sp-numSelectors+i]
			}

			val := v.stack[v.sp-numSelectors-1]
			v.sp -= numSelectors + 1

			dst := v.stack[v.curFrame.basePointer+localIndex]
			if obj, ok := dst.(*tengo.ObjectPtr); ok {
				dst = *obj.Value
			}

			if err := indexAssign(dst, val, selectors); err != nil {
				return err
			}

		case compiler.OpGetLocal:
			v.ip++
			localIndex := int(v.curInsts[v.ip])

			val := v.stack[v.curFrame.basePointer+localIndex]

			if obj, ok := val.(*tengo.ObjectPtr); ok {
				val = *obj.Value
			}

			v.stack[v.sp] = val
			v.sp++

		case compiler.OpGetBuiltin:
			v.ip++
			builtinIndex := int(v.curInsts[v.ip])

			v.stack[v.sp] = tengo.Builtins[builtinIndex]
			v.sp++

		case compiler.OpClosure:
			v.ip += 3
			constIndex := int(v.curInsts[v.ip-1]) | int(v.curInsts[v.ip-2])<<8
			numFree := int(v.curInsts[v.ip])

			fn, ok := v.constants[constIndex].(*tengo.CompiledFunction)
			if !ok {
				return fmt.Errorf("not function: %s", fn.TypeName())
			}

			free := make([]*tengo.ObjectPtr, numFree)
			for i := 0; i < numFree; i++ {
				switch freeVar := (v.stack[v.sp-numFree+i]).(type) {
				case *tengo.ObjectPtr:
					free[i] = freeVar
				default:
					free[i] = &tengo.ObjectPtr{Value: &v.stack[v.sp-numFree+i]}
				}
			}

			v.sp -= numFree

			cl := &tengo.CompiledFunction{
				Instructions:  fn.Instructions,
				NumLocals:     fn.NumLocals,
				NumParameters: fn.NumParameters,
				VarArgs:       fn.VarArgs,
				Free:          free,
			}

			v.allocs--
			if v.allocs == 0 {
				return ErrObjectAllocLimit
			}

			v.stack[v.sp] = cl
			v.sp++

		case compiler.OpGetFreePtr:
			v.ip++
			freeIndex := int(v.curInsts[v.ip])

			val := v.curFrame.freeVars[freeIndex]

			v.stack[v.sp] = val
			v.sp++

		case compiler.OpGetFree:
			v.ip++
			freeIndex := int(v.curInsts[v.ip])

			val := *v.curFrame.freeVars[freeIndex].Value

			v.stack[v.sp] = val
			v.sp++

		case compiler.OpSetFree:
			v.ip++
			freeIndex := int(v.curInsts[v.ip])

			*v.curFrame.freeVars[freeIndex].Value = v.stack[v.sp-1]

			v.sp--

		case compiler.OpGetLocalPtr:
			v.ip++
			localIndex := int(v.curInsts[v.ip])

			sp := v.curFrame.basePointer + localIndex
			val := v.stack[sp]

			var freeVar *tengo.ObjectPtr
			if obj, ok := val.(*tengo.ObjectPtr); ok {
				freeVar = obj
			} else {
				freeVar = &tengo.ObjectPtr{Value: &val}
				v.stack[sp] = freeVar
			}

			v.stack[v.sp] = freeVar
			v.sp++

		case compiler.OpSetSelFree:
			v.ip += 2
			freeIndex := int(v.curInsts[v.ip-1])
			numSelectors := int(v.curInsts[v.ip])

			// selectors and RHS value
			selectors := make([]tengo.Object, numSelectors)
			for i := 0; i < numSelectors; i++ {
				selectors[i] = v.stack[v.sp-numSelectors+i]
			}
			val := v.stack[v.sp-numSelectors-1]
			v.sp -= numSelectors + 1

			if err := indexAssign(*v.curFrame.freeVars[freeIndex].Value, val, selectors); err != nil {
				return err
			}

		case compiler.OpIteratorInit:
			var iterator tengo.Object

			dst := v.stack[v.sp-1]
			v.sp--

			if !dst.CanIterate() {
				return fmt.Errorf("not iterable: %s", dst.TypeName())
			}

			iterator = dst.Iterate()
			v.allocs--
			if v.allocs == 0 {
				return ErrObjectAllocLimit
			}

			v.stack[v.sp] = iterator
			v.sp++

		case compiler.OpIteratorNext:
			iterator := v.stack[v.sp-1]
			v.sp--

			hasMore := iterator.(tengo.Iterator).Next()

			if hasMore {
				v.stack[v.sp] = tengo.TrueValue
			} else {
				v.stack[v.sp] = tengo.FalseValue
			}
			v.sp++

		case compiler.OpIteratorKey:
			iterator := v.stack[v.sp-1]
			v.sp--

			val := iterator.(tengo.Iterator).Key()

			v.stack[v.sp] = val
			v.sp++

		case compiler.OpIteratorValue:
			iterator := v.stack[v.sp-1]
			v.sp--

			val := iterator.(tengo.Iterator).Value()

			v.stack[v.sp] = val
			v.sp++

		default:
			return fmt.Errorf("unknown opcode: %d", v.curInsts[v.ip])
		}
	}

	return nil
}

// IsStackEmpty tests if the stack is empty or not.
func (v *VM) IsStackEmpty() bool {
	return v.sp == 0
}

func indexAssign(dst, src tengo.Object, selectors []tengo.Object) error {
	numSel := len(selectors)

	for sidx := numSel - 1; sidx > 0; sidx-- {
		next, err := dst.IndexGet(selectors[sidx])
		if err != nil {
			if err == tengo.ErrNotIndexable {
				return fmt.Errorf("not indexable: %s", dst.TypeName())
			}

			if err == tengo.ErrInvalidIndexType {
				return fmt.Errorf("invalid index type: %s", selectors[sidx].TypeName())
			}

			return err
		}

		dst = next
	}

	if err := dst.IndexSet(selectors[0], src); err != nil {
		if err == tengo.ErrNotIndexAssignable {
			return fmt.Errorf("not index-assignable: %s", dst.TypeName())
		}

		if err == tengo.ErrInvalidIndexValueType {
			return fmt.Errorf("invaid index value type: %s", src.TypeName())
		}

		return err
	}

	return nil
}

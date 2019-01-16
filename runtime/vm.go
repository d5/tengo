package runtime

import (
	"errors"
	"fmt"

	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/compiler/token"
	"github.com/d5/tengo/objects"
)

const (
	// StackSize is the maximum stack size.
	StackSize = 2048

	// GlobalsSize is the maximum number of global variables.
	GlobalsSize = 1024

	// MaxFrames is the maximum number of function frames.
	MaxFrames = 1024
)

var (
	truePtr      = &objects.TrueValue
	falsePtr     = &objects.FalseValue
	undefinedPtr = &objects.UndefinedValue
	builtinFuncs []objects.Object
)

// VM is a virtual machine that executes the bytecode compiled by Compiler.
type VM struct {
	constants   []objects.Object
	stack       []*objects.Object
	sp          int
	globals     []*objects.Object
	frames      []Frame
	framesIndex int
	curFrame    *Frame
	curInsts    []byte
	curIPLimit  int
	aborting    bool
}

// NewVM creates a VM.
func NewVM(bytecode *compiler.Bytecode, globals []*objects.Object) *VM {
	if globals == nil {
		globals = make([]*objects.Object, GlobalsSize)
	} else if len(globals) < GlobalsSize {
		g := make([]*objects.Object, GlobalsSize)
		copy(g, globals)
		globals = g
	}

	frames := make([]Frame, MaxFrames)
	frames[0].fn = &objects.CompiledFunction{Instructions: bytecode.Instructions}
	frames[0].freeVars = nil
	frames[0].ip = -1
	frames[0].basePointer = 0

	return &VM{
		constants:   bytecode.Constants,
		stack:       make([]*objects.Object, StackSize),
		sp:          0,
		globals:     globals,
		frames:      frames,
		framesIndex: 1,
		curFrame:    &(frames[0]),
		curInsts:    frames[0].fn.Instructions,
		curIPLimit:  len(frames[0].fn.Instructions) - 1,
	}
}

// Abort aborts the execution.
func (v *VM) Abort() {
	v.aborting = true
}

// Run starts the execution.
func (v *VM) Run() error {
	var ip int

	for v.curFrame.ip < v.curIPLimit && !v.aborting {
		v.curFrame.ip++

		ip = v.curFrame.ip

		switch compiler.Opcode(v.curInsts[ip]) {
		case compiler.OpConstant:
			cidx := compiler.ReadUint16(v.curInsts[ip+1:])
			v.curFrame.ip += 2

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &v.constants[cidx]
			v.sp++

		case compiler.OpNull:
			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = undefinedPtr
			v.sp++

		case compiler.OpAdd:
			right := v.stack[v.sp-1]
			left := v.stack[v.sp-2]
			v.sp -= 2

			res, err := (*left).BinaryOp(token.Add, *right)
			if err != nil {
				return err
			}

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &res
			v.sp++

		case compiler.OpSub:
			right := v.stack[v.sp-1]
			left := v.stack[v.sp-2]
			v.sp -= 2

			res, err := (*left).BinaryOp(token.Sub, *right)
			if err != nil {
				return err
			}

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &res
			v.sp++

		case compiler.OpMul:
			right := v.stack[v.sp-1]
			left := v.stack[v.sp-2]
			v.sp -= 2

			res, err := (*left).BinaryOp(token.Mul, *right)
			if err != nil {
				return err
			}

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &res
			v.sp++

		case compiler.OpDiv:
			right := v.stack[v.sp-1]
			left := v.stack[v.sp-2]
			v.sp -= 2

			res, err := (*left).BinaryOp(token.Quo, *right)
			if err != nil {
				return err
			}

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &res
			v.sp++

		case compiler.OpRem:
			right := v.stack[v.sp-1]
			left := v.stack[v.sp-2]
			v.sp -= 2

			res, err := (*left).BinaryOp(token.Rem, *right)
			if err != nil {
				return err
			}

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &res
			v.sp++

		case compiler.OpBAnd:
			right := v.stack[v.sp-1]
			left := v.stack[v.sp-2]
			v.sp -= 2

			res, err := (*left).BinaryOp(token.And, *right)
			if err != nil {
				return err
			}

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &res
			v.sp++

		case compiler.OpBOr:
			right := v.stack[v.sp-1]
			left := v.stack[v.sp-2]
			v.sp -= 2

			res, err := (*left).BinaryOp(token.Or, *right)
			if err != nil {
				return err
			}

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &res
			v.sp++

		case compiler.OpBXor:
			right := v.stack[v.sp-1]
			left := v.stack[v.sp-2]
			v.sp -= 2

			res, err := (*left).BinaryOp(token.Xor, *right)
			if err != nil {
				return err
			}

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &res
			v.sp++

		case compiler.OpBAndNot:
			right := v.stack[v.sp-1]
			left := v.stack[v.sp-2]
			v.sp -= 2

			res, err := (*left).BinaryOp(token.AndNot, *right)
			if err != nil {
				return err
			}

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &res
			v.sp++

		case compiler.OpBShiftLeft:
			right := v.stack[v.sp-1]
			left := v.stack[v.sp-2]
			v.sp -= 2

			res, err := (*left).BinaryOp(token.Shl, *right)
			if err != nil {
				return err
			}

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &res
			v.sp++

		case compiler.OpBShiftRight:
			right := v.stack[v.sp-1]
			left := v.stack[v.sp-2]
			v.sp -= 2

			res, err := (*left).BinaryOp(token.Shr, *right)
			if err != nil {
				return err
			}

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &res
			v.sp++

		case compiler.OpEqual:
			right := v.stack[v.sp-1]
			left := v.stack[v.sp-2]
			v.sp -= 2

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			if (*right).Equals(*left) {
				v.stack[v.sp] = truePtr
			} else {
				v.stack[v.sp] = falsePtr
			}
			v.sp++

		case compiler.OpNotEqual:
			right := v.stack[v.sp-1]
			left := v.stack[v.sp-2]
			v.sp -= 2

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			if (*right).Equals(*left) {
				v.stack[v.sp] = falsePtr
			} else {
				v.stack[v.sp] = truePtr
			}
			v.sp++

		case compiler.OpGreaterThan:
			right := v.stack[v.sp-1]
			left := v.stack[v.sp-2]
			v.sp -= 2

			res, err := (*left).BinaryOp(token.Greater, *right)
			if err != nil {
				return err
			}

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &res
			v.sp++

		case compiler.OpGreaterThanEqual:
			right := v.stack[v.sp-1]
			left := v.stack[v.sp-2]
			v.sp -= 2

			res, err := (*left).BinaryOp(token.GreaterEq, *right)
			if err != nil {
				return err
			}

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &res
			v.sp++

		case compiler.OpPop:
			v.sp--

		case compiler.OpTrue:
			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = truePtr
			v.sp++

		case compiler.OpFalse:
			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = falsePtr
			v.sp++

		case compiler.OpLNot:
			operand := v.stack[v.sp-1]
			v.sp--

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			if (*operand).IsFalsy() {
				v.stack[v.sp] = truePtr
			} else {
				v.stack[v.sp] = falsePtr
			}
			v.sp++

		case compiler.OpBComplement:
			operand := v.stack[v.sp-1]
			v.sp--

			switch x := (*operand).(type) {
			case *objects.Int:
				if v.sp >= StackSize {
					return ErrStackOverflow
				}

				var res objects.Object = &objects.Int{Value: ^x.Value}

				v.stack[v.sp] = &res
				v.sp++
			default:
				return fmt.Errorf("invalid operation on %s", (*operand).TypeName())
			}

		case compiler.OpMinus:
			operand := v.stack[v.sp-1]
			v.sp--

			switch x := (*operand).(type) {
			case *objects.Int:
				if v.sp >= StackSize {
					return ErrStackOverflow
				}

				var res objects.Object = &objects.Int{Value: -x.Value}

				v.stack[v.sp] = &res
				v.sp++
			case *objects.Float:
				if v.sp >= StackSize {
					return ErrStackOverflow
				}

				var res objects.Object = &objects.Float{Value: -x.Value}

				v.stack[v.sp] = &res
				v.sp++
			default:
				return fmt.Errorf("invalid operation on %s", (*operand).TypeName())
			}

		case compiler.OpJumpFalsy:
			pos := int(compiler.ReadUint16(v.curInsts[ip+1:]))
			v.curFrame.ip += 2

			condition := v.stack[v.sp-1]
			v.sp--

			if (*condition).IsFalsy() {
				v.curFrame.ip = pos - 1
			}

		case compiler.OpAndJump:
			pos := int(compiler.ReadUint16(v.curInsts[ip+1:]))
			v.curFrame.ip += 2

			condition := *v.stack[v.sp-1]
			if condition.IsFalsy() {
				v.curFrame.ip = pos - 1
			} else {
				v.sp--
			}

		case compiler.OpOrJump:
			pos := int(compiler.ReadUint16(v.curInsts[ip+1:]))
			v.curFrame.ip += 2

			condition := *v.stack[v.sp-1]
			if !condition.IsFalsy() {
				v.curFrame.ip = pos - 1
			} else {
				v.sp--
			}

		case compiler.OpJump:
			pos := int(compiler.ReadUint16(v.curInsts[ip+1:]))
			v.curFrame.ip = pos - 1

		case compiler.OpSetGlobal:
			globalIndex := compiler.ReadUint16(v.curInsts[ip+1:])
			v.curFrame.ip += 2

			v.sp--

			v.globals[globalIndex] = v.stack[v.sp]

		case compiler.OpSetSelGlobal:
			globalIndex := compiler.ReadUint16(v.curInsts[ip+1:])
			numSelectors := int(compiler.ReadUint8(v.curInsts[ip+3:]))
			v.curFrame.ip += 3

			// pop selector outcomes (left to right)
			selectors := make([]interface{}, numSelectors, numSelectors)
			for i := 0; i < numSelectors; i++ {
				sel := v.stack[v.sp-1]
				v.sp--

				switch sel := (*sel).(type) {
				case *objects.String: // map key
					selectors[i] = sel.Value
				case *objects.Int: // array index
					selectors[i] = int(sel.Value)
				default:
					return fmt.Errorf("invalid selector type: %s", sel.TypeName())
				}
			}

			// RHS value
			val := v.stack[v.sp-1]
			v.sp--

			if err := selectorAssign(v.globals[globalIndex], val, selectors); err != nil {
				return err
			}

		case compiler.OpGetGlobal:
			globalIndex := compiler.ReadUint16(v.curInsts[ip+1:])
			v.curFrame.ip += 2

			val := v.globals[globalIndex]

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = val
			v.sp++

		case compiler.OpArray:
			numElements := int(compiler.ReadUint16(v.curInsts[ip+1:]))
			v.curFrame.ip += 2

			var elements []objects.Object
			for i := v.sp - numElements; i < v.sp; i++ {
				elements = append(elements, *v.stack[i])
			}
			v.sp -= numElements

			var arr objects.Object = &objects.Array{Value: elements}

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &arr
			v.sp++

		case compiler.OpMap:
			numElements := int(compiler.ReadUint16(v.curInsts[ip+1:]))
			v.curFrame.ip += 2

			kv := make(map[string]objects.Object)
			for i := v.sp - numElements; i < v.sp; i += 2 {
				key := *v.stack[i]
				value := *v.stack[i+1]
				kv[key.(*objects.String).Value] = value
			}
			v.sp -= numElements

			var m objects.Object = &objects.Map{Value: kv}

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &m
			v.sp++

		case compiler.OpError:
			value := v.stack[v.sp-1]
			v.sp--

			var err objects.Object = &objects.Error{
				Value: *value,
			}

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &err
			v.sp++

		case compiler.OpIndex:
			index := v.stack[v.sp-1]
			left := v.stack[v.sp-2]
			v.sp -= 2

			switch left := (*left).(type) {
			case *objects.Array:
				idx, ok := (*index).(*objects.Int)
				if !ok {
					return fmt.Errorf("non-integer array index: %s", left.TypeName())
				}

				if idx.Value < 0 || idx.Value >= int64(len(left.Value)) {
					return fmt.Errorf("index out of bounds: %d", index)
				}

				if v.sp >= StackSize {
					return ErrStackOverflow
				}

				v.stack[v.sp] = &left.Value[idx.Value]
				v.sp++

			case *objects.String:
				idx, ok := (*index).(*objects.Int)
				if !ok {
					return fmt.Errorf("non-integer array index: %s", left.TypeName())
				}

				str := []rune(left.Value)

				if idx.Value < 0 || idx.Value >= int64(len(str)) {
					return fmt.Errorf("index out of bounds: %d", index)
				}

				if v.sp >= StackSize {
					return ErrStackOverflow
				}

				var val objects.Object = &objects.Char{Value: str[idx.Value]}

				v.stack[v.sp] = &val
				v.sp++

			case *objects.Map:
				key, ok := (*index).(*objects.String)
				if !ok {
					return fmt.Errorf("non-string map key: %s", left.TypeName())
				}

				var res = objects.UndefinedValue
				val, ok := left.Value[key.Value]
				if ok {
					res = val
				}

				if v.sp >= StackSize {
					return ErrStackOverflow
				}

				v.stack[v.sp] = &res
				v.sp++

			case *objects.Error: // err.value
				key, ok := (*index).(*objects.String)
				if !ok || key.Value != "value" {
					return errors.New("invalid selector on error")
				}

				if v.sp >= StackSize {
					return ErrStackOverflow
				}

				v.stack[v.sp] = &left.Value
				v.sp++

			default:
				return fmt.Errorf("type %s does not support indexing", left.TypeName())
			}

		case compiler.OpSliceIndex:
			high := v.stack[v.sp-1]
			low := v.stack[v.sp-2]
			left := v.stack[v.sp-3]
			v.sp -= 3

			var lowIdx, highIdx int64

			switch low := (*low).(type) {
			case *objects.Undefined:
				//lowIdx = 0
			case *objects.Int:
				lowIdx = low.Value
			default:
				return fmt.Errorf("non-integer slice index: %s", low.TypeName())
			}

			switch high := (*high).(type) {
			case *objects.Undefined:
				highIdx = -1 // will be replaced by number of elements
			case *objects.Int:
				highIdx = high.Value
			default:
				return fmt.Errorf("non-integer slice index: %s", high.TypeName())
			}

			switch left := (*left).(type) {
			case *objects.Array:
				numElements := int64(len(left.Value))

				if lowIdx < 0 || lowIdx >= numElements {
					return fmt.Errorf("index out of bounds: %d", lowIdx)
				}
				if highIdx < 0 {
					highIdx = numElements
				} else if highIdx < 0 || highIdx > numElements {
					return fmt.Errorf("index out of bounds: %d", highIdx)
				}

				if lowIdx > highIdx {
					return fmt.Errorf("invalid slice index: %d > %d", lowIdx, highIdx)
				}

				if v.sp >= StackSize {
					return ErrStackOverflow
				}

				var val objects.Object = &objects.Array{Value: left.Value[lowIdx:highIdx]}

				v.stack[v.sp] = &val
				v.sp++

			case *objects.String:
				numElements := int64(len(left.Value))

				if lowIdx < 0 || lowIdx >= numElements {
					return fmt.Errorf("index out of bounds: %d", lowIdx)
				}
				if highIdx < 0 {
					highIdx = numElements
				} else if highIdx < 0 || highIdx > numElements {
					return fmt.Errorf("index out of bounds: %d", highIdx)
				}

				if lowIdx > highIdx {
					return fmt.Errorf("invalid slice index: %d > %d", lowIdx, highIdx)
				}

				if v.sp >= StackSize {
					return ErrStackOverflow
				}

				var val objects.Object = &objects.String{Value: left.Value[lowIdx:highIdx]}

				v.stack[v.sp] = &val
				v.sp++

			default:
				return fmt.Errorf("cannot slice %s", left.TypeName())
			}

		case compiler.OpCall:
			numArgs := int(compiler.ReadUint8(v.curInsts[ip+1:]))
			v.curFrame.ip++

			callee := *v.stack[v.sp-1-numArgs]

			switch callee := callee.(type) {
			case *objects.Closure:
				if err := v.callFunction(callee.Fn, callee.Free, numArgs); err != nil {
					return err
				}
			case *objects.CompiledFunction:
				if err := v.callFunction(callee, nil, numArgs); err != nil {
					return err
				}
			case objects.Callable:
				var args []objects.Object
				for _, arg := range v.stack[v.sp-numArgs : v.sp] {
					args = append(args, *arg)
				}

				ret, err := callee.Call(args...)
				v.sp -= numArgs + 1

				// runtime error
				if err != nil {
					return err
				}

				// nil return -> undefined
				if ret == nil {
					ret = objects.UndefinedValue
				}

				if v.sp >= StackSize {
					return ErrStackOverflow
				}

				v.stack[v.sp] = &ret
				v.sp++
			default:
				return fmt.Errorf("calling non-function: %s", callee.TypeName())
			}

		case compiler.OpReturnValue:
			//numRets := int(compiler.ReadUint8(v.curInsts[ip+1:]))
			_ = int(compiler.ReadUint8(v.curInsts[ip+1:]))
			v.curFrame.ip++

			// TODO: multi-value return is not fully implemented yet
			//var rets []*objects.Object
			//for i := 0; i < numRets; i++ {
			//	val := v.pop()
			//	rets = append(rets, val)
			//}
			retVal := v.stack[v.sp-1]
			//v.sp--

			v.framesIndex--
			lastFrame := v.frames[v.framesIndex]
			v.curFrame = &v.frames[v.framesIndex-1]
			v.curInsts = v.curFrame.fn.Instructions
			v.curIPLimit = len(v.curInsts) - 1

			//v.sp = lastFrame.basePointer - 1
			v.sp = lastFrame.basePointer

			//for _, retVal := range rets {
			//	if err := v.push(retVal); err != nil {
			//		return err
			//	}
			//}
			if v.sp-1 >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp-1] = retVal
			//v.sp++

		case compiler.OpReturn:
			v.framesIndex--
			lastFrame := v.frames[v.framesIndex]
			v.curFrame = &v.frames[v.framesIndex-1]
			v.curInsts = v.curFrame.fn.Instructions
			v.curIPLimit = len(v.curInsts) - 1

			v.sp = lastFrame.basePointer - 1

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = undefinedPtr
			v.sp++

		case compiler.OpDefineLocal:
			localIndex := compiler.ReadUint8(v.curInsts[ip+1:])
			v.curFrame.ip++

			sp := v.curFrame.basePointer + int(localIndex)

			// local variables can be mutated by other actions
			// so always store the copy of popped value
			val := *v.stack[v.sp-1]
			v.sp--

			v.stack[sp] = &val

		case compiler.OpSetLocal:
			localIndex := compiler.ReadUint8(v.curInsts[ip+1:])
			v.curFrame.ip++

			sp := v.curFrame.basePointer + int(localIndex)

			// update pointee of v.stack[sp] instead of replacing the pointer itself.
			// this is needed because there can be free variables referencing the same local variables.
			val := v.stack[v.sp-1]
			v.sp--

			*v.stack[sp] = *val // also use a copy of popped value

		case compiler.OpSetSelLocal:
			localIndex := compiler.ReadUint8(v.curInsts[ip+1:])
			numSelectors := int(compiler.ReadUint8(v.curInsts[ip+2:]))
			v.curFrame.ip += 2

			// pop selector outcomes (left to right)
			selectors := make([]interface{}, numSelectors, numSelectors)
			for i := 0; i < numSelectors; i++ {
				sel := v.stack[v.sp-1]
				v.sp--

				switch sel := (*sel).(type) {
				case *objects.String: // map key
					selectors[i] = sel.Value
				case *objects.Int: // array index
					selectors[i] = int(sel.Value)
				default:
					return fmt.Errorf("invalid selector type: %s", sel.TypeName())
				}
			}

			// RHS value
			val := v.stack[v.sp-1] // no need to copy value here; selectorAssign uses copy of value
			v.sp--

			sp := v.curFrame.basePointer + int(localIndex)

			if err := selectorAssign(v.stack[sp], val, selectors); err != nil {
				return err
			}

		case compiler.OpGetLocal:
			localIndex := compiler.ReadUint8(v.curInsts[ip+1:])
			v.curFrame.ip++

			val := v.stack[v.curFrame.basePointer+int(localIndex)]

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = val
			v.sp++

		case compiler.OpGetBuiltin:
			builtinIndex := compiler.ReadUint8(v.curInsts[ip+1:])
			v.curFrame.ip++

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &builtinFuncs[builtinIndex]
			v.sp++

		case compiler.OpClosure:
			constIndex := compiler.ReadUint16(v.curInsts[ip+1:])
			numFree := compiler.ReadUint8(v.curInsts[ip+3:])
			v.curFrame.ip += 3

			if err := v.pushClosure(int(constIndex), int(numFree)); err != nil {
				return err
			}

		case compiler.OpGetFree:
			freeIndex := compiler.ReadUint8(v.curInsts[ip+1:])
			v.curFrame.ip++

			val := v.curFrame.freeVars[freeIndex]

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = val
			v.sp++

		case compiler.OpSetSelFree:
			freeIndex := compiler.ReadUint8(v.curInsts[ip+1:])
			numSelectors := int(compiler.ReadUint8(v.curInsts[ip+2:]))
			v.curFrame.ip += 2

			// pop selector outcomes (left to right)
			selectors := make([]interface{}, numSelectors, numSelectors)
			for i := 0; i < numSelectors; i++ {
				sel := v.stack[v.sp-1]
				v.sp--

				switch sel := (*sel).(type) {
				case *objects.String: // map key
					selectors[i] = sel.Value
				case *objects.Int: // array index
					selectors[i] = int(sel.Value)
				default:
					return fmt.Errorf("invalid selector type: %s", sel.TypeName())
				}
			}

			// RHS value
			val := v.stack[v.sp-1]
			v.sp--

			if err := selectorAssign(v.curFrame.freeVars[freeIndex], val, selectors); err != nil {
				return err
			}

		case compiler.OpSetFree:
			freeIndex := compiler.ReadUint8(v.curInsts[ip+1:])
			v.curFrame.ip++

			val := v.stack[v.sp-1]
			v.sp--

			*v.curFrame.freeVars[freeIndex] = *val

		case compiler.OpIteratorInit:
			var iterator objects.Object

			dst := v.stack[v.sp-1]
			v.sp--

			switch dst := (*dst).(type) {
			case *objects.Array:
				iterator = objects.NewArrayIterator(dst)
			case *objects.Map:
				iterator = objects.NewMapIterator(dst)
			case *objects.String:
				iterator = objects.NewStringIterator(dst)
			default:
				return fmt.Errorf("non-iterable type: %s", dst.TypeName())
			}

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &iterator
			v.sp++

		case compiler.OpIteratorNext:
			iterator := v.stack[v.sp-1]
			v.sp--

			b := (*iterator).(objects.Iterator).Next()
			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			if b {
				v.stack[v.sp] = truePtr
			} else {
				v.stack[v.sp] = falsePtr
			}
			v.sp++

		case compiler.OpIteratorKey:
			iterator := v.stack[v.sp-1]
			v.sp--

			val := (*iterator).(objects.Iterator).Key()

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &val
			v.sp++

		case compiler.OpIteratorValue:
			iterator := v.stack[v.sp-1]
			v.sp--

			val := (*iterator).(objects.Iterator).Value()

			if v.sp >= StackSize {
				return ErrStackOverflow
			}

			v.stack[v.sp] = &val
			v.sp++

		default:
			return fmt.Errorf("unknown opcode: %d", v.curInsts[ip])
		}
	}

	// check if stack still has some objects left
	if v.sp > 0 && !v.aborting {
		return fmt.Errorf("non empty stack after execution")
	}

	return nil
}

// Globals returns the global variables.
func (v *VM) Globals() []*objects.Object {
	return v.globals
}

// FrameInfo returns the current function call frame information.
func (v *VM) FrameInfo() (frameIndex int, ip int) {
	return v.framesIndex - 1, v.frames[v.framesIndex-1].ip
}

func (v *VM) pushClosure(constIndex, numFree int) error {
	c := v.constants[constIndex]

	fn, ok := c.(*objects.CompiledFunction)
	if !ok {
		return fmt.Errorf("not a function: %s", fn.TypeName())
	}

	free := make([]*objects.Object, numFree)
	for i := 0; i < numFree; i++ {
		free[i] = v.stack[v.sp-numFree+i]
	}
	v.sp -= numFree

	if v.sp >= StackSize {
		return ErrStackOverflow
	}

	var cl objects.Object = &objects.Closure{
		Fn:   fn,
		Free: free,
	}

	v.stack[v.sp] = &cl
	v.sp++

	return nil
}

func (v *VM) callFunction(fn *objects.CompiledFunction, freeVars []*objects.Object, numArgs int) error {
	if numArgs != fn.NumParameters {
		return fmt.Errorf("wrong number of arguments: want=%d, got=%d",
			fn.NumParameters, numArgs)
	}

	// check if this is a tail-call (recursive call right before return)
	if fn == v.curFrame.fn { // recursion
		nextOp := compiler.Opcode(v.curInsts[v.curFrame.ip+1])
		if nextOp == compiler.OpReturnValue || // tail call
			(nextOp == compiler.OpPop &&
				compiler.OpReturn == compiler.Opcode(v.curInsts[v.curFrame.ip+2])) {

			//  stack before tail-call
			//
			//  |--------|
			//  |        | <- SP  current
			//  |--------|
			//  | *ARG2  |        for next function (tail-call)
			//  |--------|
			//  | *ARG1  |        for next function (tail-call)
			//  |--------|
			//  |  FUNC  |        function itself
			//  |--------|
			//  | LOCAL3 |        for current function
			//  |--------|
			//  | LOCAL2 |        for current function
			//  |--------|
			//  |  ARG2  |        for current function
			//  |--------|
			//  |  ARG1  | <- BP  for current function
			//  |--------|

			for p := 0; p < numArgs; p++ {
				v.stack[v.curFrame.basePointer+p] = v.stack[v.sp-numArgs+p]
			}
			v.sp -= numArgs + 1
			v.curFrame.ip = -1 // reset IP to beginning of the frame

			//  stack after tail-call
			//
			//  |--------|
			//  |        |
			//  |--------|
			//  | *ARG2  |
			//  |--------|
			//  | *ARG1  |
			//  |--------|
			//  |  FUNC  | <- SP  current
			//  |--------|
			//  | LOCAL3 |        for current function
			//  |--------|
			//  | LOCAL2 |        for current function
			//  |--------|
			//  | *ARG2  |        (copied)
			//  |--------|
			//  | *ARG1  | <- BP  (copied)
			//  |--------|

			return nil
		}
	}

	v.curFrame = &(v.frames[v.framesIndex])
	v.curFrame.fn = fn
	v.curFrame.freeVars = freeVars
	v.curFrame.ip = -1
	v.curFrame.basePointer = v.sp - numArgs
	v.curInsts = fn.Instructions
	v.curIPLimit = len(v.curInsts) - 1
	v.framesIndex++

	v.sp = v.sp - numArgs + fn.NumLocals

	//  stack after the function call
	//
	//  |--------|
	//  |        |        <- SP after function call
	//  |--------|
	//  | LOCAL4 | (BP+3)
	//  |--------|
	//  | LOCAL3 | (BP+2) <- SP before function call
	//  |--------|
	//  |  ARG2  | (BP+1)
	//  |--------|
	//  |  ARG1  | (BP+0) <- BP
	//  |--------|

	return nil
}

func selectorAssign(dst, src *objects.Object, selectors []interface{}) error {
	numSel := len(selectors)

	for idx := 0; idx < numSel; idx++ {
		switch sel := selectors[idx].(type) {
		case string:
			m, isMap := (*dst).(*objects.Map)
			if !isMap {
				return fmt.Errorf("invalid map object for selector '%s'", sel)
			}

			if idx == numSel-1 {
				m.Set(sel, *src)
				return nil
			}

			nxt, found := m.Get(sel)
			if !found {
				return fmt.Errorf("key not found '%s'", sel)
			}

			dst = &nxt
		case int:
			arr, isArray := (*dst).(*objects.Array)
			if !isArray {
				return fmt.Errorf("invalid array object for select '[%d]'", sel)
			}

			if idx == numSel-1 {
				return arr.Set(sel, *src)
			}

			nxt, err := arr.Get(sel)
			if err != nil {
				return err
			}

			dst = &nxt
		default:
			panic(fmt.Errorf("invalid selector term: %T", sel))
		}
	}

	return nil
}

func init() {
	builtinFuncs = make([]objects.Object, len(objects.Builtins))
	for i, b := range objects.Builtins {
		builtinFuncs[i] = &objects.BuiltinFunction{Value: b.Func}
	}
}

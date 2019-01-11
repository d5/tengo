package vm

import (
	"fmt"

	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/token"
)

const (
	StackSize   = 2048
	GlobalsSize = 1024
	MaxFrames   = 1024
)

var (
	trueObj      objects.Object = &objects.Bool{Value: true}
	falseObj     objects.Object = &objects.Bool{Value: false}
	undefinedObj objects.Object = &objects.Undefined{}
	builtinFuncs []objects.Object
)

type VM struct {
	constants   []objects.Object
	stack       []*objects.Object
	sp          int
	globals     []*objects.Object
	frames      []Frame
	framesIndex int
	aborting    bool
}

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
	}
}

func (v *VM) Abort() {
	v.aborting = true
}

func (v *VM) Run() error {
	var ip int
	var ins []byte
	var op compiler.Opcode

	for v.frames[v.framesIndex-1].ip < len(v.frames[v.framesIndex-1].fn.Instructions)-1 && !v.aborting {
		curFrame := &(v.frames[v.framesIndex-1])
		curFrame.ip++

		ip = curFrame.ip
		ins = curFrame.fn.Instructions
		op = compiler.Opcode(ins[ip])

		switch op {
		case compiler.OpConstant:
			cidx := compiler.ReadUint16(ins[ip+1:])
			curFrame.ip += 2

			if err := v.push(&v.constants[cidx]); err != nil {
				return err
			}

		case compiler.OpNull:
			if err := v.push(&undefinedObj); err != nil {
				return err
			}
		case compiler.OpAdd:
			right := v.pop()
			left := v.pop()

			res, err := (*left).BinaryOp(token.Add, *right)
			if err != nil {
				return err
			}

			if err := v.push(&res); err != nil {
				return err
			}
		case compiler.OpSub:
			right := v.pop()
			left := v.pop()

			res, err := (*left).BinaryOp(token.Sub, *right)
			if err != nil {
				return err
			}

			if err := v.push(&res); err != nil {
				return err
			}
		case compiler.OpMul:
			right := v.pop()
			left := v.pop()

			res, err := (*left).BinaryOp(token.Mul, *right)
			if err != nil {
				return err
			}

			if err := v.push(&res); err != nil {
				return err
			}
		case compiler.OpDiv:
			right := v.pop()
			left := v.pop()

			res, err := (*left).BinaryOp(token.Quo, *right)
			if err != nil {
				return err
			}

			if err := v.push(&res); err != nil {
				return err
			}
		case compiler.OpRem:
			right := v.pop()
			left := v.pop()

			res, err := (*left).BinaryOp(token.Rem, *right)
			if err != nil {
				return err
			}

			if err := v.push(&res); err != nil {
				return err
			}
		case compiler.OpBAnd:
			right := v.pop()
			left := v.pop()

			res, err := (*left).BinaryOp(token.And, *right)
			if err != nil {
				return err
			}

			if err := v.push(&res); err != nil {
				return err
			}
		case compiler.OpBOr:
			right := v.pop()
			left := v.pop()

			res, err := (*left).BinaryOp(token.Or, *right)
			if err != nil {
				return err
			}

			if err := v.push(&res); err != nil {
				return err
			}
		case compiler.OpBXor:
			right := v.pop()
			left := v.pop()

			res, err := (*left).BinaryOp(token.Xor, *right)
			if err != nil {
				return err
			}

			if err := v.push(&res); err != nil {
				return err
			}
		case compiler.OpBAndNot:
			right := v.pop()
			left := v.pop()

			res, err := (*left).BinaryOp(token.AndNot, *right)
			if err != nil {
				return err
			}

			if err := v.push(&res); err != nil {
				return err
			}
		case compiler.OpBShiftLeft:
			right := v.pop()
			left := v.pop()

			res, err := (*left).BinaryOp(token.Shl, *right)
			if err != nil {
				return err
			}

			if err := v.push(&res); err != nil {
				return err
			}
		case compiler.OpBShiftRight:
			right := v.pop()
			left := v.pop()

			res, err := (*left).BinaryOp(token.Shr, *right)
			if err != nil {
				return err
			}

			if err := v.push(&res); err != nil {
				return err
			}
		case compiler.OpEqual:
			right := v.pop()
			left := v.pop()

			if (*right).Equals(*left) {
				if err := v.push(&trueObj); err != nil {
					return err
				}
			} else {
				if err := v.push(&falseObj); err != nil {
					return err
				}
			}
		case compiler.OpNotEqual:
			right := v.pop()
			left := v.pop()

			if (*right).Equals(*left) {
				if err := v.push(&falseObj); err != nil {
					return err
				}
			} else {
				if err := v.push(&trueObj); err != nil {
					return err
				}
			}
		case compiler.OpGreaterThan:
			right := v.pop()
			left := v.pop()

			res, err := (*left).BinaryOp(token.Greater, *right)
			if err != nil {
				return err
			}

			if err := v.push(&res); err != nil {
				return err
			}
		case compiler.OpGreaterThanEqual:
			right := v.pop()
			left := v.pop()

			res, err := (*left).BinaryOp(token.GreaterEq, *right)
			if err != nil {
				return err
			}

			if err := v.push(&res); err != nil {
				return err
			}
		case compiler.OpPop:
			_ = v.pop()
		case compiler.OpTrue:
			if err := v.push(&trueObj); err != nil {
				return err
			}
		case compiler.OpFalse:
			if err := v.push(&falseObj); err != nil {
				return err
			}
		case compiler.OpLNot:
			operand := v.pop()

			if (*operand).IsFalsy() {
				if err := v.push(&trueObj); err != nil {
					return err
				}
			} else {
				if err := v.push(&falseObj); err != nil {
					return err
				}
			}
		case compiler.OpBComplement:
			operand := v.pop()

			switch x := (*operand).(type) {
			case *objects.Int:
				var res objects.Object = &objects.Int{Value: ^x.Value}
				if err := v.push(&res); err != nil {
					return err
				}
			default:
				return fmt.Errorf("invalid operation on %s", (*operand).TypeName())
			}
		case compiler.OpMinus:
			operand := v.pop()

			switch x := (*operand).(type) {
			case *objects.Int:
				var res objects.Object = &objects.Int{Value: -x.Value}
				if err := v.push(&res); err != nil {
					return err
				}
			case *objects.Float:
				var res objects.Object = &objects.Float{Value: -x.Value}
				if err := v.push(&res); err != nil {
					return err
				}
			default:
				return fmt.Errorf("invalid operation on %s", (*operand).TypeName())
			}
		case compiler.OpJumpFalsy:
			pos := int(compiler.ReadUint16(ins[ip+1:]))
			curFrame.ip += 2

			condition := v.pop()
			if (*condition).IsFalsy() {
				curFrame.ip = pos - 1
			}
		case compiler.OpAndJump:
			pos := int(compiler.ReadUint16(ins[ip+1:]))
			curFrame.ip += 2

			condition := *v.stack[v.sp-1]
			if condition.IsFalsy() {
				curFrame.ip = pos - 1
			} else {
				_ = v.pop()
			}
		case compiler.OpOrJump:
			pos := int(compiler.ReadUint16(ins[ip+1:]))
			curFrame.ip += 2

			condition := *v.stack[v.sp-1]
			if !condition.IsFalsy() {
				curFrame.ip = pos - 1
			} else {
				_ = v.pop()
			}
		case compiler.OpJump:
			pos := int(compiler.ReadUint16(ins[ip+1:]))
			curFrame.ip = pos - 1
		case compiler.OpSetGlobal:
			globalIndex := compiler.ReadUint16(ins[ip+1:])
			curFrame.ip += 2

			val := v.pop()

			v.globals[globalIndex] = val
		case compiler.OpSetSelGlobal:
			globalIndex := compiler.ReadUint16(ins[ip+1:])
			numSelectors := int(compiler.ReadUint8(ins[ip+3:]))
			curFrame.ip += 3

			// pop selector outcomes (left to right)
			selectors := make([]interface{}, numSelectors, numSelectors)
			for i := 0; i < numSelectors; i++ {
				sel := v.pop()

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
			val := v.pop()

			if err := selectorAssign(v.globals[globalIndex], val, selectors); err != nil {
				return err
			}
		case compiler.OpGetGlobal:
			globalIndex := compiler.ReadUint16(ins[ip+1:])
			curFrame.ip += 2

			val := v.globals[globalIndex]

			if err := v.push(val); err != nil {
				return err
			}
		case compiler.OpArray:
			numElements := int(compiler.ReadUint16(ins[ip+1:]))
			curFrame.ip += 2

			var elements []objects.Object
			for i := v.sp - numElements; i < v.sp; i++ {
				elements = append(elements, *v.stack[i])
			}
			v.sp -= numElements

			var arr objects.Object = &objects.Array{Value: elements}

			if err := v.push(&arr); err != nil {
				return err
			}
		case compiler.OpMap:
			numElements := int(compiler.ReadUint16(ins[ip+1:]))
			curFrame.ip += 2

			kv := make(map[string]objects.Object)
			for i := v.sp - numElements; i < v.sp; i += 2 {
				key := *v.stack[i]
				value := *v.stack[i+1]
				kv[key.(*objects.String).Value] = value
			}
			v.sp -= numElements

			var map_ objects.Object = &objects.Map{Value: kv}

			if err := v.push(&map_); err != nil {
				return err
			}
		case compiler.OpIndex:
			index := v.pop()
			left := v.pop()

			switch left := (*left).(type) {
			case *objects.Array:
				idx, ok := (*index).(*objects.Int)
				if !ok {
					return fmt.Errorf("non-integer array index: %s", left.TypeName())
				}

				if err := v.executeArrayIndex(left, idx.Value); err != nil {
					return err
				}
			case *objects.String:
				idx, ok := (*index).(*objects.Int)
				if !ok {
					return fmt.Errorf("non-integer array index: %s", left.TypeName())
				}

				if err := v.executeStringIndex(left, idx.Value); err != nil {
					return err
				}
			case *objects.Map:
				key, ok := (*index).(*objects.String)
				if !ok {
					return fmt.Errorf("non-string map key: %s", left.TypeName())
				}

				if err := v.executeMapIndex(left, key.Value); err != nil {
					return err
				}
			default:
				return fmt.Errorf("type %s does not support indexing", left.TypeName())
			}
		case compiler.OpSliceIndex:
			high := v.pop()
			low := v.pop()
			left := v.pop()

			var lowIdx *int64
			switch low := (*low).(type) {
			case *objects.Undefined:
			case *objects.Int:
				lowIdx = &low.Value
			default:
				return fmt.Errorf("non-integer slice index: %s", low.TypeName())
			}

			var highIdx *int64
			switch high := (*high).(type) {
			case *objects.Undefined:
			case *objects.Int:
				highIdx = &high.Value
			default:
				return fmt.Errorf("non-integer slice index: %s", high.TypeName())
			}

			switch left := (*left).(type) {
			case *objects.Array:
				if err := v.executeArraySliceIndex(left, lowIdx, highIdx); err != nil {
					return err
				}
			case *objects.String:

				if err := v.executeStringSliceIndex(left, lowIdx, highIdx); err != nil {
					return err
				}
			default:
				return fmt.Errorf("cannot slice %s", left.TypeName())
			}
		case compiler.OpCall:
			numArgs := compiler.ReadUint8(ins[ip+1:])
			curFrame.ip += 1

			if err := v.executeCall(int(numArgs)); err != nil {
				return err
			}
		case compiler.OpReturnValue:
			numRets := int(compiler.ReadUint8(ins[ip+1:]))
			curFrame.ip += 1

			var rets []*objects.Object
			for i := 0; i < numRets; i++ {
				val := v.pop()
				rets = append(rets, val)
			}

			v.framesIndex--
			frame := v.frames[v.framesIndex]

			v.sp = frame.basePointer - 1

			for _, retVal := range rets {
				if err := v.push(retVal); err != nil {
					return err
				}
			}
		case compiler.OpReturn:
			v.framesIndex--
			frame := v.frames[v.framesIndex]

			v.sp = frame.basePointer - 1

			if err := v.push(&undefinedObj); err != nil {
				return err
			}

		case compiler.OpDefineLocal:
			localIndex := compiler.ReadUint8(ins[ip+1:])
			curFrame.ip += 1

			sp := curFrame.basePointer + int(localIndex)

			// local variables can be mutated by other actions
			// so always store the copy of popped value
			val := v.popValue()
			v.stack[sp] = &val

		case compiler.OpSetLocal:
			localIndex := compiler.ReadUint8(ins[ip+1:])
			curFrame.ip += 1

			sp := curFrame.basePointer + int(localIndex)

			// update pointee of v.stack[sp] instead of replacing the pointer itself.
			// this is needed because there can be free variables referencing the same local variables.
			val := v.pop()
			*v.stack[sp] = *val // also use a copy of popped value

		case compiler.OpSetSelLocal:
			localIndex := compiler.ReadUint8(ins[ip+1:])
			numSelectors := int(compiler.ReadUint8(ins[ip+2:]))
			curFrame.ip += 2

			// pop selector outcomes (left to right)
			selectors := make([]interface{}, numSelectors, numSelectors)
			for i := 0; i < numSelectors; i++ {
				sel := v.pop()

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
			val := v.pop() // no need to copy value here; selectorAssign uses copy of value

			sp := curFrame.basePointer + int(localIndex)

			if err := selectorAssign(v.stack[sp], val, selectors); err != nil {
				return err
			}
		case compiler.OpGetLocal:
			localIndex := compiler.ReadUint8(ins[ip+1:])
			curFrame.ip += 1

			val := v.stack[curFrame.basePointer+int(localIndex)]

			if err := v.push(val); err != nil {
				return err
			}
		case compiler.OpGetBuiltin:
			builtinIndex := compiler.ReadUint8(ins[ip+1:])
			curFrame.ip += 1

			if err := v.push(&builtinFuncs[builtinIndex]); err != nil {
				return err
			}
		case compiler.OpClosure:
			constIndex := compiler.ReadUint16(ins[ip+1:])
			numFree := compiler.ReadUint8(ins[ip+3:])
			curFrame.ip += 3

			if err := v.pushClosure(int(constIndex), int(numFree)); err != nil {
				return err
			}
		case compiler.OpGetFree:
			freeIndex := compiler.ReadUint8(ins[ip+1:])
			curFrame.ip += 1

			val := curFrame.freeVars[freeIndex]

			if err := v.push(val); err != nil {
				return err
			}
		case compiler.OpSetSelFree:
			freeIndex := compiler.ReadUint8(ins[ip+1:])
			numSelectors := int(compiler.ReadUint8(ins[ip+2:]))
			curFrame.ip += 2

			// pop selector outcomes (left to right)
			selectors := make([]interface{}, numSelectors, numSelectors)
			for i := 0; i < numSelectors; i++ {
				sel := v.pop()

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
			val := v.pop()

			if err := selectorAssign(curFrame.freeVars[freeIndex], val, selectors); err != nil {
				return err
			}
		case compiler.OpSetFree:
			freeIndex := compiler.ReadUint8(ins[ip+1:])
			curFrame.ip += 1

			val := v.pop()

			*v.frames[v.framesIndex-1].freeVars[freeIndex] = *val

		case compiler.OpIteratorInit:
			var iterator objects.Object

			dst := v.pop()
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

			if err := v.push(&iterator); err != nil {
				return err
			}

		case compiler.OpIteratorNext:
			iterator := v.pop()
			b := (*iterator).(objects.Iterator).Next()
			if b {
				if err := v.push(&trueObj); err != nil {
					return err
				}
			} else {
				if err := v.push(&falseObj); err != nil {
					return err
				}
			}
		case compiler.OpIteratorKey:
			iterator := v.pop()
			val := (*iterator).(objects.Iterator).Key()

			if err := v.push(&val); err != nil {
				return err
			}
		case compiler.OpIteratorValue:
			iterator := v.pop()
			val := (*iterator).(objects.Iterator).Value()

			if err := v.push(&val); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown opcode: %d", op)
		}
	}

	return nil
}

func (v *VM) Globals() []*objects.Object {
	return v.globals
}

// for tests
func (v *VM) Stack() []*objects.Object {
	return v.stack[:v.sp]
}

// for tests
func (v *VM) FrameDebug() (frameIndex int, ip int) {
	return v.framesIndex - 1, v.frames[v.framesIndex-1].ip
}

func (v *VM) push(o *objects.Object) error {
	if v.sp >= StackSize {
		return ErrStackOverflow
	}

	v.stack[v.sp] = o
	v.sp++

	return nil
}

func (v *VM) pop() *objects.Object {
	o := v.stack[v.sp-1]
	v.sp--

	return o
}

func (v *VM) popValue() objects.Object {
	o := v.stack[v.sp-1]
	v.sp--

	return *o
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

	var cl objects.Object = &objects.Closure{
		Fn:   fn,
		Free: free,
	}

	return v.push(&cl)
}

func (v *VM) executeStringIndex(str *objects.String, index int64) error {
	rs := []rune(str.Value)

	if index < 0 || index >= int64(len(rs)) {
		return fmt.Errorf("index out of bounds: %d", index)
	}

	var val objects.Object = &objects.Char{Value: rs[index]}

	return v.push(&val)
}

func (v *VM) executeArrayIndex(arr *objects.Array, index int64) error {
	if index < 0 || index >= int64(len(arr.Value)) {
		return fmt.Errorf("index out of bounds: %d", index)
	}

	return v.push(&arr.Value[index])
}

func (v *VM) executeMapIndex(map_ *objects.Map, key string) error {
	var res = undefinedObj
	val, ok := map_.Value[key]
	if ok {
		res = val
	}

	return v.push(&res)
}

func (v *VM) executeArraySliceIndex(arr *objects.Array, low, high *int64) error {
	numElements := int64(len(arr.Value))

	var lowIdx, highIdx int64

	if low != nil {
		lowIdx = *low
		if lowIdx < 0 || lowIdx >= numElements {
			return fmt.Errorf("index out of bounds: %d", lowIdx)
		}
	}
	//} else {
	//	lowIdx = 0
	//}

	if high != nil {
		highIdx = *high
		if highIdx < 0 || highIdx > numElements {
			return fmt.Errorf("index out of bounds: %d", highIdx)
		}
	} else {
		highIdx = numElements
	}

	if lowIdx > highIdx {
		return fmt.Errorf("invalid slice index: %d > %d", lowIdx, highIdx)
	}

	var val objects.Object = &objects.Array{Value: arr.Value[lowIdx:highIdx]}

	return v.push(&val)
}

func (v *VM) executeStringSliceIndex(left *objects.String, low, high *int64) error {
	var lowIdx, highIdx int64

	if low != nil {
		lowIdx = *low
		if lowIdx < 0 || lowIdx >= int64(len(left.Value)) {
			return fmt.Errorf("index out of bounds: %d", lowIdx)
		}
	}
	//} else {
	//	lowIdx = 0
	//}

	if high != nil {
		highIdx = *high
		if highIdx < 0 || highIdx > int64(len(left.Value)) {
			return fmt.Errorf("index out of bounds: %d", highIdx)
		}
	} else {
		highIdx = int64(len(left.Value))
	}

	if lowIdx > highIdx {
		return fmt.Errorf("invalid slice index: %d > %d", lowIdx, highIdx)
	}

	var val objects.Object = &objects.String{Value: left.Value[lowIdx:highIdx]}

	return v.push(&val)
}

func (v *VM) executeCall(numArgs int) error {
	callee := *v.stack[v.sp-1-numArgs]

	switch callee := callee.(type) {
	case *objects.Closure:
		return v.callFunction(callee.Fn, callee.Free, numArgs)
	case *objects.CompiledFunction:
		return v.callFunction(callee, nil, numArgs)
	case objects.Callable:
		return v.callCallable(callee, numArgs)
	default:
		return fmt.Errorf("calling non-function: %s", callee.TypeName())
	}
}

func (v *VM) callFunction(fn *objects.CompiledFunction, freeVars []*objects.Object, numArgs int) error {
	if numArgs != fn.NumParameters {
		return fmt.Errorf("wrong number of arguments: want=%d, got=%d",
			fn.NumParameters, numArgs)
	}

	// check if this is a tail-call (recursive call right before return)
	curFrame := &(v.frames[v.framesIndex-1])
	if fn == curFrame.fn { // recursion
		nextOp := compiler.Opcode(curFrame.fn.Instructions[curFrame.ip+1])
		if nextOp == compiler.OpReturnValue || // tail call
			(nextOp == compiler.OpPop &&
				compiler.OpReturn == compiler.Opcode(curFrame.fn.Instructions[curFrame.ip+2])) {

			//  stack before tail-call
			//
			//  |--------|
			//  |        | <- SP  current
			//  |--------|
			//  | *ARG2  |        for next function (tail-call)
			//  |--------|
			//  | *ARG1  |        for next function (tail-call)
			//  |--------|
			//  | LOCAL3 |        for current function
			//  |--------|
			//  | LOCAL2 |        for current function
			//  |--------|
			//  |  ARG2  |        for current function
			//  |--------|
			//  |  ARG1  | <- BP  for current function
			//  |--------|

			copy(v.stack[curFrame.basePointer:], v.stack[v.sp-numArgs:v.sp])
			v.sp -= numArgs
			curFrame.ip = -1

			//  stack after tail-call
			//
			//  |--------|
			//  |        |
			//  |--------|
			//  | *ARG2  |
			//  |--------|
			//  | *ARG1  | <- SP  current
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

	v.frames[v.framesIndex].fn = fn
	v.frames[v.framesIndex].freeVars = freeVars
	v.frames[v.framesIndex].ip = -1
	v.frames[v.framesIndex].basePointer = v.sp - numArgs
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

func (v *VM) callCallable(callable objects.Callable, numArgs int) error {
	var args []objects.Object
	for _, arg := range v.stack[v.sp-numArgs : v.sp] {
		args = append(args, *arg)
	}

	res, err := callable.Call(args...)
	v.sp -= numArgs + 1

	// runtime error
	if err != nil {
		return err
	}

	// nil return -> undefined
	if res == nil {
		res = undefinedObj
	}

	return v.push(&res)
}

func selectorAssign(dst, src *objects.Object, selectors []interface{}) error {
	numSel := len(selectors)

	for idx := 0; idx < numSel; idx++ {
		switch sel := selectors[idx].(type) {
		case string:
			map_, isMap := (*dst).(*objects.Map)
			if !isMap {
				return fmt.Errorf("invalid map object for selector '%s'", sel)
			}

			if idx == numSel-1 {
				map_.Set(sel, *src)
				return nil
			}

			nxt, found := map_.Get(sel)
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

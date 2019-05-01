package runtime

import (
	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/compiler/source"
	"github.com/d5/tengo/objects"
)

type hooks struct {
	v *VM
}

func (h hooks) Call(val objects.Object, args ...objects.Object) (objects.Object, error) {
	fn := objects.CompiledFunction{
		Instructions: []byte{
			compiler.OpCall,
			byte(len(args)),
			compiler.OpEscape,
		},
		SourceMap: map[int]source.Pos{},
	}

	h.v.curFrame.ip = h.v.ip // store current ip before call
	h.v.curFrame = &(h.v.frames[h.v.framesIndex])
	h.v.curFrame.fn = &fn
	h.v.curFrame.freeVars = nil
	h.v.curFrame.basePointer = h.v.sp
	h.v.curInsts = fn.Instructions
	h.v.ip = -1
	h.v.framesIndex++

	h.v.stack[h.v.sp] = val
	h.v.sp++
	for _, a := range args {
		h.v.stack[h.v.sp] = a
		h.v.sp++
	}

	h.v.run()

	retVal := h.v.stack[h.v.sp-1]
	h.v.framesIndex--
	h.v.curFrame = &h.v.frames[h.v.framesIndex-1]
	h.v.curInsts = h.v.curFrame.fn.Instructions
	h.v.ip = h.v.curFrame.ip
	h.v.sp = h.v.frames[h.v.framesIndex].basePointer
	// skip stack overflow check because (newSP) <= (oldSP)
	//v.stack[v.sp-1] = retVal

	if h.v.err != nil {
		return nil, h.v.err
	}

	return retVal, nil
}

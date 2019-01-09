package vm

import (
	"github.com/d5/tengo/objects"
)

type Frame struct {
	fn          *objects.CompiledFunction
	freeVars    []*objects.Object
	ip          int
	basePointer int
}

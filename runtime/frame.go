package runtime

import (
	"github.com/d5/tengo/objects"
)

// Frame represents a function call frame.
type Frame struct {
	fn            *objects.CompiledFunction
	freeVars      []*objects.FreeVar
	localFreeVars map[int]*objects.FreeVar
	ip            int
	basePointer   int
}

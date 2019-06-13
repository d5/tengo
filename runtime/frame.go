package runtime

import "github.com/d5/tengo"

// Frame represents a function call frame.
type Frame struct {
	fn          *tengo.CompiledFunction
	freeVars    []*tengo.ObjectPtr
	ip          int
	basePointer int
}

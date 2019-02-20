package compiler

import "github.com/d5/tengo/compiler/source"

// CompilationScope represents a compiled instructions
// and the last two instructions that were emitted.
type CompilationScope struct {
	instructions     []byte
	lastInstructions [2]EmittedInstruction
	symbolInit       map[string]bool
	sourceMap        map[int]source.Pos
}

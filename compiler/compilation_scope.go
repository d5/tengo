package compiler

// CompilationScope represents a compiled instructions
// and the last two instructions that were emitted.
type CompilationScope struct {
	instructions     []byte
	lastInstructions [2]EmittedInstruction
}

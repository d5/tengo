package compiler

type CompilationScope struct {
	instructions     []byte
	lastInstructions [2]EmittedInstruction
}

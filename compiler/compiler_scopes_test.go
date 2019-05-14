package compiler_test

import (
	"testing"

	"github.com/d5/tengo/compiler"
)

func TestCompilerScopes(t *testing.T) {
	expect(t, `
if a := 1; a {
    a = 2
	b := a
} else {
    a = 3
	b := a
}`, bytecode(
		concat(
			compiler.MakeInstruction(compiler.OpConstant, 0),
			compiler.MakeInstruction(compiler.OpSetGlobal, 0),
			compiler.MakeInstruction(compiler.OpGetGlobal, 0),
			compiler.MakeInstruction(compiler.OpJumpFalsy, 27),
			compiler.MakeInstruction(compiler.OpConstant, 1),
			compiler.MakeInstruction(compiler.OpSetGlobal, 0),
			compiler.MakeInstruction(compiler.OpGetGlobal, 0),
			compiler.MakeInstruction(compiler.OpSetGlobal, 1),
			compiler.MakeInstruction(compiler.OpJump, 39),
			compiler.MakeInstruction(compiler.OpConstant, 2),
			compiler.MakeInstruction(compiler.OpSetGlobal, 0),
			compiler.MakeInstruction(compiler.OpGetGlobal, 0),
			compiler.MakeInstruction(compiler.OpSetGlobal, 1),
			compiler.MakeInstruction(compiler.OpSuspend)),
		objectsArray(
			intObject(1),
			intObject(2),
			intObject(3))))

	expect(t, `
func() {
	if a := 1; a {
    	a = 2
		b := a
	} else {
    	a = 3
		b := a
	}
}`, bytecode(
		concat(
			compiler.MakeInstruction(compiler.OpConstant, 3),
			compiler.MakeInstruction(compiler.OpPop),
			compiler.MakeInstruction(compiler.OpSuspend)),
		objectsArray(
			intObject(1),
			intObject(2),
			intObject(3),
			compiledFunction(0, 0,
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpDefineLocal, 0),
				compiler.MakeInstruction(compiler.OpGetLocal, 0),
				compiler.MakeInstruction(compiler.OpJumpFalsy, 22),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpSetLocal, 0),
				compiler.MakeInstruction(compiler.OpGetLocal, 0),
				compiler.MakeInstruction(compiler.OpDefineLocal, 1),
				compiler.MakeInstruction(compiler.OpJump, 31),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpSetLocal, 0),
				compiler.MakeInstruction(compiler.OpGetLocal, 0),
				compiler.MakeInstruction(compiler.OpDefineLocal, 1),
				compiler.MakeInstruction(compiler.OpReturn, 0)))))
}

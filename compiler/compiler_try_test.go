package compiler_test

import (
	"testing"

	"github.com/d5/tengo/compiler"
)

func TestTry(t *testing.T) {
	expect(t, `try(1)`, bytecode(
		concat(
			compiler.MakeInstruction(compiler.OpConstant, 0),
			compiler.MakeInstruction(compiler.OpReturnOnError, 1),
			compiler.MakeInstruction(compiler.OpPop),
			compiler.MakeInstruction(compiler.OpSuspend),
		),
		objectsArray(
			intObject(1),
		),
	))

	expect(t, `try(1+1)`, bytecode(
		concat(
			compiler.MakeInstruction(compiler.OpConstant, 0),
			compiler.MakeInstruction(compiler.OpConstant, 0),
			compiler.MakeInstruction(compiler.OpBinaryOp, 11),
			compiler.MakeInstruction(compiler.OpReturnOnError, 1),
			compiler.MakeInstruction(compiler.OpPop),
			compiler.MakeInstruction(compiler.OpSuspend),
		),
		objectsArray(
			intObject(1),
		),
	))

	expect(t, `x := try(1+1)`, bytecode(
		concat(
			compiler.MakeInstruction(compiler.OpConstant, 0),
			compiler.MakeInstruction(compiler.OpConstant, 0),
			compiler.MakeInstruction(compiler.OpBinaryOp, 11),
			compiler.MakeInstruction(compiler.OpReturnOnError, 1),
			compiler.MakeInstruction(compiler.OpSetGlobal, 0),
			compiler.MakeInstruction(compiler.OpSuspend),
		),
		objectsArray(
			intObject(1),
		),
	))
}

package compiler_test

import (
	"testing"

	"github.com/d5/tengo/compiler"
)

func TestSpreadArray(t *testing.T) {
	expect(t, `[1, 2, [3]...]`, bytecode(
		concat(
			compiler.MakeInstruction(compiler.OpConstant, 0),
			compiler.MakeInstruction(compiler.OpConstant, 1),
			compiler.MakeInstruction(compiler.OpConstant, 2),
			compiler.MakeInstruction(compiler.OpArray, 1),
			compiler.MakeInstruction(compiler.OpSpread),
			compiler.MakeInstruction(compiler.OpArray, 3),
			compiler.MakeInstruction(compiler.OpPop),
		),
		objectsArray(
			intObject(1),
			intObject(2),
			intObject(3),
		),
	))

	expect(t, `a := [1,2]; b := [3]; [a..., b...]`, bytecode(
		concat(
			compiler.MakeInstruction(compiler.OpConstant, 0),
			compiler.MakeInstruction(compiler.OpConstant, 1),
			compiler.MakeInstruction(compiler.OpArray, 2),
			compiler.MakeInstruction(compiler.OpSetGlobal, 0),
			compiler.MakeInstruction(compiler.OpConstant, 2),
			compiler.MakeInstruction(compiler.OpArray, 1),
			compiler.MakeInstruction(compiler.OpSetGlobal, 1),
			compiler.MakeInstruction(compiler.OpGetGlobal, 0),
			compiler.MakeInstruction(compiler.OpSpread),
			compiler.MakeInstruction(compiler.OpGetGlobal, 1),
			compiler.MakeInstruction(compiler.OpSpread),
			compiler.MakeInstruction(compiler.OpArray, 2),
			compiler.MakeInstruction(compiler.OpPop),
		),
		objectsArray(
			intObject(1),
			intObject(2),
			intObject(3),
		),
	))

	expect(t, `fn := undefined; fn(1, 2, [3]...)`, bytecode(
		concat(
			compiler.MakeInstruction(compiler.OpNull),
			compiler.MakeInstruction(compiler.OpSetGlobal, 0),
			compiler.MakeInstruction(compiler.OpGetGlobal, 0),
			compiler.MakeInstruction(compiler.OpConstant, 0),
			compiler.MakeInstruction(compiler.OpConstant, 1),
			compiler.MakeInstruction(compiler.OpConstant, 2),
			compiler.MakeInstruction(compiler.OpArray, 1),
			compiler.MakeInstruction(compiler.OpSpread),
			compiler.MakeInstruction(compiler.OpCall, 3),
			compiler.MakeInstruction(compiler.OpPop),
		),
		objectsArray(
			intObject(1),
			intObject(2),
			intObject(3),
		),
	))
}

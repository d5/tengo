package compiler_test

import (
	"testing"

	"github.com/d5/tengo/compiler"
)

func TestCompilerDeadCode(t *testing.T) {
	expect(t, `
func() {
	a := 4
	return a

	b := 5 // dead code from here
	c := a
	return b
}`,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpPop),
				compiler.MakeInstruction(compiler.OpSuspend)),
			objectsArray(
				intObject(4),
				intObject(5),
				compiledFunction(0, 0,
					compiler.MakeInstruction(compiler.OpConstant, 0),
					compiler.MakeInstruction(compiler.OpDefineLocal, 0),
					compiler.MakeInstruction(compiler.OpGetLocal, 0),
					compiler.MakeInstruction(compiler.OpReturn, 1)))))

	expect(t, `
func() {
	if true {
		return 5
		a := 4  // dead code from here
		b := a
		return b
	} else {
		return 4
		c := 5  // dead code from here
		d := c
		return d
	}
}`, bytecode(
		concat(
			compiler.MakeInstruction(compiler.OpConstant, 2),
			compiler.MakeInstruction(compiler.OpPop),
			compiler.MakeInstruction(compiler.OpSuspend)),
		objectsArray(
			intObject(5),
			intObject(4),
			compiledFunction(0, 0,
				compiler.MakeInstruction(compiler.OpTrue),
				compiler.MakeInstruction(compiler.OpJumpFalsy, 9),
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpReturn, 1),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpReturn, 1)))))

	expect(t, `
func() {
	a := 1
	for {
		if a == 5 {
			return 10
		}
		5 + 5
		return 20
		b := a
		return b
	}
}`, bytecode(
		concat(
			compiler.MakeInstruction(compiler.OpConstant, 4),
			compiler.MakeInstruction(compiler.OpPop),
			compiler.MakeInstruction(compiler.OpSuspend)),
		objectsArray(
			intObject(1),
			intObject(5),
			intObject(10),
			intObject(20),
			compiledFunction(0, 0,
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpDefineLocal, 0),
				compiler.MakeInstruction(compiler.OpGetLocal, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpEqual),
				compiler.MakeInstruction(compiler.OpJumpFalsy, 19),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpReturn, 1),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpBinaryOp, 11),
				compiler.MakeInstruction(compiler.OpPop),
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpReturn, 1)))))

	expect(t, `
func() {
	if true {
		return 5
		a := 4  // dead code from here
		b := a
		return b
	} else {
		return 4
		c := 5  // dead code from here
		d := c
		return d
	}
}`, bytecode(
		concat(
			compiler.MakeInstruction(compiler.OpConstant, 2),
			compiler.MakeInstruction(compiler.OpPop),
			compiler.MakeInstruction(compiler.OpSuspend)),
		objectsArray(
			intObject(5),
			intObject(4),
			compiledFunction(0, 0,
				compiler.MakeInstruction(compiler.OpTrue),
				compiler.MakeInstruction(compiler.OpJumpFalsy, 9),
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpReturn, 1),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpReturn, 1)))))
}

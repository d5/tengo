package compiler_test

import (
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/compiler"
)

func TestReadOperands(t *testing.T) {
	assertReadOperand(t, compiler.OpConstant, []int{65535}, 2)
}

func assertReadOperand(t *testing.T, opcode compiler.Opcode, operands []int, expectedBytes int) {
	inst := compiler.MakeInstruction(opcode, operands...)
	numOperands := compiler.OpcodeOperands[opcode]
	operandsRead, read := compiler.ReadOperands(numOperands, inst[1:])
	assert.Equal(t, expectedBytes, read)
	assert.Equal(t, operands, operandsRead)
}

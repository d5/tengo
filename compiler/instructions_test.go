package compiler_test

import (
	"strings"
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/compiler"
)

func TestInstructions_String(t *testing.T) {
	assertInstructionString(t,
		[][]byte{
			compiler.MakeInstruction(compiler.OpConstant, 1),
			compiler.MakeInstruction(compiler.OpConstant, 2),
			compiler.MakeInstruction(compiler.OpConstant, 65535),
		},
		`0000 CONST   1    
0003 CONST   2    
0006 CONST   65535`)

	assertInstructionString(t,
		[][]byte{
			compiler.MakeInstruction(compiler.OpBinaryOp, 11),
			compiler.MakeInstruction(compiler.OpConstant, 2),
			compiler.MakeInstruction(compiler.OpConstant, 65535),
		},
		`0000 BINARYOP 11   
0002 CONST   2    
0005 CONST   65535`)

	assertInstructionString(t,
		[][]byte{
			compiler.MakeInstruction(compiler.OpBinaryOp, 11),
			compiler.MakeInstruction(compiler.OpGetLocal, 1),
			compiler.MakeInstruction(compiler.OpConstant, 2),
			compiler.MakeInstruction(compiler.OpConstant, 65535),
		},
		`0000 BINARYOP 11   
0002 GETL    1    
0004 CONST   2    
0007 CONST   65535`)
}

func TestMakeInstruction(t *testing.T) {
	makeInstruction(t, []byte{byte(compiler.OpConstant), 0, 0}, compiler.OpConstant, 0)
	makeInstruction(t, []byte{byte(compiler.OpConstant), 0, 1}, compiler.OpConstant, 1)
	makeInstruction(t, []byte{byte(compiler.OpConstant), 255, 254}, compiler.OpConstant, 65534)
	makeInstruction(t, []byte{byte(compiler.OpPop)}, compiler.OpPop)
	makeInstruction(t, []byte{byte(compiler.OpTrue)}, compiler.OpTrue)
	makeInstruction(t, []byte{byte(compiler.OpFalse)}, compiler.OpFalse)
}

func assertInstructionString(t *testing.T, instructions [][]byte, expected string) {
	concatted := make([]byte, 0)
	for _, e := range instructions {
		concatted = append(concatted, e...)
	}
	assert.Equal(t, expected, strings.Join(compiler.FormatInstructions(concatted, 0), "\n"))
}

func makeInstruction(t *testing.T, expected []byte, opcode compiler.Opcode, operands ...int) {
	inst := compiler.MakeInstruction(opcode, operands...)
	assert.Equal(t, expected, []byte(inst))
}

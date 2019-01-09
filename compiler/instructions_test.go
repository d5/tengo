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
			compiler.MakeInstruction(compiler.OpAdd),
			compiler.MakeInstruction(compiler.OpConstant, 2),
			compiler.MakeInstruction(compiler.OpConstant, 65535),
		},
		`0000 ADD    
0001 CONST   2    
0004 CONST   65535`)

	assertInstructionString(t,
		[][]byte{
			compiler.MakeInstruction(compiler.OpAdd),
			compiler.MakeInstruction(compiler.OpGetLocal, 1),
			compiler.MakeInstruction(compiler.OpConstant, 2),
			compiler.MakeInstruction(compiler.OpConstant, 65535),
		},
		`0000 ADD    
0001 GETL    1    
0003 CONST   2    
0006 CONST   65535`)
}

func TestMakeInstruction(t *testing.T) {
	makeInstruction(t, []byte{byte(compiler.OpConstant), 0, 0}, compiler.OpConstant, 0)
	makeInstruction(t, []byte{byte(compiler.OpConstant), 0, 1}, compiler.OpConstant, 1)
	makeInstruction(t, []byte{byte(compiler.OpConstant), 255, 254}, compiler.OpConstant, 65534)
	makeInstruction(t, []byte{byte(compiler.OpPop)}, compiler.OpPop)
	makeInstruction(t, []byte{byte(compiler.OpTrue)}, compiler.OpTrue)
	makeInstruction(t, []byte{byte(compiler.OpFalse)}, compiler.OpFalse)
	makeInstruction(t, []byte{byte(compiler.OpAdd)}, compiler.OpAdd)
	makeInstruction(t, []byte{byte(compiler.OpSub)}, compiler.OpSub)
	makeInstruction(t, []byte{byte(compiler.OpMul)}, compiler.OpMul)
	makeInstruction(t, []byte{byte(compiler.OpDiv)}, compiler.OpDiv)
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

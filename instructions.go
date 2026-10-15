package tengo

import (
	"fmt"

	"github.com/d5/tengo/v2/parser"
)

var instructionLen [256]int

func init() {
	for op, widths := range parser.OpcodeOperands {
		n := 1
		for _, w := range widths {
			n += w
		}
		instructionLen[op] = n
	}
}

// MakeInstruction returns a bytecode for an opcode and the operands.
func MakeInstruction(opcode parser.Opcode, operands ...int) []byte {
	return appendInstruction(make([]byte, 0, instructionLen[opcode]),
		opcode, operands...)
}

func appendInstruction(
	dst []byte,
	opcode parser.Opcode,
	operands ...int,
) []byte {
	numOperands := parser.OpcodeOperands[opcode]
	dst = append(dst, opcode)
	for i, width := range numOperands {
		var o int
		if i < len(operands) {
			o = operands[i]
		}
		switch width {
		case 1:
			dst = append(dst, byte(o))
		case 2:
			n := uint16(o)
			dst = append(dst, byte(n>>8), byte(n))
		case 4:
			n := uint32(o)
			dst = append(dst, byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
		}
	}
	return dst
}

// FormatInstructions returns string representation of bytecode instructions.
func FormatInstructions(b []byte, posOffset int) []string {
	var out []string

	i := 0
	for i < len(b) {
		numOperands := parser.OpcodeOperands[b[i]]
		operands, read := parser.ReadOperands(numOperands, b[i+1:])

		switch len(numOperands) {
		case 0:
			out = append(out, fmt.Sprintf("%04d %-7s",
				posOffset+i, parser.OpcodeNames[b[i]]))
		case 1:
			out = append(out, fmt.Sprintf("%04d %-7s %-5d",
				posOffset+i, parser.OpcodeNames[b[i]], operands[0]))
		case 2:
			out = append(out, fmt.Sprintf("%04d %-7s %-5d %-5d",
				posOffset+i, parser.OpcodeNames[b[i]],
				operands[0], operands[1]))
		}
		i += 1 + read
	}
	return out
}

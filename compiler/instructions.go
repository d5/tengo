package compiler

import (
	"fmt"
)

// MakeInstruction returns a bytecode for an opcode and the operands.
func MakeInstruction(opcode Opcode, operands ...int) []byte {
	def, ok := Lookup(opcode)
	if !ok {
		return nil
	}

	totalLen := 1
	for _, w := range def.Operands {
		totalLen += w
	}

	instruction := make([]byte, totalLen, totalLen)
	instruction[0] = byte(opcode)

	offset := 1
	for i, o := range operands {
		width := def.Operands[i]
		switch width {
		case 1:
			instruction[offset] = byte(o)
		case 2:
			n := uint16(o)
			instruction[offset] = byte(n >> 8)
			instruction[offset+1] = byte(n)
		}
		offset += width
	}

	return instruction
}

// FormatInstructions returns string representation of
// bytecode instructions.
func FormatInstructions(b []byte, posOffset int) []string {
	var out []string

	i := 0
	for i < len(b) {
		def, ok := Lookup(Opcode(b[i]))
		if !ok {
			out = append(out, fmt.Sprintf("error: unknown Opcode %d", b[i]))
			continue
		}

		operands, read := ReadOperands(def, b[i+1:])

		switch len(def.Operands) {
		case 0:
			out = append(out, fmt.Sprintf("%04d %-7s", posOffset+i, def.Name))
		case 1:
			out = append(out, fmt.Sprintf("%04d %-7s %-5d", posOffset+i, def.Name, operands[0]))
		case 2:
			out = append(out, fmt.Sprintf("%04d %-7s %-5d %-5d", posOffset+i, def.Name, operands[0], operands[1]))
		}

		i += 1 + read
	}

	return out
}

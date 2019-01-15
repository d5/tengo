package compiler

// ReadOperands reads operands from the bytecode.
func ReadOperands(def *Definition, ins []byte) (operands []int, offset int) {
	for _, width := range def.Operands {
		switch width {
		case 1:
			operands = append(operands, int(ReadUint8(ins[offset:])))
		case 2:
			operands = append(operands, int(ReadUint16(ins[offset:])))
		}

		offset += width
	}

	return
}

// ReadUint16 reads uint16 from the byte slice.
func ReadUint16(b []byte) uint16 {
	return uint16(b[1]) | uint16(b[0])<<8
}

// ReadUint8 reads uint8 from the byte slice.
func ReadUint8(b []byte) uint8 {
	return uint8(b[0])
}

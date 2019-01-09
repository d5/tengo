package compiler

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

func ReadUint16(b []byte) uint16 {
	return uint16(b[1]) | uint16(b[0])<<8
}

func ReadUint8(b []byte) uint8 {
	return uint8(b[0])
}

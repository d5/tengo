package parser_test

import (
	"testing"

	"github.com/d5/tengo/v2/parser"
	"github.com/d5/tengo/v2/require"
)

func TestReadOperands(t *testing.T) {
	for _, tc := range []struct {
		op       parser.Opcode
		ins      []byte
		operands []int
		offset   int
	}{
		{parser.OpPop, []byte{0xff}, nil, 0},
		{parser.OpGetLocal, []byte{0xff, 0x01}, []int{255}, 1},
		{parser.OpConstant, []byte{0x01, 0x02, 0xff}, []int{0x0102}, 2},
		{parser.OpJump, []byte{0x01, 0x02, 0x03, 0x04}, []int{0x01020304}, 4},
		{parser.OpCall, []byte{0x03, 0x01}, []int{3, 1}, 2},
		{parser.OpClosure, []byte{0xff, 0xfe, 0x02}, []int{0xfffe, 2}, 3},
	} {
		name := parser.OpcodeNames[tc.op]
		operands, offset := parser.ReadOperands(
			parser.OpcodeOperands[tc.op], tc.ins)
		require.Equal(t, tc.operands, operands, name)
		require.Equal(t, tc.offset, offset, name)
	}
}

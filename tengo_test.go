package tengo_test

import (
	"strings"
	"testing"
	"time"

	"github.com/d5/tengo"
	"github.com/d5/tengo/internal"
	"github.com/d5/tengo/internal/require"
)

func TestInstructions_String(t *testing.T) {
	assertInstructionString(t,
		[][]byte{
			internal.MakeInstruction(internal.OpConstant, 1),
			internal.MakeInstruction(internal.OpConstant, 2),
			internal.MakeInstruction(internal.OpConstant, 65535),
		},
		`0000 CONST   1    
0003 CONST   2    
0006 CONST   65535`)

	assertInstructionString(t,
		[][]byte{
			internal.MakeInstruction(internal.OpBinaryOp, 11),
			internal.MakeInstruction(internal.OpConstant, 2),
			internal.MakeInstruction(internal.OpConstant, 65535),
		},
		`0000 BINARYOP 11   
0002 CONST   2    
0005 CONST   65535`)

	assertInstructionString(t,
		[][]byte{
			internal.MakeInstruction(internal.OpBinaryOp, 11),
			internal.MakeInstruction(internal.OpGetLocal, 1),
			internal.MakeInstruction(internal.OpConstant, 2),
			internal.MakeInstruction(internal.OpConstant, 65535),
		},
		`0000 BINARYOP 11   
0002 GETL    1    
0004 CONST   2    
0007 CONST   65535`)
}

func TestMakeInstruction(t *testing.T) {
	makeInstruction(t, []byte{internal.OpConstant, 0, 0},
		internal.OpConstant, 0)
	makeInstruction(t, []byte{internal.OpConstant, 0, 1},
		internal.OpConstant, 1)
	makeInstruction(t, []byte{internal.OpConstant, 255, 254},
		internal.OpConstant, 65534)
	makeInstruction(t, []byte{internal.OpPop}, internal.OpPop)
	makeInstruction(t, []byte{internal.OpTrue}, internal.OpTrue)
	makeInstruction(t, []byte{internal.OpFalse}, internal.OpFalse)
}

func TestNumObjects(t *testing.T) {
	testCountObjects(t, &tengo.Array{}, 1)
	testCountObjects(t, &tengo.Array{Value: []tengo.Object{
		&tengo.Int{Value: 1},
		&tengo.Int{Value: 2},
		&tengo.Array{Value: []tengo.Object{
			&tengo.Int{Value: 3},
			&tengo.Int{Value: 4},
			&tengo.Int{Value: 5},
		}},
	}}, 7)
	testCountObjects(t, tengo.TrueValue, 1)
	testCountObjects(t, tengo.FalseValue, 1)
	testCountObjects(t, &tengo.BuiltinFunction{}, 1)
	testCountObjects(t, &tengo.Bytes{Value: []byte("foobar")}, 1)
	testCountObjects(t, &tengo.Char{Value: '가'}, 1)
	testCountObjects(t, &tengo.CompiledFunction{}, 1)
	testCountObjects(t, &tengo.Error{Value: &tengo.Int{Value: 5}}, 2)
	testCountObjects(t, &tengo.Float{Value: 19.84}, 1)
	testCountObjects(t, &tengo.ImmutableArray{Value: []tengo.Object{
		&tengo.Int{Value: 1},
		&tengo.Int{Value: 2},
		&tengo.ImmutableArray{Value: []tengo.Object{
			&tengo.Int{Value: 3},
			&tengo.Int{Value: 4},
			&tengo.Int{Value: 5},
		}},
	}}, 7)
	testCountObjects(t, &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"k1": &tengo.Int{Value: 1},
			"k2": &tengo.Int{Value: 2},
			"k3": &tengo.Array{Value: []tengo.Object{
				&tengo.Int{Value: 3},
				&tengo.Int{Value: 4},
				&tengo.Int{Value: 5},
			}},
		}}, 7)
	testCountObjects(t, &tengo.Int{Value: 1984}, 1)
	testCountObjects(t, &tengo.Map{Value: map[string]tengo.Object{
		"k1": &tengo.Int{Value: 1},
		"k2": &tengo.Int{Value: 2},
		"k3": &tengo.Array{Value: []tengo.Object{
			&tengo.Int{Value: 3},
			&tengo.Int{Value: 4},
			&tengo.Int{Value: 5},
		}},
	}}, 7)
	testCountObjects(t, &tengo.String{Value: "foo bar"}, 1)
	testCountObjects(t, &tengo.Time{Value: time.Now()}, 1)
	testCountObjects(t, tengo.UndefinedValue, 1)
}

func testCountObjects(t *testing.T, o tengo.Object, expected int) {
	require.Equal(t, expected, tengo.CountObjects(o))
}

func assertInstructionString(
	t *testing.T,
	instructions [][]byte,
	expected string,
) {
	concatted := make([]byte, 0)
	for _, e := range instructions {
		concatted = append(concatted, e...)
	}
	require.Equal(t, expected, strings.Join(
		internal.FormatInstructions(concatted, 0), "\n"))
}

func makeInstruction(
	t *testing.T,
	expected []byte,
	opcode internal.Opcode,
	operands ...int,
) {
	inst := internal.MakeInstruction(opcode, operands...)
	require.Equal(t, expected, inst)
}

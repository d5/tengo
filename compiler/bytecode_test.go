package compiler_test

import (
	"bytes"
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/objects"
)

func TestBytecode(t *testing.T) {
	testBytecodeSerialization(t, &compiler.Bytecode{})

	testBytecodeSerialization(t, bytecode(
		concat(), objectsArray(
			&objects.Array{
				Value: objectsArray(
					&objects.Int{Value: 12},
					&objects.String{Value: "foo"},
					&objects.Bool{Value: true},
					&objects.Float{Value: 93.11},
					&objects.Char{Value: 'x'},
				),
			},
			&objects.Bool{Value: false},
			&objects.Char{Value: 'y'},
			&objects.Float{Value: 93.11},
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpSetLocal, 0),
				compiler.MakeInstruction(compiler.OpGetGlobal, 0),
				compiler.MakeInstruction(compiler.OpGetFree, 0)),
			&objects.Float{Value: 39.2},
			&objects.Int{Value: 192},
			&objects.Map{
				Value: map[string]objects.Object{
					"a": &objects.Float{Value: -93.1},
					"b": &objects.Bool{Value: false},
				},
			},
			&objects.String{Value: "bar"},
			&objects.Undefined{})))

	testBytecodeSerialization(t, bytecode(
		concat(
			compiler.MakeInstruction(compiler.OpConstant, 0),
			compiler.MakeInstruction(compiler.OpSetGlobal, 0),
			compiler.MakeInstruction(compiler.OpConstant, 6),
			compiler.MakeInstruction(compiler.OpPop)),
		objectsArray(
			intObject(55),
			intObject(66),
			intObject(77),
			intObject(88),
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpSetLocal, 0),
				compiler.MakeInstruction(compiler.OpGetGlobal, 0),
				compiler.MakeInstruction(compiler.OpGetFree, 0),
				compiler.MakeInstruction(compiler.OpAdd),
				compiler.MakeInstruction(compiler.OpGetFree, 1),
				compiler.MakeInstruction(compiler.OpAdd),
				compiler.MakeInstruction(compiler.OpGetLocal, 0),
				compiler.MakeInstruction(compiler.OpAdd),
				compiler.MakeInstruction(compiler.OpReturnValue, 1)),
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpSetLocal, 0),
				compiler.MakeInstruction(compiler.OpGetFree, 0),
				compiler.MakeInstruction(compiler.OpGetLocal, 0),
				compiler.MakeInstruction(compiler.OpClosure, 4, 2),
				compiler.MakeInstruction(compiler.OpReturnValue, 1)),
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpSetLocal, 0),
				compiler.MakeInstruction(compiler.OpGetLocal, 0),
				compiler.MakeInstruction(compiler.OpClosure, 5, 1),
				compiler.MakeInstruction(compiler.OpReturnValue, 1)))))
}

func testBytecodeSerialization(t *testing.T, b *compiler.Bytecode) {
	var buf bytes.Buffer
	err := b.Encode(&buf)
	assert.NoError(t, err)

	r := &compiler.Bytecode{}
	err = r.Decode(bytes.NewReader(buf.Bytes()))
	assert.NoError(t, err)

	assert.Equal(t, b.Instructions, r.Instructions)
	assert.Equal(t, b.Constants, r.Constants)
}

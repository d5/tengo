package compiler_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/d5/tengo"
	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/compiler/source"
)

type srcfile struct {
	name string
	size int
}

func TestBytecode(t *testing.T) {
	testBytecodeSerialization(t, bytecode(concat(), objectsArray()))

	testBytecodeSerialization(t, bytecode(
		concat(), objectsArray(
			&tengo.Char{Value: 'y'},
			&tengo.Float{Value: 93.11},
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpSetLocal, 0),
				compiler.MakeInstruction(compiler.OpGetGlobal, 0),
				compiler.MakeInstruction(compiler.OpGetFree, 0)),
			&tengo.Float{Value: 39.2},
			&tengo.Int{Value: 192},
			&tengo.String{Value: "bar"})))

	testBytecodeSerialization(t, bytecodeFileSet(
		concat(
			compiler.MakeInstruction(compiler.OpConstant, 0),
			compiler.MakeInstruction(compiler.OpSetGlobal, 0),
			compiler.MakeInstruction(compiler.OpConstant, 6),
			compiler.MakeInstruction(compiler.OpPop)),
		objectsArray(
			&tengo.Int{Value: 55},
			&tengo.Int{Value: 66},
			&tengo.Int{Value: 77},
			&tengo.Int{Value: 88},
			&tengo.ImmutableMap{
				Value: map[string]tengo.Object{
					"array": &tengo.ImmutableArray{
						Value: []tengo.Object{
							&tengo.Int{Value: 1},
							&tengo.Int{Value: 2},
							&tengo.Int{Value: 3},
							tengo.TrueValue,
							tengo.FalseValue,
							tengo.UndefinedValue,
						},
					},
					"true":  tengo.TrueValue,
					"false": tengo.FalseValue,
					"bytes": &tengo.Bytes{Value: make([]byte, 16)},
					"char":  &tengo.Char{Value: 'Y'},
					"error": &tengo.Error{Value: &tengo.String{Value: "some error"}},
					"float": &tengo.Float{Value: -19.84},
					"immutable_array": &tengo.ImmutableArray{
						Value: []tengo.Object{
							&tengo.Int{Value: 1},
							&tengo.Int{Value: 2},
							&tengo.Int{Value: 3},
							tengo.TrueValue,
							tengo.FalseValue,
							tengo.UndefinedValue,
						},
					},
					"immutable_map": &tengo.ImmutableMap{
						Value: map[string]tengo.Object{
							"a": &tengo.Int{Value: 1},
							"b": &tengo.Int{Value: 2},
							"c": &tengo.Int{Value: 3},
							"d": tengo.TrueValue,
							"e": tengo.FalseValue,
							"f": tengo.UndefinedValue,
						},
					},
					"int": &tengo.Int{Value: 91},
					"map": &tengo.Map{
						Value: map[string]tengo.Object{
							"a": &tengo.Int{Value: 1},
							"b": &tengo.Int{Value: 2},
							"c": &tengo.Int{Value: 3},
							"d": tengo.TrueValue,
							"e": tengo.FalseValue,
							"f": tengo.UndefinedValue,
						},
					},
					"string":    &tengo.String{Value: "foo bar"},
					"time":      &tengo.Time{Value: time.Now()},
					"undefined": tengo.UndefinedValue,
				},
			},
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpSetLocal, 0),
				compiler.MakeInstruction(compiler.OpGetGlobal, 0),
				compiler.MakeInstruction(compiler.OpGetFree, 0),
				compiler.MakeInstruction(compiler.OpBinaryOp, 11),
				compiler.MakeInstruction(compiler.OpGetFree, 1),
				compiler.MakeInstruction(compiler.OpBinaryOp, 11),
				compiler.MakeInstruction(compiler.OpGetLocal, 0),
				compiler.MakeInstruction(compiler.OpBinaryOp, 11),
				compiler.MakeInstruction(compiler.OpReturn, 1)),
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpSetLocal, 0),
				compiler.MakeInstruction(compiler.OpGetFree, 0),
				compiler.MakeInstruction(compiler.OpGetLocal, 0),
				compiler.MakeInstruction(compiler.OpClosure, 4, 2),
				compiler.MakeInstruction(compiler.OpReturn, 1)),
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpSetLocal, 0),
				compiler.MakeInstruction(compiler.OpGetLocal, 0),
				compiler.MakeInstruction(compiler.OpClosure, 5, 1),
				compiler.MakeInstruction(compiler.OpReturn, 1))),
		fileSet(srcfile{name: "file1", size: 100}, srcfile{name: "file2", size: 200})))
}

func TestBytecode_RemoveDuplicates(t *testing.T) {
	testBytecodeRemoveDuplicates(t,
		bytecode(
			concat(), objectsArray(
				&tengo.Char{Value: 'y'},
				&tengo.Float{Value: 93.11},
				compiledFunction(1, 0,
					compiler.MakeInstruction(compiler.OpConstant, 3),
					compiler.MakeInstruction(compiler.OpSetLocal, 0),
					compiler.MakeInstruction(compiler.OpGetGlobal, 0),
					compiler.MakeInstruction(compiler.OpGetFree, 0)),
				&tengo.Float{Value: 39.2},
				&tengo.Int{Value: 192},
				&tengo.String{Value: "bar"})),
		bytecode(
			concat(), objectsArray(
				&tengo.Char{Value: 'y'},
				&tengo.Float{Value: 93.11},
				compiledFunction(1, 0,
					compiler.MakeInstruction(compiler.OpConstant, 3),
					compiler.MakeInstruction(compiler.OpSetLocal, 0),
					compiler.MakeInstruction(compiler.OpGetGlobal, 0),
					compiler.MakeInstruction(compiler.OpGetFree, 0)),
				&tengo.Float{Value: 39.2},
				&tengo.Int{Value: 192},
				&tengo.String{Value: "bar"})))

	testBytecodeRemoveDuplicates(t,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpConstant, 4),
				compiler.MakeInstruction(compiler.OpConstant, 5),
				compiler.MakeInstruction(compiler.OpConstant, 6),
				compiler.MakeInstruction(compiler.OpConstant, 7),
				compiler.MakeInstruction(compiler.OpConstant, 8),
				compiler.MakeInstruction(compiler.OpClosure, 4, 1)),
			objectsArray(
				&tengo.Int{Value: 1},
				&tengo.Float{Value: 2.0},
				&tengo.Char{Value: '3'},
				&tengo.String{Value: "four"},
				compiledFunction(1, 0,
					compiler.MakeInstruction(compiler.OpConstant, 3),
					compiler.MakeInstruction(compiler.OpConstant, 7),
					compiler.MakeInstruction(compiler.OpSetLocal, 0),
					compiler.MakeInstruction(compiler.OpGetGlobal, 0),
					compiler.MakeInstruction(compiler.OpGetFree, 0)),
				&tengo.Int{Value: 1},
				&tengo.Float{Value: 2.0},
				&tengo.Char{Value: '3'},
				&tengo.String{Value: "four"})),
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpConstant, 4),
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpClosure, 4, 1)),
			objectsArray(
				&tengo.Int{Value: 1},
				&tengo.Float{Value: 2.0},
				&tengo.Char{Value: '3'},
				&tengo.String{Value: "four"},
				compiledFunction(1, 0,
					compiler.MakeInstruction(compiler.OpConstant, 3),
					compiler.MakeInstruction(compiler.OpConstant, 2),
					compiler.MakeInstruction(compiler.OpSetLocal, 0),
					compiler.MakeInstruction(compiler.OpGetGlobal, 0),
					compiler.MakeInstruction(compiler.OpGetFree, 0)))))

	testBytecodeRemoveDuplicates(t,
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpConstant, 4)),
			objectsArray(
				&tengo.Int{Value: 1},
				&tengo.Int{Value: 2},
				&tengo.Int{Value: 3},
				&tengo.Int{Value: 1},
				&tengo.Int{Value: 3})),
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 2)),
			objectsArray(
				&tengo.Int{Value: 1},
				&tengo.Int{Value: 2},
				&tengo.Int{Value: 3})))
}

func TestBytecode_CountObjects(t *testing.T) {
	b := bytecode(
		concat(),
		objectsArray(
			&tengo.Int{Value: 55},
			&tengo.Int{Value: 66},
			&tengo.Int{Value: 77},
			&tengo.Int{Value: 88},
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpReturn, 1)),
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpReturn, 1)),
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpReturn, 1))))
	assert.Equal(t, 7, b.CountObjects())
}

func fileSet(files ...srcfile) *source.FileSet {
	fileSet := source.NewFileSet()
	for _, f := range files {
		fileSet.AddFile(f.name, -1, f.size)
	}
	return fileSet
}

func bytecodeFileSet(instructions []byte, constants []tengo.Object, fileSet *source.FileSet) *compiler.Bytecode {
	return &compiler.Bytecode{
		FileSet:      fileSet,
		MainFunction: &tengo.CompiledFunction{Instructions: instructions},
		Constants:    constants,
	}
}

func testBytecodeRemoveDuplicates(t *testing.T, input, expected *compiler.Bytecode) {
	input.RemoveDuplicates()

	assert.Equal(t, expected.FileSet, input.FileSet)
	assert.Equal(t, expected.MainFunction, input.MainFunction)
	assert.Equal(t, expected.Constants, input.Constants)
}

func testBytecodeSerialization(t *testing.T, b *compiler.Bytecode) {
	var buf bytes.Buffer
	err := b.Encode(&buf)
	assert.NoError(t, err)

	r := &compiler.Bytecode{}
	err = r.Decode(bytes.NewReader(buf.Bytes()), nil)
	assert.NoError(t, err)

	assert.Equal(t, b.FileSet, r.FileSet)
	assert.Equal(t, b.MainFunction, r.MainFunction)
	assert.Equal(t, b.Constants, r.Constants)
}

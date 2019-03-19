package compiler_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/compiler/source"
	"github.com/d5/tengo/objects"
)

type srcfile struct {
	name string
	size int
}

func TestBytecode(t *testing.T) {
	testBytecodeSerialization(t, bytecode(concat(), objectsArray()))

	testBytecodeSerialization(t, bytecode(
		concat(), objectsArray(
			&objects.Char{Value: 'y'},
			&objects.Float{Value: 93.11},
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpSetLocal, 0),
				compiler.MakeInstruction(compiler.OpGetGlobal, 0),
				compiler.MakeInstruction(compiler.OpGetFree, 0)),
			&objects.Float{Value: 39.2},
			&objects.Int{Value: 192},
			&objects.String{Value: "bar"})))

	testBytecodeSerialization(t, bytecodeFileSet(
		concat(
			compiler.MakeInstruction(compiler.OpConstant, 0),
			compiler.MakeInstruction(compiler.OpSetGlobal, 0),
			compiler.MakeInstruction(compiler.OpConstant, 6),
			compiler.MakeInstruction(compiler.OpPop)),
		objectsArray(
			&objects.Int{Value: 55},
			&objects.Int{Value: 66},
			&objects.Int{Value: 77},
			&objects.Int{Value: 88},
			&objects.ImmutableMap{
				Value: map[string]objects.Object{
					"array": &objects.ImmutableArray{
						Value: []objects.Object{
							&objects.Int{Value: 1},
							&objects.Int{Value: 2},
							&objects.Int{Value: 3},
							objects.TrueValue,
							objects.FalseValue,
							objects.UndefinedValue,
						},
					},
					"true":  objects.TrueValue,
					"false": objects.FalseValue,
					"bytes": &objects.Bytes{Value: make([]byte, 16)},
					"char":  &objects.Char{Value: 'Y'},
					"error": &objects.Error{Value: &objects.String{Value: "some error"}},
					"float": &objects.Float{Value: -19.84},
					"immutable_array": &objects.ImmutableArray{
						Value: []objects.Object{
							&objects.Int{Value: 1},
							&objects.Int{Value: 2},
							&objects.Int{Value: 3},
							objects.TrueValue,
							objects.FalseValue,
							objects.UndefinedValue,
						},
					},
					"immutable_map": &objects.ImmutableMap{
						Value: map[string]objects.Object{
							"a": &objects.Int{Value: 1},
							"b": &objects.Int{Value: 2},
							"c": &objects.Int{Value: 3},
							"d": objects.TrueValue,
							"e": objects.FalseValue,
							"f": objects.UndefinedValue,
						},
					},
					"int": &objects.Int{Value: 91},
					"map": &objects.Map{
						Value: map[string]objects.Object{
							"a": &objects.Int{Value: 1},
							"b": &objects.Int{Value: 2},
							"c": &objects.Int{Value: 3},
							"d": objects.TrueValue,
							"e": objects.FalseValue,
							"f": objects.UndefinedValue,
						},
					},
					"string":    &objects.String{Value: "foo bar"},
					"time":      &objects.Time{Value: time.Now()},
					"undefined": objects.UndefinedValue,
				},
			},
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
				compiler.MakeInstruction(compiler.OpReturnValue)),
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpSetLocal, 0),
				compiler.MakeInstruction(compiler.OpGetFree, 0),
				compiler.MakeInstruction(compiler.OpGetLocal, 0),
				compiler.MakeInstruction(compiler.OpClosure, 4, 2),
				compiler.MakeInstruction(compiler.OpReturnValue)),
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpSetLocal, 0),
				compiler.MakeInstruction(compiler.OpGetLocal, 0),
				compiler.MakeInstruction(compiler.OpClosure, 5, 1),
				compiler.MakeInstruction(compiler.OpReturnValue))),
		fileSet(srcfile{name: "file1", size: 100}, srcfile{name: "file2", size: 200})))
}

func TestBytecode_RemoveDuplicates(t *testing.T) {
	testBytecodeRemoveDuplicates(t,
		bytecode(
			concat(), objectsArray(
				&objects.Char{Value: 'y'},
				&objects.Float{Value: 93.11},
				compiledFunction(1, 0,
					compiler.MakeInstruction(compiler.OpConstant, 3),
					compiler.MakeInstruction(compiler.OpSetLocal, 0),
					compiler.MakeInstruction(compiler.OpGetGlobal, 0),
					compiler.MakeInstruction(compiler.OpGetFree, 0)),
				&objects.Float{Value: 39.2},
				&objects.Int{Value: 192},
				&objects.String{Value: "bar"})),
		bytecode(
			concat(), objectsArray(
				&objects.Char{Value: 'y'},
				&objects.Float{Value: 93.11},
				compiledFunction(1, 0,
					compiler.MakeInstruction(compiler.OpConstant, 3),
					compiler.MakeInstruction(compiler.OpSetLocal, 0),
					compiler.MakeInstruction(compiler.OpGetGlobal, 0),
					compiler.MakeInstruction(compiler.OpGetFree, 0)),
				&objects.Float{Value: 39.2},
				&objects.Int{Value: 192},
				&objects.String{Value: "bar"})))

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
				&objects.Int{Value: 1},
				&objects.Float{Value: 2.0},
				&objects.Char{Value: '3'},
				&objects.String{Value: "four"},
				compiledFunction(1, 0,
					compiler.MakeInstruction(compiler.OpConstant, 3),
					compiler.MakeInstruction(compiler.OpConstant, 7),
					compiler.MakeInstruction(compiler.OpSetLocal, 0),
					compiler.MakeInstruction(compiler.OpGetGlobal, 0),
					compiler.MakeInstruction(compiler.OpGetFree, 0)),
				&objects.Int{Value: 1},
				&objects.Float{Value: 2.0},
				&objects.Char{Value: '3'},
				&objects.String{Value: "four"})),
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
				&objects.Int{Value: 1},
				&objects.Float{Value: 2.0},
				&objects.Char{Value: '3'},
				&objects.String{Value: "four"},
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
				&objects.Int{Value: 1},
				&objects.Int{Value: 2},
				&objects.Int{Value: 3},
				&objects.Int{Value: 1},
				&objects.Int{Value: 3})),
		bytecode(
			concat(
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpConstant, 0),
				compiler.MakeInstruction(compiler.OpConstant, 2)),
			objectsArray(
				&objects.Int{Value: 1},
				&objects.Int{Value: 2},
				&objects.Int{Value: 3})))
}

func TestBytecode_CountObjects(t *testing.T) {
	b := bytecode(
		concat(),
		objectsArray(
			&objects.Int{Value: 55},
			&objects.Int{Value: 66},
			&objects.Int{Value: 77},
			&objects.Int{Value: 88},
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 3),
				compiler.MakeInstruction(compiler.OpReturnValue)),
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 2),
				compiler.MakeInstruction(compiler.OpReturnValue)),
			compiledFunction(1, 0,
				compiler.MakeInstruction(compiler.OpConstant, 1),
				compiler.MakeInstruction(compiler.OpReturnValue))))
	assert.Equal(t, 7, b.CountObjects())
}

func fileSet(files ...srcfile) *source.FileSet {
	fileSet := source.NewFileSet()
	for _, f := range files {
		fileSet.AddFile(f.name, -1, f.size)
	}
	return fileSet
}

func bytecodeFileSet(instructions []byte, constants []objects.Object, fileSet *source.FileSet) *compiler.Bytecode {
	return &compiler.Bytecode{
		FileSet:      fileSet,
		MainFunction: &objects.CompiledFunction{Instructions: instructions},
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
	err = r.Decode(bytes.NewReader(buf.Bytes()))
	assert.NoError(t, err)

	assert.Equal(t, b.FileSet, r.FileSet)
	assert.Equal(t, b.MainFunction, r.MainFunction)
	assert.Equal(t, b.Constants, r.Constants)
}

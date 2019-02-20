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
			objects.UndefinedValue,
			&objects.Time{Value: time.Now()},
			&objects.Array{
				Value: objectsArray(
					&objects.Int{Value: 12},
					&objects.String{Value: "foo"},
					objects.TrueValue,
					objects.FalseValue,
					&objects.Float{Value: 93.11},
					&objects.Char{Value: 'x'},
					objects.UndefinedValue,
				),
			},
			objects.FalseValue,
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
					"b": objects.FalseValue,
					"c": objects.UndefinedValue,
				},
			},
			&objects.String{Value: "bar"},
			objects.UndefinedValue)))

	testBytecodeSerialization(t, bytecodeFileSet(
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

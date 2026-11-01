package tengo_test

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/parser"
	"github.com/d5/tengo/v2/require"
)

type srcfile struct {
	name string
	size int
}

func TestBytecode(t *testing.T) {
	testBytecodeSerialization(t, bytecode(concatInsts(), objectsArray()))

	testBytecodeSerialization(t, bytecode(
		concatInsts(), objectsArray(
			&tengo.Char{Value: 'y'},
			&tengo.Float{Value: 93.11},
			compiledFunction(1, 0,
				tengo.MakeInstruction(parser.OpConstant, 3),
				tengo.MakeInstruction(parser.OpSetLocal, 0),
				tengo.MakeInstruction(parser.OpGetGlobal, 0),
				tengo.MakeInstruction(parser.OpGetFree, 0)),
			&tengo.Float{Value: 39.2},
			&tengo.Int{Value: 192},
			&tengo.String{Value: "bar"})))

	testBytecodeSerialization(t, bytecodeFileSet(
		concatInsts(
			tengo.MakeInstruction(parser.OpConstant, 0),
			tengo.MakeInstruction(parser.OpSetGlobal, 0),
			tengo.MakeInstruction(parser.OpConstant, 6),
			tengo.MakeInstruction(parser.OpPop)),
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
					"error": &tengo.Error{Value: &tengo.String{
						Value: "some error",
					}},
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
				tengo.MakeInstruction(parser.OpConstant, 3),
				tengo.MakeInstruction(parser.OpSetLocal, 0),
				tengo.MakeInstruction(parser.OpGetGlobal, 0),
				tengo.MakeInstruction(parser.OpGetFree, 0),
				tengo.MakeInstruction(parser.OpBinaryOp, 11),
				tengo.MakeInstruction(parser.OpGetFree, 1),
				tengo.MakeInstruction(parser.OpBinaryOp, 11),
				tengo.MakeInstruction(parser.OpGetLocal, 0),
				tengo.MakeInstruction(parser.OpBinaryOp, 11),
				tengo.MakeInstruction(parser.OpReturn, 1)),
			compiledFunction(1, 0,
				tengo.MakeInstruction(parser.OpConstant, 2),
				tengo.MakeInstruction(parser.OpSetLocal, 0),
				tengo.MakeInstruction(parser.OpGetFree, 0),
				tengo.MakeInstruction(parser.OpGetLocal, 0),
				tengo.MakeInstruction(parser.OpClosure, 4, 2),
				tengo.MakeInstruction(parser.OpReturn, 1)),
			compiledFunction(1, 0,
				tengo.MakeInstruction(parser.OpConstant, 1),
				tengo.MakeInstruction(parser.OpSetLocal, 0),
				tengo.MakeInstruction(parser.OpGetLocal, 0),
				tengo.MakeInstruction(parser.OpClosure, 5, 1),
				tengo.MakeInstruction(parser.OpReturn, 1))),
		fileSet(srcfile{name: "file1", size: 100},
			srcfile{name: "file2", size: 200})))
}

func TestBytecode_RemoveDuplicates(t *testing.T) {
	testBytecodeRemoveDuplicates(t,
		bytecode(
			concatInsts(), objectsArray(
				&tengo.Char{Value: 'y'},
				&tengo.Float{Value: 93.11},
				compiledFunction(1, 0,
					tengo.MakeInstruction(parser.OpConstant, 3),
					tengo.MakeInstruction(parser.OpSetLocal, 0),
					tengo.MakeInstruction(parser.OpGetGlobal, 0),
					tengo.MakeInstruction(parser.OpGetFree, 0)),
				&tengo.Float{Value: 39.2},
				&tengo.Int{Value: 192},
				&tengo.String{Value: "bar"})),
		bytecode(
			concatInsts(), objectsArray(
				&tengo.Char{Value: 'y'},
				&tengo.Float{Value: 93.11},
				compiledFunction(1, 0,
					tengo.MakeInstruction(parser.OpConstant, 3),
					tengo.MakeInstruction(parser.OpSetLocal, 0),
					tengo.MakeInstruction(parser.OpGetGlobal, 0),
					tengo.MakeInstruction(parser.OpGetFree, 0)),
				&tengo.Float{Value: 39.2},
				&tengo.Int{Value: 192},
				&tengo.String{Value: "bar"})))

	testBytecodeRemoveDuplicates(t,
		bytecode(
			concatInsts(
				tengo.MakeInstruction(parser.OpConstant, 0),
				tengo.MakeInstruction(parser.OpConstant, 1),
				tengo.MakeInstruction(parser.OpConstant, 2),
				tengo.MakeInstruction(parser.OpConstant, 3),
				tengo.MakeInstruction(parser.OpConstant, 4),
				tengo.MakeInstruction(parser.OpConstant, 5),
				tengo.MakeInstruction(parser.OpConstant, 6),
				tengo.MakeInstruction(parser.OpConstant, 7),
				tengo.MakeInstruction(parser.OpConstant, 8),
				tengo.MakeInstruction(parser.OpClosure, 4, 1)),
			objectsArray(
				&tengo.Int{Value: 1},
				&tengo.Float{Value: 2.0},
				&tengo.Char{Value: '3'},
				&tengo.String{Value: "four"},
				compiledFunction(1, 0,
					tengo.MakeInstruction(parser.OpConstant, 3),
					tengo.MakeInstruction(parser.OpConstant, 7),
					tengo.MakeInstruction(parser.OpSetLocal, 0),
					tengo.MakeInstruction(parser.OpGetGlobal, 0),
					tengo.MakeInstruction(parser.OpGetFree, 0)),
				&tengo.Int{Value: 1},
				&tengo.Float{Value: 2.0},
				&tengo.Char{Value: '3'},
				&tengo.String{Value: "four"})),
		bytecode(
			concatInsts(
				tengo.MakeInstruction(parser.OpConstant, 0),
				tengo.MakeInstruction(parser.OpConstant, 1),
				tengo.MakeInstruction(parser.OpConstant, 2),
				tengo.MakeInstruction(parser.OpConstant, 3),
				tengo.MakeInstruction(parser.OpConstant, 4),
				tengo.MakeInstruction(parser.OpConstant, 0),
				tengo.MakeInstruction(parser.OpConstant, 1),
				tengo.MakeInstruction(parser.OpConstant, 2),
				tengo.MakeInstruction(parser.OpConstant, 3),
				tengo.MakeInstruction(parser.OpClosure, 4, 1)),
			objectsArray(
				&tengo.Int{Value: 1},
				&tengo.Float{Value: 2.0},
				&tengo.Char{Value: '3'},
				&tengo.String{Value: "four"},
				compiledFunction(1, 0,
					tengo.MakeInstruction(parser.OpConstant, 3),
					tengo.MakeInstruction(parser.OpConstant, 2),
					tengo.MakeInstruction(parser.OpSetLocal, 0),
					tengo.MakeInstruction(parser.OpGetGlobal, 0),
					tengo.MakeInstruction(parser.OpGetFree, 0)))))

	testBytecodeRemoveDuplicates(t,
		bytecode(
			concatInsts(
				tengo.MakeInstruction(parser.OpConstant, 0),
				tengo.MakeInstruction(parser.OpConstant, 1),
				tengo.MakeInstruction(parser.OpConstant, 2),
				tengo.MakeInstruction(parser.OpConstant, 3),
				tengo.MakeInstruction(parser.OpConstant, 4)),
			objectsArray(
				&tengo.Int{Value: 1},
				&tengo.Int{Value: 2},
				&tengo.Int{Value: 3},
				&tengo.Int{Value: 1},
				&tengo.Int{Value: 3})),
		bytecode(
			concatInsts(
				tengo.MakeInstruction(parser.OpConstant, 0),
				tengo.MakeInstruction(parser.OpConstant, 1),
				tengo.MakeInstruction(parser.OpConstant, 2),
				tengo.MakeInstruction(parser.OpConstant, 0),
				tengo.MakeInstruction(parser.OpConstant, 2)),
			objectsArray(
				&tengo.Int{Value: 1},
				&tengo.Int{Value: 2},
				&tengo.Int{Value: 3})))
}

func TestBytecode_CountObjects(t *testing.T) {
	b := bytecode(
		concatInsts(),
		objectsArray(
			&tengo.Int{Value: 55},
			&tengo.Int{Value: 66},
			&tengo.Int{Value: 77},
			&tengo.Int{Value: 88},
			compiledFunction(1, 0,
				tengo.MakeInstruction(parser.OpConstant, 3),
				tengo.MakeInstruction(parser.OpReturn, 1)),
			compiledFunction(1, 0,
				tengo.MakeInstruction(parser.OpConstant, 2),
				tengo.MakeInstruction(parser.OpReturn, 1)),
			compiledFunction(1, 0,
				tengo.MakeInstruction(parser.OpConstant, 1),
				tengo.MakeInstruction(parser.OpReturn, 1))))
	require.Equal(t, 7, b.CountObjects())
}

func fileSet(files ...srcfile) *parser.SourceFileSet {
	fileSet := parser.NewFileSet()
	for _, f := range files {
		fileSet.AddFile(f.name, -1, f.size)
	}
	return fileSet
}

func bytecodeFileSet(
	instructions []byte,
	constants []tengo.Object,
	fileSet *parser.SourceFileSet,
) *tengo.Bytecode {
	return &tengo.Bytecode{
		FileSet:      fileSet,
		MainFunction: &tengo.CompiledFunction{Instructions: instructions},
		Constants:    constants,
	}
}

func testBytecodeRemoveDuplicates(
	t *testing.T,
	input, expected *tengo.Bytecode,
) {
	input.RemoveDuplicates()

	require.Equal(t, expected.FileSet, input.FileSet)
	require.Equal(t, expected.MainFunction, input.MainFunction)
	require.Equal(t, expected.Constants, input.Constants)
}

func testBytecodeSerialization(t *testing.T, b *tengo.Bytecode) {
	var buf bytes.Buffer
	err := b.Encode(&buf)
	require.NoError(t, err)

	r := &tengo.Bytecode{}
	err = r.Decode(bytes.NewReader(buf.Bytes()), nil)
	require.NoError(t, err)

	require.Equal(t, b.FileSet, r.FileSet)
	require.Equal(t, b.MainFunction, r.MainFunction)
	require.Equal(t, b.Constants, r.Constants)
}

func TestBytecode_Size(t *testing.T) {
	b := &tengo.Bytecode{
		FileSet: fileSet(srcfile{name: "a", size: 10},
			srcfile{name: "b", size: 20}),
		MainFunction: &tengo.CompiledFunction{
			Instructions: concatInsts(
				tengo.MakeInstruction(parser.OpConstant, 0),
				tengo.MakeInstruction(parser.OpPop)),
			SourceMap: map[int]parser.Pos{0: 1},
		},
		Constants: objectsArray(&tengo.Int{Value: 1}, &tengo.Int{Value: 2}),
	}
	require.Equal(t, int64(4+1+30+2), b.Size())
}

func TestBytecode_DecodeCanonicalizesSingletons(t *testing.T) {
	b := bytecode(concatInsts(), objectsArray(
		&tengo.ImmutableMap{Value: map[string]tengo.Object{
			"array": &tengo.Array{Value: []tengo.Object{
				tengo.TrueValue,
				tengo.FalseValue,
				tengo.UndefinedValue,
				&tengo.Map{Value: map[string]tengo.Object{
					"t": tengo.TrueValue,
				}},
			}},
		}}))

	var buf bytes.Buffer
	require.NoError(t, b.Encode(&buf))
	r := &tengo.Bytecode{}
	require.NoError(t, r.Decode(bytes.NewReader(buf.Bytes()), nil))
	arr := r.Constants[0].(*tengo.ImmutableMap).Value["array"].(*tengo.Array)
	require.True(t, arr.Value[0] == tengo.TrueValue)
	require.True(t, arr.Value[1] == tengo.FalseValue)
	require.True(t, arr.Value[2] == tengo.UndefinedValue)
	require.True(t, arr.Value[3].(*tengo.Map).Value["t"] == tengo.TrueValue)
}

func TestBytecode_DecodeBuiltinModule(t *testing.T) {
	mods := tengo.NewModuleMap()
	mods.AddBuiltinModule("mod", map[string]tengo.Object{
		"double": &tengo.UserFunction{
			Name: "double",
			Value: func(args ...tengo.Object) (tengo.Object, error) {
				return &tengo.Int{Value: args[0].(*tengo.Int).Value * 2}, nil
			},
		},
	})

	src := []byte(`out = import("mod").double(21)`)
	fileSet := parser.NewFileSet()
	srcFile := fileSet.AddFile("test", -1, len(src))
	file, err := parser.NewParser(srcFile, src, nil).ParseFile()
	require.NoError(t, err)
	symbols := tengo.NewSymbolTable()
	out := symbols.Define("out")
	c := tengo.NewCompiler(srcFile, symbols, nil, mods, nil)
	require.NoError(t, c.Compile(file))

	var buf bytes.Buffer
	require.NoError(t, c.Bytecode().Encode(&buf))

	err = (&tengo.Bytecode{}).Decode(bytes.NewReader(buf.Bytes()), nil)
	require.Error(t, err)
	require.Equal(t, "user function not decodable", err.Error())

	decoded := &tengo.Bytecode{}
	require.NoError(t, decoded.Decode(bytes.NewReader(buf.Bytes()), mods))
	globals := make([]tengo.Object, tengo.GlobalsSize)
	require.NoError(t, tengo.NewVM(decoded, globals, -1).Run())
	require.Equal(t, &tengo.Int{Value: 42}, globals[out.Index])
}

func TestBytecode_DecodeNestedUserFunction(t *testing.T) {
	module := func() tengo.Object {
		return &tengo.ImmutableMap{Value: map[string]tengo.Object{
			"fn": &tengo.UserFunction{Name: "fn"},
		}}
	}
	for _, c := range []tengo.Object{
		module(),
		&tengo.Array{Value: []tengo.Object{module()}},
		&tengo.ImmutableArray{Value: []tengo.Object{module()}},
		&tengo.Map{Value: map[string]tengo.Object{"m": module()}},
		&tengo.ImmutableMap{Value: map[string]tengo.Object{"m": module()}},
	} {
		var buf bytes.Buffer
		require.NoError(t, bytecode(concatInsts(), objectsArray(c)).Encode(&buf))
		err := (&tengo.Bytecode{}).Decode(bytes.NewReader(buf.Bytes()), nil)
		require.Error(t, err, c.TypeName())
		require.Equal(t, "user function not decodable", err.Error())
	}
}

func TestBytecode_DecodeTruncated(t *testing.T) {
	b := bytecode(
		concatInsts(
			tengo.MakeInstruction(parser.OpConstant, 0),
			tengo.MakeInstruction(parser.OpPop)),
		objectsArray(&tengo.String{Value: "foo"}))
	var buf bytes.Buffer
	require.NoError(t, b.Encode(&buf))
	data := buf.Bytes()

	for n := 0; n < len(data); n++ {
		err := (&tengo.Bytecode{}).Decode(bytes.NewReader(data[:n]), nil)
		require.Error(t, err, "decoding %d of %d bytes", n, len(data))
	}
	require.NoError(t,
		(&tengo.Bytecode{}).Decode(bytes.NewReader(data), nil))
}

func TestBytecode_EncodeWriteError(t *testing.T) {
	writeErr := errors.New("write failed")
	b := bytecode(concatInsts(), objectsArray())
	err := b.Encode(failingWriter{err: writeErr})
	require.True(t, errors.Is(err, writeErr), "%v", err)
}

type failingWriter struct {
	err error
}

func (w failingWriter) Write([]byte) (int, error) {
	return 0, w.err
}

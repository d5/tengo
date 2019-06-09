package runtime_test

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/d5/tengo"
	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/compiler/ast"
	"github.com/d5/tengo/compiler/parser"
	"github.com/d5/tengo/compiler/source"
	"github.com/d5/tengo/runtime"
	"github.com/d5/tengo/stdlib"
)

const testOut = "out"

type IARR []interface{}
type IMAP map[string]interface{}
type MAP = map[string]interface{}
type ARR = []interface{}

type testopts struct {
	modules     *tengo.ModuleMap
	symbols     map[string]tengo.Object
	maxAllocs   int64
	skip2ndPass bool
}

func Opts() *testopts {
	return &testopts{
		modules:     tengo.NewModuleMap(),
		symbols:     make(map[string]tengo.Object),
		maxAllocs:   -1,
		skip2ndPass: false,
	}
}

func (o *testopts) copy() *testopts {
	c := &testopts{
		modules:     o.modules.Copy(),
		symbols:     make(map[string]tengo.Object),
		maxAllocs:   o.maxAllocs,
		skip2ndPass: o.skip2ndPass,
	}
	for k, v := range o.symbols {
		c.symbols[k] = v
	}
	return c
}

func (o *testopts) Stdlib() *testopts {
	o.modules.AddMap(stdlib.GetModuleMap(stdlib.AllModuleNames()...))
	return o
}

func (o *testopts) Module(name string, mod interface{}) *testopts {
	c := o.copy()
	switch mod := mod.(type) {
	case tengo.Importable:
		c.modules.Add(name, mod)
	case string:
		c.modules.AddSourceModule(name, []byte(mod))
	case []byte:
		c.modules.AddSourceModule(name, mod)
	default:
		panic(fmt.Errorf("invalid module type: %T", mod))
	}
	return c
}

func (o *testopts) Symbol(name string, value tengo.Object) *testopts {
	c := o.copy()
	c.symbols[name] = value
	return c
}

func (o *testopts) MaxAllocs(limit int64) *testopts {
	c := o.copy()
	c.maxAllocs = limit
	return c
}

func (o *testopts) Skip2ndPass() *testopts {
	c := o.copy()
	c.skip2ndPass = true
	return c
}

func expect(t *testing.T, input string, opts *testopts, expected interface{}) {
	if opts == nil {
		opts = Opts()
	}

	symbols := opts.symbols
	modules := opts.modules
	maxAllocs := opts.maxAllocs

	expectedObj := toObject(expected)

	if symbols == nil {
		symbols = make(map[string]tengo.Object)
	}
	symbols[testOut] = objectZeroCopy(expectedObj)

	// first pass: run the code normally
	{
		// parse
		file := parse(t, input)
		if file == nil {
			return
		}

		// compiler/VM
		res, trace, err := traceCompileRun(file, symbols, modules, maxAllocs)
		if !assert.NoError(t, err) ||
			!assert.Equal(t, expectedObj, res[testOut]) {
			t.Log("\n" + strings.Join(trace, "\n") + "\n")
		}
	}

	// second pass: run the code as import module
	if !opts.skip2ndPass {
		file := parse(t, `out = import("__code__")`)
		if file == nil {
			return
		}

		expectedObj := toObject(expected)
		switch eo := expectedObj.(type) {
		case *tengo.Array:
			expectedObj = &tengo.ImmutableArray{Value: eo.Value}
		case *tengo.Map:
			expectedObj = &tengo.ImmutableMap{Value: eo.Value}
		}

		modules.AddSourceModule("__code__", []byte(fmt.Sprintf("out := undefined; %s; export out", input)))

		res, trace, err := traceCompileRun(file, symbols, modules, maxAllocs)
		if !assert.NoError(t, err) ||
			!assert.Equal(t, expectedObj, res[testOut]) {
			t.Log("\n" + strings.Join(trace, "\n") + "\n")
		}
	}
}

func expectError(t *testing.T, input string, opts *testopts, expected string) {
	if opts == nil {
		opts = Opts()
	}

	symbols := opts.symbols
	modules := opts.modules
	maxAllocs := opts.maxAllocs

	expected = strings.TrimSpace(expected)
	if expected == "" {
		panic("expected must not be empty")
	}

	// parse
	program := parse(t, input)
	if program == nil {
		return
	}

	// compiler/VM
	_, trace, err := traceCompileRun(program, symbols, modules, maxAllocs)
	if !assert.Error(t, err) ||
		!assert.True(t, strings.Contains(err.Error(), expected), "expected error string: %s, got: %s", expected, err.Error()) {
		t.Log("\n" + strings.Join(trace, "\n") + "\n")
	}
}

func expectPanic(t *testing.T, input string, opts *testopts, expected string) {
	if opts == nil {
		opts = Opts()
	}

	symbols := opts.symbols
	modules := opts.modules
	maxAllocs := opts.maxAllocs

	expected = strings.TrimSpace(expected)
	if expected == "" {
		panic("expected must not be empty")
	}

	// parse
	program := parse(t, input)
	if program == nil {
		return
	}

	// compiler/VM
	var trace []string
	defer func() {
		e := recover()
		if !assert.True(t, strings.Contains(fmt.Sprint(e), expected),
			"expected panic string: %s, got: %s", expected, fmt.Sprint(e)) {
			t.Log("\n" + strings.Join(trace, "\n") + "\n")
		}
	}()
	_, trace, _ = traceCompileRun(program, symbols, modules, maxAllocs)
}

type tracer struct {
	Out []string
}

func (o *tracer) Write(p []byte) (n int, err error) {
	o.Out = append(o.Out, string(p))
	return len(p), nil
}

func traceCompileRun(file *ast.File, symbols map[string]tengo.Object, modules *tengo.ModuleMap, maxAllocs int64) (res map[string]tengo.Object, trace []string, err error) {
	var v *runtime.VM

	//defer func() {
	//	if e := recover(); e != nil {
	//		err = fmt.Errorf("panic: %v", e)
	//
	//		// stack trace
	//		var stackTrace []string
	//		for i := 2; ; i += 1 {
	//			_, file, line, ok := _runtime.Caller(i)
	//			if !ok {
	//				break
	//			}
	//			stackTrace = append(stackTrace, fmt.Sprintf("  %s:%d", file, line))
	//		}
	//
	//		trace = append(trace, fmt.Sprintf("\n[Error Trace]\n\n  %s\n", strings.Join(stackTrace, "\n  ")))
	//	}
	//}()

	globals := make([]tengo.Object, runtime.GlobalsSize)

	symTable := compiler.NewSymbolTable()
	for name, value := range symbols {
		sym := symTable.Define(name)

		// should not store pointer to 'value' variable
		// which is re-used in each iteration.
		valueCopy := value
		globals[sym.Index] = valueCopy
	}
	for idx, fn := range tengo.Builtins {
		symTable.DefineBuiltin(idx, fn.Name)
	}

	//tr := &tracer{}
	c := compiler.NewCompiler(file.InputFile, symTable, nil, modules, nil)
	err = c.Compile(file)
	//trace = append(trace, fmt.Sprintf("\n[Compiler Trace]\n\n%s", strings.Join(tr.Out, "")))
	if err != nil {
		return
	}

	bytecode := c.Bytecode()
	bytecode.RemoveDuplicates()
	trace = append(trace, fmt.Sprintf("\n[Compiled Constants]\n\n%s", strings.Join(bytecode.FormatConstants(), "\n")))
	trace = append(trace, fmt.Sprintf("\n[Compiled Instructions]\n\n%s", strings.Join(bytecode.FormatInstructions(), "\n")))

	v = runtime.NewVM(bytecode, globals, maxAllocs)

	err = v.Run()
	{
		res = make(map[string]tengo.Object)
		for name := range symbols {
			sym, depth, ok := symTable.Resolve(name)
			if !ok || depth != 0 {
				err = fmt.Errorf("symbol not found: %s", name)
				return
			}

			res[name] = globals[sym.Index]
		}
		trace = append(trace, fmt.Sprintf("\n[Globals]\n\n%s", strings.Join(formatGlobals(globals), "\n")))
	}
	if err == nil && !v.IsStackEmpty() {
		err = errors.New("non empty stack after execution")
	}

	return
}

func formatGlobals(globals []tengo.Object) (formatted []string) {
	for idx, global := range globals {
		if global == nil {
			return
		}

		formatted = append(formatted, fmt.Sprintf("[% 3d] %s (%s|%p)", idx, global.String(), reflect.TypeOf(global).Elem().Name(), global))
	}

	return
}

func parse(t *testing.T, input string) *ast.File {
	testFileSet := source.NewFileSet()
	testFile := testFileSet.AddFile("test", -1, len(input))

	p := parser.NewParser(testFile, []byte(input), nil)
	file, err := p.ParseFile()
	if !assert.NoError(t, err) {
		return nil
	}

	return file
}

func errorObject(v interface{}) *tengo.Error {
	return &tengo.Error{Value: toObject(v)}
}

func toObject(v interface{}) tengo.Object {
	switch v := v.(type) {
	case tengo.Object:
		return v
	case string:
		return &tengo.String{Value: v}
	case int64:
		return &tengo.Int{Value: v}
	case int: // for convenience
		return &tengo.Int{Value: int64(v)}
	case bool:
		if v {
			return tengo.TrueValue
		}
		return tengo.FalseValue
	case rune:
		return &tengo.Char{Value: v}
	case byte: // for convenience
		return &tengo.Char{Value: rune(v)}
	case float64:
		return &tengo.Float{Value: v}
	case []byte:
		return &tengo.Bytes{Value: v}
	case MAP:
		objs := make(map[string]tengo.Object)
		for k, v := range v {
			objs[k] = toObject(v)
		}

		return &tengo.Map{Value: objs}
	case ARR:
		var objs []tengo.Object
		for _, e := range v {
			objs = append(objs, toObject(e))
		}

		return &tengo.Array{Value: objs}
	case IMAP:
		objs := make(map[string]tengo.Object)
		for k, v := range v {
			objs[k] = toObject(v)
		}

		return &tengo.ImmutableMap{Value: objs}
	case IARR:
		var objs []tengo.Object
		for _, e := range v {
			objs = append(objs, toObject(e))
		}

		return &tengo.ImmutableArray{Value: objs}
	}

	panic(fmt.Errorf("unknown type: %T", v))
}

func objectZeroCopy(o tengo.Object) tengo.Object {
	switch o.(type) {
	case *tengo.Int:
		return &tengo.Int{}
	case *tengo.Float:
		return &tengo.Float{}
	case *tengo.Bool:
		return &tengo.Bool{}
	case *tengo.Char:
		return &tengo.Char{}
	case *tengo.String:
		return &tengo.String{}
	case *tengo.Array:
		return &tengo.Array{}
	case *tengo.Map:
		return &tengo.Map{}
	case *tengo.Undefined:
		return tengo.UndefinedValue
	case *tengo.Error:
		return &tengo.Error{}
	case *tengo.Bytes:
		return &tengo.Bytes{}
	case *tengo.ImmutableArray:
		return &tengo.ImmutableArray{}
	case *tengo.ImmutableMap:
		return &tengo.ImmutableMap{}
	case nil:
		panic("nil")
	default:
		panic(fmt.Errorf("unknown object type: %s", o.TypeName()))
	}
}

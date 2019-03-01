package runtime_test

import (
	"fmt"
	"reflect"
	_runtime "runtime"
	"strings"
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/compiler/ast"
	"github.com/d5/tengo/compiler/parser"
	"github.com/d5/tengo/compiler/source"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/runtime"
)

const (
	testOut = "out"
)

type IARR []interface{}
type IMAP map[string]interface{}
type MAP = map[string]interface{}
type ARR = []interface{}
type SYM = map[string]objects.Object

func expect(t *testing.T, input string, expected interface{}) {
	expectWithUserModules(t, input, expected, nil)
}

func expectNoMod(t *testing.T, input string, expected interface{}) {
	runVM(t, input, expected, nil, nil, true)
}

func expectWithSymbols(t *testing.T, input string, expected interface{}, symbols map[string]objects.Object) {
	runVM(t, input, expected, symbols, nil, true)
}

func expectWithUserModules(t *testing.T, input string, expected interface{}, userModules map[string]string) {
	runVM(t, input, expected, nil, userModules, false)
}

func expectError(t *testing.T, input, expected string) {
	runVMError(t, input, nil, nil, expected)
}

func expectErrorWithUserModules(t *testing.T, input string, userModules map[string]string, expected string) {
	runVMError(t, input, nil, userModules, expected)
}

func expectErrorWithSymbols(t *testing.T, input string, symbols map[string]objects.Object, expected string) {
	runVMError(t, input, symbols, nil, expected)
}

func runVM(t *testing.T, input string, expected interface{}, symbols map[string]objects.Object, userModules map[string]string, skipModuleTest bool) {
	expectedObj := toObject(expected)

	if symbols == nil {
		symbols = make(map[string]objects.Object)
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
		res, trace, err := traceCompileRun(file, symbols, userModules)
		if !assert.NoError(t, err) ||
			!assert.Equal(t, expectedObj, res[testOut]) {
			t.Log("\n" + strings.Join(trace, "\n"))
		}
	}

	// second pass: run the code as import module
	if !skipModuleTest {
		file := parse(t, `out = import("__code__")`)
		if file == nil {
			return
		}

		expectedObj := toObject(expected)
		switch eo := expectedObj.(type) {
		case *objects.Array:
			expectedObj = &objects.ImmutableArray{Value: eo.Value}
		case *objects.Map:
			expectedObj = &objects.ImmutableMap{Value: eo.Value}
		}

		if userModules == nil {
			userModules = make(map[string]string)
		}
		userModules["__code__"] = fmt.Sprintf("out := undefined; %s; export out", input)

		res, trace, err := traceCompileRun(file, symbols, userModules)
		if !assert.NoError(t, err) ||
			!assert.Equal(t, expectedObj, res[testOut]) {
			t.Log("\n" + strings.Join(trace, "\n"))
		}
	}
}

func runVMError(t *testing.T, input string, symbols map[string]objects.Object, userModules map[string]string, expected string) {
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
	_, trace, err := traceCompileRun(program, symbols, userModules)
	if !assert.Error(t, err) ||
		!assert.True(t, strings.Contains(err.Error(), expected), "expected error string: %s, got: %s", expected, err.Error()) {
		t.Log("\n" + strings.Join(trace, "\n"))
	}
}

type tracer struct {
	Out []string
}

func (o *tracer) Write(p []byte) (n int, err error) {
	o.Out = append(o.Out, string(p))
	return len(p), nil
}

func traceCompileRun(file *ast.File, symbols map[string]objects.Object, userModules map[string]string) (res map[string]objects.Object, trace []string, err error) {
	var v *runtime.VM

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic: %v", e)

			// stack trace
			var stackTrace []string
			for i := 2; ; i += 1 {
				_, file, line, ok := _runtime.Caller(i)
				if !ok {
					break
				}
				stackTrace = append(stackTrace, fmt.Sprintf("  %s:%d", file, line))
			}

			trace = append(trace, fmt.Sprintf("[Error Trace]\n\n  %s\n", strings.Join(stackTrace, "\n  ")))
		}
	}()

	globals := make([]*objects.Object, runtime.GlobalsSize)

	symTable := compiler.NewSymbolTable()
	for name, value := range symbols {
		sym := symTable.Define(name)

		// should not store pointer to 'value' variable
		// which is re-used in each iteration.
		valueCopy := value
		globals[sym.Index] = &valueCopy
	}
	for idx, fn := range objects.Builtins {
		symTable.DefineBuiltin(idx, fn.Name)
	}

	tr := &tracer{}
	c := compiler.NewCompiler(file.InputFile, symTable, nil, nil, tr)
	c.SetModuleLoader(func(moduleName string) ([]byte, error) {
		if src, ok := userModules[moduleName]; ok {
			return []byte(src), nil
		}

		return nil, fmt.Errorf("module '%s' not found", moduleName)
	})
	err = c.Compile(file)
	trace = append(trace, fmt.Sprintf("\n[Compiler Trace]\n\n%s", strings.Join(tr.Out, "")))
	if err != nil {
		return
	}

	bytecode := c.Bytecode()
	trace = append(trace, fmt.Sprintf("\n[Compiled Constants]\n\n%s", strings.Join(bytecode.FormatConstants(), "\n")))
	trace = append(trace, fmt.Sprintf("\n[Compiled Instructions]\n\n%s\n", strings.Join(bytecode.FormatInstructions(), "\n")))

	v = runtime.NewVM(bytecode, globals, nil, nil)

	err = v.Run()
	{
		res = make(map[string]objects.Object)
		for name := range symbols {
			sym, depth, ok := symTable.Resolve(name)
			if !ok || depth != 0 {
				err = fmt.Errorf("symbol not found: %s", name)
				return
			}

			res[name] = *globals[sym.Index]
		}
		trace = append(trace, fmt.Sprintf("\n[Globals]\n\n%s", strings.Join(formatGlobals(globals), "\n")))
	}
	if err != nil {
		return
	}

	return
}

func formatGlobals(globals []*objects.Object) (formatted []string) {
	for idx, global := range globals {
		if global == nil {
			return
		}

		switch global := (*global).(type) {
		case *objects.Closure:
			formatted = append(formatted, fmt.Sprintf("[% 3d] (Closure|%p)", idx, global))
			for _, l := range compiler.FormatInstructions(global.Fn.Instructions, 0) {
				formatted = append(formatted, fmt.Sprintf("     %s", l))
			}
		default:
			formatted = append(formatted, fmt.Sprintf("[% 3d] %s (%s|%p)", idx, global.String(), reflect.TypeOf(global).Elem().Name(), global))
		}
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

func errorObject(v interface{}) *objects.Error {
	return &objects.Error{Value: toObject(v)}
}

func toObject(v interface{}) objects.Object {
	switch v := v.(type) {
	case objects.Object:
		return v
	case string:
		return &objects.String{Value: v}
	case int64:
		return &objects.Int{Value: v}
	case int: // for convenience
		return &objects.Int{Value: int64(v)}
	case bool:
		if v {
			return objects.TrueValue
		}
		return objects.FalseValue
	case rune:
		return &objects.Char{Value: v}
	case byte: // for convenience
		return &objects.Char{Value: rune(v)}
	case float64:
		return &objects.Float{Value: v}
	case []byte:
		return &objects.Bytes{Value: v}
	case MAP:
		objs := make(map[string]objects.Object)
		for k, v := range v {
			objs[k] = toObject(v)
		}

		return &objects.Map{Value: objs}
	case ARR:
		var objs []objects.Object
		for _, e := range v {
			objs = append(objs, toObject(e))
		}

		return &objects.Array{Value: objs}
	case IMAP:
		objs := make(map[string]objects.Object)
		for k, v := range v {
			objs[k] = toObject(v)
		}

		return &objects.ImmutableMap{Value: objs}
	case IARR:
		var objs []objects.Object
		for _, e := range v {
			objs = append(objs, toObject(e))
		}

		return &objects.ImmutableArray{Value: objs}
	}

	panic(fmt.Errorf("unknown type: %T", v))
}

func objectZeroCopy(o objects.Object) objects.Object {
	switch o.(type) {
	case *objects.Int:
		return &objects.Int{}
	case *objects.Float:
		return &objects.Float{}
	case *objects.Bool:
		return &objects.Bool{}
	case *objects.Char:
		return &objects.Char{}
	case *objects.String:
		return &objects.String{}
	case *objects.Array:
		return &objects.Array{}
	case *objects.Map:
		return &objects.Map{}
	case *objects.Undefined:
		return objects.UndefinedValue
	case *objects.Error:
		return &objects.Error{}
	case *objects.Bytes:
		return &objects.Bytes{}
	case *objects.ImmutableArray:
		return &objects.ImmutableArray{}
	case *objects.ImmutableMap:
		return &objects.ImmutableMap{}
	case nil:
		panic("nil")
	default:
		panic(fmt.Errorf("unknown object type: %s", o.TypeName()))
	}
}

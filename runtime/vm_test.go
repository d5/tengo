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

func expectWithSymbols(t *testing.T, input string, expected interface{}, symbols map[string]objects.Object) {
	// parse
	file := parse(t, input)
	if file == nil {
		return
	}

	// compiler/VM
	runVM(t, file, expected, symbols, nil)
}

func expectWithUserModules(t *testing.T, input string, expected interface{}, userModules map[string]string) {
	// parse
	file := parse(t, input)
	if file == nil {
		return
	}

	// compiler/VM
	runVM(t, file, expected, nil, userModules)
}

func expectError(t *testing.T, input string) {
	expectErrorWithUserModules(t, input, nil)
}

func expectErrorWithUserModules(t *testing.T, input string, userModules map[string]string) {
	// parse
	program := parse(t, input)
	if program == nil {
		return
	}

	// compiler/VM
	runVMError(t, program, nil, userModules)
}

func expectErrorWithSymbols(t *testing.T, input string, symbols map[string]objects.Object) {
	// parse
	program := parse(t, input)
	if program == nil {
		return
	}

	// compiler/VM
	runVMError(t, program, symbols, nil)
}

func runVM(t *testing.T, file *ast.File, expected interface{}, symbols map[string]objects.Object, userModules map[string]string) (ok bool) {
	expectedObj := toObject(expected)

	if symbols == nil {
		symbols = make(map[string]objects.Object)
	}
	symbols[testOut] = objectZeroCopy(expectedObj)

	res, trace, err := traceCompileRun(file, symbols, userModules)

	defer func() {
		if !ok {
			t.Log("\n" + strings.Join(trace, "\n"))
		}
	}()

	if !assert.NoError(t, err) {
		return
	}

	ok = assert.Equal(t, expectedObj, res[testOut])

	return
}

// TODO: should differentiate compile-time error, runtime error, and, error object returned
func runVMError(t *testing.T, file *ast.File, symbols map[string]objects.Object, userModules map[string]string) (ok bool) {
	_, trace, err := traceCompileRun(file, symbols, userModules)

	defer func() {
		if !ok {
			t.Log("\n" + strings.Join(trace, "\n"))
		}
	}()

	ok = assert.Error(t, err)

	return
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
		return &objects.Bool{Value: v}
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

			var ipstr string
			if v != nil {
				frameIdx, ip := v.FrameInfo()
				ipstr = fmt.Sprintf("\n  (Frame=%d, IP=%d)", frameIdx, ip+1)
			}
			trace = append(trace, fmt.Sprintf("[Panic]\n\n  %v%s\n", e, ipstr))
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
	c := compiler.NewCompiler(symTable, nil, tr)
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
	var constStr []string
	for cidx, cn := range bytecode.Constants {
		if cmFn, ok := cn.(*objects.CompiledFunction); ok {
			constStr = append(constStr, fmt.Sprintf("[% 3d] (Compiled Function|%p)", cidx, &cn))
			for _, l := range compiler.FormatInstructions(cmFn.Instructions, 0) {
				constStr = append(constStr, fmt.Sprintf("     %s", l))
			}
		} else {
			constStr = append(constStr, fmt.Sprintf("[% 3d] %s (%s|%p)", cidx, cn, reflect.TypeOf(cn).Elem().Name(), &cn))
		}
	}
	trace = append(trace, fmt.Sprintf("\n[Compiled Constants]\n\n%s", strings.Join(constStr, "\n")))
	trace = append(trace, fmt.Sprintf("\n[Compiled Instructions]\n\n%s\n", strings.Join(compiler.FormatInstructions(bytecode.Instructions, 0), "\n")))

	v = runtime.NewVM(bytecode, globals)

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
		var globalsStr []string
		for gidx, g := range globals {
			if g == nil {
				break
			}

			if cmFn, ok := (*g).(*objects.Closure); ok {
				globalsStr = append(globalsStr, fmt.Sprintf("[% 3d] (Closure|%p)", gidx, g))
				for _, l := range compiler.FormatInstructions(cmFn.Fn.Instructions, 0) {
					globalsStr = append(globalsStr, fmt.Sprintf("     %s", l))
				}
			} else {
				globalsStr = append(globalsStr, fmt.Sprintf("[% 3d] %s (%s|%p)", gidx, (*g).String(), reflect.TypeOf(*g).Elem().Name(), g))
			}
		}
		trace = append(trace, fmt.Sprintf("\n[Globals]\n\n%s", strings.Join(globalsStr, "\n")))

		frameIdx, ip := v.FrameInfo()
		trace = append(trace, fmt.Sprintf("\n[IP]\n\nFrame=%d, IP=%d", frameIdx, ip+1))
	}
	if err != nil {
		return
	}

	return
}

func parse(t *testing.T, input string) *ast.File {
	testFileSet := source.NewFileSet()
	testFile := testFileSet.AddFile("", -1, len(input))

	file, err := parser.ParseFile(testFile, []byte(input), nil)
	if !assert.NoError(t, err) {
		return nil
	}

	return file
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
		return &objects.Undefined{}
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

func undefined() *objects.Undefined {
	return &objects.Undefined{}
}

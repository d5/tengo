package stdlib_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/d5/tengo"
	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/script"
	"github.com/d5/tengo/stdlib"
)

type ARR = []interface{}
type MAP = map[string]interface{}
type IARR []interface{}
type IMAP map[string]interface{}

func TestAllModuleNames(t *testing.T) {
	names := stdlib.AllModuleNames()
	if !assert.Equal(t, len(stdlib.BuiltinModules), len(names)) {
		return
	}
}

func TestModulesRun(t *testing.T) {
	// os.File
	expect(t, `
os := import("os")
out := ""

write_file := func(filename, data) {
	file := os.create(filename)
	if !file { return file }

	if res := file.write(bytes(data)); is_error(res) {
		return res
	}

	return file.close()
}

read_file := func(filename) {
	file := os.open(filename)
	if !file { return file }

	data := bytes(100)
	cnt := file.read(data)
	if  is_error(cnt) {
		return cnt
	}

	file.close()
	return data[:cnt]
}

if write_file("./temp", "foobar") {
	out = string(read_file("./temp"))
}

os.remove("./temp")
`, "foobar")

	// exec.command
	expect(t, `
out := ""
os := import("os")
cmd := os.exec("echo", "foo", "bar")
if !is_error(cmd) {
	out = cmd.output()
}
`, []byte("foo bar\n"))

}

func TestGetModules(t *testing.T) {
	mods := stdlib.GetModuleMap()
	assert.Equal(t, 0, mods.Len())

	mods = stdlib.GetModuleMap("os")
	assert.Equal(t, 1, mods.Len())
	assert.NotNil(t, mods.Get("os"))

	mods = stdlib.GetModuleMap("os", "rand")
	assert.Equal(t, 2, mods.Len())
	assert.NotNil(t, mods.Get("os"))
	assert.NotNil(t, mods.Get("rand"))

	mods = stdlib.GetModuleMap("text", "text")
	assert.Equal(t, 1, mods.Len())
	assert.NotNil(t, mods.Get("text"))

	mods = stdlib.GetModuleMap("nonexisting", "text")
	assert.Equal(t, 1, mods.Len())
	assert.NotNil(t, mods.Get("text"))
}

type callres struct {
	t *testing.T
	o interface{}
	e error
}

func (c callres) call(funcName string, rt tengo.Interop, args ...interface{}) callres {
	if c.e != nil {
		return c
	}

	var oargs []tengo.Object
	for _, v := range args {
		oargs = append(oargs, object(v))
	}

	switch o := c.o.(type) {
	case *tengo.BuiltinModule:
		m, ok := o.Attrs[funcName]
		if !ok {
			return callres{t: c.t, e: fmt.Errorf("function not found: %s", funcName)}
		}

		f, ok := m.(*tengo.GoFunction)
		if !ok {
			return callres{t: c.t, e: fmt.Errorf("non-callable: %s", funcName)}
		}

		res, err := f.Value(rt, oargs...)
		return callres{t: c.t, o: res, e: err}
	case *tengo.GoFunction:
		res, err := o.Value(rt, oargs...)
		return callres{t: c.t, o: res, e: err}
	case *tengo.ImmutableMap:
		m, ok := o.Value[funcName]
		if !ok {
			return callres{t: c.t, e: fmt.Errorf("function not found: %s", funcName)}
		}

		f, ok := m.(*tengo.GoFunction)
		if !ok {
			return callres{t: c.t, e: fmt.Errorf("non-callable: %s", funcName)}
		}

		res, err := f.Value(rt, oargs...)
		return callres{t: c.t, o: res, e: err}
	default:
		panic(fmt.Errorf("unexpected object: %v (%T)", o, o))
	}
}

func (c callres) expect(expected interface{}, msgAndArgs ...interface{}) bool {
	return assert.NoError(c.t, c.e, msgAndArgs...) &&
		assert.Equal(c.t, object(expected), c.o, msgAndArgs...)
}

func (c callres) expectError() bool {
	return assert.Error(c.t, c.e)
}

func module(t *testing.T, moduleName string) callres {
	mod := stdlib.GetModuleMap(moduleName).GetBuiltinModule(moduleName)
	if mod == nil {
		return callres{t: t, e: fmt.Errorf("module not found: %s", moduleName)}
	}

	return callres{t: t, o: mod}
}

func object(v interface{}) tengo.Object {
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
			objs[k] = object(v)
		}

		return &tengo.Map{Value: objs}
	case ARR:
		var objs []tengo.Object
		for _, e := range v {
			objs = append(objs, object(e))
		}

		return &tengo.Array{Value: objs}
	case IMAP:
		objs := make(map[string]tengo.Object)
		for k, v := range v {
			objs[k] = object(v)
		}

		return &tengo.ImmutableMap{Value: objs}
	case IARR:
		var objs []tengo.Object
		for _, e := range v {
			objs = append(objs, object(e))
		}

		return &tengo.ImmutableArray{Value: objs}
	case time.Time:
		return &tengo.Time{Value: v}
	case []int:
		var objs []tengo.Object
		for _, e := range v {
			objs = append(objs, &tengo.Int{Value: int64(e)})
		}

		return &tengo.Array{Value: objs}
	}

	panic(fmt.Errorf("unknown type: %T", v))
}

func expect(t *testing.T, input string, expected interface{}) {
	s := script.New([]byte(input))
	s.SetImports(stdlib.GetModuleMap(stdlib.AllModuleNames()...))
	c, err := s.Run()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	v := c.Get("out")
	if !assert.NotNil(t, v) {
		return
	}

	assert.Equal(t, expected, v.Value())
}

type mockInterop struct {
	val tengo.Object
	err error
}

func (i mockInterop) InteropCall(
	callable tengo.Object,
	args ...tengo.Object,
) (tengo.Object, error) {
	return i.val, i.err
}

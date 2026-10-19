package stdlib_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/require"
	"github.com/d5/tengo/v2/stdlib"
)

type ARR = []interface{}
type MAP = map[string]interface{}
type IARR []interface{}
type IMAP map[string]interface{}

func TestAllModuleNames(t *testing.T) {
	names := stdlib.AllModuleNames()
	require.Equal(t,
		len(stdlib.BuiltinModules)+len(stdlib.SourceModules),
		len(names))
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

func TestSortModule(t *testing.T) {
	// normal
	expect(t, `
sort := import("sort")
a := [4, 5, 3, 1, 2]
a = sort.sort(a, func(i, j) {
	return a[i] < a[j]
})
out := import("json").encode(a)
`, []byte("[1,2,3,4,5]"))

	// generate random sequences for sorting
	rand.Seed(time.Now().UnixNano())
	generateRandomSequence := func() []int {
		length := 1 + rand.Intn(1000)
		arr := make([]int, length)
		for i := 0; i < length; i++ {
			arr[i] = rand.Intn(1000000)
		}
		return arr
	}
	for i := 0; i < 500; i++ {
		seq := generateRandomSequence()
		seqBytes, _ := json.Marshal(seq)
		sort.Ints(seq)
		sortResult, _ := json.Marshal(seq)
		expect(t, fmt.Sprintf(`
sort := import("sort")
a := %s
a = sort.sort(a, func(i, j) {
	return a[i] < a[j]
})
out := import("json").encode(a)
`, string(seqBytes)), sortResult)
	}

	// less is not a function
	expectErr(t, `
sort := import("sort")
a := [4, 5, 3, 1, 2]
a = sort.sort(a, 0)
out := import("json").encode(a)`, "Runtime Error: not callable: int")

	// arr is not an array
	expectErr(t, `
sort := import("sort")
a := 12345
a = sort.sort(a, func(i, j) {
	return a[i] < a[j]
})
out := import("json").encode(a)`, "Runtime Error: invalid type for argument 'first' in call to 'builtin-function:len': expected array/s")

	// empty array
	expect(t, `
sort := import("sort")
a := []
a = sort.sort(a, func(i, j) {
	return a[i] < a[j]
})
out := import("json").encode(a)`, []byte("[]"))

	// sort json
	expect(t, `
sort := import("sort")
a := [{"age": 12, "name": "A"}, {"age": 18, "name": "B"}, {"age": 9, "name": "C"}, {"age": 10, "name": "D"}, {"age": 21, "name": "E"}]
a = sort.sort(a, func(i, j) {
	return a[i].age < a[j].age
})
out := []
for item in a {
	out = append(out, item.name)
}
out = import("json").encode(out)`, []byte(`["C","D","A","B","E"]`))
}

func TestGetModules(t *testing.T) {
	mods := stdlib.GetModuleMap()
	require.Equal(t, 0, mods.Len())

	mods = stdlib.GetModuleMap("os")
	require.Equal(t, 1, mods.Len())
	require.NotNil(t, mods.Get("os"))

	mods = stdlib.GetModuleMap("os", "rand")
	require.Equal(t, 2, mods.Len())
	require.NotNil(t, mods.Get("os"))
	require.NotNil(t, mods.Get("rand"))

	mods = stdlib.GetModuleMap("text", "text")
	require.Equal(t, 1, mods.Len())
	require.NotNil(t, mods.Get("text"))

	mods = stdlib.GetModuleMap("nonexisting", "text")
	require.Equal(t, 1, mods.Len())
	require.NotNil(t, mods.Get("text"))
}

type callres struct {
	t *testing.T
	o interface{}
	e error
}

func (c callres) call(funcName string, args ...interface{}) callres {
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
			return callres{t: c.t, e: fmt.Errorf(
				"function not found: %s", funcName)}
		}

		f, ok := m.(*tengo.UserFunction)
		if !ok {
			return callres{t: c.t, e: fmt.Errorf(
				"non-callable: %s", funcName)}
		}

		res, err := f.Value(oargs...)
		return callres{t: c.t, o: res, e: err}
	case *tengo.UserFunction:
		res, err := o.Value(oargs...)
		return callres{t: c.t, o: res, e: err}
	case *tengo.ImmutableMap:
		m, ok := o.Value[funcName]
		if !ok {
			return callres{t: c.t, e: fmt.Errorf("function not found: %s", funcName)}
		}

		f, ok := m.(*tengo.UserFunction)
		if !ok {
			return callres{t: c.t, e: fmt.Errorf("non-callable: %s", funcName)}
		}

		res, err := f.Value(oargs...)
		return callres{t: c.t, o: res, e: err}
	default:
		panic(fmt.Errorf("unexpected object: %v (%T)", o, o))
	}
}

func (c callres) expect(expected interface{}, msgAndArgs ...interface{}) {
	require.NoError(c.t, c.e, msgAndArgs...)
	require.Equal(c.t, object(expected), c.o, msgAndArgs...)
}

func (c callres) expectError() {
	require.Error(c.t, c.e)
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
	s := tengo.NewScript([]byte(input))
	s.SetImports(stdlib.GetModuleMap(stdlib.AllModuleNames()...))
	c, err := s.Run()
	require.NoError(t, err)
	require.NotNil(t, c)
	v := c.Get("out")
	require.NotNil(t, v)
	require.Equal(t, expected, v.Value())
}

func expectErr(t *testing.T, input string, errMsg string) {
	s := tengo.NewScript([]byte(input))
	s.SetImports(stdlib.GetModuleMap(stdlib.AllModuleNames()...))
	_, err := s.Run()
	require.True(t, strings.Contains(err.Error(), errMsg))
}

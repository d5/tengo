package runtime_test

import (
	"testing"

	"github.com/d5/tengo"
	"github.com/d5/tengo/compiler/token"
)

func TestBuiltinLen(t *testing.T) {
	expect(t, `out = len("")`, nil, 0)
	expect(t, `out = len("four")`, nil, 4)
	expect(t, `out = len("hello world")`, nil, 11)
	expect(t, `out = len([])`, nil, 0)
	expect(t, `out = len([1, 2, 3])`, nil, 3)
	expect(t, `out = len({})`, nil, 0)
	expect(t, `out = len({a:1, b:2})`, nil, 2)
	expect(t, `out = len(immutable([]))`, nil, 0)
	expect(t, `out = len(immutable([1, 2, 3]))`, nil, 3)
	expect(t, `out = len(immutable({}))`, nil, 0)
	expect(t, `out = len(immutable({a:1, b:2}))`, nil, 2)
	expectError(t, `len(1)`, nil, "invalid type for argument")
	expectError(t, `len("one", "two")`, nil, "wrong number of arguments")
}

func TestBuiltinCopy(t *testing.T) {
	expect(t, `out = copy(1)`, nil, 1)
	expectError(t, `copy(1, 2)`, nil, "wrong number of arguments")
}

func TestBuiltinAppend(t *testing.T) {
	expect(t, `out = append([1, 2, 3], 4)`, nil, ARR{1, 2, 3, 4})
	expect(t, `out = append([1, 2, 3], 4, 5, 6)`, nil, ARR{1, 2, 3, 4, 5, 6})
	expect(t, `out = append([1, 2, 3], "foo", false)`, nil, ARR{1, 2, 3, "foo", false})
}

func TestBuiltinInt(t *testing.T) {
	expect(t, `out = int(1)`, nil, 1)
	expect(t, `out = int(1.8)`, nil, 1)
	expect(t, `out = int("-522")`, nil, -522)
	expect(t, `out = int(true)`, nil, 1)
	expect(t, `out = int(false)`, nil, 0)
	expect(t, `out = int('8')`, nil, 56)
	expect(t, `out = int([1])`, nil, tengo.UndefinedValue)
	expect(t, `out = int({a: 1})`, nil, tengo.UndefinedValue)
	expect(t, `out = int(undefined)`, nil, tengo.UndefinedValue)
	expect(t, `out = int("-522", 1)`, nil, -522)
	expect(t, `out = int(undefined, 1)`, nil, 1)
	expect(t, `out = int(undefined, 1.8)`, nil, 1.8)
	expect(t, `out = int(undefined, string(1))`, nil, "1")
	expect(t, `out = int(undefined, undefined)`, nil, tengo.UndefinedValue)
}

func TestBuiltinString(t *testing.T) {
	expect(t, `out = string(1)`, nil, "1")
	expect(t, `out = string(1.8)`, nil, "1.8")
	expect(t, `out = string("-522")`, nil, "-522")
	expect(t, `out = string(true)`, nil, "true")
	expect(t, `out = string(false)`, nil, "false")
	expect(t, `out = string('8')`, nil, "8")
	expect(t, `out = string([1,8.1,true,3])`, nil, "[1, 8.1, true, 3]")
	expect(t, `out = string({b: "foo"})`, nil, `{b: "foo"}`)
	expect(t, `out = string(undefined)`, nil, tengo.UndefinedValue) // not "undefined"
	expect(t, `out = string(1, "-522")`, nil, "1")
	expect(t, `out = string(undefined, "-522")`, nil, "-522") // not "undefined"
}

func TestBuiltinFloat(t *testing.T) {
	expect(t, `out = float(1)`, nil, 1.0)
	expect(t, `out = float(1.8)`, nil, 1.8)
	expect(t, `out = float("-52.2")`, nil, -52.2)
	expect(t, `out = float(true)`, nil, tengo.UndefinedValue)
	expect(t, `out = float(false)`, nil, tengo.UndefinedValue)
	expect(t, `out = float('8')`, nil, tengo.UndefinedValue)
	expect(t, `out = float([1,8.1,true,3])`, nil, tengo.UndefinedValue)
	expect(t, `out = float({a: 1, b: "foo"})`, nil, tengo.UndefinedValue)
	expect(t, `out = float(undefined)`, nil, tengo.UndefinedValue)
	expect(t, `out = float("-52.2", 1.8)`, nil, -52.2)
	expect(t, `out = float(undefined, 1)`, nil, 1)
	expect(t, `out = float(undefined, 1.8)`, nil, 1.8)
	expect(t, `out = float(undefined, "-52.2")`, nil, "-52.2")
	expect(t, `out = float(undefined, char(56))`, nil, '8')
	expect(t, `out = float(undefined, undefined)`, nil, tengo.UndefinedValue)
}

func TestBuiltinChar(t *testing.T) {
	expect(t, `out = char(56)`, nil, '8')
	expect(t, `out = char(1.8)`, nil, tengo.UndefinedValue)
	expect(t, `out = char("-52.2")`, nil, tengo.UndefinedValue)
	expect(t, `out = char(true)`, nil, tengo.UndefinedValue)
	expect(t, `out = char(false)`, nil, tengo.UndefinedValue)
	expect(t, `out = char('8')`, nil, '8')
	expect(t, `out = char([1,8.1,true,3])`, nil, tengo.UndefinedValue)
	expect(t, `out = char({a: 1, b: "foo"})`, nil, tengo.UndefinedValue)
	expect(t, `out = char(undefined)`, nil, tengo.UndefinedValue)
	expect(t, `out = char(56, 'a')`, nil, '8')
	expect(t, `out = char(undefined, '8')`, nil, '8')
	expect(t, `out = char(undefined, 56)`, nil, 56)
	expect(t, `out = char(undefined, "-52.2")`, nil, "-52.2")
	expect(t, `out = char(undefined, undefined)`, nil, tengo.UndefinedValue)
}

func TestBuiltinBool(t *testing.T) {
	expect(t, `out = bool(1)`, nil, true)          // non-zero integer: true
	expect(t, `out = bool(0)`, nil, false)         // zero: true
	expect(t, `out = bool(1.8)`, nil, true)        // all floats (except for NaN): true
	expect(t, `out = bool(0.0)`, nil, true)        // all floats (except for NaN): true
	expect(t, `out = bool("false")`, nil, true)    // non-empty string: true
	expect(t, `out = bool("")`, nil, false)        // empty string: false
	expect(t, `out = bool(true)`, nil, true)       // true: true
	expect(t, `out = bool(false)`, nil, false)     // false: false
	expect(t, `out = bool('8')`, nil, true)        // non-zero chars: true
	expect(t, `out = bool(char(0))`, nil, false)   // zero char: false
	expect(t, `out = bool([1])`, nil, true)        // non-empty arrays: true
	expect(t, `out = bool([])`, nil, false)        // empty array: false
	expect(t, `out = bool({a: 1})`, nil, true)     // non-empty maps: true
	expect(t, `out = bool({})`, nil, false)        // empty maps: false
	expect(t, `out = bool(undefined)`, nil, false) // undefined: false
}

func TestBuiltinBytes(t *testing.T) {
	expect(t, `out = bytes(1)`, nil, []byte{0})
	expect(t, `out = bytes(1.8)`, nil, tengo.UndefinedValue)
	expect(t, `out = bytes("-522")`, nil, []byte{'-', '5', '2', '2'})
	expect(t, `out = bytes(true)`, nil, tengo.UndefinedValue)
	expect(t, `out = bytes(false)`, nil, tengo.UndefinedValue)
	expect(t, `out = bytes('8')`, nil, tengo.UndefinedValue)
	expect(t, `out = bytes([1])`, nil, tengo.UndefinedValue)
	expect(t, `out = bytes({a: 1})`, nil, tengo.UndefinedValue)
	expect(t, `out = bytes(undefined)`, nil, tengo.UndefinedValue)
	expect(t, `out = bytes("-522", ['8'])`, nil, []byte{'-', '5', '2', '2'})
	expect(t, `out = bytes(undefined, "-522")`, nil, "-522")
	expect(t, `out = bytes(undefined, 1)`, nil, 1)
	expect(t, `out = bytes(undefined, 1.8)`, nil, 1.8)
	expect(t, `out = bytes(undefined, int("-522"))`, nil, -522)
	expect(t, `out = bytes(undefined, undefined)`, nil, tengo.UndefinedValue)
}

func TestBuiltinBytesN(t *testing.T) {
	curMaxBytesLen := tengo.MaxBytesLen
	defer func() { tengo.MaxBytesLen = curMaxBytesLen }()
	tengo.MaxBytesLen = 10

	expect(t, `out = bytes(0)`, nil, make([]byte, 0))
	expect(t, `out = bytes(10)`, nil, make([]byte, 10))
	expectError(t, `bytes(11)`, nil, "bytes size limit")

	tengo.MaxBytesLen = 1000
	expect(t, `out = bytes(1000)`, nil, make([]byte, 1000))
	expectError(t, `bytes(1001)`, nil, "bytes size limit")
}

func TestBuiltinIsError(t *testing.T) {
	expect(t, `out = is_error(error(1))`, nil, true)
	expect(t, `out = is_error(1)`, nil, false)
}

func TestBuiltinIsUndefined(t *testing.T) {
	expect(t, `out = is_undefined(undefined)`, nil, true)
	expect(t, `out = is_undefined(error(1))`, nil, false)
}

func TestBuiltinTypeName(t *testing.T) {
	expect(t, `out = type_name(1)`, nil, "int")
	expect(t, `out = type_name(1.1)`, nil, "float")
	expect(t, `out = type_name("a")`, nil, "string")
	expect(t, `out = type_name([1,2,3])`, nil, "array")
	expect(t, `out = type_name({k:1})`, nil, "map")
	expect(t, `out = type_name('a')`, nil, "char")
	expect(t, `out = type_name(true)`, nil, "bool")
	expect(t, `out = type_name(false)`, nil, "bool")
	expect(t, `out = type_name(bytes( 1))`, nil, "bytes")
	expect(t, `out = type_name(undefined)`, nil, "undefined")
	expect(t, `out = type_name(error("err"))`, nil, "error")
	expect(t, `out = type_name(func() {})`, nil, "function")
	expect(t, `a := func(x) { return func() { return x } }; out = type_name(a(5))`, nil, "function")
}

func TestBuiltinIsFunction(t *testing.T) {
	expect(t, `out = is_function(1)`, nil, false)
	expect(t, `out = is_function(func() {})`, nil, true)
	expect(t, `out = is_function(func(x) { return x })`, nil, true)
	expect(t, `out = is_function(len)`, nil, false)                                                                         // builtin function
	expect(t, `a := func(x) { return func() { return x } }; out = is_function(a)`, nil, true)                               // function
	expect(t, `a := func(x) { return func() { return x } }; out = is_function(a(5))`, nil, true)                            // closure
	expect(t, `out = is_function(x)`, Opts().Symbol("x", &StringArray{Value: []string{"foo", "bar"}}).Skip2ndPass(), false) // user object
}

func TestBuiltinIsCallable(t *testing.T) {
	expect(t, `out = is_callable(1)`, nil, false)
	expect(t, `out = is_callable(func() {})`, nil, true)
	expect(t, `out = is_callable(func(x) { return x })`, nil, true)
	expect(t, `out = is_callable(len)`, nil, true)                                                                         // builtin function
	expect(t, `a := func(x) { return func() { return x } }; out = is_callable(a)`, nil, true)                              // function
	expect(t, `a := func(x) { return func() { return x } }; out = is_callable(a(5))`, nil, true)                           // closure
	expect(t, `out = is_callable(x)`, Opts().Symbol("x", &StringArray{Value: []string{"foo", "bar"}}).Skip2ndPass(), true) // user object
}

func TestBuiltinFormat(t *testing.T) {
	expect(t, `out = format("")`, nil, "")
	expect(t, `out = format("foo")`, nil, "foo")
	expect(t, `out = format("foo %d %v %s", 1, 2, "bar")`, nil, "foo 1 2 bar")
	expect(t, `out = format("foo %v", [1, "bar", true])`, nil, `foo [1, "bar", true]`)
	expect(t, `out = format("foo %v %d", [1, "bar", true], 19)`, nil, `foo [1, "bar", true] 19`)
	expect(t, `out = format("foo %v", {"a": {"b": {"c": [1, 2, 3]}}})`, nil, `foo {a: {b: {c: [1, 2, 3]}}}`)
	expect(t, `out = format("%v", [1, [2, [3, 4]]])`, nil, `[1, [2, [3, 4]]]`)

	tengo.MaxStringLen = 9
	expectError(t, `format("%s", "1234567890")`, nil, "exceeding string size limit")
	tengo.MaxStringLen = 2147483647
}

func TestBuiltinBind(t *testing.T) {
	expectError(t, `bind()`, nil, "wrong number of arguments")

	// binding non-callable: bind() will not fail, but, calling the bound
	// function will fail.
	expect(t, `out = type_name(bind(1))`, nil, "go-function")
	expectError(t, `f := bind(1); f()`, nil, "not callable")
	expect(t, `out = type_name(bind(1, 2))`, nil, "go-function")
	expectError(t, `f := bind(1); f()`, nil, "not callable")

	// binding builtin functions
	expect(t, `f := bind(format, "%v"); out = f(123)`, nil, "123")
	expect(t, `f := bind(format, "%v %v"); out = f("foo", false)`, nil, `"foo" false`)
	expect(t, `f := bind(format, "%v %v", "bar"); out = f(999)`, nil, `"bar" 999`)
	expect(t, `f := bind(format, "%v %v", "bar"); out = f()`, nil, `"bar" %!v(MISSING)`)
	expect(t, `f := bind(format, "%v %v", "bar"); out = f(123, 345)`, nil, `"bar" 123%!(EXTRA int=345)`)
	expect(t, `f := bind(format, "%v %v %v"); out = f([1, 2, 3]...)`, nil, "1 2 3")
	expect(t, `f := bind(format, "%v %v %v", 1); out = f([2, 3]...)`, nil, "1 2 3")
	expect(t, `f := bind(format, "%v %v %v", 1, 2); out = f([3]...)`, nil, "1 2 3")
	expect(t, `f := bind(format, "%v %v %v", 1, 2, 3); out = f([]...)`, nil, "1 2 3")
	expect(t, `f := bind(format, "%v %v %v", [1, 2]...); out = f(3)`, nil, "1 2 3")
	expect(t, `f := bind(format, "%v %v %v", [1, 2, 3]...); out = f()`, nil, "1 2 3")

	// binding compiled functions
	expect(t, `
add := func(a, b) { return a + b }
add4 := bind(add, 4)
out = add4(6)
`, nil, 10)
	expect(t, `
add := func(a, b) { return a + b }
add4_6 := bind(add, 4, 6)
out = add4_6()
`, nil, 10)
	expectError(t, `
add := func(a, b) { return a + b }
add4 := bind(add, 4)
add4()
`, nil, "wrong number of arguments")
	expectError(t, `
add := func(a, b) { return a + b }
add4 := bind(add, 4)
add4(5, 6)
`, nil, "wrong number of arguments")
	expect(t, `
add := func(...a) {
	s := 0
	for _, n in a { s += n }
	return s
}
add4 := bind(add, 4)
out = [add4(), add4(6), add4(5, 6), add4([1, 2, 3]...)]
`, nil, ARR{4, 10, 15, 10})
	expect(t, `
add_fn := func(a, b) { return func() { return a + b } }
add4 := bind(add_fn, 4)
out = add4(6)()
`, nil, 10)
	expect(t, `
add_fn := func(a, b) { return func(c) { return a + b + c } }
add4 := bind(add_fn, 4)
out = add4(6)(5)
`, nil, 15)

	// binding Go functions
	opts := Opts().Skip2ndPass().
		Symbol("add", &tengo.GoFunction{
			Value: func(_ tengo.Interop, args ...tengo.Object) (tengo.Object, error) {
				if len(args) == 0 {
					return tengo.UndefinedValue, nil
				} else if len(args) == 1 {
					return args[0], nil
				}

				var err error
				sum := args[0]
				for _, arg := range args[1:] {
					sum, err = sum.BinaryOp(token.Add, arg)
					if err != nil {
						return nil, err
					}
				}
				return sum, nil
			},
		})
	expect(t, `
add4 := bind(add, 4)
out = add4()
`, opts, 4)
	expect(t, `
add4 := bind(add, 4)
out = add4(6)
`, opts, 10)
	expect(t, `
add4 := bind(add, 4)
out = add4(5, 6)
`, opts, 15)
	expect(t, `
add4 := bind(add, 4)
out = add4([1, 2, 3]...)
`, opts, 10)

	// nested
	expect(t, `
add := func(...a) {
	s := 0
	for _, n in a { s += n }
	return s
}
f := bind(bind(bind(add, 1), 2), 3)
out = f(4)
`, nil, 10)
}

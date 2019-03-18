package runtime_test

import (
	"testing"

	"github.com/d5/tengo"
	"github.com/d5/tengo/objects"
)

func TestBuiltinFunction(t *testing.T) {
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

	expect(t, `out = copy(1)`, nil, 1)
	expectError(t, `copy(1, 2)`, nil, "wrong number of arguments")

	expect(t, `out = append([1, 2, 3], 4)`, nil, ARR{1, 2, 3, 4})
	expect(t, `out = append([1, 2, 3], 4, 5, 6)`, nil, ARR{1, 2, 3, 4, 5, 6})
	expect(t, `out = append([1, 2, 3], "foo", false)`, nil, ARR{1, 2, 3, "foo", false})

	expect(t, `out = int(1)`, nil, 1)
	expect(t, `out = int(1.8)`, nil, 1)
	expect(t, `out = int("-522")`, nil, -522)
	expect(t, `out = int(true)`, nil, 1)
	expect(t, `out = int(false)`, nil, 0)
	expect(t, `out = int('8')`, nil, 56)
	expect(t, `out = int([1])`, nil, objects.UndefinedValue)
	expect(t, `out = int({a: 1})`, nil, objects.UndefinedValue)
	expect(t, `out = int(undefined)`, nil, objects.UndefinedValue)
	expect(t, `out = int("-522", 1)`, nil, -522)
	expect(t, `out = int(undefined, 1)`, nil, 1)
	expect(t, `out = int(undefined, 1.8)`, nil, 1.8)
	expect(t, `out = int(undefined, string(1))`, nil, "1")
	expect(t, `out = int(undefined, undefined)`, nil, objects.UndefinedValue)

	expect(t, `out = string(1)`, nil, "1")
	expect(t, `out = string(1.8)`, nil, "1.8")
	expect(t, `out = string("-522")`, nil, "-522")
	expect(t, `out = string(true)`, nil, "true")
	expect(t, `out = string(false)`, nil, "false")
	expect(t, `out = string('8')`, nil, "8")
	expect(t, `out = string([1,8.1,true,3])`, nil, "[1, 8.1, true, 3]")
	expect(t, `out = string({b: "foo"})`, nil, `{b: "foo"}`)
	expect(t, `out = string(undefined)`, nil, objects.UndefinedValue) // not "undefined"
	expect(t, `out = string(1, "-522")`, nil, "1")
	expect(t, `out = string(undefined, "-522")`, nil, "-522") // not "undefined"

	expect(t, `out = float(1)`, nil, 1.0)
	expect(t, `out = float(1.8)`, nil, 1.8)
	expect(t, `out = float("-52.2")`, nil, -52.2)
	expect(t, `out = float(true)`, nil, objects.UndefinedValue)
	expect(t, `out = float(false)`, nil, objects.UndefinedValue)
	expect(t, `out = float('8')`, nil, objects.UndefinedValue)
	expect(t, `out = float([1,8.1,true,3])`, nil, objects.UndefinedValue)
	expect(t, `out = float({a: 1, b: "foo"})`, nil, objects.UndefinedValue)
	expect(t, `out = float(undefined)`, nil, objects.UndefinedValue)
	expect(t, `out = float("-52.2", 1.8)`, nil, -52.2)
	expect(t, `out = float(undefined, 1)`, nil, 1)
	expect(t, `out = float(undefined, 1.8)`, nil, 1.8)
	expect(t, `out = float(undefined, "-52.2")`, nil, "-52.2")
	expect(t, `out = float(undefined, char(56))`, nil, '8')
	expect(t, `out = float(undefined, undefined)`, nil, objects.UndefinedValue)

	expect(t, `out = char(56)`, nil, '8')
	expect(t, `out = char(1.8)`, nil, objects.UndefinedValue)
	expect(t, `out = char("-52.2")`, nil, objects.UndefinedValue)
	expect(t, `out = char(true)`, nil, objects.UndefinedValue)
	expect(t, `out = char(false)`, nil, objects.UndefinedValue)
	expect(t, `out = char('8')`, nil, '8')
	expect(t, `out = char([1,8.1,true,3])`, nil, objects.UndefinedValue)
	expect(t, `out = char({a: 1, b: "foo"})`, nil, objects.UndefinedValue)
	expect(t, `out = char(undefined)`, nil, objects.UndefinedValue)
	expect(t, `out = char(56, 'a')`, nil, '8')
	expect(t, `out = char(undefined, '8')`, nil, '8')
	expect(t, `out = char(undefined, 56)`, nil, 56)
	expect(t, `out = char(undefined, "-52.2")`, nil, "-52.2")
	expect(t, `out = char(undefined, undefined)`, nil, objects.UndefinedValue)

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

	expect(t, `out = bytes(1)`, nil, []byte{0})
	expect(t, `out = bytes(1.8)`, nil, objects.UndefinedValue)
	expect(t, `out = bytes("-522")`, nil, []byte{'-', '5', '2', '2'})
	expect(t, `out = bytes(true)`, nil, objects.UndefinedValue)
	expect(t, `out = bytes(false)`, nil, objects.UndefinedValue)
	expect(t, `out = bytes('8')`, nil, objects.UndefinedValue)
	expect(t, `out = bytes([1])`, nil, objects.UndefinedValue)
	expect(t, `out = bytes({a: 1})`, nil, objects.UndefinedValue)
	expect(t, `out = bytes(undefined)`, nil, objects.UndefinedValue)
	expect(t, `out = bytes("-522", ['8'])`, nil, []byte{'-', '5', '2', '2'})
	expect(t, `out = bytes(undefined, "-522")`, nil, "-522")
	expect(t, `out = bytes(undefined, 1)`, nil, 1)
	expect(t, `out = bytes(undefined, 1.8)`, nil, 1.8)
	expect(t, `out = bytes(undefined, int("-522"))`, nil, -522)
	expect(t, `out = bytes(undefined, undefined)`, nil, objects.UndefinedValue)

	expect(t, `out = is_error(error(1))`, nil, true)
	expect(t, `out = is_error(1)`, nil, false)

	expect(t, `out = is_undefined(undefined)`, nil, true)
	expect(t, `out = is_undefined(error(1))`, nil, false)

	// type_name
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
	expect(t, `out = type_name(func() {})`, nil, "compiled-function")
	expect(t, `a := func(x) { return func() { return x } }; out = type_name(a(5))`, nil, "closure") // closure

	// is_function
	expect(t, `out = is_function(1)`, nil, false)
	expect(t, `out = is_function(func() {})`, nil, true)
	expect(t, `out = is_function(func(x) { return x })`, nil, true)
	expect(t, `out = is_function(len)`, nil, false)                                                                         // builtin function
	expect(t, `a := func(x) { return func() { return x } }; out = is_function(a)`, nil, true)                               // function
	expect(t, `a := func(x) { return func() { return x } }; out = is_function(a(5))`, nil, true)                            // closure
	expect(t, `out = is_function(x)`, Opts().Symbol("x", &StringArray{Value: []string{"foo", "bar"}}).Skip2ndPass(), false) // user object

	// is_callable
	expect(t, `out = is_callable(1)`, nil, false)
	expect(t, `out = is_callable(func() {})`, nil, true)
	expect(t, `out = is_callable(func(x) { return x })`, nil, true)
	expect(t, `out = is_callable(len)`, nil, true)                                                                         // builtin function
	expect(t, `a := func(x) { return func() { return x } }; out = is_callable(a)`, nil, true)                              // function
	expect(t, `a := func(x) { return func() { return x } }; out = is_callable(a(5))`, nil, true)                           // closure
	expect(t, `out = is_callable(x)`, Opts().Symbol("x", &StringArray{Value: []string{"foo", "bar"}}).Skip2ndPass(), true) // user object
}

func TestBytesN(t *testing.T) {
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

package runtime_test

import (
	"testing"

	"github.com/d5/tengo/objects"
)

func TestBuiltinFunction(t *testing.T) {
	expect(t, `out = len("")`, 0)
	expect(t, `out = len("four")`, 4)
	expect(t, `out = len("hello world")`, 11)
	expect(t, `out = len([])`, 0)
	expect(t, `out = len([1, 2, 3])`, 3)
	expect(t, `out = len({})`, 0)
	expect(t, `out = len({a:1, b:2})`, 2)
	expect(t, `out = len(immutable([]))`, 0)
	expect(t, `out = len(immutable([1, 2, 3]))`, 3)
	expect(t, `out = len(immutable({}))`, 0)
	expect(t, `out = len(immutable({a:1, b:2}))`, 2)
	expectError(t, `len(1)`)
	expectError(t, `len("one", "two")`)

	expect(t, `out = copy(1)`, 1)
	expectError(t, `out = copy(1, 2)`)

	expect(t, `out = append([1, 2, 3], 4)`, ARR{1, 2, 3, 4})
	expect(t, `out = append([1, 2, 3], 4, 5, 6)`, ARR{1, 2, 3, 4, 5, 6})
	expect(t, `out = append([1, 2, 3], "foo", false)`, ARR{1, 2, 3, "foo", false})

	expect(t, `out = int(1)`, 1)
	expect(t, `out = int(1.8)`, 1)
	expect(t, `out = int("-522")`, -522)
	expect(t, `out = int(true)`, 1)
	expect(t, `out = int(false)`, 0)
	expect(t, `out = int('8')`, 56)
	expect(t, `out = int([1])`, objects.UndefinedValue)
	expect(t, `out = int({a: 1})`, objects.UndefinedValue)
	expect(t, `out = int(undefined)`, objects.UndefinedValue)
	expect(t, `out = int("-522", 1)`, -522)
	expect(t, `out = int(undefined, 1)`, 1)
	expect(t, `out = int(undefined, 1.8)`, 1.8)
	expect(t, `out = int(undefined, string(1))`, "1")
	expect(t, `out = int(undefined, undefined)`, objects.UndefinedValue)

	expect(t, `out = string(1)`, "1")
	expect(t, `out = string(1.8)`, "1.8")
	expect(t, `out = string("-522")`, "-522")
	expect(t, `out = string(true)`, "true")
	expect(t, `out = string(false)`, "false")
	expect(t, `out = string('8')`, "8")
	expect(t, `out = string([1,8.1,true,3])`, "[1, 8.1, true, 3]")
	expect(t, `out = string({b: "foo"})`, `{b: "foo"}`)
	expect(t, `out = string(undefined)`, objects.UndefinedValue) // not "undefined"
	expect(t, `out = string(1, "-522")`, "1")
	expect(t, `out = string(undefined, "-522")`, "-522") // not "undefined"

	expect(t, `out = float(1)`, 1.0)
	expect(t, `out = float(1.8)`, 1.8)
	expect(t, `out = float("-52.2")`, -52.2)
	expect(t, `out = float(true)`, objects.UndefinedValue)
	expect(t, `out = float(false)`, objects.UndefinedValue)
	expect(t, `out = float('8')`, objects.UndefinedValue)
	expect(t, `out = float([1,8.1,true,3])`, objects.UndefinedValue)
	expect(t, `out = float({a: 1, b: "foo"})`, objects.UndefinedValue)
	expect(t, `out = float(undefined)`, objects.UndefinedValue)
	expect(t, `out = float("-52.2", 1.8)`, -52.2)
	expect(t, `out = float(undefined, 1)`, 1)
	expect(t, `out = float(undefined, 1.8)`, 1.8)
	expect(t, `out = float(undefined, "-52.2")`, "-52.2")
	expect(t, `out = float(undefined, char(56))`, '8')
	expect(t, `out = float(undefined, undefined)`, objects.UndefinedValue)

	expect(t, `out = char(56)`, '8')
	expect(t, `out = char(1.8)`, objects.UndefinedValue)
	expect(t, `out = char("-52.2")`, objects.UndefinedValue)
	expect(t, `out = char(true)`, objects.UndefinedValue)
	expect(t, `out = char(false)`, objects.UndefinedValue)
	expect(t, `out = char('8')`, '8')
	expect(t, `out = char([1,8.1,true,3])`, objects.UndefinedValue)
	expect(t, `out = char({a: 1, b: "foo"})`, objects.UndefinedValue)
	expect(t, `out = char(undefined)`, objects.UndefinedValue)
	expect(t, `out = char(56, 'a')`, '8')
	expect(t, `out = char(undefined, '8')`, '8')
	expect(t, `out = char(undefined, 56)`, 56)
	expect(t, `out = char(undefined, "-52.2")`, "-52.2")
	expect(t, `out = char(undefined, undefined)`, objects.UndefinedValue)

	expect(t, `out = bool(1)`, true)          // non-zero integer: true
	expect(t, `out = bool(0)`, false)         // zero: true
	expect(t, `out = bool(1.8)`, true)        // all floats (except for NaN): true
	expect(t, `out = bool(0.0)`, true)        // all floats (except for NaN): true
	expect(t, `out = bool("false")`, true)    // non-empty string: true
	expect(t, `out = bool("")`, false)        // empty string: false
	expect(t, `out = bool(true)`, true)       // true: true
	expect(t, `out = bool(false)`, false)     // false: false
	expect(t, `out = bool('8')`, true)        // non-zero chars: true
	expect(t, `out = bool(char(0))`, false)   // zero char: false
	expect(t, `out = bool([1])`, true)        // non-empty arrays: true
	expect(t, `out = bool([])`, false)        // empty array: false
	expect(t, `out = bool({a: 1})`, true)     // non-empty maps: true
	expect(t, `out = bool({})`, false)        // empty maps: false
	expect(t, `out = bool(undefined)`, false) // undefined: false

	expect(t, `out = bytes(1)`, []byte{0})
	expect(t, `out = bytes(1.8)`, objects.UndefinedValue)
	expect(t, `out = bytes("-522")`, []byte{'-', '5', '2', '2'})
	expect(t, `out = bytes(true)`, objects.UndefinedValue)
	expect(t, `out = bytes(false)`, objects.UndefinedValue)
	expect(t, `out = bytes('8')`, objects.UndefinedValue)
	expect(t, `out = bytes([1])`, objects.UndefinedValue)
	expect(t, `out = bytes({a: 1})`, objects.UndefinedValue)
	expect(t, `out = bytes(undefined)`, objects.UndefinedValue)
	expect(t, `out = bytes("-522", ['8'])`, []byte{'-', '5', '2', '2'})
	expect(t, `out = bytes(undefined, "-522")`, "-522")
	expect(t, `out = bytes(undefined, 1)`, 1)
	expect(t, `out = bytes(undefined, 1.8)`, 1.8)
	expect(t, `out = bytes(undefined, int("-522"))`, -522)
	expect(t, `out = bytes(undefined, undefined)`, objects.UndefinedValue)

	expect(t, `out = is_error(error(1))`, true)
	expect(t, `out = is_error(1)`, false)

	expect(t, `out = is_undefined(undefined)`, true)
	expect(t, `out = is_undefined(error(1))`, false)

	// to_json
	expect(t, `out = to_json(5)`, []byte("5"))
	expect(t, `out = to_json({foo: 5})`, []byte("{\"foo\":5}"))
	expect(t, `out = to_json({foo: "bar"})`, []byte("{\"foo\":\"bar\"}"))
	expect(t, `out = to_json({foo: 1.8})`, []byte("{\"foo\":1.8}"))
	expect(t, `out = to_json({foo: true})`, []byte("{\"foo\":true}"))
	expect(t, `out = to_json({foo: '8'})`, []byte("{\"foo\":56}"))
	expect(t, `out = to_json({foo: bytes("foo")})`, []byte("{\"foo\":\"Zm9v\"}")) // json encoding returns []byte as base64 encoded string
	expect(t, `out = to_json({foo: ["bar", 1, 1.8, '8', true]})`, []byte("{\"foo\":[\"bar\",1,1.8,56,true]}"))
	expect(t, `out = to_json({foo: [["bar", 1], ["bar", 1]]})`, []byte("{\"foo\":[[\"bar\",1],[\"bar\",1]]}"))
	expect(t, `out = to_json({foo: {string: "bar", int: 1, float: 1.8, char: '8', bool: true}})`, []byte("{\"foo\":{\"bool\":true,\"char\":56,\"float\":1.8,\"int\":1,\"string\":\"bar\"}}"))
	expect(t, `out = to_json({foo: {map1: {string: "bar"}, map2: {int: "1"}}})`, []byte("{\"foo\":{\"map1\":{\"string\":\"bar\"},\"map2\":{\"int\":\"1\"}}}"))
	expect(t, `out = to_json([["bar", 1], ["bar", 1]])`, []byte("[[\"bar\",1],[\"bar\",1]]"))

	// from_json
	expect(t, `out = from_json("{\"foo\":5}").foo`, 5.0)
	expect(t, `out = from_json("{\"foo\":\"bar\"}").foo`, "bar")
	expect(t, `out = from_json("{\"foo\":1.8}").foo`, 1.8)
	expect(t, `out = from_json("{\"foo\":true}").foo`, true)
	expect(t, `out = from_json("{\"foo\":[\"bar\",1,1.8,56,true]}").foo`, ARR{"bar", 1.0, 1.8, 56.0, true})
	expect(t, `out = from_json("{\"foo\":[[\"bar\",1],[\"bar\",1]]}").foo[0]`, ARR{"bar", 1.0})
	expect(t, `out = from_json("{\"foo\":{\"bool\":true,\"char\":56,\"float\":1.8,\"int\":1,\"string\":\"bar\"}}").foo.bool`, true)
	expect(t, `out = from_json("{\"foo\":{\"map1\":{\"string\":\"bar\"},\"map2\":{\"int\":\"1\"}}}").foo.map1.string`, "bar")

	expect(t, `out = from_json("5")`, 5.0)
	expect(t, `out = from_json("[\"bar\",1,1.8,56,true]")`, ARR{"bar", 1.0, 1.8, 56.0, true})

	// sprintf
	expect(t, `out = sprintf("")`, "")
	expect(t, `out = sprintf("foo")`, "foo")
	expect(t, `out = sprintf("foo %d %v %s", 1, 2, "bar")`, "foo 1 2 bar")
	expect(t, `out = sprintf("foo %v", [1, "bar", true])`, "foo [1 bar true]")
	expect(t, `out = sprintf("foo %v %d", [1, "bar", true], 19)`, "foo [1 bar true] 19")
	expectError(t, `sprintf(1)`)   // format has to be String
	expectError(t, `sprintf('c')`) // format has to be String

	// type_name
	expect(t, `out = type_name(1)`, "int")
	expect(t, `out = type_name(1.1)`, "float")
	expect(t, `out = type_name("a")`, "string")
	expect(t, `out = type_name([1,2,3])`, "array")
	expect(t, `out = type_name({k:1})`, "map")
	expect(t, `out = type_name('a')`, "char")
	expect(t, `out = type_name(true)`, "bool")
	expect(t, `out = type_name(false)`, "bool")
	expect(t, `out = type_name(bytes( 1))`, "bytes")
	expect(t, `out = type_name(undefined)`, "undefined")
	expect(t, `out = type_name(error("err"))`, "error")
}

package runtime_test

import (
	"testing"
)

func TestBuiltinFunction(t *testing.T) {
	expect(t, `out = len("")`, 0)
	expect(t, `out = len("four")`, 4)
	expect(t, `out = len("hello world")`, 11)
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
	expect(t, `out = int([1])`, undefined())
	expect(t, `out = int({a: 1})`, undefined())
	expect(t, `out = int(undefined)`, undefined())
	expect(t, `out = int("-522", 1)`, -522)
	expect(t, `out = int(undefined, 1)`, 1)
	expect(t, `out = int(undefined, 1.8)`, undefined())
	expect(t, `out = int(undefined, undefined)`, undefined())

	expect(t, `out = string(1)`, "1")
	expect(t, `out = string(1.8)`, "1.8")
	expect(t, `out = string("-522")`, "-522")
	expect(t, `out = string(true)`, "true")
	expect(t, `out = string(false)`, "false")
	expect(t, `out = string('8')`, "8")
	expect(t, `out = string([1,8.1,true,3])`, "[1, 8.1, true, 3]")
	expect(t, `out = string({b: "foo"})`, `{b: "foo"}`)
	expect(t, `out = string(undefined)`, undefined()) // not "undefined"
	expect(t, `out = string(1, "-522")`, "1")
	expect(t, `out = string(undefined, "-522")`, "-522") // not "undefined"

	expect(t, `out = float(1)`, 1.0)
	expect(t, `out = float(1.8)`, 1.8)
	expect(t, `out = float("-52.2")`, -52.2)
	expect(t, `out = float(true)`, undefined())
	expect(t, `out = float(false)`, undefined())
	expect(t, `out = float('8')`, undefined())
	expect(t, `out = float([1,8.1,true,3])`, undefined())
	expect(t, `out = float({a: 1, b: "foo"})`, undefined())
	expect(t, `out = float(undefined)`, undefined())
	expect(t, `out = float("-52.2", 1.8)`, -52.2)
	expect(t, `out = float(undefined, 1)`, 1.0)
	expect(t, `out = float(undefined, 1.8)`, 1.8)
	expect(t, `out = float(undefined, "-52.2")`, undefined())
	expect(t, `out = float(undefined, undefined)`, undefined())

	expect(t, `out = char(56)`, '8')
	expect(t, `out = char(1.8)`, undefined())
	expect(t, `out = char("-52.2")`, undefined())
	expect(t, `out = char(true)`, undefined())
	expect(t, `out = char(false)`, undefined())
	expect(t, `out = char('8')`, '8')
	expect(t, `out = char([1,8.1,true,3])`, undefined())
	expect(t, `out = char({a: 1, b: "foo"})`, undefined())
	expect(t, `out = char(undefined)`, undefined())
	expect(t, `out = char(56, 'a')`, '8')
	expect(t, `out = char(undefined, '8')`, '8')
	expect(t, `out = char(undefined, 56)`, undefined())
	expect(t, `out = char(undefined, "-52.2")`, undefined())

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
	expect(t, `out = bytes(1.8)`, undefined())
	expect(t, `out = bytes("-522")`, []byte{'-', '5', '2', '2'})
	expect(t, `out = bytes(true)`, undefined())
	expect(t, `out = bytes(false)`, undefined())
	expect(t, `out = bytes('8')`, undefined())
	expect(t, `out = bytes([1])`, undefined())
	expect(t, `out = bytes({a: 1})`, undefined())
	expect(t, `out = bytes(undefined)`, undefined())
	expect(t, `out = bytes("-522", ['8'])`, []byte{'-', '5', '2', '2'})
	expect(t, `out = bytes(undefined, "-522")`, []byte{'-', '5', '2', '2'})
	expect(t, `out = bytes(undefined, 1)`, []byte{0})
	expect(t, `out = bytes(undefined, 1.8)`, undefined())
	expect(t, `out = bytes(undefined, undefined)`, undefined())

	expect(t, `out = is_error(error(1))`, true)
	expect(t, `out = is_error(1)`, false)

	expect(t, `out = is_undefined(undefined)`, true)
	expect(t, `out = is_undefined(error(1))`, false)
}

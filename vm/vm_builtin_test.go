package vm_test

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
}

package runtime_test

import (
	"testing"

	"github.com/d5/tengo/objects"
)

func TestObjectsLimit(t *testing.T) {
	testAllocsLimit(t, `5`, 0)
	testAllocsLimit(t, `5 + 5`, 1)
	testAllocsLimit(t, `a := [1, 2, 3]`, 1)
	testAllocsLimit(t, `a := 1; b := 2; c := 3; d := [a, b, c]`, 1)
	testAllocsLimit(t, `a := {foo: 1, bar: 2}`, 1)
	testAllocsLimit(t, `a := 1; b := 2; c := {foo: a, bar: b}`, 1)
	testAllocsLimit(t, `
f := func() {
	return 5 + 5
}
a := f() + 5
`, 2)
	testAllocsLimit(t, `
f := func() {
	return 5 + 5
}
a := f()
`, 1)
	testAllocsLimit(t, `
a := []
f := func() {
	a = append(a, 5)
}
f()
f()
f()
`, 4)
}

func testAllocsLimit(t *testing.T, src string, limit int64) {
	expectAllocsLimit(t, src, -1, objects.UndefinedValue) // no limit
	expectAllocsLimit(t, src, limit, objects.UndefinedValue)
	expectAllocsLimit(t, src, limit+1, objects.UndefinedValue)
	if limit > 1 {
		expectErrorAllocsLimit(t, src, limit-1, "allocation limit exceeded")
	}
	if limit > 2 {
		expectErrorAllocsLimit(t, src, limit-2, "allocation limit exceeded")
	}
}

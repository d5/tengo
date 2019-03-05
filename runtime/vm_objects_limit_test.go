package runtime_test

import "testing"

func TestObjectsLimit(t *testing.T) {
	expectAllocsLimit(t, `out = 5`, 0, 5)
	expectAllocsLimit(t, `out = 5 + 5`, 1, 10)
	expectErrorAllocsLimit(t, `5 + 5`, 0, "allocation limit exceeded")

	// compound types
	expectAllocsLimit(t, `out = [1, 2, 3]`, 1, ARR{1, 2, 3})
	expectAllocsLimit(t, `a := 1; b := 2; c := 3; out = [a, b, c]`, 1, ARR{1, 2, 3})
	expectAllocsLimit(t, `out = {foo: 1, bar: 2}`, 1, MAP{"foo": 1, "bar": 2})
	expectAllocsLimit(t, `a := 1; b := 2; out = {foo: a, bar: b}`, 1, MAP{"foo": 1, "bar": 2})

	expectAllocsLimit(t, `
f := func() {
	return 5 + 5
}
out = f() + 5
`, 1, 15)
	expectErrorAllocsLimit(t, `
f := func() {
	return 5 + 5
}
f()
`, 0, "allocation limit exceeded")

	expectAllocsLimit(t, `
a := []
f := func() {
	a = append(a, 5)
}
f()
f()
f()
out = a
`, 1, ARR{5, 5, 5})
}

package runtime_test

import "testing"

func TestObjectsLimit(t *testing.T) {
	expectAllocsLimit(t, `out = 5`, 0, 5)
	expectAllocsLimit(t, `out = 5 + 5`, 1, 10)
	expectErrorAllocsLimit(t, `5 + 5`, 0, "objects limit exceeded")

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
`, 0, "objects limit exceeded")

	expectAllocsLimit(t, `
a := []
f := func() {
	a = append(a, 5)
}
f()
f()
f()
out = a
`, 4, ARR{5, 5, 5})
}

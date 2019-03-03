package runtime_test

import "testing"

func TestObjectsLimit(t *testing.T) {
	expectObjectsLimit(t, `out = 5`, 0, 5)
	expectObjectsLimit(t, `out = 5 + 5`, 1, 10)
	expectErrorObjectsLimit(t, `5 + 5`, 0, "objects limit exceeded")

	expectObjectsLimit(t, `
f := func() {
	return 5 + 5
}
out = f() + 5
`, 1, 15)
	expectErrorObjectsLimit(t, `
f := func() {
	return 5 + 5
}
f()
`, 0, "objects limit exceeded")

	expectObjectsLimit(t, `
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

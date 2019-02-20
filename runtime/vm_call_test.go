package runtime_test

import "testing"

func TestCall(t *testing.T) {
	expect(t, `a := { b: func(x) { return x + 2 } }; out = a.b(5)`, 7)
	expect(t, `a := { b: { c: func(x) { return x + 2 } } }; out = a.b.c(5)`, 7)
	expect(t, `a := { b: { c: func(x) { return x + 2 } } }; out = a["b"].c(5)`, 7)
	expectError(t, `a := 1
b := func(a, c) {
   c(a)
}

c := func(a) {
   a()
}
b(a, c)
`, `test:7:4: not callable: int
in test:3:4
in test:9:1`)
}

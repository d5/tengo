package runtime_test

import "testing"

func TestCall(t *testing.T) {
	expect(t, `a := { b: func(x) { return x + 2 } }; out = a.b(5)`, nil, 7)
	expect(t, `a := { b: { c: func(x) { return x + 2 } } }; out = a.b.c(5)`, nil, 7)
	expect(t, `a := { b: { c: func(x) { return x + 2 } } }; out = a["b"].c(5)`, nil, 7)
	expectError(t, `a := 1
b := func(a, c) {
   c(a)
}

c := func(a) {
   a()
}
b(a, c)
`, nil, "Runtime Error: not callable: int\n\tat test:7:4\n\tat test:3:4\n\tat test:9:1")
}

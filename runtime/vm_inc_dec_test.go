package runtime_test

import (
	"testing"
)

func TestIncDec(t *testing.T) {
	expect(t, `out = 0; out++`, nil, 1)
	expect(t, `out = 0; out--`, nil, -1)
	expect(t, `a := 0; a++; out = a`, nil, 1)
	expect(t, `a := 0; a++; a--; out = a`, nil, 0)

	// this seems strange but it works because 'a += b' is
	// translated into 'a = a + b' and string type takes other types for + operator.
	expect(t, `a := "foo"; a++; out = a`, nil, "foo1")
	expectError(t, `a := "foo"; a--`, nil, "invalid operation")

	expectError(t, `a++`, nil, "unresolved reference") // not declared
	expectError(t, `a--`, nil, "unresolved reference") // not declared
	expectError(t, `4++`, nil, "unresolved reference")
}

package runtime_test

import (
	"testing"
)

func TestIncDec(t *testing.T) {
	expect(t, `out++`, 1)
	expect(t, `out--`, -1)
	expect(t, `a := 0; a++; out = a`, 1)
	expect(t, `a := 0; a++; a--; out = a`, 0)

	// this seems strange but it works because 'a += b' is
	// translated into 'a = a + b' and string type takes other types for + operator.
	expect(t, `a := "foo"; a++; out = a`, "foo1")
	expectError(t, `a := "foo"; a--`)

	expectError(t, `a++`) // not declared
	expectError(t, `a--`) // not declared
	//expectError(t, `a := 0; b := a++`) // inc-dec is statement not expression <- parser error
	expectError(t, `4++`)
}

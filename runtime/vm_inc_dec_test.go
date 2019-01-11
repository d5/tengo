package runtime_test

import (
	"testing"
)

func TestIncDec(t *testing.T) {
	expect(t, `out++`, 1)
	expect(t, `out--`, -1)
	expect(t, `a := 0; a++; out = a`, 1)
	expect(t, `a := 0; a++; a--; out = a`, 0)

	expectError(t, `a++`)             // not declared
	expectError(t, `a--`)             // not declared
	expectError(t, `a := "foo"; a++`) // invalid operand
	//expectError(t, `a := 0; b := a++`) // inc-dec is statement not expression <- parser error

	expectError(t, `4++`)
}

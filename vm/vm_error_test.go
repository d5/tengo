package vm_test

import (
	"testing"
)

func TestError(t *testing.T) {
	expectError(t, `5 + true`)

	expectError(t, `5 + true; 5`)

	expectError(t, `-true`)

	expectError(t, `true + false`)

	expectError(t, `5; true + false; 5`)

	expectError(t, `if (10 > 1) { true + false; }`)

	expectError(t, `
if (10 > 1) {
	if (10 > 1) {
		return true + false;
	}

	return 1;
}
`)

	expectError(t, `if (true + false) { 10 }`)

	expectError(t, `10 + (true + false)`)

	expectError(t, `(true + false) + 20`)

	expectError(t, `!(true + false)`)

	expectError(t, `foobar`)

	expectError(t, `"foo" - "bar"`)
}

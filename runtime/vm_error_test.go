package runtime_test

import (
	"testing"
)

func TestErrorObject(t *testing.T) {
	expect(t, `out = error(1)`, errorObject(1))
	expect(t, `out = error("some error")`, errorObject("some error"))
	expect(t, `out = error("some" + " error")`, errorObject("some error"))
	expect(t, `out = func() { return error(5) }()`, errorObject(5))
}

func TestError(t *testing.T) {
	// TODO: these tests should probably be moved to other more relevant tests.
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

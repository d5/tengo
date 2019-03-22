package runtime_test

import "testing"

func TestVMStackOverflow(t *testing.T) {
	expectError(t, `f := func() { return f() + 1 }; f()`, nil, "stack overflow")
}

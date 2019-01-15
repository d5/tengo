package runtime_test

import "testing"

func TestModule(t *testing.T) {
	expect(t, `math := import("math"); out = math.abs(1)`, 1.0)
	expect(t, `math := import("math"); out = math.abs(-1)`, 1.0)
	expect(t, `math := import("math"); out = math.abs(1.0)`, 1.0)
	expect(t, `math := import("math"); out = math.abs(-1.0)`, 1.0)

	// TODO: test for-in statement with module map

	// TODO: test sharing of global variables

	// TODO: test module function that mutate other variables in the same module
}

package runtime_test

import (
	"testing"
)

func TestBoolean(t *testing.T) {
	expect(t, `out = true`, true)
	expect(t, `out = false`, false)

	expect(t, `out = 1 < 2`, true)
	expect(t, `out = 1 > 2`, false)
	expect(t, `out = 1 < 1`, false)
	expect(t, `out = 1 > 2`, false)
	expect(t, `out = 1 == 1`, true)
	expect(t, `out = 1 != 1`, false)
	expect(t, `out = 1 == 2`, false)
	expect(t, `out = 1 != 2`, true)
	expect(t, `out = 1 <= 2`, true)
	expect(t, `out = 1 >= 2`, false)
	expect(t, `out = 1 <= 1`, true)
	expect(t, `out = 1 >= 2`, false)

	expect(t, `out = true == true`, true)
	expect(t, `out = false == false`, true)
	expect(t, `out = true == false`, false)
	expect(t, `out = true != false`, true)
	expect(t, `out = false != true`, true)
	expect(t, `out = (1 < 2) == true`, true)
	expect(t, `out = (1 < 2) == false`, false)
	expect(t, `out = (1 > 2) == true`, false)
	expect(t, `out = (1 > 2) == false`, true)

	expectError(t, `5 + true`, "invalid operation")
	expectError(t, `5 + true; 5`, "invalid operation")
	expectError(t, `-true`, "invalid operation")
	expectError(t, `true + false`, "invalid operation")
	expectError(t, `5; true + false; 5`, "invalid operation")
	expectError(t, `if (10 > 1) { true + false; }`, "invalid operation")
	expectError(t, `
func() {
	if (10 > 1) {
		if (10 > 1) {
			return true + false;
		}

		return 1;
	}
}()
`, "invalid operation")
	expectError(t, `if (true + false) { 10 }`, "invalid operation")
	expectError(t, `10 + (true + false)`, "invalid operation")
	expectError(t, `(true + false) + 20`, "invalid operation")
	expectError(t, `!(true + false)`, "invalid operation")
}

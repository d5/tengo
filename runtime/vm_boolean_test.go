package runtime_test

import (
	"testing"
)

func TestBoolean(t *testing.T) {
	expect(t, `out = true`, nil, true)
	expect(t, `out = false`, nil, false)

	expect(t, `out = 1 < 2`, nil, true)
	expect(t, `out = 1 > 2`, nil, false)
	expect(t, `out = 1 < 1`, nil, false)
	expect(t, `out = 1 > 2`, nil, false)
	expect(t, `out = 1 == 1`, nil, true)
	expect(t, `out = 1 != 1`, nil, false)
	expect(t, `out = 1 == 2`, nil, false)
	expect(t, `out = 1 != 2`, nil, true)
	expect(t, `out = 1 <= 2`, nil, true)
	expect(t, `out = 1 >= 2`, nil, false)
	expect(t, `out = 1 <= 1`, nil, true)
	expect(t, `out = 1 >= 2`, nil, false)

	expect(t, `out = true == true`, nil, true)
	expect(t, `out = false == false`, nil, true)
	expect(t, `out = true == false`, nil, false)
	expect(t, `out = true != false`, nil, true)
	expect(t, `out = false != true`, nil, true)
	expect(t, `out = (1 < 2) == true`, nil, true)
	expect(t, `out = (1 < 2) == false`, nil, false)
	expect(t, `out = (1 > 2) == true`, nil, false)
	expect(t, `out = (1 > 2) == false`, nil, true)

	expectError(t, `5 + true`, nil, "invalid operation")
	expectError(t, `5 + true; 5`, nil, "invalid operation")
	expectError(t, `-true`, nil, "invalid operation")
	expectError(t, `true + false`, nil, "invalid operation")
	expectError(t, `5; true + false; 5`, nil, "invalid operation")
	expectError(t, `if (10 > 1) { true + false; }`, nil, "invalid operation")
	expectError(t, `
func() {
	if (10 > 1) {
		if (10 > 1) {
			return true + false;
		}

		return 1;
	}
}()
`, nil, "invalid operation")
	expectError(t, `if (true + false) { 10 }`, nil, "invalid operation")
	expectError(t, `10 + (true + false)`, nil, "invalid operation")
	expectError(t, `(true + false) + 20`, nil, "invalid operation")
	expectError(t, `!(true + false)`, nil, "invalid operation")
}

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
}

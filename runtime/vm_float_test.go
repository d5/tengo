package runtime_test

import (
	"testing"
)

func TestFloat(t *testing.T) {
	expect(t, `out = 0.0`, 0.0)
	expect(t, `out = -10.3`, -10.3)
	expect(t, `out = 3.2 + 2.0 * -4.0`, -4.8)
	expect(t, `out = 4 + 2.3`, 6.3)
	expect(t, `out = 2.3 + 4`, 6.3)
	expect(t, `out = +5.0`, 5.0)
	expect(t, `out = -5.0 + +5.0`, 0.0)
}

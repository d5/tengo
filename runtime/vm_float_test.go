package runtime_test

import (
	"testing"
)

func TestFloat(t *testing.T) {
	expect(t, `out = 0.0`, nil, 0.0)
	expect(t, `out = -10.3`, nil, -10.3)
	expect(t, `out = 3.2 + 2.0 * -4.0`, nil, -4.8)
	expect(t, `out = 4 + 2.3`, nil, 6.3)
	expect(t, `out = 2.3 + 4`, nil, 6.3)
	expect(t, `out = +5.0`, nil, 5.0)
	expect(t, `out = -5.0 + +5.0`, nil, 0.0)
}

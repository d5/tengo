package runtime_test

import "testing"

func TestBitwise(t *testing.T) {
	expect(t, `out = 1 & 1`, nil, 1)
	expect(t, `out = 1 & 0`, nil, 0)
	expect(t, `out = 0 & 1`, nil, 0)
	expect(t, `out = 0 & 0`, nil, 0)
	expect(t, `out = 1 | 1`, nil, 1)
	expect(t, `out = 1 | 0`, nil, 1)
	expect(t, `out = 0 | 1`, nil, 1)
	expect(t, `out = 0 | 0`, nil, 0)
	expect(t, `out = 1 ^ 1`, nil, 0)
	expect(t, `out = 1 ^ 0`, nil, 1)
	expect(t, `out = 0 ^ 1`, nil, 1)
	expect(t, `out = 0 ^ 0`, nil, 0)
	expect(t, `out = 1 &^ 1`, nil, 0)
	expect(t, `out = 1 &^ 0`, nil, 1)
	expect(t, `out = 0 &^ 1`, nil, 0)
	expect(t, `out = 0 &^ 0`, nil, 0)
	expect(t, `out = 1 << 2`, nil, 4)
	expect(t, `out = 16 >> 2`, nil, 4)

	expect(t, `out = 1; out &= 1`, nil, 1)
	expect(t, `out = 1; out |= 0`, nil, 1)
	expect(t, `out = 1; out ^= 0`, nil, 1)
	expect(t, `out = 1; out &^= 0`, nil, 1)
	expect(t, `out = 1; out <<= 2`, nil, 4)
	expect(t, `out = 16; out >>= 2`, nil, 4)

	expect(t, `out = ^0`, nil, ^0)
	expect(t, `out = ^1`, nil, ^1)
	expect(t, `out = ^55`, nil, ^55)
	expect(t, `out = ^-55`, nil, ^-55)
}

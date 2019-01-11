package runtime_test

import "testing"

func TestBitwise(t *testing.T) {
	expect(t, `out = 1 & 1`, 1&1)
	expect(t, `out = 1 & 0`, 1&0)
	expect(t, `out = 0 & 1`, 0&1)
	expect(t, `out = 0 & 0`, 0&0)
	expect(t, `out = 1 | 1`, 1|1)
	expect(t, `out = 1 | 0`, 1|0)
	expect(t, `out = 0 | 1`, 0|1)
	expect(t, `out = 0 | 0`, 0|0)
	expect(t, `out = 1 ^ 1`, 1^1)
	expect(t, `out = 1 ^ 0`, 1^0)
	expect(t, `out = 0 ^ 1`, 0^1)
	expect(t, `out = 0 ^ 0`, 0^0)
	expect(t, `out = 1 &^ 1`, 1&^1)
	expect(t, `out = 1 &^ 0`, 1&^0)
	expect(t, `out = 0 &^ 1`, 0&^1)
	expect(t, `out = 0 &^ 0`, 0&^0)
	expect(t, `out = 1 << 2`, 1<<2)
	expect(t, `out = 16 >> 2`, 16>>2)

	expect(t, `out = 1; out &= 1`, 1)
	expect(t, `out = 1; out |= 0`, 1)
	expect(t, `out = 1; out ^= 0`, 1)
	expect(t, `out = 1; out &^= 0`, 1)
	expect(t, `out = 1; out <<= 2`, 4)
	expect(t, `out = 16; out >>= 2`, 4)

	expect(t, `out = ^0`, ^0)
	expect(t, `out = ^1`, ^1)
	expect(t, `out = ^55`, ^55)
	expect(t, `out = ^-55`, ^-55)
}

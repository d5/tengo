package runtime_test

import (
	"testing"
)

func TestInteger(t *testing.T) {
	expect(t, `out = 5`, 5)
	expect(t, `out = 10`, 10)
	expect(t, `out = -5`, -5)
	expect(t, `out = -10`, -10)
	expect(t, `out = 5 + 5 + 5 + 5 - 10`, 10)
	expect(t, `out = 2 * 2 * 2 * 2 * 2`, 32)
	expect(t, `out = -50 + 100 + -50`, 0)
	expect(t, `out = 5 * 2 + 10`, 20)
	expect(t, `out = 5 + 2 * 10`, 25)
	expect(t, `out = 20 + 2 * -10`, 0)
	expect(t, `out = 50 / 2 * 2 + 10`, 60)
	expect(t, `out = 2 * (5 + 10)`, 30)
	expect(t, `out = 3 * 3 * 3 + 10`, 37)
	expect(t, `out = 3 * (3 * 3) + 10`, 37)
	expect(t, `out = (5 + 10 * 2 + 15 /3) * 2 + -10`, 50)
	expect(t, `out = 5 % 3`, 2)
	expect(t, `out = 5 % 3 + 4`, 6)
	expect(t, `out = +5`, 5)
	expect(t, `out = +5 + -5`, 0)

	expect(t, `out = 9 + '0'`, '9')
	expect(t, `out = '9' - 5`, '4')
}

package runtime_test

import (
	"testing"
)

func TestInteger(t *testing.T) {
	expect(t, `out = 5`, nil, 5)
	expect(t, `out = 10`, nil, 10)
	expect(t, `out = -5`, nil, -5)
	expect(t, `out = -10`, nil, -10)
	expect(t, `out = 5 + 5 + 5 + 5 - 10`, nil, 10)
	expect(t, `out = 2 * 2 * 2 * 2 * 2`, nil, 32)
	expect(t, `out = -50 + 100 + -50`, nil, 0)
	expect(t, `out = 5 * 2 + 10`, nil, 20)
	expect(t, `out = 5 + 2 * 10`, nil, 25)
	expect(t, `out = 20 + 2 * -10`, nil, 0)
	expect(t, `out = 50 / 2 * 2 + 10`, nil, 60)
	expect(t, `out = 2 * (5 + 10)`, nil, 30)
	expect(t, `out = 3 * 3 * 3 + 10`, nil, 37)
	expect(t, `out = 3 * (3 * 3) + 10`, nil, 37)
	expect(t, `out = (5 + 10 * 2 + 15 /3) * 2 + -10`, nil, 50)
	expect(t, `out = 5 % 3`, nil, 2)
	expect(t, `out = 5 % 3 + 4`, nil, 6)
	expect(t, `out = +5`, nil, 5)
	expect(t, `out = +5 + -5`, nil, 0)

	expect(t, `out = 9 + '0'`, nil, '9')
	expect(t, `out = '9' - 5`, nil, '4')
}

package runtime_test

import "testing"

func TestChar(t *testing.T) {
	expect(t, `out = 'a'`, nil, 'a')
	expect(t, `out = '九'`, nil, rune(20061))
	expect(t, `out = 'Æ'`, nil, rune(198))

	expect(t, `out = '0' + '9'`, nil, rune(105))
	expect(t, `out = '0' + 9`, nil, '9')
	expect(t, `out = '9' - 4`, nil, '5')
	expect(t, `out = '0' == '0'`, nil, true)
	expect(t, `out = '0' != '0'`, nil, false)
	expect(t, `out = '2' < '4'`, nil, true)
	expect(t, `out = '2' > '4'`, nil, false)
	expect(t, `out = '2' <= '4'`, nil, true)
	expect(t, `out = '2' >= '4'`, nil, false)
	expect(t, `out = '4' < '4'`, nil, false)
	expect(t, `out = '4' > '4'`, nil, false)
	expect(t, `out = '4' <= '4'`, nil, true)
	expect(t, `out = '4' >= '4'`, nil, true)
}

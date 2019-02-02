package runtime_test

import "testing"

func TestChar(t *testing.T) {
	expect(t, `out = 'a'`, 'a')
	expect(t, `out = '九'`, rune(20061))
	expect(t, `out = 'Æ'`, rune(198))

	expect(t, `out = '0' + '9'`, rune(105))
	expect(t, `out = '0' + 9`, '9')
	expect(t, `out = '9' - 4`, '5')
	expect(t, `out = '0' == '0'`, true)
	expect(t, `out = '0' != '0'`, false)
	expect(t, `out = '2' < '4'`, true)
	expect(t, `out = '2' > '4'`, false)
	expect(t, `out = '2' <= '4'`, true)
	expect(t, `out = '2' >= '4'`, false)
	expect(t, `out = '4' < '4'`, false)
	expect(t, `out = '4' > '4'`, false)
	expect(t, `out = '4' <= '4'`, true)
	expect(t, `out = '4' >= '4'`, true)
}

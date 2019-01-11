package runtime_test

import "testing"

func TestChar(t *testing.T) {
	expect(t, `out = 'a'`, 'a')
	expect(t, `out = '九'`, rune(20061))
	expect(t, `out = 'Æ'`, rune(198))
}

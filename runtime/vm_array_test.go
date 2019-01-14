package runtime_test

import (
	"fmt"
	"testing"
)

func TestArray(t *testing.T) {
	expect(t, `out = [1, 2 * 2, 3 + 3]`, ARR{1, 4, 6})

	// array copy-by-reference
	expect(t, `a1 := [1, 2, 3]; a2 := a1; a1[0] = 5; out = a2`, ARR{5, 2, 3})
	expect(t, `func () { a1 := [1, 2, 3]; a2 := a1; a1[0] = 5; out = a2 }()`, ARR{5, 2, 3})

	// index operator
	arr := ARR{1, 2, 3, 4, 5, 6}
	arrStr := `[1, 2, 3, 4, 5, 6]`
	arrLen := 6
	for idx := 0; idx < arrLen; idx++ {
		expect(t, fmt.Sprintf("out = %s[%d]", arrStr, idx), arr[idx])
		expect(t, fmt.Sprintf("out = %s[0 + %d]", arrStr, idx), arr[idx])
		expect(t, fmt.Sprintf("out = %s[1 + %d - 1]", arrStr, idx), arr[idx])
		expect(t, fmt.Sprintf("idx := %d; out = %s[idx]", idx, arrStr), arr[idx])
	}
	expectError(t, fmt.Sprintf("%s[%d]", arrStr, -1))
	expectError(t, fmt.Sprintf("%s[%d]", arrStr, arrLen))

	// slice operator
	for low := 0; low < arrLen; low++ {
		for high := low; high <= arrLen; high++ {
			expect(t, fmt.Sprintf("out = %s[%d:%d]", arrStr, low, high), arr[low:high])
			expect(t, fmt.Sprintf("out = %s[0 + %d : 0 + %d]", arrStr, low, high), arr[low:high])
			expect(t, fmt.Sprintf("out = %s[1 + %d - 1 : 1 + %d - 1]", arrStr, low, high), arr[low:high])
			expect(t, fmt.Sprintf("out = %s[:%d]", arrStr, high), arr[:high])
			expect(t, fmt.Sprintf("out = %s[%d:]", arrStr, low), arr[low:])
			expect(t, fmt.Sprintf("out = %s[:]", arrStr), arr[:])
		}
	}
	expectError(t, fmt.Sprintf("%s[%d:]", arrStr, -1))
	expectError(t, fmt.Sprintf("%s[:%d]", arrStr, arrLen+1))
	expectError(t, fmt.Sprintf("%s[%d:%d]", arrStr, 2, 1))
}

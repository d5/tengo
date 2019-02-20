package runtime_test

import (
	"fmt"
	"testing"

	"github.com/d5/tengo/objects"
)

func TestArray(t *testing.T) {
	expect(t, `out = [1, 2 * 2, 3 + 3]`, ARR{1, 4, 6})

	// array copy-by-reference
	expect(t, `a1 := [1, 2, 3]; a2 := a1; a1[0] = 5; out = a2`, ARR{5, 2, 3})
	expect(t, `func () { a1 := [1, 2, 3]; a2 := a1; a1[0] = 5; out = a2 }()`, ARR{5, 2, 3})

	// array index set
	expectError(t, `a1 := [1, 2, 3]; a1[3] = 5`, "index out of bounds")

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

	expect(t, fmt.Sprintf("%s[%d]", arrStr, -1), objects.UndefinedValue)
	expect(t, fmt.Sprintf("%s[%d]", arrStr, arrLen), objects.UndefinedValue)

	// slice operator
	for low := 0; low < arrLen; low++ {
		expect(t, fmt.Sprintf("out = %s[%d:%d]", arrStr, low, low), ARR{})
		for high := low; high <= arrLen; high++ {
			expect(t, fmt.Sprintf("out = %s[%d:%d]", arrStr, low, high), arr[low:high])
			expect(t, fmt.Sprintf("out = %s[0 + %d : 0 + %d]", arrStr, low, high), arr[low:high])
			expect(t, fmt.Sprintf("out = %s[1 + %d - 1 : 1 + %d - 1]", arrStr, low, high), arr[low:high])
			expect(t, fmt.Sprintf("out = %s[:%d]", arrStr, high), arr[:high])
			expect(t, fmt.Sprintf("out = %s[%d:]", arrStr, low), arr[low:])
		}
	}

	expect(t, fmt.Sprintf("out = %s[:]", arrStr), arr)
	expect(t, fmt.Sprintf("out = %s[%d:]", arrStr, -1), arr)
	expect(t, fmt.Sprintf("out = %s[:%d]", arrStr, arrLen+1), arr)
	expect(t, fmt.Sprintf("out = %s[%d:%d]", arrStr, 2, 2), ARR{})

	expectError(t, fmt.Sprintf("%s[:%d]", arrStr, -1), "invalid slice index")
	expectError(t, fmt.Sprintf("%s[%d:]", arrStr, arrLen+1), "invalid slice index")
	expectError(t, fmt.Sprintf("%s[%d:%d]", arrStr, 0, -1), "invalid slice index")
	expectError(t, fmt.Sprintf("%s[%d:%d]", arrStr, 2, 1), "invalid slice index")
}

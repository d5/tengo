package runtime_test

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	expect(t, `out = "Hello World!"`, "Hello World!")
	expect(t, `out = "Hello" + " " + "World!"`, "Hello World!")

	expect(t, `out = "Hello" == "Hello"`, true)
	expect(t, `out = "Hello" == "World"`, false)
	expect(t, `out = "Hello" != "Hello"`, false)
	expect(t, `out = "Hello" != "World"`, true)

	// index operator
	str := "abcdef"
	strStr := `"abcdef"`
	strLen := 6
	for idx := 0; idx < strLen; idx++ {
		expect(t, fmt.Sprintf("out = %s[%d]", strStr, idx), str[idx])
		expect(t, fmt.Sprintf("out = %s[0 + %d]", strStr, idx), str[idx])
		expect(t, fmt.Sprintf("out = %s[1 + %d - 1]", strStr, idx), str[idx])
		expect(t, fmt.Sprintf("idx := %d; out = %s[idx]", idx, strStr), str[idx])
	}
	expectError(t, fmt.Sprintf("%s[%d]", strStr, -1))
	expectError(t, fmt.Sprintf("%s[%d]", strStr, strLen))

	// slice operator
	for low := 0; low < strLen; low++ {
		for high := low; high <= strLen; high++ {
			expect(t, fmt.Sprintf("out = %s[%d:%d]", strStr, low, high), str[low:high])
			expect(t, fmt.Sprintf("out = %s[0 + %d : 0 + %d]", strStr, low, high), str[low:high])
			expect(t, fmt.Sprintf("out = %s[1 + %d - 1 : 1 + %d - 1]", strStr, low, high), str[low:high])
			expect(t, fmt.Sprintf("out = %s[:%d]", strStr, high), str[:high])
			expect(t, fmt.Sprintf("out = %s[%d:]", strStr, low), str[low:])
			expect(t, fmt.Sprintf("out = %s[:]", strStr), str[:])
		}
	}
	expectError(t, fmt.Sprintf("%s[%d:]", strStr, -1))
	expectError(t, fmt.Sprintf("%s[:%d]", strStr, strLen+1))
	expectError(t, fmt.Sprintf("%s[%d:%d]", strStr, 2, 1))
}

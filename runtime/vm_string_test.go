package runtime_test

import (
	"fmt"
	"testing"

	"github.com/d5/tengo/objects"
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

	expect(t, fmt.Sprintf("%s[%d]", strStr, -1), objects.UndefinedValue)
	expect(t, fmt.Sprintf("%s[%d]", strStr, strLen), objects.UndefinedValue)

	// slice operator
	for low := 0; low <= strLen; low++ {
		expect(t, fmt.Sprintf("out = %s[%d:%d]", strStr, low, low), "")
		for high := low; high <= strLen; high++ {
			expect(t, fmt.Sprintf("out = %s[%d:%d]", strStr, low, high), str[low:high])
			expect(t, fmt.Sprintf("out = %s[0 + %d : 0 + %d]", strStr, low, high), str[low:high])
			expect(t, fmt.Sprintf("out = %s[1 + %d - 1 : 1 + %d - 1]", strStr, low, high), str[low:high])
			expect(t, fmt.Sprintf("out = %s[:%d]", strStr, high), str[:high])
			expect(t, fmt.Sprintf("out = %s[%d:]", strStr, low), str[low:])
		}
	}

	expect(t, fmt.Sprintf("out = %s[:]", strStr), str[:])
	expect(t, fmt.Sprintf("out = %s[:]", strStr), str)
	expect(t, fmt.Sprintf("out = %s[%d:]", strStr, -1), str)
	expect(t, fmt.Sprintf("out = %s[:%d]", strStr, strLen+1), str)
	expect(t, fmt.Sprintf("out = %s[%d:%d]", strStr, 2, 2), "")

	expectError(t, fmt.Sprintf("%s[:%d]", strStr, -1), "invalid slice index")
	expectError(t, fmt.Sprintf("%s[%d:]", strStr, strLen+1), "invalid slice index")
	expectError(t, fmt.Sprintf("%s[%d:%d]", strStr, 0, -1), "invalid slice index")
	expectError(t, fmt.Sprintf("%s[%d:%d]", strStr, 2, 1), "invalid slice index")

	// string concatenation with other types
	expect(t, `out = "foo" + 1`, "foo1")
	// Float.String() returns the smallest number of digits
	// necessary such that ParseFloat will return f exactly.
	expect(t, `out = "foo" + 1.0`, "foo1") // <- note '1' instead of '1.0'
	expect(t, `out = "foo" + 1.5`, "foo1.5")
	expect(t, `out = "foo" + true`, "footrue")
	expect(t, `out = "foo" + 'X'`, "fooX")
	expect(t, `out = "foo" + error(5)`, "fooerror: 5")
	expect(t, `out = "foo" + undefined`, "foo<undefined>")
	expect(t, `out = "foo" + [1,2,3]`, "foo[1, 2, 3]")
	//expect(t, `out = "foo" + {a: 1, b: 2}`, "foo{a: 1, b: 2}") // TODO: commented because order of key is not consistent
	// also works with "+=" operator
	expect(t, `out = "foo"; out += 1.5`, "foo1.5")
	// string concats works only when string is LHS
	expectError(t, `1 + "foo"`, "invalid operation")

	expectError(t, `"foo" - "bar"`, "invalid operation")
}

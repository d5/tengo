package runtime_test

import (
	"fmt"
	"testing"

	"github.com/d5/tengo/objects"
)

func TestString(t *testing.T) {
	expect(t, `out = "Hello World!"`, nil, "Hello World!")
	expect(t, `out = "Hello" + " " + "World!"`, nil, "Hello World!")

	expect(t, `out = "Hello" == "Hello"`, nil, true)
	expect(t, `out = "Hello" == "World"`, nil, false)
	expect(t, `out = "Hello" != "Hello"`, nil, false)
	expect(t, `out = "Hello" != "World"`, nil, true)

	// index operator
	str := "abcdef"
	strStr := `"abcdef"`
	strLen := 6
	for idx := 0; idx < strLen; idx++ {
		expect(t, fmt.Sprintf("out = %s[%d]", strStr, idx), nil, str[idx])
		expect(t, fmt.Sprintf("out = %s[0 + %d]", strStr, idx), nil, str[idx])
		expect(t, fmt.Sprintf("out = %s[1 + %d - 1]", strStr, idx), nil, str[idx])
		expect(t, fmt.Sprintf("idx := %d; out = %s[idx]", idx, strStr), nil, str[idx])
	}

	expect(t, fmt.Sprintf("%s[%d]", strStr, -1), nil, objects.UndefinedValue)
	expect(t, fmt.Sprintf("%s[%d]", strStr, strLen), nil, objects.UndefinedValue)

	// slice operator
	for low := 0; low <= strLen; low++ {
		expect(t, fmt.Sprintf("out = %s[%d:%d]", strStr, low, low), nil, "")
		for high := low; high <= strLen; high++ {
			expect(t, fmt.Sprintf("out = %s[%d:%d]", strStr, low, high), nil, str[low:high])
			expect(t, fmt.Sprintf("out = %s[0 + %d : 0 + %d]", strStr, low, high), nil, str[low:high])
			expect(t, fmt.Sprintf("out = %s[1 + %d - 1 : 1 + %d - 1]", strStr, low, high), nil, str[low:high])
			expect(t, fmt.Sprintf("out = %s[:%d]", strStr, high), nil, str[:high])
			expect(t, fmt.Sprintf("out = %s[%d:]", strStr, low), nil, str[low:])
		}
	}

	expect(t, fmt.Sprintf("out = %s[:]", strStr), nil, str[:])
	expect(t, fmt.Sprintf("out = %s[:]", strStr), nil, str)
	expect(t, fmt.Sprintf("out = %s[%d:]", strStr, -1), nil, str)
	expect(t, fmt.Sprintf("out = %s[:%d]", strStr, strLen+1), nil, str)
	expect(t, fmt.Sprintf("out = %s[%d:%d]", strStr, 2, 2), nil, "")

	expectError(t, fmt.Sprintf("%s[:%d]", strStr, -1), nil, "invalid slice index")
	expectError(t, fmt.Sprintf("%s[%d:]", strStr, strLen+1), nil, "invalid slice index")
	expectError(t, fmt.Sprintf("%s[%d:%d]", strStr, 0, -1), nil, "invalid slice index")
	expectError(t, fmt.Sprintf("%s[%d:%d]", strStr, 2, 1), nil, "invalid slice index")

	// string concatenation with other types
	expect(t, `out = "foo" + 1`, nil, "foo1")
	// Float.String() returns the smallest number of digits
	// necessary such that ParseFloat will return f exactly.
	expect(t, `out = "foo" + 1.0`, nil, "foo1") // <- note '1' instead of '1.0'
	expect(t, `out = "foo" + 1.5`, nil, "foo1.5")
	expect(t, `out = "foo" + true`, nil, "footrue")
	expect(t, `out = "foo" + 'X'`, nil, "fooX")
	expect(t, `out = "foo" + error(5)`, nil, "fooerror: 5")
	expect(t, `out = "foo" + undefined`, nil, "foo<undefined>")
	expect(t, `out = "foo" + [1,2,3]`, nil, "foo[1, 2, 3]")
	// also works with "+=" operator
	expect(t, `out = "foo"; out += 1.5`, nil, "foo1.5")
	// string concats works only when string is LHS
	expectError(t, `1 + "foo"`, nil, "invalid operation")

	expectError(t, `"foo" - "bar"`, nil, "invalid operation")
}

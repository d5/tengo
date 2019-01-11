package runtime_test

import (
	"testing"
)

func TestString(t *testing.T) {
	expect(t, `out = "Hello World!"`, "Hello World!")
	expect(t, `out = "Hello" + " " + "World!"`, "Hello World!")

	expect(t, `out = "Hello" == "Hello"`, true)
	expect(t, `out = "Hello" == "World"`, false)
	expect(t, `out = "Hello" != "Hello"`, false)
	expect(t, `out = "Hello" != "World"`, true)

	expect(t, `out = "abcde"[0]`, 'a')
	expect(t, `out = "abcde"[4]`, 'e')
	expectError(t, `out = "abcde"[-1]`)
	expectError(t, `out = "abcde"[5]`)

	expect(t, `out = "abcde"[:]`, "abcde")
	expect(t, `out = "abcde"[0:5]`, "abcde")
	expect(t, `out = "abcde"[1:]`, "bcde")
	expect(t, `out = "abcde"[:4]`, "abcd")
	expect(t, `out = "abcde"[1:4]`, "bcd")
	expect(t, `out = "abcde"[2:3]`, "c")
	expect(t, `out = "abcde"[2:2]`, "")
	expectError(t, `out = "abcde"[-1:]`)
	expectError(t, `out = "abcde"[:6]`)
	expectError(t, `out = "abcde"[-1:6]`)
	expectError(t, `out = "abcde"[3:2]`)
}

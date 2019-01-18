package runtime_test

import (
	"testing"
)

func TestBytes(t *testing.T) {
	expect(t, `out = bytes("Hello World!")`, []byte("Hello World!"))
	expect(t, `out = bytes("Hello") + bytes(" ") + bytes("World!")`, []byte("Hello World!"))

	// bytes[] -> int
	expect(t, `out = bytes("abcde")[0]`, 97)
	expect(t, `out = bytes("abcde")[1]`, 98)
	expect(t, `out = bytes("abcde")[4]`, 101)
}

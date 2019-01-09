package vm_test

import (
	"testing"
)

func TestForIn(t *testing.T) {
	expect(t, `for i, x in [1, 2, 3] { out += i + x }`, 9)
	expect(t, `func() { for i, x in [1, 2, 3] { out += i + x } }()`, 9)

	expect(t, `for i, _ in [1, 2, 3] { out += i }`, 3)
	expect(t, `func() { for i, _ in [1, 2, 3] { out += i  } }()`, 3)

	expect(t, `for k, v in {a:2,b:3,c:4} { out = k; if v==3 { break } }`, "b")
	expect(t, `func() { for k, v in {a:2,b:3,c:4} { out = k; if v==3 { break } } }()`, "b")

	expect(t, `for c in "abcde" { out += c }`, "abcde")
	expect(t, `for i, c in "abcde" { if i == 2 { continue }; out += c }`, "abde")
}

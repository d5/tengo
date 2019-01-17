package runtime_test

import (
	"testing"
)

func TestForIn(t *testing.T) {
	// array
	expect(t, `for x in [1, 2, 3] { out += x }`, 6)                     // value
	expect(t, `for i, x in [1, 2, 3] { out += i + x }`, 9)              // index, value
	expect(t, `func() { for i, x in [1, 2, 3] { out += i + x } }()`, 9) // index, value
	expect(t, `for i, _ in [1, 2, 3] { out += i }`, 3)                  // index, _
	expect(t, `func() { for i, _ in [1, 2, 3] { out += i  } }()`, 3)    // index, _

	// map
	expect(t, `for v in {a:2,b:3,c:4} { out += v }`, 9)                                     // value
	expect(t, `for k, v in {a:2,b:3,c:4} { out = k; if v==3 { break } }`, "b")              // key, value
	expect(t, `for k, _ in {a:2} { out += k }`, "a")                                        // key, _
	expect(t, `for _, v in {a:2,b:3,c:4} { out += v }`, 9)                                  // _, value
	expect(t, `func() { for k, v in {a:2,b:3,c:4} { out = k; if v==3 { break } } }()`, "b") // key, value

	// string
	expect(t, `for c in "abcde" { out += c }`, "abcde")
	expect(t, `for i, c in "abcde" { if i == 2 { continue }; out += c }`, "abde")
}

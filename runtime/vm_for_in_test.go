package runtime_test

import (
	"testing"
)

func TestForIn(t *testing.T) {
	// array
	expect(t, `out = 0; for x in [1, 2, 3] { out += x }`, nil, 6)                     // value
	expect(t, `out = 0; for i, x in [1, 2, 3] { out += i + x }`, nil, 9)              // index, value
	expect(t, `out = 0; func() { for i, x in [1, 2, 3] { out += i + x } }()`, nil, 9) // index, value
	expect(t, `out = 0; for i, _ in [1, 2, 3] { out += i }`, nil, 3)                  // index, _
	expect(t, `out = 0; func() { for i, _ in [1, 2, 3] { out += i  } }()`, nil, 3)    // index, _

	// map
	expect(t, `out = 0; for v in {a:2,b:3,c:4} { out += v }`, nil, 9)                                      // value
	expect(t, `out = ""; for k, v in {a:2,b:3,c:4} { out = k; if v==3 { break } }`, nil, "b")              // key, value
	expect(t, `out = ""; for k, _ in {a:2} { out += k }`, nil, "a")                                        // key, _
	expect(t, `out = 0; for _, v in {a:2,b:3,c:4} { out += v }`, nil, 9)                                   // _, value
	expect(t, `out = ""; func() { for k, v in {a:2,b:3,c:4} { out = k; if v==3 { break } } }()`, nil, "b") // key, value

	// string
	expect(t, `out = ""; for c in "abcde" { out += c }`, nil, "abcde")
	expect(t, `out = ""; for i, c in "abcde" { if i == 2 { continue }; out += c }`, nil, "abde")
}

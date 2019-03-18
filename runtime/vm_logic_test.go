package runtime_test

import "testing"

func TestLogical(t *testing.T) {
	expect(t, `out = true && true`, nil, true)
	expect(t, `out = true && false`, nil, false)
	expect(t, `out = false && true`, nil, false)
	expect(t, `out = false && false`, nil, false)
	expect(t, `out = !true && true`, nil, false)
	expect(t, `out = !true && false`, nil, false)
	expect(t, `out = !false && true`, nil, true)
	expect(t, `out = !false && false`, nil, false)

	expect(t, `out = true || true`, nil, true)
	expect(t, `out = true || false`, nil, true)
	expect(t, `out = false || true`, nil, true)
	expect(t, `out = false || false`, nil, false)
	expect(t, `out = !true || true`, nil, true)
	expect(t, `out = !true || false`, nil, false)
	expect(t, `out = !false || true`, nil, true)
	expect(t, `out = !false || false`, nil, true)

	expect(t, `out = 1 && 2`, nil, 2)
	expect(t, `out = 1 || 2`, nil, 1)
	expect(t, `out = 1 && 0`, nil, 0)
	expect(t, `out = 1 || 0`, nil, 1)
	expect(t, `out = 1 && (0 || 2)`, nil, 2)
	expect(t, `out = 0 || (0 || 2)`, nil, 2)
	expect(t, `out = 0 || (0 && 2)`, nil, 0)
	expect(t, `out = 0 || (2 && 0)`, nil, 0)

	expect(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; t() && f()`, nil, 7)
	expect(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; f() && t()`, nil, 7)
	expect(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; f() || t()`, nil, 3)
	expect(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; t() || f()`, nil, 3)
	expect(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; !t() && f()`, nil, 3)
	expect(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; !f() && t()`, nil, 3)
	expect(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; !f() || t()`, nil, 7)
	expect(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; !t() || f()`, nil, 7)
}

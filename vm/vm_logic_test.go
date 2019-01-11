package vm_test

import "testing"

func TestLogical(t *testing.T) {
	expect(t, `out = true && true`, true)
	expect(t, `out = true && false`, false)
	expect(t, `out = false && true`, false)
	expect(t, `out = false && false`, false)
	expect(t, `out = !true && true`, false)
	expect(t, `out = !true && false`, false)
	expect(t, `out = !false && true`, true)
	expect(t, `out = !false && false`, false)

	expect(t, `out = true || true`, true)
	expect(t, `out = true || false`, true)
	expect(t, `out = false || true`, true)
	expect(t, `out = false || false`, false)
	expect(t, `out = !true || true`, true)
	expect(t, `out = !true || false`, false)
	expect(t, `out = !false || true`, true)
	expect(t, `out = !false || false`, true)

	expect(t, `out = 1 && 2`, 2)
	expect(t, `out = 1 || 2`, 1)
	expect(t, `out = 1 && 0`, 0)
	expect(t, `out = 1 || 0`, 1)
	expect(t, `out = 1 && (0 || 2)`, 2)
	expect(t, `out = 0 || (0 || 2)`, 2)
	expect(t, `out = 0 || (0 && 2)`, 0)
	expect(t, `out = 0 || (2 && 0)`, 0)

	expect(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; t() && f()`, 7)
	expect(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; f() && t()`, 7)
	expect(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; f() || t()`, 3)
	expect(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; t() || f()`, 3)
	expect(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; !t() && f()`, 3)
	expect(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; !f() && t()`, 3)
	expect(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; !f() || t()`, 7)
	expect(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; !t() || f()`, 7)
}

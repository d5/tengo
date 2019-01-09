package vm_test

import (
	"testing"
)

func TestBangOperator(t *testing.T) {
	expect(t, `out = !true`, false)
	expect(t, `out = !false`, true)
	expect(t, `out = !0`, true)
	expect(t, `out = !5`, false)
	expect(t, `out = !!true`, true)
	expect(t, `out = !!false`, false)
	expect(t, `out = !!5`, true)
}

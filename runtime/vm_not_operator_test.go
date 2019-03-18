package runtime_test

import (
	"testing"
)

func TestBangOperator(t *testing.T) {
	expect(t, `out = !true`, nil, false)
	expect(t, `out = !false`, nil, true)
	expect(t, `out = !0`, nil, true)
	expect(t, `out = !5`, nil, false)
	expect(t, `out = !!true`, nil, true)
	expect(t, `out = !!false`, nil, false)
	expect(t, `out = !!5`, nil, true)
}

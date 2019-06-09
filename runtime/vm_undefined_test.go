package runtime_test

import (
	"testing"

	"github.com/d5/tengo"
)

func TestUndefined(t *testing.T) {
	expect(t, `out = undefined`, nil, tengo.UndefinedValue)
	expect(t, `out = undefined.a`, nil, tengo.UndefinedValue)
	expect(t, `out = undefined[1]`, nil, tengo.UndefinedValue)
	expect(t, `out = undefined.a.b`, nil, tengo.UndefinedValue)
	expect(t, `out = undefined[1][2]`, nil, tengo.UndefinedValue)
	expect(t, `out = undefined ? 1 : 2`, nil, 2)
	expect(t, `out = undefined == undefined`, nil, true)
	expect(t, `out = undefined == 1`, nil, false)
	expect(t, `out = 1 == undefined`, nil, false)
	expect(t, `out = undefined == float([])`, nil, true)
	expect(t, `out = float([]) == undefined`, nil, true)
}

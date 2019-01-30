package runtime_test

import (
	"testing"

	"github.com/d5/tengo/objects"
)

func TestUndefined(t *testing.T) {
	expect(t, `out = undefined`, objects.UndefinedValue)
	expect(t, `out = undefined == undefined`, true)
	expect(t, `out = undefined == 1`, false)
	expect(t, `out = 1 == undefined`, false)
	expect(t, `out = undefined == float([])`, true)
	expect(t, `out = float([]) == undefined`, true)
}

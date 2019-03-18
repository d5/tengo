package runtime_test

import (
	"testing"

	"github.com/d5/tengo/stdlib"
)

func TestSourceModules(t *testing.T) {
	expect(t, `enum := import("enum"); out = enum.any([1,2,3], func(i, v) { return v == 2 })`,
		Opts().Module("enum", stdlib.SourceModules["enum"]),
		true)
}

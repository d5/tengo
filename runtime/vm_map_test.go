package runtime_test

import (
	"testing"

	"github.com/d5/tengo"
)

func TestMap(t *testing.T) {
	expect(t, `
out = {
	one: 10 - 9,
	two: 1 + 1,
	three: 6 / 2
}`, nil, MAP{
		"one":   1,
		"two":   2,
		"three": 3,
	})

	expect(t, `
out = {
	"one": 10 - 9,
	"two": 1 + 1,
	"three": 6 / 2
}`, nil, MAP{
		"one":   1,
		"two":   2,
		"three": 3,
	})

	expect(t, `out = {foo: 5}["foo"]`, nil, 5)
	expect(t, `out = {foo: 5}["bar"]`, nil, tengo.UndefinedValue)
	expect(t, `key := "foo"; out = {foo: 5}[key]`, nil, 5)
	expect(t, `out = {}["foo"]`, nil, tengo.UndefinedValue)

	expect(t, `
m := {
	foo: func(x) {
		return x * 2
	}
}
out = m["foo"](2) + m["foo"](3)
`, nil, 10)

	// map assignment is copy-by-reference
	expect(t, `m1 := {k1: 1, k2: "foo"}; m2 := m1; m1.k1 = 5; out = m2.k1`, nil, 5)
	expect(t, `m1 := {k1: 1, k2: "foo"}; m2 := m1; m2.k1 = 3; out = m1.k1`, nil, 3)
	expect(t, `func() { m1 := {k1: 1, k2: "foo"}; m2 := m1; m1.k1 = 5; out = m2.k1 }()`, nil, 5)
	expect(t, `func() { m1 := {k1: 1, k2: "foo"}; m2 := m1; m2.k1 = 3; out = m1.k1 }()`, nil, 3)
}

package runtime_test

import (
	"testing"

	"github.com/d5/tengo/objects"
)

func TestSelector(t *testing.T) {
	expect(t, `a := {k1: 5, k2: "foo"}; out = a.k1`, 5)
	expect(t, `a := {k1: 5, k2: "foo"}; out = a.k2`, "foo")
	expect(t, `a := {k1: 5, k2: "foo"}; out = a.k3`, objects.UndefinedValue)

	expect(t, `
a := {
	b: {
		c: 4,
		a: false
	},
	c: "foo bar"
}
out = a.b.c`, 4)

	expect(t, `
a := {
	b: {
		c: 4,
		a: false
	},
	c: "foo bar"
}
b := a.x.c`, objects.UndefinedValue)

	expect(t, `
a := {
	b: {
		c: 4,
		a: false
	},
	c: "foo bar"
}
b := a.x.y`, objects.UndefinedValue)

	expect(t, `a := {b: 1, c: "foo"}; a.b = 2; out = a.b`, 2)
	expect(t, `a := {b: 1, c: "foo"}; a.c = 2; out = a.c`, 2) // type not checked on sub-field
	expect(t, `a := {b: {c: 1}}; a.b.c = 2; out = a.b.c`, 2)
	expect(t, `a := {b: 1}; a.c = 2; out = a`, MAP{"b": 1, "c": 2})
	expect(t, `a := {b: {c: 1}}; a.b.d = 2; out = a`, MAP{"b": MAP{"c": 1, "d": 2}})

	expect(t, `func() { a := {b: 1, c: "foo"}; a.b = 2; out = a.b }()`, 2)
	expect(t, `func() { a := {b: 1, c: "foo"}; a.c = 2; out = a.c }()`, 2) // type not checked on sub-field
	expect(t, `func() { a := {b: {c: 1}}; a.b.c = 2; out = a.b.c }()`, 2)
	expect(t, `func() { a := {b: 1}; a.c = 2; out = a }()`, MAP{"b": 1, "c": 2})
	expect(t, `func() { a := {b: {c: 1}}; a.b.d = 2; out = a }()`, MAP{"b": MAP{"c": 1, "d": 2}})

	expect(t, `func() { a := {b: 1, c: "foo"}; func() { a.b = 2 }(); out = a.b }()`, 2)
	expect(t, `func() { a := {b: 1, c: "foo"}; func() { a.c = 2 }(); out = a.c }()`, 2) // type not checked on sub-field
	expect(t, `func() { a := {b: {c: 1}}; func() { a.b.c = 2 }(); out = a.b.c }()`, 2)
	expect(t, `func() { a := {b: 1}; func() { a.c = 2 }(); out = a }()`, MAP{"b": 1, "c": 2})
	expect(t, `func() { a := {b: {c: 1}}; func() { a.b.d = 2 }(); out = a }()`, MAP{"b": MAP{"c": 1, "d": 2}})

	expect(t, `
a := {
	b: [1, 2, 3],
	c: {
		d: 8,
		e: "foo",
		f: [9, 8]
	}
}
out = [a.b[2], a.c.d, a.c.e, a.c.f[1]]
`, ARR{3, 8, "foo", 8})

	expect(t, `
func() {
	a := [1, 2, 3]
	b := 9
	a[1] = b
	b = 7     // make sure a[1] has a COPY of value of 'b'
	out = a[1]
}()
`, 9)

	expectErrorString(t, `a := {b: {c: 1}}; a.d.c = 2`, "not index-assignable")
	expectErrorString(t, `a := [1, 2, 3]; a.b = 2`, "invalid index type")
	expectErrorString(t, `a := "foo"; a.b = 2`, "not index-assignable")
	expectErrorString(t, `func() { a := {b: {c: 1}}; a.d.c = 2 }()`, "not index-assignable")
	expectErrorString(t, `func() { a := [1, 2, 3]; a.b = 2 }()`, "invalid index type")
	expectErrorString(t, `func() { a := "foo"; a.b = 2 }()`, "not index-assignable")
}

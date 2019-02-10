package runtime_test

import (
	"testing"

	"github.com/d5/tengo/objects"
)

func TestImmutable(t *testing.T) {
	// primitive types are already immutable values
	// immutable expression has no effects.
	expect(t, `a := immutable(1); out = a`, 1)
	expect(t, `a := 5; b := immutable(a); out = b`, 5)
	expect(t, `a := immutable(1); a = 5; out = a`, 5)

	// array
	expectError(t, `a := immutable([1, 2, 3]); a[1] = 5`)
	expectError(t, `a := immutable(["foo", [1,2,3]]); a[1] = "bar"`)
	expect(t, `a := immutable(["foo", [1,2,3]]); a[1][1] = "bar"; out = a`, IARR{"foo", ARR{1, "bar", 3}})
	expectError(t, `a := immutable(["foo", immutable([1,2,3])]); a[1][1] = "bar"`)
	expectError(t, `a := ["foo", immutable([1,2,3])]; a[1][1] = "bar"`)
	expect(t, `a := immutable([1,2,3]); b := copy(a); b[1] = 5; out = b`, ARR{1, 5, 3})
	expect(t, `a := immutable([1,2,3]); b := copy(a); b[1] = 5; out = a`, IARR{1, 2, 3})
	expect(t, `out = immutable([1,2,3]) == [1,2,3]`, true)
	expect(t, `out = immutable([1,2,3]) == immutable([1,2,3])`, true)
	expect(t, `out = [1,2,3] == immutable([1,2,3])`, true)
	expect(t, `out = immutable([1,2,3]) == [1,2]`, false)
	expect(t, `out = immutable([1,2,3]) == immutable([1,2])`, false)
	expect(t, `out = [1,2,3] == immutable([1,2])`, false)
	expect(t, `out = immutable([1, 2, 3, 4])[1]`, 2)
	expect(t, `out = immutable([1, 2, 3, 4])[1:3]`, ARR{2, 3})
	expect(t, `a := immutable([1,2,3]); a = 5; out = a`, 5)
	expect(t, `a := immutable([1, 2, 3]); out = a[5]`, objects.UndefinedValue)

	// map
	expectError(t, `a := immutable({b: 1, c: 2}); a.b = 5`)
	expectError(t, `a := immutable({b: 1, c: 2}); a["b"] = "bar"`)
	expect(t, `a := immutable({b: 1, c: [1,2,3]}); a.c[1] = "bar"; out = a`, IMAP{"b": 1, "c": ARR{1, "bar", 3}})
	expectError(t, `a := immutable({b: 1, c: immutable([1,2,3])}); a.c[1] = "bar"`)
	expectError(t, `a := {b: 1, c: immutable([1,2,3])}; a.c[1] = "bar"`)
	expect(t, `out = immutable({a:1,b:2}) == {a:1,b:2}`, true)
	expect(t, `out = immutable({a:1,b:2}) == immutable({a:1,b:2})`, true)
	expect(t, `out = {a:1,b:2} == immutable({a:1,b:2})`, true)
	expect(t, `out = immutable({a:1,b:2}) == {a:1,b:3}`, false)
	expect(t, `out = immutable({a:1,b:2}) == immutable({a:1,b:3})`, false)
	expect(t, `out = {a:1,b:2} == immutable({a:1,b:3})`, false)
	expect(t, `out = immutable({a:1,b:2}).b`, 2)
	expect(t, `out = immutable({a:1,b:2})["b"]`, 2)
	expect(t, `a := immutable({a:1,b:2}); a = 5; out = 5`, 5)
	expect(t, `a := immutable({a:1,b:2}); out = a.c`, objects.UndefinedValue)

	expect(t, `a := immutable({b: 5, c: "foo"}); out = a.b`, 5)
	expectError(t, `a := immutable({b: 5, c: "foo"}); a.b = 10`)
}

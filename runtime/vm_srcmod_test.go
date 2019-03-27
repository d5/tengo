package runtime_test

import "testing"

func TestSrcModEnum(t *testing.T) {
	expect(t, `
x := import("enum")
out = x.all([1, 2, 3], func(_, v) { return v >= 1 }) 
`, Opts().Stdlib(), true)
	expect(t, `
x := import("enum")
out = x.all([1, 2, 3], func(_, v) { return v >= 2 }) 
`, Opts().Stdlib(), false)

	expect(t, `
x := import("enum")
out = x.any([1, 2, 3], func(_, v) { return v >= 1 }) 
`, Opts().Stdlib(), true)
	expect(t, `
x := import("enum")
out = x.any([1, 2, 3], func(_, v) { return v >= 2 }) 
`, Opts().Stdlib(), true)

	expect(t, `
x := import("enum")
out = x.chunk([1, 2, 3], 1) 
`, Opts().Stdlib(), ARR{ARR{1}, ARR{2}, ARR{3}})
	expect(t, `
x := import("enum")
out = x.chunk([1, 2, 3], 2) 
`, Opts().Stdlib(), ARR{ARR{1, 2}, ARR{3}})
	expect(t, `
x := import("enum")
out = x.chunk([1, 2, 3], 3) 
`, Opts().Stdlib(), ARR{ARR{1, 2, 3}})
	expect(t, `
x := import("enum")
out = x.chunk([1, 2, 3], 4) 
`, Opts().Stdlib(), ARR{ARR{1, 2, 3}})
	expect(t, `
x := import("enum")
out = x.chunk([1, 2, 3, 4, 5, 6], 2) 
`, Opts().Stdlib(), ARR{ARR{1, 2}, ARR{3, 4}, ARR{5, 6}})

	expect(t, `
x := import("enum")
out = x.at([1, 2, 3], 0) 
`, Opts().Stdlib(), 1)
}

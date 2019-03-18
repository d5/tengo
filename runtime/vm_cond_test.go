package runtime_test

import "testing"

func TestCondExpr(t *testing.T) {
	expect(t, `out = true ? 5 : 10`, nil, 5)
	expect(t, `out = false ? 5 : 10`, nil, 10)
	expect(t, `out = (1 == 1) ? 2 + 3 : 12 - 2`, nil, 5)
	expect(t, `out = (1 != 1) ? 2 + 3 : 12 - 2`, nil, 10)
	expect(t, `out = (1 == 1) ? true ? 10 - 8 : 1 + 3 : 12 - 2`, nil, 2)
	expect(t, `out = (1 == 1) ? false ? 10 - 8 : 1 + 3 : 12 - 2`, nil, 4)

	expect(t, `
out = 0
f1 := func() { out += 10 }
f2 := func() { out = -out }
true ? f1() : f2()
`, nil, 10)
	expect(t, `
out = 5
f1 := func() { out += 10 }
f2 := func() { out = -out }
false ? f1() : f2()
`, nil, -5)
	expect(t, `
f1 := func(a) { return a + 2 }
f2 := func(a) { return a - 2 }
f3 := func(a) { return a + 10 }
f4 := func(a) { return -a }

f := func(c) {
	return c == 0 ? f1(c) : f2(c) ? f3(c) : f4(c)
}

out = [f(0), f(1), f(2)]
`, nil, ARR{2, 11, -2})

	expect(t, `f := func(a) { return -a }; out = f(true ? 5 : 3)`, nil, -5)
	expect(t, `out = [false?5:10, true?1:2]`, nil, ARR{10, 1})

	expect(t, `
out = 1 > 2 ?
	1 + 2 + 3 :
	10 - 5`, nil, 5)
}

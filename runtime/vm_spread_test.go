package runtime_test

import "testing"

func TestSpreadArray(t *testing.T) {
	expect(t, `out = [[1,2,3]...]`, nil, ARR{1, 2, 3})
	expect(t, `a := [1,2,3]; out = [a...]`, nil, ARR{1, 2, 3})
	expect(t, `a := [1,2]; b := [3]; out = [a..., b...]`, nil, ARR{1, 2, 3})
	expect(t, `a := [1]; b := [3]; out = [a..., 2, b...]`, nil, ARR{1, 2, 3})
	expect(t, `a := [1]; b := [2]; out = [a..., b..., 3]`, nil, ARR{1, 2, 3})
	expect(t, `a := [2]; b := [3]; out = [1, a..., b...]`, nil, ARR{1, 2, 3})
	expect(t, `out = [[[1]..., [2, 3]...]...]`, nil, ARR{1, 2, 3})
}

func TestSpreadCall(t *testing.T) {

	const defVars = `x := [1,2]; ` +
		`y := [3,4]; ` +
		`z := [5,6,7]; ` +
		`fn1 := func(...a) { return a; }; ` +
		`fn2 := func(a, ...b) { return [a, b]; }; ` +
		`fn3 := func(a, b, c) { return [a, b, c]; }; `

	expect(t, defVars+`out = fn1([1,2,3]...)`, nil, ARR{1, 2, 3})
	expect(t, defVars+`out = fn1(x..., 3)`, nil, ARR{1, 2, 3})
	expect(t, defVars+`out = fn1(x..., y..., z...)`, nil, ARR{1, 2, 3, 4, 5, 6, 7})
	expect(t, defVars+`out = fn1(1,2,3,4,z...)`, nil, ARR{1, 2, 3, 4, 5, 6, 7})

	expect(t, defVars+`out = fn2([1,2,3]...)`, nil, ARR{1, ARR{2, 3}})
	expect(t, defVars+`out = fn2(x..., 3)`, nil, ARR{1, ARR{2, 3}})
	expect(t, defVars+`out = fn2(x..., y..., z...)`, nil, ARR{1, ARR{2, 3, 4, 5, 6, 7}})
	expect(t, defVars+`out = fn2(1,2,3,4,z...)`, nil, ARR{1, ARR{2, 3, 4, 5, 6, 7}})

	expect(t, defVars+`out = fn3([1,2,3]...)`, nil, ARR{1, 2, 3})
	expect(t, defVars+`out = fn3(x..., 3)`, nil, ARR{1, 2, 3})
	expect(t, defVars+`out = fn3([x..., y...][:3]...)`, nil, ARR{1, 2, 3})
	expectError(t, defVars+`fn3(x..., y..., z...)`, nil, "Runtime Error: wrong number of arguments: want=3, got=7")
	expectError(t, defVars+`fn3(1,2,3,4,z...)`, nil, "Runtime Error: wrong number of arguments: want=3, got=7")
}

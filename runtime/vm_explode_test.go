package runtime_test

import (
	"testing"
)

func TestExplodeArray(t *testing.T) {
	expect(t, `a := [1, 2]; out = [a...]`, nil, ARR{1, 2})
	expect(t, `a := immutable([1, 2]); out = [a...]`, nil, ARR{1, 2})
	expect(t, `a := [1, 2]; out = [a..., 3]`, nil, ARR{1, 2, 3})
	expect(t, `a := [1, 2]; b := [3, 4]; out = [a..., b...]`, nil, ARR{1, 2, 3, 4})
	expect(t, `out = [[1, 2]...]`, nil, ARR{1, 2})
	expect(t, `out = [[1, 2]..., 3]`, nil, ARR{1, 2, 3})
	expect(t, `out = [[1, 2]..., [3, 4]...]`, nil, ARR{1, 2, 3, 4})
}

func TestExplodeCall(t *testing.T) {
	expect(t, `fn := func(a, b, c) { return [a, b, c] }; out = fn([1, 2, 3]...)`, nil, ARR{1, 2, 3})
	expect(t, `a := [1, 2]; fn := func(a, b, c) { return [a, b, c] }; out = fn(a..., 3)`, nil, ARR{1, 2, 3})
	expect(t, `a := [1, 2]; b := [3]; fn := func(a, b, c) { return [a, b, c] }; out = fn(a..., b...)`, nil, ARR{1, 2, 3})

	expectError(t, `a := [1, 2]; b := [3, 4]; fn := func(a, b, c) { return [a, b, c] }; fn(a..., b...)`, nil,
		"Runtime Error: wrong number of arguments: want=3, got=4")

	expectError(t, `a := [1, 2]; b := []; fn := func(a, b, c) { return [a, b, c] }; fn(a..., b...)`, nil,
		"Runtime Error: wrong number of arguments: want=3, got=2")
}

func TestExplodableTypes(t *testing.T) {
	expectError(t, `a := 1; [a...]`, nil, "Runtime Error: cannot explode object of type int")
	expectError(t, `a := 1.5; [a...]`, nil, "Runtime Error: cannot explode object of type float")
	expectError(t, `a := "abc"; [a...]`, nil, "Runtime Error: cannot explode object of type string")
	expectError(t, `a := bytes("abc"); [a...]`, nil, "Runtime Error: cannot explode object of type bytes")
	expectError(t, `a := {}; [a...]`, nil, "Runtime Error: cannot explode object of type map")
	expectError(t, `a := immutable({}); [a...]`, nil, "Runtime Error: cannot explode object of type immutable-map")
	expectError(t, `a := func(){return []}; [a...]`, nil, "Runtime Error: cannot explode object of type compiled-function")
	expectError(t, `a := (func() { inner := 0; return func(){return [inner]}; })(); [a...]`, nil, "Runtime Error: cannot explode object of type closure")
}

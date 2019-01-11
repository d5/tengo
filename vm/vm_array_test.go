package vm_test

import (
	"testing"
)

func TestArray(t *testing.T) {
	expect(t, `out = [1, 2 * 2, 3 + 3]`, ARR{1, 4, 6})

	expect(t, `out = [1, 2, 3][0]`, 1)
	expect(t, `out = [1, 2, 3][1]`, 2)
	expect(t, `out = [1, 2, 3][2]`, 3)
	expect(t, `i := 0; out = [1, 2, 3][i]`, 1)
	expect(t, `out = [1, 2, 3][1 + 1]`, 3)
	expect(t, `arr := [1, 2, 3]; out = arr[2]`, 3)
	expect(t, `arr := [1, 2, 3]; out = arr[0] + arr[1] + arr[2]`, 6)
	expect(t, `arr := [1, 2, 3]; i := arr[0]; out = arr[i]`, 2)

	expect(t, `out = [1, 2, 3][1+1]`, 3)
	expect(t, `a := 1; out = [1, 2, 3][a+1]`, 3)

	expect(t, `out = [1, 2, 3][:]`, ARR{1, 2, 3})
	expect(t, `out = [1, 2, 3][0:3]`, ARR{1, 2, 3})
	expect(t, `out = [1, 2, 3][1:]`, ARR{2, 3})
	expect(t, `out = [1, 2, 3][1:2]`, ARR{2})
	expect(t, `out = [1, 2, 3][:2]`, ARR{1, 2})
	expect(t, `out = [1, 2, 3][1:1]`, ARR{})

	expect(t, `out = [1, 2, 3][3-2:1+1]`, ARR{2})
	expect(t, `a := 1; out = [1, 2, 3][a-1:a+1]`, ARR{1, 2})

	// array copy-by-reference
	expect(t, `a1 := [1, 2, 3]; a2 := a1; a1[0] = 5; out = a2`, ARR{5, 2, 3})
	expect(t, `func () { a1 := [1, 2, 3]; a2 := a1; a1[0] = 5; out = a2 }()`, ARR{5, 2, 3})

	expectError(t, `[1, 2, 3][3]`)
	expectError(t, `[1, 2, 3][-1]`)

	expectError(t, `[1, 2, 3][-1:]`)
	expectError(t, `[1, 2, 3][:4]`)
	expectError(t, `[1, 2, 3][-1:3]`)
	expectError(t, `[1, 2, 3][0:4]`)
	expectError(t, `[1, 2, 3][2:1]`)
}

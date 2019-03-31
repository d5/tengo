package runtime_test

import (
	"testing"

	"github.com/d5/tengo/objects"
)

func TestIf(t *testing.T) {
	expect(t, `if (true) { out = 10 }`, nil, 10)
	expect(t, `if (false) { out = 10 }`, nil, objects.UndefinedValue)
	expect(t, `if (false) { out = 10 } else { out = 20 }`, nil, 20)
	expect(t, `if (1) { out = 10 }`, nil, 10)
	expect(t, `if (0) { out = 10 } else { out = 20 }`, nil, 20)
	expect(t, `if (1 < 2) { out = 10 }`, nil, 10)
	expect(t, `if (1 > 2) { out = 10 }`, nil, objects.UndefinedValue)
	expect(t, `if (1 < 2) { out = 10 } else { out = 20 }`, nil, 10)
	expect(t, `if (1 > 2) { out = 10 } else { out = 20 }`, nil, 20)

	expect(t, `if (1 < 2) { out = 10 } else if (1 > 2) { out = 20 } else { out = 30 }`, nil, 10)
	expect(t, `if (1 > 2) { out = 10 } else if (1 < 2) { out = 20 } else { out = 30 }`, nil, 20)
	expect(t, `if (1 > 2) { out = 10 } else if (1 == 2) { out = 20 } else { out = 30 }`, nil, 30)
	expect(t, `if (1 > 2) { out = 10 } else if (1 == 2) { out = 20 } else if (1 < 2) { out = 30 } else { out = 40 }`, nil, 30)
	expect(t, `if (1 > 2) { out = 10 } else if (1 < 2) { out = 20; out = 21; out = 22 } else { out = 30 }`, nil, 22)
	expect(t, `if (1 > 2) { out = 10 } else if (1 == 2) { out = 20 } else { out = 30; out = 31; out = 32}`, nil, 32)
	expect(t, `if (1 > 2) { out = 10 } else if (1 < 2) { if (1 == 2) { out = 21 } else { out = 22 } } else { out = 30 }`, nil, 22)
	expect(t, `if (1 > 2) { out = 10 } else if (1 < 2) { if (1 == 2) { out = 21 } else if (2 == 3) { out = 22 } else { out = 23 } } else { out = 30 }`, nil, 23)
	expect(t, `if (1 > 2) { out = 10 } else if (1 == 2) { if (1 == 2) { out = 21 } else if (2 == 3) { out = 22 } else { out = 23 } } else { out = 30 }`, nil, 30)
	expect(t, `if (1 > 2) { out = 10 } else if (1 == 2) { out = 20 } else { if (1 == 2) { out = 31 } else if (2 == 3) { out = 32 } else { out = 33 } }`, nil, 33)

	expect(t, `if a:=0; a<1 { out = 10 }`, nil, 10)
	expect(t, `a:=0; if a++; a==1 { out = 10 }`, nil, 10)
	expect(t, `
func() {
	a := 1
	if a++; a > 1 {
		out = a
	}
}()
`, nil, 2)
	expect(t, `
func() {
	a := 1
	if a++; a == 1 {
		out = 10
	} else {
		out = 20
	}
}()
`, nil, 20)
	expect(t, `
func() {
	a := 1

	func() {
		if a++; a > 1 {
			a++
		}
	}()

	out = a
}()
`, nil, 3)

	// expression statement in init (should not leave objects on stack)
	expect(t, `a := 1; if a; a { out = a }`, nil, 1)
	expect(t, `a := 1; if a + 4; a { out = a }`, nil, 1)

	// dead code elimination
	expect(t, `
out = func() {
	if false { return 1 }

	a := undefined

	a = 2
	if !a {
		b := func() {
			return is_callable(a) ? a(8) : a
		}()
		if is_error(b) { 
			return b 
		} else if !is_undefined(b) { 
			return immutable(b)
		}
	}
	
	a = 3
	if a {
		b := func() {
			return is_callable(a) ? a(9) : a
		}()
		if is_error(b) { 
			return b 
		} else if !is_undefined(b) { 
			return immutable(b)
		}
	}

	return a
}()
`, nil, 3)
}

package vm_test

import "testing"

func TestTailCall(t *testing.T) {
	expect(t, `
fac := func(n, a) {
	if n == 1 {
		return a
	}
	return fac(n-1, n*a)
}
out = fac(5, 1)`, 120)

	expect(t, `
fac := func(n, a) {
	if n == 1 {
		return a
	}
	x := {foo: fac} // indirection for test
	return x.foo(n-1, n*a)
}
out = fac(5, 1)`, 120)

	expect(t, `
fib := func(x, s) {
	if x == 0 {
		return 0 + s
	} else if x == 1 {
		return 1 + s
	}
	return fib(x-1, fib(x-2, s))
}
out = fib(15, 0)`, 610)

	// global variable and no return value
	expect(t, `
	out = 0
	foo := func(a) {
	   if a == 0 {
	       return
	   }
	   out += a
	   foo(a-1)
	}
	foo(10)`, 55)

	// TODO: see comment inside the code
	expect(t, `
f1 := func() {
	f2 := undefined  // TODO: this is really inconvenient 
	f2 = func(n, s) {
		if n == 0 { return s }
		return f2(n-1, n + s)
	}

	return f2(5, 0)
}

out = f1()`, 15)

	// tail call with free vars
	//	expect(t, `
	//f1 := func() {
	//	a := 10
	//	f2 := undefined
	//	f2 = func(n, s) {
	//		a += s
	//		if n == 0 {
	//			return
	//		}
	//		f2(n-1, n + s)
	//	}
	//	return f2
	//}
	//out = f1()(5, 0)`, 25)
}

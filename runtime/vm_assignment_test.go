package runtime_test

import (
	"testing"
)

func TestAssignment(t *testing.T) {
	expect(t, `a := 1; a = 2; out = a`, nil, 2)
	expect(t, `a := 1; a = 2; out = a`, nil, 2)
	expect(t, `a := 1; a = a + 4; out = a`, nil, 5)
	expect(t, `a := 1; f1 := func() { a = 2; return a }; out = f1()`, nil, 2)
	expect(t, `a := 1; f1 := func() { a := 3; a = 2; return a }; out = f1()`, nil, 2)

	expect(t, `a := 1; out = a`, nil, 1)
	expect(t, `a := 1; a = 2; out = a`, nil, 2)
	expect(t, `a := 1; func() { a = 2 }(); out = a`, nil, 2)
	expect(t, `a := 1; func() { a := 2 }(); out = a`, nil, 1) // "a := 2" defines a new local variable 'a'
	expect(t, `a := 1; func() { b := 2; out = b }()`, nil, 2)
	expect(t, `
out = func() { 
	a := 2
	func() {
		a = 3 // captured from outer scope
	}()
	return a
}()
`, nil, 3)

	expect(t, `
func() {
	a := 5
	out = func() {  	
		a := 4						
		return a
	}()
}()`, nil, 4)

	expectError(t, `a := 1; a := 2`, nil, "redeclared")              // redeclared in the same scope
	expectError(t, `func() { a := 1; a := 2 }()`, nil, "redeclared") // redeclared in the same scope

	expect(t, `a := 1; a += 2; out = a`, nil, 3)
	expect(t, `a := 1; a += 4 - 2;; out = a`, nil, 3)
	expect(t, `a := 3; a -= 1;; out = a`, nil, 2)
	expect(t, `a := 3; a -= 5 - 4;; out = a`, nil, 2)
	expect(t, `a := 2; a *= 4;; out = a`, nil, 8)
	expect(t, `a := 2; a *= 1 + 3;; out = a`, nil, 8)
	expect(t, `a := 10; a /= 2;; out = a`, nil, 5)
	expect(t, `a := 10; a /= 5 - 3;; out = a`, nil, 5)

	// compound assignment operator does not define new variable
	expectError(t, `a += 4`, nil, "unresolved reference")
	expectError(t, `a -= 4`, nil, "unresolved reference")
	expectError(t, `a *= 4`, nil, "unresolved reference")
	expectError(t, `a /= 4`, nil, "unresolved reference")

	expect(t, `
f1 := func() { 
	f2 := func() { 
		a := 1
		a += 2    // it's a statement, not an expression
		return a
	}; 
	
	return f2(); 
}; 

out = f1();`, nil, 3)
	expect(t, `f1 := func() { f2 := func() { a := 1; a += 4 - 2; return a }; return f2(); }; out = f1()`, nil, 3)
	expect(t, `f1 := func() { f2 := func() { a := 3; a -= 1; return a }; return f2(); }; out = f1()`, nil, 2)
	expect(t, `f1 := func() { f2 := func() { a := 3; a -= 5 - 4; return a }; return f2(); }; out = f1()`, nil, 2)
	expect(t, `f1 := func() { f2 := func() { a := 2; a *= 4; return a }; return f2(); }; out = f1()`, nil, 8)
	expect(t, `f1 := func() { f2 := func() { a := 2; a *= 1 + 3; return a }; return f2(); }; out = f1()`, nil, 8)
	expect(t, `f1 := func() { f2 := func() { a := 10; a /= 2; return a }; return f2(); }; out = f1()`, nil, 5)
	expect(t, `f1 := func() { f2 := func() { a := 10; a /= 5 - 3; return a }; return f2(); }; out = f1()`, nil, 5)

	expect(t, `a := 1; f1 := func() { f2 := func() { a += 2; return a }; return f2(); }; out = f1()`, nil, 3)

	expect(t, `
	f1 := func(a) {
		return func(b) {
			c := a
			c += b * 2
			return c
		}
	}
	
	out = f1(3)(4)
	`, nil, 11)

	expect(t, `
	out = func() {
		a := 1
		func() {
			a = 2
			func() {
				a = 3
				func() {
					a := 4 // declared new
				}()
			}()
		}()
		return a
	}()
	`, nil, 3)

	// write on free variables
	expect(t, `
	f1 := func() {
		a := 5
	
		return func() {
			a += 3
			return a
		}()
	}
	out = f1()
	`, nil, 8)

	expect(t, `
    out = func() {
        f1 := func() {
            a := 5
            add1 := func() { a += 1 }
            add2 := func() { a += 2 }
            a += 3
            return func() { a += 4; add1(); add2(); a += 5; return a }
        }
        return f1()
    }()()
    `, nil, 20)

	expect(t, `
		it := func(seq, fn) {
			fn(seq[0])
			fn(seq[1])
			fn(seq[2])
		}
	
		foo := func(a) {
			b := 0
			it([1, 2, 3], func(x) {
				b = x + a
			})
			return b
		}
	
		out = foo(2)
		`, nil, 5)

	expect(t, `
		it := func(seq, fn) {
			fn(seq[0])
			fn(seq[1])
			fn(seq[2])
		}
	
		foo := func(a) {
			b := 0
			it([1, 2, 3], func(x) {
				b += x + a
			})
			return b
		}
	
		out = foo(2)
		`, nil, 12)

	expect(t, `
out = func() {
	a := 1
	func() {
		a = 2
	}()
	return a
}()
`, nil, 2)

	expect(t, `
f := func() {
	a := 1
	return {
		b: func() { a += 3 },
		c: func() { a += 2 },
		d: func() { return a }
	}
}
m := f()
m.b()
m.c()
out = m.d()
`, nil, 6)

	expect(t, `
each := func(s, x) { for i:=0; i<len(s); i++ { x(s[i]) } }

out = func() {
	a := 100
	each([1, 2, 3], func(x) {
		a += x
	})
	a += 10
	return func(b) {
		return a + b
	}
}()(20)
`, nil, 136)

	// assigning different type value
	expect(t, `a := 1; a = "foo"; out = a`, nil, "foo")              // global
	expect(t, `func() { a := 1; a = "foo"; out = a }()`, nil, "foo") // local
	expect(t, `
out = func() { 
	a := 5
	return func() { 
		a = "foo"
		return a
	}()
}()`, nil, "foo") // free

	// variables declared in if/for blocks
	expect(t, `for a:=0; a<5; a++ {}; a := "foo"; out = a`, nil, "foo")
	expect(t, `func() { for a:=0; a<5; a++ {}; a := "foo"; out = a }()`, nil, "foo")

	// selectors
	expect(t, `a:=[1,2,3]; a[1] = 5; out = a[1]`, nil, 5)
	expect(t, `a:=[1,2,3]; a[1] += 5; out = a[1]`, nil, 7)
	expect(t, `a:={b:1,c:2}; a.b = 5; out = a.b`, nil, 5)
	expect(t, `a:={b:1,c:2}; a.b += 5; out = a.b`, nil, 6)
	expect(t, `a:={b:1,c:2}; a.b += a.c; out = a.b`, nil, 3)
	expect(t, `a:={b:1,c:2}; a.b += a.c; out = a.c`, nil, 2)
	expect(t, `
a := {
	b: [1, 2, 3],
	c: {
		d: 8,
		e: "foo",
		f: [9, 8]
	}
}
a.c.f[1] += 2
out = a["c"]["f"][1]
`, nil, 10)

	expect(t, `
a := {
	b: [1, 2, 3],
	c: {
		d: 8,
		e: "foo",
		f: [9, 8]
	}
}
a.c.h = "bar"
out = a.c.h
`, nil, "bar")

	expectError(t, `
a := {
	b: [1, 2, 3],
	c: {
		d: 8,
		e: "foo",
		f: [9, 8]
	}
}
a.x.e = "bar"`, nil, "not index-assignable")
}

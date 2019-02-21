package runtime_test

import (
	"testing"

	"github.com/d5/tengo/objects"
)

func TestFunction(t *testing.T) {
	// function with no "return" statement returns "invalid" value.
	expect(t, `f1 := func() {}; out = f1();`, objects.UndefinedValue)
	expect(t, `f1 := func() {}; f2 := func() { return f1(); }; f1(); out = f2();`, objects.UndefinedValue)
	expect(t, `f := func(x) { x; }; out = f(5);`, objects.UndefinedValue)

	expect(t, `f := func(x) { return x; }; out = f(5);`, 5)
	expect(t, `f := func(x) { return x * 2; }; out = f(5);`, 10)
	expect(t, `f := func(x, y) { return x + y; }; out = f(5, 5);`, 10)
	expect(t, `f := func(x, y) { return x + y; }; out = f(5 + 5, f(5, 5));`, 20)
	expect(t, `out = func(x) { return x; }(5)`, 5)
	expect(t, `x := 10; f := func(x) { return x; }; f(5); out = x;`, 10)

	expect(t, `
	f2 := func(a) {
		f1 := func(a) {
			return a * 2;
		};
	
		return f1(a) * 3;
	};
	
	out = f2(10);
	`, 60)

	// closures
	expect(t, `
		newAdder := func(x) {
			return func(y) { return x + y };
		};
	
		add2 := newAdder(2);
		out = add2(5);
		`, 7)

	// function as a argument
	expect(t, `
	add := func(a, b) { return a + b };
	sub := func(a, b) { return a - b };
	applyFunc := func(a, b, f) { return f(a, b) };
	
	out = applyFunc(applyFunc(2, 2, add), 3, sub);
	`, 1)

	expect(t, `f1 := func() { return 5 + 10; }; out = f1();`, 15)
	expect(t, `f1 := func() { return 1 }; f2 := func() { return 2 }; out = f1() + f2()`, 3)
	expect(t, `f1 := func() { return 1 }; f2 := func() { return f1() + 2 }; f3 := func() { return f2() + 3 }; out = f3()`, 6)
	expect(t, `f1 := func() { return 99; 100 }; out = f1();`, 99)
	expect(t, `f1 := func() { return 99; return 100 }; out = f1();`, 99)
	expect(t, `f1 := func() { return 33; }; f2 := func() { return f1 }; out = f2()();`, 33)
	expect(t, `one := func() { one = 1; return one }; out = one()`, 1)
	expect(t, `three := func() { one := 1; two := 2; return one + two }; out = three()`, 3)
	expect(t, `three := func() { one := 1; two := 2; return one + two }; seven := func() { three := 3; four := 4; return three + four }; out = three() + seven()`, 10)
	expect(t, `
	foo1 := func() {
		foo := 50
		return foo
	}
	foo2 := func() {
		foo := 100
		return foo
	}
	out = foo1() + foo2()`, 150)
	expect(t, `
	g := 50;
	minusOne := func() {
		n := 1;
		return g - n;
	};
	minusTwo := func() {
		n := 2;
		return g - n;
	};
	out = minusOne() + minusTwo()
	`, 97)
	expect(t, `
	f1 := func() {
		f2 := func() { return 1; }
		return f2
	};
	out = f1()()
	`, 1)

	expect(t, `
	f1 := func(a) { return a; };
	out = f1(4)`, 4)
	expect(t, `
	f1 := func(a, b) { return a + b; };
	out = f1(1, 2)`, 3)

	expect(t, `
	sum := func(a, b) {
		c := a + b;
		return c;
	};
	out = sum(1, 2);`, 3)

	expect(t, `
	sum := func(a, b) {
		c := a + b;
		return c;
	};
	out = sum(1, 2) + sum(3, 4);`, 10)

	expect(t, `
	sum := func(a, b) {
		c := a + b
		return c
	};
	outer := func() {
		return sum(1, 2) + sum(3, 4)
	};
	out = outer();`, 10)

	expect(t, `
	g := 10;
	
	sum := func(a, b) {
		c := a + b;
		return c + g;
	}
	
	outer := func() {
		return sum(1, 2) + sum(3, 4) + g;
	}
	
	out = outer() + g
	`, 50)

	expectError(t, `func() { return 1; }(1)`, "wrong number of arguments")
	expectError(t, `func(a) { return a; }()`, "wrong number of arguments")
	expectError(t, `func(a, b) { return a + b; }(1)`, "wrong number of arguments")

	expect(t, `
		f1 := func(a) {
			return func() { return a; };
		};
		f2 := f1(99);
		out = f2()
		`, 99)

	expect(t, `
		f1 := func(a, b) {
			return func(c) { return a + b + c };
		};
	
		f2 := f1(1, 2);
		out = f2(8);
		`, 11)
	expect(t, `
		f1 := func(a, b) {
			c := a + b;
			return func(d) { return c + d };
		};
		f2 := f1(1, 2);
		out = f2(8);
		`, 11)
	expect(t, `
		f1 := func(a, b) {
			c := a + b;
			return func(d) {
				e := d + c;
				return func(f) { return e + f };
			}
		};
		f2 := f1(1, 2);
		f3 := f2(3);
		out = f3(8);
		`, 14)
	expect(t, `
		a := 1;
		f1 := func(b) {
			return func(c) {
				return func(d) { return a + b + c + d }
			};
		};
		f2 := f1(2);
		f3 := f2(3);
		out = f3(8);
		`, 14)
	expect(t, `
		f1 := func(a, b) {
			one := func() { return a; };
			two := func() { return b; };
			return func() { return one() + two(); }
		};
		f2 := f1(9, 90);
		out = f2();
		`, 99)

	// global function recursion
	expect(t, `
		fib := func(x) {
			if x == 0 {
				return 0
			} else if x == 1 {
				return 1
			} else {
				return fib(x-1) + fib(x-2)
			}
		}
		out = fib(15)`, 610)

	// local function recursion
	expect(t, `
out = func() {
	sum := func(x) {
		return x == 0 ? 0 : x + sum(x-1)
	}
	return sum(5)
}()`, 15)

	expectError(t, `return 5`, "return not allowed outside function")

	// closure and block scopes
	expect(t, `
func() {
	a := 10
	func() {
		b := 5
		if true {
			out = a + 5
		}
	}()
}()`, 15)
	expect(t, `
func() {
	a := 10
	b := func() { return 5 }
	func() {
		if b() {
			out = a + b()
		}
	}()
}()`, 15)
	expect(t, `
func() {
	a := 10
	func() {
		b := func() { return 5 }
		func() {
			if true {
				out = a + b()
			}
		}()
	}()
}()`, 15)
}

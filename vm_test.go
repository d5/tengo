package tengo_test

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	_runtime "runtime"
	"strings"
	"testing"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/parser"
	"github.com/d5/tengo/v2/require"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/d5/tengo/v2/token"
)

const testOut = "out"

type IARR []interface{}
type IMAP map[string]interface{}
type MAP = map[string]interface{}
type ARR = []interface{}

type testopts struct {
	modules     *tengo.ModuleMap
	symbols     map[string]tengo.Object
	maxAllocs   int64
	skip2ndPass bool
}

func Opts() *testopts {
	return &testopts{
		modules:     tengo.NewModuleMap(),
		symbols:     make(map[string]tengo.Object),
		maxAllocs:   -1,
		skip2ndPass: false,
	}
}

func (o *testopts) copy() *testopts {
	c := &testopts{
		modules:     o.modules.Copy(),
		symbols:     make(map[string]tengo.Object),
		maxAllocs:   o.maxAllocs,
		skip2ndPass: o.skip2ndPass,
	}
	for k, v := range o.symbols {
		c.symbols[k] = v
	}
	return c
}

func (o *testopts) Stdlib() *testopts {
	o.modules.AddMap(stdlib.GetModuleMap(stdlib.AllModuleNames()...))
	return o
}

func (o *testopts) Module(name string, mod interface{}) *testopts {
	c := o.copy()
	switch mod := mod.(type) {
	case tengo.Importable:
		c.modules.Add(name, mod)
	case string:
		c.modules.AddSourceModule(name, []byte(mod))
	case []byte:
		c.modules.AddSourceModule(name, mod)
	default:
		panic(fmt.Errorf("invalid module type: %T", mod))
	}
	return c
}

func (o *testopts) Symbol(name string, value tengo.Object) *testopts {
	c := o.copy()
	c.symbols[name] = value
	return c
}

func (o *testopts) MaxAllocs(limit int64) *testopts {
	c := o.copy()
	c.maxAllocs = limit
	return c
}

func (o *testopts) Skip2ndPass() *testopts {
	c := o.copy()
	c.skip2ndPass = true
	return c
}

type customError struct {
	err error
	str string
}

func (c *customError) Error() string {
	return c.str
}

func (c *customError) Unwrap() error {
	return c.err
}

func TestArray(t *testing.T) {
	expectRun(t, `out = [1, 2 * 2, 3 + 3]`, nil, ARR{1, 4, 6})

	// array copy-by-reference
	expectRun(t, `a1 := [1, 2, 3]; a2 := a1; a1[0] = 5; out = a2`,
		nil, ARR{5, 2, 3})
	expectRun(t, `func () { a1 := [1, 2, 3]; a2 := a1; a1[0] = 5; out = a2 }()`,
		nil, ARR{5, 2, 3})

	// array index set
	expectError(t, `a1 := [1, 2, 3]; a1[3] = 5`,
		nil, "index out of bounds")

	// index operator
	arr := ARR{1, 2, 3, 4, 5, 6}
	arrStr := `[1, 2, 3, 4, 5, 6]`
	arrLen := 6
	for idx := 0; idx < arrLen; idx++ {
		expectRun(t, fmt.Sprintf("out = %s[%d]", arrStr, idx),
			nil, arr[idx])
		expectRun(t, fmt.Sprintf("out = %s[0 + %d]", arrStr, idx),
			nil, arr[idx])
		expectRun(t, fmt.Sprintf("out = %s[1 + %d - 1]", arrStr, idx),
			nil, arr[idx])
		expectRun(t, fmt.Sprintf("idx := %d; out = %s[idx]", idx, arrStr),
			nil, arr[idx])
	}

	expectRun(t, fmt.Sprintf("%s[%d]", arrStr, -1),
		nil, tengo.UndefinedValue)
	expectRun(t, fmt.Sprintf("%s[%d]", arrStr, arrLen),
		nil, tengo.UndefinedValue)

	// slice operator
	for low := 0; low < arrLen; low++ {
		expectRun(t, fmt.Sprintf("out = %s[%d:%d]", arrStr, low, low),
			nil, ARR{})
		for high := low; high <= arrLen; high++ {
			expectRun(t, fmt.Sprintf("out = %s[%d:%d]", arrStr, low, high),
				nil, arr[low:high])
			expectRun(t, fmt.Sprintf("out = %s[0 + %d : 0 + %d]",
				arrStr, low, high), nil, arr[low:high])
			expectRun(t, fmt.Sprintf("out = %s[1 + %d - 1 : 1 + %d - 1]",
				arrStr, low, high), nil, arr[low:high])
			expectRun(t, fmt.Sprintf("out = %s[:%d]", arrStr, high),
				nil, arr[:high])
			expectRun(t, fmt.Sprintf("out = %s[%d:]", arrStr, low),
				nil, arr[low:])
		}
	}

	expectRun(t, fmt.Sprintf("out = %s[:]", arrStr),
		nil, arr)
	expectRun(t, fmt.Sprintf("out = %s[%d:]", arrStr, -1),
		nil, arr)
	expectRun(t, fmt.Sprintf("out = %s[:%d]", arrStr, arrLen+1),
		nil, arr)
	expectRun(t, fmt.Sprintf("out = %s[%d:%d]", arrStr, 2, 2),
		nil, ARR{})

	expectError(t, fmt.Sprintf("%s[:%d]", arrStr, -1),
		nil, "invalid slice index")
	expectError(t, fmt.Sprintf("%s[%d:]", arrStr, arrLen+1),
		nil, "invalid slice index")
	expectError(t, fmt.Sprintf("%s[%d:%d]", arrStr, 0, -1),
		nil, "invalid slice index")
	expectError(t, fmt.Sprintf("%s[%d:%d]", arrStr, 2, 1),
		nil, "invalid slice index")
}

func TestAssignment(t *testing.T) {
	expectRun(t, `a := 1; a = 2; out = a`, nil, 2)
	expectRun(t, `a := 1; a = 2; out = a`, nil, 2)
	expectRun(t, `a := 1; a = a + 4; out = a`, nil, 5)
	expectRun(t, `a := 1; f1 := func() { a = 2; return a }; out = f1()`,
		nil, 2)
	expectRun(t, `a := 1; f1 := func() { a := 3; a = 2; return a }; out = f1()`,
		nil, 2)

	expectRun(t, `a := 1; out = a`, nil, 1)
	expectRun(t, `a := 1; a = 2; out = a`, nil, 2)
	expectRun(t, `a := 1; func() { a = 2 }(); out = a`, nil, 2)
	expectRun(t, `a := 1; func() { a := 2 }(); out = a`, nil, 1) // "a := 2" defines a new local variable 'a'
	expectRun(t, `a := 1; func() { b := 2; out = b }()`, nil, 2)
	expectRun(t, `
out = func() { 
	a := 2
	func() {
		a = 3 // captured from outer scope
	}()
	return a
}()
`, nil, 3)

	expectRun(t, `
func() {
	a := 5
	out = func() {  	
		a := 4						
		return a
	}()
}()`, nil, 4)

	expectError(t, `a := 1; a := 2`, nil, "redeclared")              // redeclared in the same scope
	expectError(t, `func() { a := 1; a := 2 }()`, nil, "redeclared") // redeclared in the same scope

	expectRun(t, `a := 1; a += 2; out = a`, nil, 3)
	expectRun(t, `a := 1; a += 4 - 2;; out = a`, nil, 3)
	expectRun(t, `a := 3; a -= 1;; out = a`, nil, 2)
	expectRun(t, `a := 3; a -= 5 - 4;; out = a`, nil, 2)
	expectRun(t, `a := 2; a *= 4;; out = a`, nil, 8)
	expectRun(t, `a := 2; a *= 1 + 3;; out = a`, nil, 8)
	expectRun(t, `a := 10; a /= 2;; out = a`, nil, 5)
	expectRun(t, `a := 10; a /= 5 - 3;; out = a`, nil, 5)

	// compound assignment operator does not define new variable
	expectError(t, `a += 4`, nil, "unresolved reference")
	expectError(t, `a -= 4`, nil, "unresolved reference")
	expectError(t, `a *= 4`, nil, "unresolved reference")
	expectError(t, `a /= 4`, nil, "unresolved reference")

	expectRun(t, `
f1 := func() { 
	f2 := func() { 
		a := 1
		a += 2    // it's a statement, not an expression
		return a
	}; 
	
	return f2(); 
}; 

out = f1();`, nil, 3)
	expectRun(t, `f1 := func() { f2 := func() { a := 1; a += 4 - 2; return a }; return f2(); }; out = f1()`,
		nil, 3)
	expectRun(t, `f1 := func() { f2 := func() { a := 3; a -= 1; return a }; return f2(); }; out = f1()`,
		nil, 2)
	expectRun(t, `f1 := func() { f2 := func() { a := 3; a -= 5 - 4; return a }; return f2(); }; out = f1()`,
		nil, 2)
	expectRun(t, `f1 := func() { f2 := func() { a := 2; a *= 4; return a }; return f2(); }; out = f1()`,
		nil, 8)
	expectRun(t, `f1 := func() { f2 := func() { a := 2; a *= 1 + 3; return a }; return f2(); }; out = f1()`,
		nil, 8)
	expectRun(t, `f1 := func() { f2 := func() { a := 10; a /= 2; return a }; return f2(); }; out = f1()`,
		nil, 5)
	expectRun(t, `f1 := func() { f2 := func() { a := 10; a /= 5 - 3; return a }; return f2(); }; out = f1()`,
		nil, 5)

	expectRun(t, `a := 1; f1 := func() { f2 := func() { a += 2; return a }; return f2(); }; out = f1()`,
		nil, 3)

	expectRun(t, `
	f1 := func(a) {
		return func(b) {
			c := a
			c += b * 2
			return c
		}
	}
	
	out = f1(3)(4)
	`, nil, 11)

	expectRun(t, `
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
	expectRun(t, `
	f1 := func() {
		a := 5
	
		return func() {
			a += 3
			return a
		}()
	}
	out = f1()
	`, nil, 8)

	expectRun(t, `
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

	expectRun(t, `
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

	expectRun(t, `
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

	expectRun(t, `
out = func() {
	a := 1
	func() {
		a = 2
	}()
	return a
}()
`, nil, 2)

	expectRun(t, `
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

	expectRun(t, `
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
	expectRun(t, `a := 1; a = "foo"; out = a`, nil, "foo")              // global
	expectRun(t, `func() { a := 1; a = "foo"; out = a }()`, nil, "foo") // local
	expectRun(t, `
out = func() { 
	a := 5
	return func() { 
		a = "foo"
		return a
	}()
}()`, nil, "foo") // free

	// variables declared in if/for blocks
	expectRun(t, `for a:=0; a<5; a++ {}; a := "foo"; out = a`,
		nil, "foo")
	expectRun(t, `func() { for a:=0; a<5; a++ {}; a := "foo"; out = a }()`,
		nil, "foo")

	// selectors
	expectRun(t, `a:=[1,2,3]; a[1] = 5; out = a[1]`, nil, 5)
	expectRun(t, `a:=[1,2,3]; a[1] += 5; out = a[1]`, nil, 7)
	expectRun(t, `a:={b:1,c:2}; a.b = 5; out = a.b`, nil, 5)
	expectRun(t, `a:={b:1,c:2}; a.b += 5; out = a.b`, nil, 6)
	expectRun(t, `a:={b:1,c:2}; a.b += a.c; out = a.b`, nil, 3)
	expectRun(t, `a:={b:1,c:2}; a.b += a.c; out = a.c`, nil, 2)
	expectRun(t, `
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

	expectRun(t, `
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

func TestBitwise(t *testing.T) {
	expectRun(t, `out = 1 & 1`, nil, 1)
	expectRun(t, `out = 1 & 0`, nil, 0)
	expectRun(t, `out = 0 & 1`, nil, 0)
	expectRun(t, `out = 0 & 0`, nil, 0)
	expectRun(t, `out = 1 | 1`, nil, 1)
	expectRun(t, `out = 1 | 0`, nil, 1)
	expectRun(t, `out = 0 | 1`, nil, 1)
	expectRun(t, `out = 0 | 0`, nil, 0)
	expectRun(t, `out = 1 ^ 1`, nil, 0)
	expectRun(t, `out = 1 ^ 0`, nil, 1)
	expectRun(t, `out = 0 ^ 1`, nil, 1)
	expectRun(t, `out = 0 ^ 0`, nil, 0)
	expectRun(t, `out = 1 &^ 1`, nil, 0)
	expectRun(t, `out = 1 &^ 0`, nil, 1)
	expectRun(t, `out = 0 &^ 1`, nil, 0)
	expectRun(t, `out = 0 &^ 0`, nil, 0)
	expectRun(t, `out = 1 << 2`, nil, 4)
	expectRun(t, `out = 16 >> 2`, nil, 4)

	expectRun(t, `out = 1; out &= 1`, nil, 1)
	expectRun(t, `out = 1; out |= 0`, nil, 1)
	expectRun(t, `out = 1; out ^= 0`, nil, 1)
	expectRun(t, `out = 1; out &^= 0`, nil, 1)
	expectRun(t, `out = 1; out <<= 2`, nil, 4)
	expectRun(t, `out = 16; out >>= 2`, nil, 4)

	expectRun(t, `out = ^0`, nil, ^0)
	expectRun(t, `out = ^1`, nil, ^1)
	expectRun(t, `out = ^55`, nil, ^55)
	expectRun(t, `out = ^-55`, nil, ^-55)
}

func TestBoolean(t *testing.T) {
	expectRun(t, `out = true`, nil, true)
	expectRun(t, `out = false`, nil, false)

	expectRun(t, `out = 1 < 2`, nil, true)
	expectRun(t, `out = 1 > 2`, nil, false)
	expectRun(t, `out = 1 < 1`, nil, false)
	expectRun(t, `out = 1 > 2`, nil, false)
	expectRun(t, `out = 1 == 1`, nil, true)
	expectRun(t, `out = 1 != 1`, nil, false)
	expectRun(t, `out = 1 == 2`, nil, false)
	expectRun(t, `out = 1 != 2`, nil, true)
	expectRun(t, `out = 1 <= 2`, nil, true)
	expectRun(t, `out = 1 >= 2`, nil, false)
	expectRun(t, `out = 1 <= 1`, nil, true)
	expectRun(t, `out = 1 >= 2`, nil, false)

	expectRun(t, `out = true == true`, nil, true)
	expectRun(t, `out = false == false`, nil, true)
	expectRun(t, `out = true == false`, nil, false)
	expectRun(t, `out = true != false`, nil, true)
	expectRun(t, `out = false != true`, nil, true)
	expectRun(t, `out = (1 < 2) == true`, nil, true)
	expectRun(t, `out = (1 < 2) == false`, nil, false)
	expectRun(t, `out = (1 > 2) == true`, nil, false)
	expectRun(t, `out = (1 > 2) == false`, nil, true)

	expectError(t, `5 + true`, nil, "invalid operation")
	expectError(t, `5 + true; 5`, nil, "invalid operation")
	expectError(t, `-true`, nil, "invalid operation")
	expectError(t, `true + false`, nil, "invalid operation")
	expectError(t, `5; true + false; 5`, nil, "invalid operation")
	expectError(t, `if (10 > 1) { true + false; }`, nil, "invalid operation")
	expectError(t, `
func() {
	if (10 > 1) {
		if (10 > 1) {
			return true + false;
		}

		return 1;
	}
}()
`, nil, "invalid operation")
	expectError(t, `if (true + false) { 10 }`, nil, "invalid operation")
	expectError(t, `10 + (true + false)`, nil, "invalid operation")
	expectError(t, `(true + false) + 20`, nil, "invalid operation")
	expectError(t, `!(true + false)`, nil, "invalid operation")
}

func TestUndefined(t *testing.T) {
	expectRun(t, `out = undefined`, nil, tengo.UndefinedValue)
	expectRun(t, `out = undefined.a`, nil, tengo.UndefinedValue)
	expectRun(t, `out = undefined[1]`, nil, tengo.UndefinedValue)
	expectRun(t, `out = undefined.a.b`, nil, tengo.UndefinedValue)
	expectRun(t, `out = undefined[1][2]`, nil, tengo.UndefinedValue)
	expectRun(t, `out = undefined ? 1 : 2`, nil, 2)
	expectRun(t, `out = undefined == undefined`, nil, true)
	expectRun(t, `out = undefined == 1`, nil, false)
	expectRun(t, `out = 1 == undefined`, nil, false)
	expectRun(t, `out = undefined == float([])`, nil, true)
	expectRun(t, `out = float([]) == undefined`, nil, true)
}

func TestBuiltinFunction(t *testing.T) {
	expectRun(t, `out = len("")`, nil, 0)
	expectRun(t, `out = len("four")`, nil, 4)
	expectRun(t, `out = len("hello world")`, nil, 11)
	expectRun(t, `out = len([])`, nil, 0)
	expectRun(t, `out = len([1, 2, 3])`, nil, 3)
	expectRun(t, `out = len({})`, nil, 0)
	expectRun(t, `out = len({a:1, b:2})`, nil, 2)
	expectRun(t, `out = len(immutable([]))`, nil, 0)
	expectRun(t, `out = len(immutable([1, 2, 3]))`, nil, 3)
	expectRun(t, `out = len(immutable({}))`, nil, 0)
	expectRun(t, `out = len(immutable({a:1, b:2}))`, nil, 2)
	expectError(t, `len(1)`, nil, "invalid type for argument")
	expectError(t, `len("one", "two")`, nil, "wrong number of arguments")

	expectRun(t, `out = copy(1)`, nil, 1)
	expectError(t, `copy(1, 2)`, nil, "wrong number of arguments")

	expectRun(t, `out = append([1, 2, 3], 4)`, nil, ARR{1, 2, 3, 4})
	expectRun(t, `out = append([1, 2, 3], 4, 5, 6)`, nil, ARR{1, 2, 3, 4, 5, 6})
	expectRun(t, `out = append([1, 2, 3], "foo", false)`,
		nil, ARR{1, 2, 3, "foo", false})

	expectRun(t, `out = int(1)`, nil, 1)
	expectRun(t, `out = int(1.8)`, nil, 1)
	expectRun(t, `out = int("-522")`, nil, -522)
	expectRun(t, `out = int(true)`, nil, 1)
	expectRun(t, `out = int(false)`, nil, 0)
	expectRun(t, `out = int('8')`, nil, 56)
	expectRun(t, `out = int([1])`, nil, tengo.UndefinedValue)
	expectRun(t, `out = int({a: 1})`, nil, tengo.UndefinedValue)
	expectRun(t, `out = int(undefined)`, nil, tengo.UndefinedValue)
	expectRun(t, `out = int("-522", 1)`, nil, -522)
	expectRun(t, `out = int(undefined, 1)`, nil, 1)
	expectRun(t, `out = int(undefined, 1.8)`, nil, 1.8)
	expectRun(t, `out = int(undefined, string(1))`, nil, "1")
	expectRun(t, `out = int(undefined, undefined)`, nil, tengo.UndefinedValue)

	expectRun(t, `out = string(1)`, nil, "1")
	expectRun(t, `out = string(1.8)`, nil, "1.8")
	expectRun(t, `out = string("-522")`, nil, "-522")
	expectRun(t, `out = string(true)`, nil, "true")
	expectRun(t, `out = string(false)`, nil, "false")
	expectRun(t, `out = string('8')`, nil, "8")
	expectRun(t, `out = string([1,8.1,true,3])`, nil, "[1, 8.1, true, 3]")
	expectRun(t, `out = string({b: "foo"})`, nil, `{b: "foo"}`)
	expectRun(t, `out = string(undefined)`, nil, tengo.UndefinedValue) // not "undefined"
	expectRun(t, `out = string(1, "-522")`, nil, "1")
	expectRun(t, `out = string(undefined, "-522")`, nil, "-522") // not "undefined"

	expectRun(t, `out = float(1)`, nil, 1.0)
	expectRun(t, `out = float(1.8)`, nil, 1.8)
	expectRun(t, `out = float("-52.2")`, nil, -52.2)
	expectRun(t, `out = float(true)`, nil, tengo.UndefinedValue)
	expectRun(t, `out = float(false)`, nil, tengo.UndefinedValue)
	expectRun(t, `out = float('8')`, nil, tengo.UndefinedValue)
	expectRun(t, `out = float([1,8.1,true,3])`, nil, tengo.UndefinedValue)
	expectRun(t, `out = float({a: 1, b: "foo"})`, nil, tengo.UndefinedValue)
	expectRun(t, `out = float(undefined)`, nil, tengo.UndefinedValue)
	expectRun(t, `out = float("-52.2", 1.8)`, nil, -52.2)
	expectRun(t, `out = float(undefined, 1)`, nil, 1)
	expectRun(t, `out = float(undefined, 1.8)`, nil, 1.8)
	expectRun(t, `out = float(undefined, "-52.2")`, nil, "-52.2")
	expectRun(t, `out = float(undefined, char(56))`, nil, '8')
	expectRun(t, `out = float(undefined, undefined)`, nil, tengo.UndefinedValue)

	expectRun(t, `out = char(56)`, nil, '8')
	expectRun(t, `out = char(1.8)`, nil, tengo.UndefinedValue)
	expectRun(t, `out = char("-52.2")`, nil, tengo.UndefinedValue)
	expectRun(t, `out = char(true)`, nil, tengo.UndefinedValue)
	expectRun(t, `out = char(false)`, nil, tengo.UndefinedValue)
	expectRun(t, `out = char('8')`, nil, '8')
	expectRun(t, `out = char([1,8.1,true,3])`, nil, tengo.UndefinedValue)
	expectRun(t, `out = char({a: 1, b: "foo"})`, nil, tengo.UndefinedValue)
	expectRun(t, `out = char(undefined)`, nil, tengo.UndefinedValue)
	expectRun(t, `out = char(56, 'a')`, nil, '8')
	expectRun(t, `out = char(undefined, '8')`, nil, '8')
	expectRun(t, `out = char(undefined, 56)`, nil, 56)
	expectRun(t, `out = char(undefined, "-52.2")`, nil, "-52.2")
	expectRun(t, `out = char(undefined, undefined)`, nil, tengo.UndefinedValue)

	expectRun(t, `out = bool(1)`, nil, true)          // non-zero integer: true
	expectRun(t, `out = bool(0)`, nil, false)         // zero: true
	expectRun(t, `out = bool(1.8)`, nil, true)        // all floats (except for NaN): true
	expectRun(t, `out = bool(0.0)`, nil, true)        // all floats (except for NaN): true
	expectRun(t, `out = bool("false")`, nil, true)    // non-empty string: true
	expectRun(t, `out = bool("")`, nil, false)        // empty string: false
	expectRun(t, `out = bool(true)`, nil, true)       // true: true
	expectRun(t, `out = bool(false)`, nil, false)     // false: false
	expectRun(t, `out = bool('8')`, nil, true)        // non-zero chars: true
	expectRun(t, `out = bool(char(0))`, nil, false)   // zero char: false
	expectRun(t, `out = bool([1])`, nil, true)        // non-empty arrays: true
	expectRun(t, `out = bool([])`, nil, false)        // empty array: false
	expectRun(t, `out = bool({a: 1})`, nil, true)     // non-empty maps: true
	expectRun(t, `out = bool({})`, nil, false)        // empty maps: false
	expectRun(t, `out = bool(undefined)`, nil, false) // undefined: false

	expectRun(t, `out = bytes(1)`, nil, []byte{0})
	expectRun(t, `out = bytes(1.8)`, nil, tengo.UndefinedValue)
	expectRun(t, `out = bytes("-522")`, nil, []byte{'-', '5', '2', '2'})
	expectRun(t, `out = bytes(true)`, nil, tengo.UndefinedValue)
	expectRun(t, `out = bytes(false)`, nil, tengo.UndefinedValue)
	expectRun(t, `out = bytes('8')`, nil, tengo.UndefinedValue)
	expectRun(t, `out = bytes([1])`, nil, tengo.UndefinedValue)
	expectRun(t, `out = bytes({a: 1})`, nil, tengo.UndefinedValue)
	expectRun(t, `out = bytes(undefined)`, nil, tengo.UndefinedValue)
	expectRun(t, `out = bytes("-522", ['8'])`, nil, []byte{'-', '5', '2', '2'})
	expectRun(t, `out = bytes(undefined, "-522")`, nil, "-522")
	expectRun(t, `out = bytes(undefined, 1)`, nil, 1)
	expectRun(t, `out = bytes(undefined, 1.8)`, nil, 1.8)
	expectRun(t, `out = bytes(undefined, int("-522"))`, nil, -522)
	expectRun(t, `out = bytes(undefined, undefined)`, nil, tengo.UndefinedValue)

	expectRun(t, `out = is_error(error(1))`, nil, true)
	expectRun(t, `out = is_error(1)`, nil, false)

	expectRun(t, `out = is_undefined(undefined)`, nil, true)
	expectRun(t, `out = is_undefined(error(1))`, nil, false)

	// type_name
	expectRun(t, `out = type_name(1)`, nil, "int")
	expectRun(t, `out = type_name(1.1)`, nil, "float")
	expectRun(t, `out = type_name("a")`, nil, "string")
	expectRun(t, `out = type_name([1,2,3])`, nil, "array")
	expectRun(t, `out = type_name({k:1})`, nil, "map")
	expectRun(t, `out = type_name('a')`, nil, "char")
	expectRun(t, `out = type_name(true)`, nil, "bool")
	expectRun(t, `out = type_name(false)`, nil, "bool")
	expectRun(t, `out = type_name(bytes( 1))`, nil, "bytes")
	expectRun(t, `out = type_name(undefined)`, nil, "undefined")
	expectRun(t, `out = type_name(error("err"))`, nil, "error")
	expectRun(t, `out = type_name(func() {})`, nil, "compiled-function")
	expectRun(t, `a := func(x) { return func() { return x } }; out = type_name(a(5))`,
		nil, "compiled-function") // closure

	// is_function
	expectRun(t, `out = is_function(1)`, nil, false)
	expectRun(t, `out = is_function(func() {})`, nil, true)
	expectRun(t, `out = is_function(func(x) { return x })`, nil, true)
	expectRun(t, `out = is_function(len)`, nil, false) // builtin function
	expectRun(t, `a := func(x) { return func() { return x } }; out = is_function(a)`,
		nil, true) // function
	expectRun(t, `a := func(x) { return func() { return x } }; out = is_function(a(5))`,
		nil, true) // closure
	expectRun(t, `out = is_function(x)`,
		Opts().Symbol("x", &StringArray{
			Value: []string{"foo", "bar"},
		}).Skip2ndPass(),
		false) // user object

	// is_callable
	expectRun(t, `out = is_callable(1)`, nil, false)
	expectRun(t, `out = is_callable(func() {})`, nil, true)
	expectRun(t, `out = is_callable(func(x) { return x })`, nil, true)
	expectRun(t, `out = is_callable(len)`, nil, true) // builtin function
	expectRun(t, `a := func(x) { return func() { return x } }; out = is_callable(a)`,
		nil, true) // function
	expectRun(t, `a := func(x) { return func() { return x } }; out = is_callable(a(5))`,
		nil, true) // closure
	expectRun(t, `out = is_callable(x)`,
		Opts().Symbol("x", &StringArray{
			Value: []string{"foo", "bar"},
		}).Skip2ndPass(), true) // user object

	expectRun(t, `out = format("")`, nil, "")
	expectRun(t, `out = format("foo")`, nil, "foo")
	expectRun(t, `out = format("foo %d %v %s", 1, 2, "bar")`,
		nil, "foo 1 2 bar")
	expectRun(t, `out = format("foo %v", [1, "bar", true])`,
		nil, `foo [1, "bar", true]`)
	expectRun(t, `out = format("foo %v %d", [1, "bar", true], 19)`,
		nil, `foo [1, "bar", true] 19`)
	expectRun(t, `out = format("foo %v", {"a": {"b": {"c": [1, 2, 3]}}})`,
		nil, `foo {a: {b: {c: [1, 2, 3]}}}`)
	expectRun(t, `out = format("%v", [1, [2, [3, 4]]])`,
		nil, `[1, [2, [3, 4]]]`)

	tengo.MaxStringLen = 9
	expectError(t, `format("%s", "1234567890")`,
		nil, "exceeding string size limit")
	tengo.MaxStringLen = 2147483647

	// delete
	expectError(t, `delete()`, nil, tengo.ErrWrongNumArguments.Error())
	expectError(t, `delete(1)`, nil, tengo.ErrWrongNumArguments.Error())
	expectError(t, `delete(1, 2, 3)`, nil, tengo.ErrWrongNumArguments.Error())
	expectError(t, `delete({}, "", 3)`, nil, tengo.ErrWrongNumArguments.Error())
	expectError(t, `delete(1, 1)`, nil, `invalid type for argument 'first'`)
	expectError(t, `delete(1.0, 1)`, nil, `invalid type for argument 'first'`)
	expectError(t, `delete("str", 1)`, nil, `invalid type for argument 'first'`)
	expectError(t, `delete(bytes("str"), 1)`, nil,
		`invalid type for argument 'first'`)
	expectError(t, `delete(error("err"), 1)`, nil,
		`invalid type for argument 'first'`)
	expectError(t, `delete(true, 1)`, nil, `invalid type for argument 'first'`)
	expectError(t, `delete(char('c'), 1)`, nil,
		`invalid type for argument 'first'`)
	expectError(t, `delete(undefined, 1)`, nil,
		`invalid type for argument 'first'`)
	expectError(t, `delete(time(1257894000), 1)`, nil,
		`invalid type for argument 'first'`)
	expectError(t, `delete(immutable({}), "key")`, nil,
		`invalid type for argument 'first'`)
	expectError(t, `delete(immutable([]), "")`, nil,
		`invalid type for argument 'first'`)
	expectError(t, `delete([], "")`, nil, `invalid type for argument 'first'`)
	expectError(t, `delete({}, 1)`, nil, `invalid type for argument 'second'`)
	expectError(t, `delete({}, 1.0)`, nil, `invalid type for argument 'second'`)
	expectError(t, `delete({}, undefined)`, nil,
		`invalid type for argument 'second'`)
	expectError(t, `delete({}, [])`, nil, `invalid type for argument 'second'`)
	expectError(t, `delete({}, {})`, nil, `invalid type for argument 'second'`)
	expectError(t, `delete({}, error("err"))`, nil,
		`invalid type for argument 'second'`)
	expectError(t, `delete({}, bytes("str"))`, nil,
		`invalid type for argument 'second'`)
	expectError(t, `delete({}, char(35))`, nil,
		`invalid type for argument 'second'`)
	expectError(t, `delete({}, time(1257894000))`, nil,
		`invalid type for argument 'second'`)
	expectError(t, `delete({}, immutable({}))`, nil,
		`invalid type for argument 'second'`)
	expectError(t, `delete({}, immutable([]))`, nil,
		`invalid type for argument 'second'`)

	expectRun(t, `out = delete({}, "")`, nil, tengo.UndefinedValue)
	expectRun(t, `out = {key1: 1}; delete(out, "key1")`, nil, MAP{})
	expectRun(t, `out = {key1: 1, key2: "2"}; delete(out, "key1")`, nil,
		MAP{"key2": "2"})
	expectRun(t, `out = [1, "2", {a: "b", c: 10}]; delete(out[2], "c")`, nil,
		ARR{1, "2", MAP{"a": "b"}})

	// splice
	expectError(t, `splice()`, nil, tengo.ErrWrongNumArguments.Error())
	expectError(t, `splice(1)`, nil, `invalid type for argument 'first'`)
	expectError(t, `splice(1.0)`, nil, `invalid type for argument 'first'`)
	expectError(t, `splice("str")`, nil, `invalid type for argument 'first'`)
	expectError(t, `splice(bytes("str"))`, nil,
		`invalid type for argument 'first'`)
	expectError(t, `splice(error("err"))`, nil,
		`invalid type for argument 'first'`)
	expectError(t, `splice(true)`, nil, `invalid type for argument 'first'`)
	expectError(t, `splice(char('c'))`, nil,
		`invalid type for argument 'first'`)
	expectError(t, `splice(undefined)`, nil,
		`invalid type for argument 'first'`)
	expectError(t, `splice(time(1257894000))`, nil,
		`invalid type for argument 'first'`)
	expectError(t, `splice(immutable({}))`, nil,
		`invalid type for argument 'first'`)
	expectError(t, `splice(immutable([]))`, nil,
		`invalid type for argument 'first'`)
	expectError(t, `splice({})`, nil, `invalid type for argument 'first'`)
	expectError(t, `splice([], 1.0)`, nil,
		`invalid type for argument 'second'`)
	expectError(t, `splice([], "str")`, nil,
		`invalid type for argument 'second'`)
	expectError(t, `splice([], bytes("str"))`, nil,
		`invalid type for argument 'second'`)
	expectError(t, `splice([], error("error"))`, nil,
		`invalid type for argument 'second'`)
	expectError(t, `splice([], false)`, nil,
		`invalid type for argument 'second'`)
	expectError(t, `splice([], char('d'))`, nil,
		`invalid type for argument 'second'`)
	expectError(t, `splice([], undefined)`, nil,
		`invalid type for argument 'second'`)
	expectError(t, `splice([], time(0))`, nil,
		`invalid type for argument 'second'`)
	expectError(t, `splice([], [])`, nil,
		`invalid type for argument 'second'`)
	expectError(t, `splice([], {})`, nil,
		`invalid type for argument 'second'`)
	expectError(t, `splice([], immutable([]))`, nil,
		`invalid type for argument 'second'`)
	expectError(t, `splice([], immutable({}))`, nil,
		`invalid type for argument 'second'`)
	expectError(t, `splice([], 0, 1.0)`, nil,
		`invalid type for argument 'third'`)
	expectError(t, `splice([], 0, "string")`, nil,
		`invalid type for argument 'third'`)
	expectError(t, `splice([], 0, bytes("string"))`, nil,
		`invalid type for argument 'third'`)
	expectError(t, `splice([], 0, error("string"))`, nil,
		`invalid type for argument 'third'`)
	expectError(t, `splice([], 0, true)`, nil,
		`invalid type for argument 'third'`)
	expectError(t, `splice([], 0, char('f'))`, nil,
		`invalid type for argument 'third'`)
	expectError(t, `splice([], 0, undefined)`, nil,
		`invalid type for argument 'third'`)
	expectError(t, `splice([], 0, time(0))`, nil,
		`invalid type for argument 'third'`)
	expectError(t, `splice([], 0, [])`, nil,
		`invalid type for argument 'third'`)
	expectError(t, `splice([], 0, {})`, nil,
		`invalid type for argument 'third'`)
	expectError(t, `splice([], 0, immutable([]))`, nil,
		`invalid type for argument 'third'`)
	expectError(t, `splice([], 0, immutable({}))`, nil,
		`invalid type for argument 'third'`)
	expectError(t, `splice([], 1)`, nil, tengo.ErrIndexOutOfBounds.Error())
	expectError(t, `splice([1, 2, 3], 0, -1)`, nil,
		tengo.ErrIndexOutOfBounds.Error())
	expectError(t, `splice([1, 2, 3], 99, 0, "a", "b")`, nil,
		tengo.ErrIndexOutOfBounds.Error())
	expectRun(t, `out = []; splice(out)`, nil, ARR{})
	expectRun(t, `out = ["a"]; splice(out, 1)`, nil, ARR{"a"})
	expectRun(t, `out = ["a"]; out = splice(out, 1)`, nil, ARR{})
	expectRun(t, `out = [1, 2, 3]; splice(out, 0, 1)`, nil, ARR{2, 3})
	expectRun(t, `out = [1, 2, 3]; out = splice(out, 0, 1)`, nil, ARR{1})
	expectRun(t, `out = [1, 2, 3]; splice(out, 0, 0, "a", "b")`, nil,
		ARR{"a", "b", 1, 2, 3})
	expectRun(t, `out = [1, 2, 3]; out = splice(out, 0, 0, "a", "b")`, nil,
		ARR{})
	expectRun(t, `out = [1, 2, 3]; splice(out, 1, 0, "a", "b")`, nil,
		ARR{1, "a", "b", 2, 3})
	expectRun(t, `out = [1, 2, 3]; out = splice(out, 1, 0, "a", "b")`, nil,
		ARR{})
	expectRun(t, `out = [1, 2, 3]; splice(out, 1, 0, "a", "b")`, nil,
		ARR{1, "a", "b", 2, 3})
	expectRun(t, `out = [1, 2, 3]; splice(out, 2, 0, "a", "b")`, nil,
		ARR{1, 2, "a", "b", 3})
	expectRun(t, `out = [1, 2, 3]; splice(out, 3, 0, "a", "b")`, nil,
		ARR{1, 2, 3, "a", "b"})
	expectRun(t, `array := [1, 2, 3]; deleted := splice(array, 1, 1, "a", "b");
				out = [deleted, array]`, nil, ARR{ARR{2}, ARR{1, "a", "b", 3}})
	expectRun(t, `array := [1, 2, 3]; deleted := splice(array, 1); 
		out = [deleted, array]`, nil, ARR{ARR{2, 3}, ARR{1}})
	expectRun(t, `out = []; splice(out, 0, 0, "a", "b")`, nil, ARR{"a", "b"})
	expectRun(t, `out = []; splice(out, 0, 1, "a", "b")`, nil, ARR{"a", "b"})
	expectRun(t, `out = []; out = splice(out, 0, 0, "a", "b")`, nil, ARR{})
	expectRun(t, `out = splice(splice([1, 2, 3], 0, 3), 1, 3)`, nil, ARR{2, 3})
	// splice doc examples
	expectRun(t, `v := [1, 2, 3]; deleted := splice(v, 0);
		out = [deleted, v]`, nil, ARR{ARR{1, 2, 3}, ARR{}})
	expectRun(t, `v := [1, 2, 3]; deleted := splice(v, 1);
		out = [deleted, v]`, nil, ARR{ARR{2, 3}, ARR{1}})
	expectRun(t, `v := [1, 2, 3]; deleted := splice(v, 0, 1);
		out = [deleted, v]`, nil, ARR{ARR{1}, ARR{2, 3}})
	expectRun(t, `v := ["a", "b", "c"]; deleted := splice(v, 1, 2);
		out = [deleted, v]`, nil, ARR{ARR{"b", "c"}, ARR{"a"}})
	expectRun(t, `v := ["a", "b", "c"]; deleted := splice(v, 2, 1, "d");
		out = [deleted, v]`, nil, ARR{ARR{"c"}, ARR{"a", "b", "d"}})
	expectRun(t, `v := ["a", "b", "c"]; deleted := splice(v, 0, 0, "d", "e");
		out = [deleted, v]`, nil, ARR{ARR{}, ARR{"d", "e", "a", "b", "c"}})
	expectRun(t, `v := ["a", "b", "c"]; deleted := splice(v, 1, 1, "d", "e");
		out = [deleted, v]`, nil, ARR{ARR{"b"}, ARR{"a", "d", "e", "c"}})
}

func TestBytesN(t *testing.T) {
	curMaxBytesLen := tengo.MaxBytesLen
	defer func() { tengo.MaxBytesLen = curMaxBytesLen }()
	tengo.MaxBytesLen = 10

	expectRun(t, `out = bytes(0)`, nil, make([]byte, 0))
	expectRun(t, `out = bytes(10)`, nil, make([]byte, 10))
	expectError(t, `bytes(11)`, nil, "bytes size limit")

	tengo.MaxBytesLen = 1000
	expectRun(t, `out = bytes(1000)`, nil, make([]byte, 1000))
	expectError(t, `bytes(1001)`, nil, "bytes size limit")
}

func TestBytes(t *testing.T) {
	expectRun(t, `out = bytes("Hello World!")`, nil, []byte("Hello World!"))
	expectRun(t, `out = bytes("Hello") + bytes(" ") + bytes("World!")`,
		nil, []byte("Hello World!"))

	// bytes[] -> int
	expectRun(t, `out = bytes("abcde")[0]`, nil, 97)
	expectRun(t, `out = bytes("abcde")[1]`, nil, 98)
	expectRun(t, `out = bytes("abcde")[4]`, nil, 101)
	expectRun(t, `out = bytes("abcde")[10]`, nil, tengo.UndefinedValue)
}

func TestCall(t *testing.T) {
	expectRun(t, `a := { b: func(x) { return x + 2 } }; out = a.b(5)`,
		nil, 7)
	expectRun(t, `a := { b: { c: func(x) { return x + 2 } } }; out = a.b.c(5)`,
		nil, 7)
	expectRun(t, `a := { b: { c: func(x) { return x + 2 } } }; out = a["b"].c(5)`,
		nil, 7)
	expectError(t, `a := 1
b := func(a, c) {
   c(a)
}

c := func(a) {
   a()
}
b(a, c)
`, nil, "Runtime Error: not callable: int\n\tat test:7:4\n\tat test:3:4\n\tat test:9:1")
}

func TestChar(t *testing.T) {
	expectRun(t, `out = 'a'`, nil, 'a')
	expectRun(t, `out = '九'`, nil, rune(20061))
	expectRun(t, `out = 'Æ'`, nil, rune(198))

	expectRun(t, `out = '0' + '9'`, nil, rune(105))
	expectRun(t, `out = '0' + 9`, nil, '9')
	expectRun(t, `out = '9' - 4`, nil, '5')
	expectRun(t, `out = '0' == '0'`, nil, true)
	expectRun(t, `out = '0' != '0'`, nil, false)
	expectRun(t, `out = '2' < '4'`, nil, true)
	expectRun(t, `out = '2' > '4'`, nil, false)
	expectRun(t, `out = '2' <= '4'`, nil, true)
	expectRun(t, `out = '2' >= '4'`, nil, false)
	expectRun(t, `out = '4' < '4'`, nil, false)
	expectRun(t, `out = '4' > '4'`, nil, false)
	expectRun(t, `out = '4' <= '4'`, nil, true)
	expectRun(t, `out = '4' >= '4'`, nil, true)
}

func TestCondExpr(t *testing.T) {
	expectRun(t, `out = true ? 5 : 10`, nil, 5)
	expectRun(t, `out = false ? 5 : 10`, nil, 10)
	expectRun(t, `out = (1 == 1) ? 2 + 3 : 12 - 2`, nil, 5)
	expectRun(t, `out = (1 != 1) ? 2 + 3 : 12 - 2`, nil, 10)
	expectRun(t, `out = (1 == 1) ? true ? 10 - 8 : 1 + 3 : 12 - 2`, nil, 2)
	expectRun(t, `out = (1 == 1) ? false ? 10 - 8 : 1 + 3 : 12 - 2`, nil, 4)

	expectRun(t, `
out = 0
f1 := func() { out += 10 }
f2 := func() { out = -out }
true ? f1() : f2()
`, nil, 10)
	expectRun(t, `
out = 5
f1 := func() { out += 10 }
f2 := func() { out = -out }
false ? f1() : f2()
`, nil, -5)
	expectRun(t, `
f1 := func(a) { return a + 2 }
f2 := func(a) { return a - 2 }
f3 := func(a) { return a + 10 }
f4 := func(a) { return -a }

f := func(c) {
	return c == 0 ? f1(c) : f2(c) ? f3(c) : f4(c)
}

out = [f(0), f(1), f(2)]
`, nil, ARR{2, 11, -2})

	expectRun(t, `f := func(a) { return -a }; out = f(true ? 5 : 3)`, nil, -5)
	expectRun(t, `out = [false?5:10, true?1:2]`, nil, ARR{10, 1})

	expectRun(t, `
out = 1 > 2 ?
	1 + 2 + 3 :
	10 - 5`, nil, 5)
}

func TestEquality(t *testing.T) {
	testEquality(t, `1`, `1`, true)
	testEquality(t, `1`, `2`, false)

	testEquality(t, `1.0`, `1.0`, true)
	testEquality(t, `1.0`, `1.1`, false)

	testEquality(t, `true`, `true`, true)
	testEquality(t, `true`, `false`, false)

	testEquality(t, `"foo"`, `"foo"`, true)
	testEquality(t, `"foo"`, `"bar"`, false)

	testEquality(t, `'f'`, `'f'`, true)
	testEquality(t, `'f'`, `'b'`, false)

	testEquality(t, `[]`, `[]`, true)
	testEquality(t, `[1]`, `[1]`, true)
	testEquality(t, `[1]`, `[1, 2]`, false)
	testEquality(t, `["foo", "bar"]`, `["foo", "bar"]`, true)
	testEquality(t, `["foo", "bar"]`, `["bar", "foo"]`, false)

	testEquality(t, `{}`, `{}`, true)
	testEquality(t, `{a: 1, b: 2}`, `{b: 2, a: 1}`, true)
	testEquality(t, `{a: 1, b: 2}`, `{b: 2}`, false)
	testEquality(t, `{a: 1, b: {}}`, `{b: {}, a: 1}`, true)

	testEquality(t, `1`, `"foo"`, false)
	testEquality(t, `1`, `true`, false)
	testEquality(t, `[1]`, `["1"]`, false)
	testEquality(t, `[1, [2]]`, `[1, ["2"]]`, false)
	testEquality(t, `{a: 1}`, `{a: "1"}`, false)
	testEquality(t, `{a: 1, b: {c: 2}}`, `{a: 1, b: {c: "2"}}`, false)
}

func testEquality(t *testing.T, lhs, rhs string, expected bool) {
	// 1. equality is commutative
	// 2. equality and inequality must be always opposite
	expectRun(t, fmt.Sprintf("out = %s == %s", lhs, rhs), nil, expected)
	expectRun(t, fmt.Sprintf("out = %s == %s", rhs, lhs), nil, expected)
	expectRun(t, fmt.Sprintf("out = %s != %s", lhs, rhs), nil, !expected)
	expectRun(t, fmt.Sprintf("out = %s != %s", rhs, lhs), nil, !expected)
}

func TestVMErrorInfo(t *testing.T) {
	expectError(t, `a := 5
a + "boo"`,
		nil, "Runtime Error: invalid operation: int + string\n\tat test:2:1")

	expectError(t, `a := 5
b := a(5)`,
		nil, "Runtime Error: not callable: int\n\tat test:2:6")

	expectError(t, `a := 5
b := {}
b.x.y = 10`,
		nil, "Runtime Error: not index-assignable: undefined\n\tat test:3:1")

	expectError(t, `
a := func() {
	b := 5
	b += "foo"
}
a()`,
		nil, "Runtime Error: invalid operation: int + string\n\tat test:4:2")

	expectError(t, `a := 5
a + import("mod1")`, Opts().Module(
		"mod1", `export "foo"`,
	), ": invalid operation: int + string\n\tat test:2:1")

	expectError(t, `a := import("mod1")()`,
		Opts().Module(
			"mod1", `
export func() {
	b := 5
	return b + "foo"
}`), "Runtime Error: invalid operation: int + string\n\tat mod1:4:9")

	expectError(t, `a := import("mod1")()`,
		Opts().Module(
			"mod1", `export import("mod2")()`).
			Module(
				"mod2", `
export func() {
	b := 5
	return b + "foo"
}`), "Runtime Error: invalid operation: int + string\n\tat mod2:4:9")
}

func TestVMErrorUnwrap(t *testing.T) {
	userErr := errors.New("user runtime error")
	userFunc := func(err error) *tengo.UserFunction {
		return &tengo.UserFunction{Name: "user_func", Value: func(args ...tengo.Object) (tengo.Object, error) {
			return nil, err
		}}
	}
	userModule := func(err error) *tengo.BuiltinModule {
		return &tengo.BuiltinModule{
			Attrs: map[string]tengo.Object{
				"afunction": &tengo.UserFunction{
					Name: "afunction",
					Value: func(a ...tengo.Object) (tengo.Object, error) {
						return nil, err
					},
				},
			},
		}
	}

	expectError(t, `user_func()`,
		Opts().Symbol("user_func", userFunc(userErr)),
		"Runtime Error: "+userErr.Error(),
	)
	expectErrorIs(t, `user_func()`,
		Opts().Symbol("user_func", userFunc(userErr)),
		userErr,
	)

	wrapUserErr := &customError{err: userErr, str: "custom error"}

	expectErrorIs(t, `user_func()`,
		Opts().Symbol("user_func", userFunc(wrapUserErr)),
		wrapUserErr,
	)
	expectErrorIs(t, `user_func()`,
		Opts().Symbol("user_func", userFunc(wrapUserErr)),
		userErr,
	)
	var asErr1 *customError
	expectErrorAs(t, `user_func()`,
		Opts().Symbol("user_func", userFunc(wrapUserErr)),
		&asErr1,
	)
	require.True(t, asErr1.Error() == wrapUserErr.Error(),
		"expected error as:%v, got:%v", wrapUserErr, asErr1)

	expectError(t, `import("mod1").afunction()`,
		Opts().Module("mod1", userModule(userErr)),
		"Runtime Error: "+userErr.Error(),
	)
	expectErrorIs(t, `import("mod1").afunction()`,
		Opts().Module("mod1", userModule(userErr)),
		userErr,
	)
	expectError(t, `import("mod1").afunction()`,
		Opts().Module("mod1", userModule(wrapUserErr)),
		"Runtime Error: "+wrapUserErr.Error(),
	)
	expectErrorIs(t, `import("mod1").afunction()`,
		Opts().Module("mod1", userModule(wrapUserErr)),
		wrapUserErr,
	)
	expectErrorIs(t, `import("mod1").afunction()`,
		Opts().Module("mod1", userModule(wrapUserErr)),
		userErr,
	)
	var asErr2 *customError
	expectErrorAs(t, `import("mod1").afunction()`,
		Opts().Module("mod1", userModule(wrapUserErr)),
		&asErr2,
	)
	require.True(t, asErr2.Error() == wrapUserErr.Error(),
		"expected error as:%v, got:%v", wrapUserErr, asErr2)
}

func TestError(t *testing.T) {
	expectRun(t, `out = error(1)`, nil, errorObject(1))
	expectRun(t, `out = error(1).value`, nil, 1)
	expectRun(t, `out = error("some error")`, nil, errorObject("some error"))
	expectRun(t, `out = error("some" + " error")`, nil, errorObject("some error"))
	expectRun(t, `out = func() { return error(5) }()`, nil, errorObject(5))
	expectRun(t, `out = error(error("foo"))`, nil, errorObject(errorObject("foo")))
	expectRun(t, `out = error("some error")`, nil, errorObject("some error"))
	expectRun(t, `out = error("some error").value`, nil, "some error")
	expectRun(t, `out = error("some error")["value"]`, nil, "some error")

	expectError(t, `error("error").err`, nil, "invalid index on error")
	expectError(t, `error("error").value_`, nil, "invalid index on error")
	expectError(t, `error([1,2,3])[1]`, nil, "invalid index on error")
}

func TestFloat(t *testing.T) {
	expectRun(t, `out = 0.0`, nil, 0.0)
	expectRun(t, `out = -10.3`, nil, -10.3)
	expectRun(t, `out = 3.2 + 2.0 * -4.0`, nil, -4.8)
	expectRun(t, `out = 4 + 2.3`, nil, 6.3)
	expectRun(t, `out = 2.3 + 4`, nil, 6.3)
	expectRun(t, `out = +5.0`, nil, 5.0)
	expectRun(t, `out = -5.0 + +5.0`, nil, 0.0)
}

func TestForIn(t *testing.T) {
	// array
	expectRun(t, `out = 0; for x in [1, 2, 3] { out += x }`,
		nil, 6) // value
	expectRun(t, `out = 0; for i, x in [1, 2, 3] { out += i + x }`,
		nil, 9) // index, value
	expectRun(t, `out = 0; func() { for i, x in [1, 2, 3] { out += i + x } }()`,
		nil, 9) // index, value
	expectRun(t, `out = 0; for i, _ in [1, 2, 3] { out += i }`,
		nil, 3) // index, _
	expectRun(t, `out = 0; func() { for i, _ in [1, 2, 3] { out += i  } }()`,
		nil, 3) // index, _

	// map
	expectRun(t, `out = 0; for v in {a:2,b:3,c:4} { out += v }`,
		nil, 9) // value
	expectRun(t, `out = ""; for k, v in {a:2,b:3,c:4} { out = k; if v==3 { break } }`,
		nil, "b") // key, value
	expectRun(t, `out = ""; for k, _ in {a:2} { out += k }`,
		nil, "a") // key, _
	expectRun(t, `out = 0; for _, v in {a:2,b:3,c:4} { out += v }`,
		nil, 9) // _, value
	expectRun(t, `out = ""; func() { for k, v in {a:2,b:3,c:4} { out = k; if v==3 { break } } }()`,
		nil, "b") // key, value

	// string
	expectRun(t, `out = ""; for c in "abcde" { out += c }`,
		nil, "abcde")
	expectRun(t, `out = ""; for i, c in "abcde" { if i == 2 { continue }; out += c }`,
		nil, "abde")
}

func TestFor(t *testing.T) {
	expectRun(t, `
	out = 0
	for {
		out++
		if out == 5 {
			break
		}
	}`, nil, 5)

	expectRun(t, `
	out = 0
	for {
		out++
		if out == 5 {
			break
		}
	}`, nil, 5)

	expectRun(t, `
	out = 0
	a := 0
	for {
		a++
		if a == 3 { continue }
		if a == 5 { break }
		out += a
	}`, nil, 7) // 1 + 2 + 4

	expectRun(t, `
	out = 0
	a := 0
	for {
		a++
		if a == 3 { continue }
		out += a
		if a == 5 { break }
	}`, nil, 12) // 1 + 2 + 4 + 5

	expectRun(t, `
	out = 0
	for true {
		out++
		if out == 5 {
			break
		}
	}`, nil, 5)

	expectRun(t, `
	a := 0
	for true {
		a++
		if a == 5 {
			break
		}
	}
	out = a`, nil, 5)

	expectRun(t, `
	out = 0
	a := 0
	for true {
		a++
		if a == 3 { continue }
		if a == 5 { break }
		out += a
	}`, nil, 7) // 1 + 2 + 4

	expectRun(t, `
	out = 0
	a := 0
	for true {
		a++
		if a == 3 { continue }
		out += a
		if a == 5 { break }
	}`, nil, 12) // 1 + 2 + 4 + 5

	expectRun(t, `
	out = 0
	func() {
		for true {
			out++
			if out == 5 {
				return
			}
		}
	}()`, nil, 5)

	expectRun(t, `
	out = 0
	for a:=1; a<=10; a++ {
		out += a
	}`, nil, 55)

	expectRun(t, `
	out = 0
	for a:=1; a<=3; a++ {
		for b:=3; b<=6; b++ {
			out += b
		}
	}`, nil, 54)

	expectRun(t, `
	out = 0
	func() {
		for {
			out++
			if out == 5 {
				break
			}
		}
	}()`, nil, 5)

	expectRun(t, `
	out = 0
	func() {
		for true {
			out++
			if out == 5 {
				break
			}
		}
	}()`, nil, 5)

	expectRun(t, `
	out = func() {
		a := 0
		for {
			a++
			if a == 5 {
				break
			}
		}
		return a
	}()`, nil, 5)

	expectRun(t, `
	out = func() {
		a := 0
		for true {
			a++
			if a== 5 {
				break
			}
		}
		return a
	}()`, nil, 5)

	expectRun(t, `
	out = func() {
		a := 0
		func() {
			for {
				a++
				if a == 5 {
					break
				}
			}
		}()
		return a
	}()`, nil, 5)

	expectRun(t, `
	out = func() {
		a := 0
		func() {
			for true {
				a++
				if a == 5 {
					break
				}
			}
		}()
		return a
	}()`, nil, 5)

	expectRun(t, `
	out = func() {
		sum := 0
		for a:=1; a<=10; a++ {
			sum += a
		}
		return sum
	}()`, nil, 55)

	expectRun(t, `
	out = func() {
		sum := 0
		for a:=1; a<=4; a++ {
			for b:=3; b<=5; b++ {
				sum += b
			}
		}
		return sum
	}()`, nil, 48) // (3+4+5) * 4

	expectRun(t, `
	a := 1
	for ; a<=10; a++ {
		if a == 5 {
			break
		}
	}
	out = a`, nil, 5)

	expectRun(t, `
	out = 0
	for a:=1; a<=10; a++ {
		if a == 3 {
			continue
		}
		out += a
		if a == 5 {
			break
		}
	}`, nil, 12) // 1 + 2 + 4 + 5

	expectRun(t, `
	out = 0
	for a:=1; a<=10; {
		if a == 3 {
			a++
			continue
		}
		out += a
		if a == 5 {
			break
		}
		a++
	}`, nil, 12) // 1 + 2 + 4 + 5
}

func TestFunction(t *testing.T) {
	// function with no "return" statement returns "invalid" value.
	expectRun(t, `f1 := func() {}; out = f1();`,
		nil, tengo.UndefinedValue)
	expectRun(t, `f1 := func() {}; f2 := func() { return f1(); }; f1(); out = f2();`,
		nil, tengo.UndefinedValue)
	expectRun(t, `f := func(x) { x; }; out = f(5);`,
		nil, tengo.UndefinedValue)

	expectRun(t, `f := func(...x) { return x; }; out = f(1,2,3);`,
		nil, ARR{1, 2, 3})

	expectRun(t, `f := func(a, b, ...x) { return [a, b, x]; }; out = f(8,9,1,2,3);`,
		nil, ARR{8, 9, ARR{1, 2, 3}})

	expectRun(t, `f := func(v) { x := 2; return func(a, ...b){ return [a, b, v+x]}; }; out = f(5)("a", "b");`,
		nil, ARR{"a", ARR{"b"}, 7})

	expectRun(t, `f := func(...x) { return x; }; out = f();`,
		nil, &tengo.Array{Value: []tengo.Object{}})

	expectRun(t, `f := func(a, b, ...x) { return [a, b, x]; }; out = f(8, 9);`,
		nil, ARR{8, 9, ARR{}})

	expectRun(t, `f := func(v) { x := 2; return func(a, ...b){ return [a, b, v+x]}; }; out = f(5)("a");`,
		nil, ARR{"a", ARR{}, 7})

	expectError(t, `f := func(a, b, ...x) { return [a, b, x]; }; f();`, nil,
		"Runtime Error: wrong number of arguments: want>=2, got=0\n\tat test:1:46")

	expectError(t, `f := func(a, b, ...x) { return [a, b, x]; }; f(1);`, nil,
		"Runtime Error: wrong number of arguments: want>=2, got=1\n\tat test:1:46")

	expectRun(t, `f := func(x) { return x; }; out = f(5);`, nil, 5)
	expectRun(t, `f := func(x) { return x * 2; }; out = f(5);`, nil, 10)
	expectRun(t, `f := func(x, y) { return x + y; }; out = f(5, 5);`, nil, 10)
	expectRun(t, `f := func(x, y) { return x + y; }; out = f(5 + 5, f(5, 5));`,
		nil, 20)
	expectRun(t, `out = func(x) { return x; }(5)`, nil, 5)
	expectRun(t, `x := 10; f := func(x) { return x; }; f(5); out = x;`, nil, 10)

	expectRun(t, `
	f2 := func(a) {
		f1 := func(a) {
			return a * 2;
		};
	
		return f1(a) * 3;
	};
	
	out = f2(10);
	`, nil, 60)

	expectRun(t, `
		f1 := func(f) {
			a := [undefined]
			a[0] = func() { return f(a) }
			return a[0]()
		}

		out = f1(func(a) { return 2 })
	`, nil, 2)

	// closures
	expectRun(t, `
		newAdder := func(x) {
			return func(y) { return x + y };
		};
	
		add2 := newAdder(2);
		out = add2(5);
		`, nil, 7)
	expectRun(t, `
		m := {a: 1}
		for k,v in m {
			func(){
				out = k
			}()
		}
		`, nil, "a")

	expectRun(t, `
		m := {a: 1}
		for k,v in m {
			func(){
				out = v
			}()
		}
		`, nil, 1)
	// function as a argument
	expectRun(t, `
	add := func(a, b) { return a + b };
	sub := func(a, b) { return a - b };
	applyFunc := func(a, b, f) { return f(a, b) };
	
	out = applyFunc(applyFunc(2, 2, add), 3, sub);
	`, nil, 1)

	expectRun(t, `f1 := func() { return 5 + 10; }; out = f1();`,
		nil, 15)
	expectRun(t, `f1 := func() { return 1 }; f2 := func() { return 2 }; out = f1() + f2()`,
		nil, 3)
	expectRun(t, `f1 := func() { return 1 }; f2 := func() { return f1() + 2 }; f3 := func() { return f2() + 3 }; out = f3()`,
		nil, 6)
	expectRun(t, `f1 := func() { return 99; 100 }; out = f1();`,
		nil, 99)
	expectRun(t, `f1 := func() { return 99; return 100 }; out = f1();`,
		nil, 99)
	expectRun(t, `f1 := func() { return 33; }; f2 := func() { return f1 }; out = f2()();`,
		nil, 33)
	expectRun(t, `one := func() { one = 1; return one }; out = one()`,
		nil, 1)
	expectRun(t, `three := func() { one := 1; two := 2; return one + two }; out = three()`,
		nil, 3)
	expectRun(t, `three := func() { one := 1; two := 2; return one + two }; seven := func() { three := 3; four := 4; return three + four }; out = three() + seven()`,
		nil, 10)
	expectRun(t, `
	foo1 := func() {
		foo := 50
		return foo
	}
	foo2 := func() {
		foo := 100
		return foo
	}
	out = foo1() + foo2()`, nil, 150)
	expectRun(t, `
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
	`, nil, 97)
	expectRun(t, `
	f1 := func() {
		f2 := func() { return 1; }
		return f2
	};
	out = f1()()
	`, nil, 1)

	expectRun(t, `
	f1 := func(a) { return a; };
	out = f1(4)`, nil, 4)
	expectRun(t, `
	f1 := func(a, b) { return a + b; };
	out = f1(1, 2)`, nil, 3)

	expectRun(t, `
	sum := func(a, b) {
		c := a + b;
		return c;
	};
	out = sum(1, 2);`, nil, 3)

	expectRun(t, `
	sum := func(a, b) {
		c := a + b;
		return c;
	};
	out = sum(1, 2) + sum(3, 4);`, nil, 10)

	expectRun(t, `
	sum := func(a, b) {
		c := a + b
		return c
	};
	outer := func() {
		return sum(1, 2) + sum(3, 4)
	};
	out = outer();`, nil, 10)

	expectRun(t, `
	g := 10;
	
	sum := func(a, b) {
		c := a + b;
		return c + g;
	}
	
	outer := func() {
		return sum(1, 2) + sum(3, 4) + g;
	}
	
	out = outer() + g
	`, nil, 50)

	expectError(t, `func() { return 1; }(1)`,
		nil, "wrong number of arguments")
	expectError(t, `func(a) { return a; }()`,
		nil, "wrong number of arguments")
	expectError(t, `func(a, b) { return a + b; }(1)`,
		nil, "wrong number of arguments")

	expectRun(t, `
		f1 := func(a) {
			return func() { return a; };
		};
		f2 := f1(99);
		out = f2()
		`, nil, 99)

	expectRun(t, `
		f1 := func(a, b) {
			return func(c) { return a + b + c };
		};
	
		f2 := f1(1, 2);
		out = f2(8);
		`, nil, 11)
	expectRun(t, `
		f1 := func(a, b) {
			c := a + b;
			return func(d) { return c + d };
		};
		f2 := f1(1, 2);
		out = f2(8);
		`, nil, 11)
	expectRun(t, `
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
		`, nil, 14)
	expectRun(t, `
		a := 1;
		f1 := func(b) {
			return func(c) {
				return func(d) { return a + b + c + d }
			};
		};
		f2 := f1(2);
		f3 := f2(3);
		out = f3(8);
		`, nil, 14)
	expectRun(t, `
		f1 := func(a, b) {
			one := func() { return a; };
			two := func() { return b; };
			return func() { return one() + two(); }
		};
		f2 := f1(9, 90);
		out = f2();
		`, nil, 99)

	// global function recursion
	expectRun(t, `
		fib := func(x) {
			if x == 0 {
				return 0
			} else if x == 1 {
				return 1
			} else {
				return fib(x-1) + fib(x-2)
			}
		}
		out = fib(15)`, nil, 610)

	// local function recursion
	expectRun(t, `
out = func() {
	sum := func(x) {
		return x == 0 ? 0 : x + sum(x-1)
	}
	return sum(5)
}()`, nil, 15)

	expectError(t, `return 5`, nil, "return not allowed outside function")

	// closure and block scopes
	expectRun(t, `
func() {
	a := 10
	func() {
		b := 5
		if true {
			out = a + 5
		}
	}()
}()`, nil, 15)
	expectRun(t, `
func() {
	a := 10
	b := func() { return 5 }
	func() {
		if b() {
			out = a + b()
		}
	}()
}()`, nil, 15)
	expectRun(t, `
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
}()`, nil, 15)

	// function skipping return
	expectRun(t, `out = func() {}()`,
		nil, tengo.UndefinedValue)
	expectRun(t, `out = func(v) { if v { return true } }(1)`,
		nil, true)
	expectRun(t, `out = func(v) { if v { return true } }(0)`,
		nil, tengo.UndefinedValue)
	expectRun(t, `out = func(v) { if v { } else { return true } }(1)`,
		nil, tengo.UndefinedValue)
	expectRun(t, `out = func(v) { if v { return } }(1)`,
		nil, tengo.UndefinedValue)
	expectRun(t, `out = func(v) { if v { return } }(0)`,
		nil, tengo.UndefinedValue)
	expectRun(t, `out = func(v) { if v { } else { return } }(1)`,
		nil, tengo.UndefinedValue)
	expectRun(t, `out = func(v) { for ;;v++ { if v == 3 { return true } } }(1)`,
		nil, true)
	expectRun(t, `out = func(v) { for ;;v++ { if v == 3 { break } } }(1)`,
		nil, tengo.UndefinedValue)

	// 'f' in RHS at line 4 must reference global variable 'f'
	// See https://github.com/d5/tengo/issues/314
	expectRun(t, `
f := func() { return 2 }
out = (func() {
	f := f()
	return f
})()
	`, nil, 2)
}

func TestBlocksInGlobalScope(t *testing.T) {
	expectRun(t, `
f := undefined
if true {
	a := 1
	f = func() {
		a = 2
	}
}
b := 3
f()
out = b`,
		nil, 3)

	expectRun(t, `
func() {
	f := undefined
	if true {
		a := 10
		f = func() {
			a = 20
		}
	}
	b := 5
	f()
	out = b
}()
	`,
		nil, 5)

	expectRun(t, `
f := undefined
if true {
	a := 1
	b := 2
	f = func() {
		a = 3
		b = 4
	}
}
c := 5
d := 6
f()
out = c + d`,
		nil, 11)

	expectRun(t, `
fn := undefined
if true {
	a := 1
	b := 2
	if true {
		c := 3
		d := 4
		fn = func() {
			a = 5
			b = 6
			c = 7
			d = 8
		}
	}
}
e := 9
f := 10
fn()
out = e + f`,
		nil, 19)

	expectRun(t, `
out = 0
func() {
	for x in [1, 2, 3] {
		out += x
	}
}()`,
		nil, 6)

	expectRun(t, `
out = 0
for x in [1, 2, 3] {
	out += x
}`,
		nil, 6)
}

func TestIf(t *testing.T) {

	expectRun(t, `if (true) { out = 10 }`, nil, 10)
	expectRun(t, `if (false) { out = 10 }`, nil, tengo.UndefinedValue)
	expectRun(t, `if (false) { out = 10 } else { out = 20 }`, nil, 20)
	expectRun(t, `if (1) { out = 10 }`, nil, 10)
	expectRun(t, `if (0) { out = 10 } else { out = 20 }`, nil, 20)
	expectRun(t, `if (1 < 2) { out = 10 }`, nil, 10)
	expectRun(t, `if (1 > 2) { out = 10 }`, nil, tengo.UndefinedValue)
	expectRun(t, `if (1 < 2) { out = 10 } else { out = 20 }`, nil, 10)
	expectRun(t, `if (1 > 2) { out = 10 } else { out = 20 }`, nil, 20)

	expectRun(t, `if (1 < 2) { out = 10 } else if (1 > 2) { out = 20 } else { out = 30 }`,
		nil, 10)
	expectRun(t, `if (1 > 2) { out = 10 } else if (1 < 2) { out = 20 } else { out = 30 }`,
		nil, 20)
	expectRun(t, `if (1 > 2) { out = 10 } else if (1 == 2) { out = 20 } else { out = 30 }`,
		nil, 30)
	expectRun(t, `if (1 > 2) { out = 10 } else if (1 == 2) { out = 20 } else if (1 < 2) { out = 30 } else { out = 40 }`,
		nil, 30)
	expectRun(t, `if (1 > 2) { out = 10 } else if (1 < 2) { out = 20; out = 21; out = 22 } else { out = 30 }`,
		nil, 22)
	expectRun(t, `if (1 > 2) { out = 10 } else if (1 == 2) { out = 20 } else { out = 30; out = 31; out = 32}`,
		nil, 32)
	expectRun(t, `if (1 > 2) { out = 10 } else if (1 < 2) { if (1 == 2) { out = 21 } else { out = 22 } } else { out = 30 }`,
		nil, 22)
	expectRun(t, `if (1 > 2) { out = 10 } else if (1 < 2) { if (1 == 2) { out = 21 } else if (2 == 3) { out = 22 } else { out = 23 } } else { out = 30 }`,
		nil, 23)
	expectRun(t, `if (1 > 2) { out = 10 } else if (1 == 2) { if (1 == 2) { out = 21 } else if (2 == 3) { out = 22 } else { out = 23 } } else { out = 30 }`,
		nil, 30)
	expectRun(t, `if (1 > 2) { out = 10 } else if (1 == 2) { out = 20 } else { if (1 == 2) { out = 31 } else if (2 == 3) { out = 32 } else { out = 33 } }`,
		nil, 33)

	expectRun(t, `if a:=0; a<1 { out = 10 }`, nil, 10)
	expectRun(t, `a:=0; if a++; a==1 { out = 10 }`, nil, 10)
	expectRun(t, `
func() {
	a := 1
	if a++; a > 1 {
		out = a
	}
}()
`, nil, 2)
	expectRun(t, `
func() {
	a := 1
	if a++; a == 1 {
		out = 10
	} else {
		out = 20
	}
}()
`, nil, 20)
	expectRun(t, `
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
	expectRun(t, `a := 1; if a; a { out = a }`, nil, 1)
	expectRun(t, `a := 1; if a + 4; a { out = a }`, nil, 1)

	// dead code elimination
	expectRun(t, `
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

func TestImmutable(t *testing.T) {
	// primitive types are already immutable values
	// immutable expression has no effects.
	expectRun(t, `a := immutable(1); out = a`, nil, 1)
	expectRun(t, `a := 5; b := immutable(a); out = b`, nil, 5)
	expectRun(t, `a := immutable(1); a = 5; out = a`, nil, 5)

	// array
	expectError(t, `a := immutable([1, 2, 3]); a[1] = 5`,
		nil, "not index-assignable")
	expectError(t, `a := immutable(["foo", [1,2,3]]); a[1] = "bar"`,
		nil, "not index-assignable")
	expectRun(t, `a := immutable(["foo", [1,2,3]]); a[1][1] = "bar"; out = a`,
		nil, IARR{"foo", ARR{1, "bar", 3}})
	expectError(t, `a := immutable(["foo", immutable([1,2,3])]); a[1][1] = "bar"`,
		nil, "not index-assignable")
	expectError(t, `a := ["foo", immutable([1,2,3])]; a[1][1] = "bar"`,
		nil, "not index-assignable")
	expectRun(t, `a := immutable([1,2,3]); b := copy(a); b[1] = 5; out = b`,
		nil, ARR{1, 5, 3})
	expectRun(t, `a := immutable([1,2,3]); b := copy(a); b[1] = 5; out = a`,
		nil, IARR{1, 2, 3})
	expectRun(t, `out = immutable([1,2,3]) == [1,2,3]`,
		nil, true)
	expectRun(t, `out = immutable([1,2,3]) == immutable([1,2,3])`,
		nil, true)
	expectRun(t, `out = [1,2,3] == immutable([1,2,3])`,
		nil, true)
	expectRun(t, `out = immutable([1,2,3]) == [1,2]`,
		nil, false)
	expectRun(t, `out = immutable([1,2,3]) == immutable([1,2])`,
		nil, false)
	expectRun(t, `out = [1,2,3] == immutable([1,2])`,
		nil, false)
	expectRun(t, `out = immutable([1, 2, 3, 4])[1]`,
		nil, 2)
	expectRun(t, `out = immutable([1, 2, 3, 4])[1:3]`,
		nil, ARR{2, 3})
	expectRun(t, `a := immutable([1,2,3]); a = 5; out = a`,
		nil, 5)
	expectRun(t, `a := immutable([1, 2, 3]); out = a[5]`,
		nil, tengo.UndefinedValue)

	// map
	expectError(t, `a := immutable({b: 1, c: 2}); a.b = 5`,
		nil, "not index-assignable")
	expectError(t, `a := immutable({b: 1, c: 2}); a["b"] = "bar"`,
		nil, "not index-assignable")
	expectRun(t, `a := immutable({b: 1, c: [1,2,3]}); a.c[1] = "bar"; out = a`,
		nil, IMAP{"b": 1, "c": ARR{1, "bar", 3}})
	expectError(t, `a := immutable({b: 1, c: immutable([1,2,3])}); a.c[1] = "bar"`,
		nil, "not index-assignable")
	expectError(t, `a := {b: 1, c: immutable([1,2,3])}; a.c[1] = "bar"`,
		nil, "not index-assignable")
	expectRun(t, `out = immutable({a:1,b:2}) == {a:1,b:2}`,
		nil, true)
	expectRun(t, `out = immutable({a:1,b:2}) == immutable({a:1,b:2})`,
		nil, true)
	expectRun(t, `out = {a:1,b:2} == immutable({a:1,b:2})`,
		nil, true)
	expectRun(t, `out = immutable({a:1,b:2}) == {a:1,b:3}`,
		nil, false)
	expectRun(t, `out = immutable({a:1,b:2}) == immutable({a:1,b:3})`,
		nil, false)
	expectRun(t, `out = {a:1,b:2} == immutable({a:1,b:3})`,
		nil, false)
	expectRun(t, `out = immutable({a:1,b:2}).b`,
		nil, 2)
	expectRun(t, `out = immutable({a:1,b:2})["b"]`,
		nil, 2)
	expectRun(t, `a := immutable({a:1,b:2}); a = 5; out = 5`,
		nil, 5)
	expectRun(t, `a := immutable({a:1,b:2}); out = a.c`,
		nil, tengo.UndefinedValue)

	expectRun(t, `a := immutable({b: 5, c: "foo"}); out = a.b`,
		nil, 5)
	expectError(t, `a := immutable({b: 5, c: "foo"}); a.b = 10`,
		nil, "not index-assignable")
}

func TestIncDec(t *testing.T) {
	expectRun(t, `out = 0; out++`, nil, 1)
	expectRun(t, `out = 0; out--`, nil, -1)
	expectRun(t, `a := 0; a++; out = a`, nil, 1)
	expectRun(t, `a := 0; a++; a--; out = a`, nil, 0)

	// this seems strange but it works because 'a += b' is
	// translated into 'a = a + b' and string type takes other types for + operator.
	expectRun(t, `a := "foo"; a++; out = a`, nil, "foo1")
	expectError(t, `a := "foo"; a--`, nil, "invalid operation")

	expectError(t, `a++`, nil, "unresolved reference") // not declared
	expectError(t, `a--`, nil, "unresolved reference") // not declared
	expectError(t, `4++`, nil, "unresolved reference")
}

type StringDict struct {
	tengo.ObjectImpl
	Value map[string]string
}

func (o *StringDict) String() string { return "" }

func (o *StringDict) TypeName() string {
	return "string-dict"
}

func (o *StringDict) IndexGet(index tengo.Object) (tengo.Object, error) {
	strIdx, ok := index.(*tengo.String)
	if !ok {
		return nil, tengo.ErrInvalidIndexType
	}

	for k, v := range o.Value {
		if strings.EqualFold(strIdx.Value, k) {
			return &tengo.String{Value: v}, nil
		}
	}

	return tengo.UndefinedValue, nil
}

func (o *StringDict) IndexSet(index, value tengo.Object) error {
	strIdx, ok := index.(*tengo.String)
	if !ok {
		return tengo.ErrInvalidIndexType
	}

	strVal, ok := tengo.ToString(value)
	if !ok {
		return tengo.ErrInvalidIndexValueType
	}

	o.Value[strings.ToLower(strIdx.Value)] = strVal

	return nil
}

type StringCircle struct {
	tengo.ObjectImpl
	Value []string
}

func (o *StringCircle) TypeName() string {
	return "string-circle"
}

func (o *StringCircle) String() string {
	return ""
}

func (o *StringCircle) IndexGet(index tengo.Object) (tengo.Object, error) {
	intIdx, ok := index.(*tengo.Int)
	if !ok {
		return nil, tengo.ErrInvalidIndexType
	}

	r := int(intIdx.Value) % len(o.Value)
	if r < 0 {
		r = len(o.Value) + r
	}

	return &tengo.String{Value: o.Value[r]}, nil
}

func (o *StringCircle) IndexSet(index, value tengo.Object) error {
	intIdx, ok := index.(*tengo.Int)
	if !ok {
		return tengo.ErrInvalidIndexType
	}

	r := int(intIdx.Value) % len(o.Value)
	if r < 0 {
		r = len(o.Value) + r
	}

	strVal, ok := tengo.ToString(value)
	if !ok {
		return tengo.ErrInvalidIndexValueType
	}

	o.Value[r] = strVal

	return nil
}

type StringArray struct {
	tengo.ObjectImpl
	Value []string
}

func (o *StringArray) String() string {
	return strings.Join(o.Value, ", ")
}

func (o *StringArray) BinaryOp(
	op token.Token,
	rhs tengo.Object,
) (tengo.Object, error) {
	if rhs, ok := rhs.(*StringArray); ok {
		switch op {
		case token.Add:
			if len(rhs.Value) == 0 {
				return o, nil
			}
			return &StringArray{Value: append(o.Value, rhs.Value...)}, nil
		}
	}

	return nil, tengo.ErrInvalidOperator
}

func (o *StringArray) IsFalsy() bool {
	return len(o.Value) == 0
}

func (o *StringArray) Equals(x tengo.Object) bool {
	if x, ok := x.(*StringArray); ok {
		if len(o.Value) != len(x.Value) {
			return false
		}

		for i, v := range o.Value {
			if v != x.Value[i] {
				return false
			}
		}

		return true
	}

	return false
}

func (o *StringArray) Copy() tengo.Object {
	return &StringArray{
		Value: append([]string{}, o.Value...),
	}
}

func (o *StringArray) TypeName() string {
	return "string-array"
}

func (o *StringArray) IndexGet(index tengo.Object) (tengo.Object, error) {
	intIdx, ok := index.(*tengo.Int)
	if ok {
		if intIdx.Value >= 0 && intIdx.Value < int64(len(o.Value)) {
			return &tengo.String{Value: o.Value[intIdx.Value]}, nil
		}

		return nil, tengo.ErrIndexOutOfBounds
	}

	strIdx, ok := index.(*tengo.String)
	if ok {
		for vidx, str := range o.Value {
			if strIdx.Value == str {
				return &tengo.Int{Value: int64(vidx)}, nil
			}
		}

		return tengo.UndefinedValue, nil
	}

	return nil, tengo.ErrInvalidIndexType
}

func (o *StringArray) IndexSet(index, value tengo.Object) error {
	strVal, ok := tengo.ToString(value)
	if !ok {
		return tengo.ErrInvalidIndexValueType
	}

	intIdx, ok := index.(*tengo.Int)
	if ok {
		if intIdx.Value >= 0 && intIdx.Value < int64(len(o.Value)) {
			o.Value[intIdx.Value] = strVal
			return nil
		}

		return tengo.ErrIndexOutOfBounds
	}

	return tengo.ErrInvalidIndexType
}

func (o *StringArray) Call(
	args ...tengo.Object,
) (ret tengo.Object, err error) {
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}

	s1, ok := tengo.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	for i, v := range o.Value {
		if v == s1 {
			return &tengo.Int{Value: int64(i)}, nil
		}
	}

	return tengo.UndefinedValue, nil
}

func (o *StringArray) CanCall() bool {
	return true
}

func TestIndexable(t *testing.T) {
	dict := func() *StringDict {
		return &StringDict{Value: map[string]string{"a": "foo", "b": "bar"}}
	}
	expectRun(t, `out = dict["a"]`,
		Opts().Symbol("dict", dict()).Skip2ndPass(), "foo")
	expectRun(t, `out = dict["B"]`,
		Opts().Symbol("dict", dict()).Skip2ndPass(), "bar")
	expectRun(t, `out = dict["x"]`,
		Opts().Symbol("dict", dict()).Skip2ndPass(), tengo.UndefinedValue)
	expectError(t, `dict[0]`,
		Opts().Symbol("dict", dict()).Skip2ndPass(), "invalid index type")

	strCir := func() *StringCircle {
		return &StringCircle{Value: []string{"one", "two", "three"}}
	}
	expectRun(t, `out = cir[0]`,
		Opts().Symbol("cir", strCir()).Skip2ndPass(), "one")
	expectRun(t, `out = cir[1]`,
		Opts().Symbol("cir", strCir()).Skip2ndPass(), "two")
	expectRun(t, `out = cir[-1]`,
		Opts().Symbol("cir", strCir()).Skip2ndPass(), "three")
	expectRun(t, `out = cir[-2]`,
		Opts().Symbol("cir", strCir()).Skip2ndPass(), "two")
	expectRun(t, `out = cir[3]`,
		Opts().Symbol("cir", strCir()).Skip2ndPass(), "one")
	expectError(t, `cir["a"]`,
		Opts().Symbol("cir", strCir()).Skip2ndPass(), "invalid index type")

	strArr := func() *StringArray {
		return &StringArray{Value: []string{"one", "two", "three"}}
	}
	expectRun(t, `out = arr["one"]`,
		Opts().Symbol("arr", strArr()).Skip2ndPass(), 0)
	expectRun(t, `out = arr["three"]`,
		Opts().Symbol("arr", strArr()).Skip2ndPass(), 2)
	expectRun(t, `out = arr["four"]`,
		Opts().Symbol("arr", strArr()).Skip2ndPass(), tengo.UndefinedValue)
	expectRun(t, `out = arr[0]`,
		Opts().Symbol("arr", strArr()).Skip2ndPass(), "one")
	expectRun(t, `out = arr[1]`,
		Opts().Symbol("arr", strArr()).Skip2ndPass(), "two")
	expectError(t, `arr[-1]`,
		Opts().Symbol("arr", strArr()).Skip2ndPass(), "index out of bounds")
}

func TestIndexAssignable(t *testing.T) {
	dict := func() *StringDict {
		return &StringDict{Value: map[string]string{"a": "foo", "b": "bar"}}
	}
	expectRun(t, `dict["a"] = "1984"; out = dict["a"]`,
		Opts().Symbol("dict", dict()).Skip2ndPass(), "1984")
	expectRun(t, `dict["c"] = "1984"; out = dict["c"]`,
		Opts().Symbol("dict", dict()).Skip2ndPass(), "1984")
	expectRun(t, `dict["c"] = 1984; out = dict["C"]`,
		Opts().Symbol("dict", dict()).Skip2ndPass(), "1984")
	expectError(t, `dict[0] = "1984"`,
		Opts().Symbol("dict", dict()).Skip2ndPass(), "invalid index type")

	strCir := func() *StringCircle {
		return &StringCircle{Value: []string{"one", "two", "three"}}
	}
	expectRun(t, `cir[0] = "ONE"; out = cir[0]`,
		Opts().Symbol("cir", strCir()).Skip2ndPass(), "ONE")
	expectRun(t, `cir[1] = "TWO"; out = cir[1]`,
		Opts().Symbol("cir", strCir()).Skip2ndPass(), "TWO")
	expectRun(t, `cir[-1] = "THREE"; out = cir[2]`,
		Opts().Symbol("cir", strCir()).Skip2ndPass(), "THREE")
	expectRun(t, `cir[0] = "ONE"; out = cir[3]`,
		Opts().Symbol("cir", strCir()).Skip2ndPass(), "ONE")
	expectError(t, `cir["a"] = "ONE"`,
		Opts().Symbol("cir", strCir()).Skip2ndPass(), "invalid index type")

	strArr := func() *StringArray {
		return &StringArray{Value: []string{"one", "two", "three"}}
	}
	expectRun(t, `arr[0] = "ONE"; out = arr[0]`,
		Opts().Symbol("arr", strArr()).Skip2ndPass(), "ONE")
	expectRun(t, `arr[1] = "TWO"; out = arr[1]`,
		Opts().Symbol("arr", strArr()).Skip2ndPass(), "TWO")
	expectError(t, `arr["one"] = "ONE"`,
		Opts().Symbol("arr", strArr()).Skip2ndPass(), "invalid index type")
}

func TestInteger(t *testing.T) {
	expectRun(t, `out = 5`, nil, 5)
	expectRun(t, `out = 10`, nil, 10)
	expectRun(t, `out = -5`, nil, -5)
	expectRun(t, `out = -10`, nil, -10)
	expectRun(t, `out = 5 + 5 + 5 + 5 - 10`, nil, 10)
	expectRun(t, `out = 2 * 2 * 2 * 2 * 2`, nil, 32)
	expectRun(t, `out = -50 + 100 + -50`, nil, 0)
	expectRun(t, `out = 5 * 2 + 10`, nil, 20)
	expectRun(t, `out = 5 + 2 * 10`, nil, 25)
	expectRun(t, `out = 20 + 2 * -10`, nil, 0)
	expectRun(t, `out = 50 / 2 * 2 + 10`, nil, 60)
	expectRun(t, `out = 2 * (5 + 10)`, nil, 30)
	expectRun(t, `out = 3 * 3 * 3 + 10`, nil, 37)
	expectRun(t, `out = 3 * (3 * 3) + 10`, nil, 37)
	expectRun(t, `out = (5 + 10 * 2 + 15 /3) * 2 + -10`, nil, 50)
	expectRun(t, `out = 5 % 3`, nil, 2)
	expectRun(t, `out = 5 % 3 + 4`, nil, 6)
	expectRun(t, `out = +5`, nil, 5)
	expectRun(t, `out = +5 + -5`, nil, 0)

	expectRun(t, `out = 9 + '0'`, nil, '9')
	expectRun(t, `out = '9' - 5`, nil, '4')
}

type StringArrayIterator struct {
	tengo.ObjectImpl
	strArr *StringArray
	idx    int
}

func (i *StringArrayIterator) TypeName() string {
	return "string-array-iterator"
}

func (i *StringArrayIterator) String() string {
	return ""
}

func (i *StringArrayIterator) Next() bool {
	i.idx++
	return i.idx <= len(i.strArr.Value)
}

func (i *StringArrayIterator) Key() tengo.Object {
	return &tengo.Int{Value: int64(i.idx - 1)}
}

func (i *StringArrayIterator) Value() tengo.Object {
	return &tengo.String{Value: i.strArr.Value[i.idx-1]}
}

func (o *StringArray) Iterate() tengo.Iterator {
	return &StringArrayIterator{
		strArr: o,
	}
}

func (o *StringArray) CanIterate() bool {
	return true
}

func TestIterable(t *testing.T) {
	strArr := func() *StringArray {
		return &StringArray{Value: []string{"one", "two", "three"}}
	}
	expectRun(t, `for i, s in arr { out += i }`,
		Opts().Symbol("arr", strArr()).Skip2ndPass(), 3)
	expectRun(t, `for i, s in arr { out += s }`,
		Opts().Symbol("arr", strArr()).Skip2ndPass(), "onetwothree")
	expectRun(t, `for i, s in arr { out += s + i }`,
		Opts().Symbol("arr", strArr()).Skip2ndPass(), "one0two1three2")
}

func TestLogical(t *testing.T) {
	expectRun(t, `out = true && true`, nil, true)
	expectRun(t, `out = true && false`, nil, false)
	expectRun(t, `out = false && true`, nil, false)
	expectRun(t, `out = false && false`, nil, false)
	expectRun(t, `out = !true && true`, nil, false)
	expectRun(t, `out = !true && false`, nil, false)
	expectRun(t, `out = !false && true`, nil, true)
	expectRun(t, `out = !false && false`, nil, false)

	expectRun(t, `out = true || true`, nil, true)
	expectRun(t, `out = true || false`, nil, true)
	expectRun(t, `out = false || true`, nil, true)
	expectRun(t, `out = false || false`, nil, false)
	expectRun(t, `out = !true || true`, nil, true)
	expectRun(t, `out = !true || false`, nil, false)
	expectRun(t, `out = !false || true`, nil, true)
	expectRun(t, `out = !false || false`, nil, true)

	expectRun(t, `out = 1 && 2`, nil, 2)
	expectRun(t, `out = 1 || 2`, nil, 1)
	expectRun(t, `out = 1 && 0`, nil, 0)
	expectRun(t, `out = 1 || 0`, nil, 1)
	expectRun(t, `out = 1 && (0 || 2)`, nil, 2)
	expectRun(t, `out = 0 || (0 || 2)`, nil, 2)
	expectRun(t, `out = 0 || (0 && 2)`, nil, 0)
	expectRun(t, `out = 0 || (2 && 0)`, nil, 0)

	expectRun(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; t() && f()`,
		nil, 7)
	expectRun(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; f() && t()`,
		nil, 7)
	expectRun(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; f() || t()`,
		nil, 3)
	expectRun(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; t() || f()`,
		nil, 3)
	expectRun(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; !t() && f()`,
		nil, 3)
	expectRun(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; !f() && t()`,
		nil, 3)
	expectRun(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; !f() || t()`,
		nil, 7)
	expectRun(t, `t:=func() {out = 3; return true}; f:=func() {out = 7; return false}; !t() || f()`,
		nil, 7)
}

func TestMap(t *testing.T) {
	expectRun(t, `
out = {
	one: 10 - 9,
	two: 1 + 1,
	three: 6 / 2
}`, nil, MAP{
		"one":   1,
		"two":   2,
		"three": 3,
	})

	expectRun(t, `
out = {
	"one": 10 - 9,
	"two": 1 + 1,
	"three": 6 / 2
}`, nil, MAP{
		"one":   1,
		"two":   2,
		"three": 3,
	})

	expectRun(t, `out = {foo: 5}["foo"]`, nil, 5)
	expectRun(t, `out = {foo: 5}["bar"]`, nil, tengo.UndefinedValue)
	expectRun(t, `key := "foo"; out = {foo: 5}[key]`, nil, 5)
	expectRun(t, `out = {}["foo"]`, nil, tengo.UndefinedValue)

	expectRun(t, `
m := {
	foo: func(x) {
		return x * 2
	}
}
out = m["foo"](2) + m["foo"](3)
`, nil, 10)

	// map assignment is copy-by-reference
	expectRun(t, `m1 := {k1: 1, k2: "foo"}; m2 := m1; m1.k1 = 5; out = m2.k1`,
		nil, 5)
	expectRun(t, `m1 := {k1: 1, k2: "foo"}; m2 := m1; m2.k1 = 3; out = m1.k1`,
		nil, 3)
	expectRun(t, `func() { m1 := {k1: 1, k2: "foo"}; m2 := m1; m1.k1 = 5; out = m2.k1 }()`,
		nil, 5)
	expectRun(t, `func() { m1 := {k1: 1, k2: "foo"}; m2 := m1; m2.k1 = 3; out = m1.k1 }()`,
		nil, 3)
}

func TestBuiltin(t *testing.T) {
	m := Opts().Module("math",
		&tengo.BuiltinModule{
			Attrs: map[string]tengo.Object{
				"abs": &tengo.UserFunction{
					Name: "abs",
					Value: func(a ...tengo.Object) (tengo.Object, error) {
						v, _ := tengo.ToFloat64(a[0])
						return &tengo.Float{Value: math.Abs(v)}, nil
					},
				},
			},
		})

	// builtin
	expectRun(t, `math := import("math"); out = math.abs(1)`, m, 1.0)
	expectRun(t, `math := import("math"); out = math.abs(-1)`, m, 1.0)
	expectRun(t, `math := import("math"); out = math.abs(1.0)`, m, 1.0)
	expectRun(t, `math := import("math"); out = math.abs(-1.0)`, m, 1.0)
}

func TestUserModules(t *testing.T) {
	// export none
	expectRun(t, `out = import("mod1")`,
		Opts().Module("mod1", `fn := func() { return 5.0 }; a := 2`),
		tengo.UndefinedValue)

	// export values
	expectRun(t, `out = import("mod1")`,
		Opts().Module("mod1", `export 5`), 5)
	expectRun(t, `out = import("mod1")`,
		Opts().Module("mod1", `export "foo"`), "foo")

	// export compound types
	expectRun(t, `out = import("mod1")`,
		Opts().Module("mod1", `export [1, 2, 3]`), IARR{1, 2, 3})
	expectRun(t, `out = import("mod1")`,
		Opts().Module("mod1", `export {a: 1, b: 2}`), IMAP{"a": 1, "b": 2})

	// export value is immutable
	expectError(t, `m1 := import("mod1"); m1.a = 5`,
		Opts().Module("mod1", `export {a: 1, b: 2}`), "not index-assignable")
	expectError(t, `m1 := import("mod1"); m1[1] = 5`,
		Opts().Module("mod1", `export [1, 2, 3]`), "not index-assignable")

	// code after export statement will not be executed
	expectRun(t, `out = import("mod1")`,
		Opts().Module("mod1", `a := 10; export a; a = 20`), 10)
	expectRun(t, `out = import("mod1")`,
		Opts().Module("mod1", `a := 10; export a; a = 20; export a`), 10)

	// export function
	expectRun(t, `out = import("mod1")()`,
		Opts().Module("mod1", `export func() { return 5.0 }`), 5.0)
	// export function that reads module-global variable
	expectRun(t, `out = import("mod1")()`,
		Opts().Module("mod1", `a := 1.5; export func() { return a + 5.0 }`), 6.5)
	// export function that read local variable
	expectRun(t, `out = import("mod1")()`,
		Opts().Module("mod1", `export func() { a := 1.5; return a + 5.0 }`), 6.5)
	// export function that read free variables
	expectRun(t, `out = import("mod1")()`,
		Opts().Module("mod1", `export func() { a := 1.5; return func() { return a + 5.0 }() }`), 6.5)

	// recursive function in module
	expectRun(t, `out = import("mod1")`,
		Opts().Module(
			"mod1", `
a := func(x) {
	return x == 0 ? 0 : x + a(x-1)
}

export a(5)
`), 15)
	expectRun(t, `out = import("mod1")`,
		Opts().Module(
			"mod1", `
export func() {
	a := func(x) {
		return x == 0 ? 0 : x + a(x-1)
	}

	return a(5)
}()
`), 15)

	// (main) -> mod1 -> mod2
	expectRun(t, `out = import("mod1")()`,
		Opts().Module("mod1", `export import("mod2")`).
			Module("mod2", `export func() { return 5.0 }`),
		5.0)
	// (main) -> mod1 -> mod2
	//        -> mod2
	expectRun(t, `import("mod1"); out = import("mod2")()`,
		Opts().Module("mod1", `export import("mod2")`).
			Module("mod2", `export func() { return 5.0 }`),
		5.0)
	// (main) -> mod1 -> mod2 -> mod3
	//        -> mod2 -> mod3
	expectRun(t, `import("mod1"); out = import("mod2")()`,
		Opts().Module("mod1", `export import("mod2")`).
			Module("mod2", `export import("mod3")`).
			Module("mod3", `export func() { return 5.0 }`),
		5.0)

	// cyclic imports
	// (main) -> mod1 -> mod2 -> mod1
	expectError(t, `import("mod1")`,
		Opts().Module("mod1", `import("mod2")`).
			Module("mod2", `import("mod1")`),
		"Compile Error: cyclic module import: mod1\n\tat mod2:1:1")
	// (main) -> mod1 -> mod2 -> mod3 -> mod1
	expectError(t, `import("mod1")`,
		Opts().Module("mod1", `import("mod2")`).
			Module("mod2", `import("mod3")`).
			Module("mod3", `import("mod1")`),
		"Compile Error: cyclic module import: mod1\n\tat mod3:1:1")
	// (main) -> mod1 -> mod2 -> mod3 -> mod2
	expectError(t, `import("mod1")`,
		Opts().Module("mod1", `import("mod2")`).
			Module("mod2", `import("mod3")`).
			Module("mod3", `import("mod2")`),
		"Compile Error: cyclic module import: mod2\n\tat mod3:1:1")

	// unknown modules
	expectError(t, `import("mod0")`,
		Opts().Module("mod1", `a := 5`), "module 'mod0' not found")
	expectError(t, `import("mod1")`,
		Opts().Module("mod1", `import("mod2")`), "module 'mod2' not found")

	// module is immutable but its variables is not necessarily immutable.
	expectRun(t, `m1 := import("mod1"); m1.a.b = 5; out = m1.a.b`,
		Opts().Module("mod1", `export {a: {b: 3}}`),
		5)

	// make sure module has same builtin functions
	expectRun(t, `out = import("mod1")`,
		Opts().Module("mod1", `export func() { return type_name(0) }()`),
		"int")

	// 'export' statement is ignored outside module
	expectRun(t, `a := 5; export func() { a = 10 }(); out = a`,
		Opts().Skip2ndPass(), 5)

	// 'export' must be in the top-level
	expectError(t, `import("mod1")`,
		Opts().Module("mod1", `func() { export 5 }()`),
		"Compile Error: export not allowed inside function\n\tat mod1:1:10")
	expectError(t, `import("mod1")`,
		Opts().Module("mod1", `func() { func() { export 5 }() }()`),
		"Compile Error: export not allowed inside function\n\tat mod1:1:19")

	// module cannot access outer scope
	expectError(t, `a := 5; import("mod1")`,
		Opts().Module("mod1", `export a`),
		"Compile Error: unresolved reference 'a'\n\tat mod1:1:8")

	// runtime error within modules
	expectError(t, `
a := 1;
b := import("mod1");
b(a)`,
		Opts().Module("mod1", `
export func(a) {
   a()
}
`), "Runtime Error: not callable: int\n\tat mod1:3:4\n\tat test:4:1")

	// module skipping export
	expectRun(t, `out = import("mod0")`,
		Opts().Module("mod0", ``), tengo.UndefinedValue)
	expectRun(t, `out = import("mod0")`,
		Opts().Module("mod0", `if 1 { export true }`), true)
	expectRun(t, `out = import("mod0")`,
		Opts().Module("mod0", `if 0 { export true }`),
		tengo.UndefinedValue)
	expectRun(t, `out = import("mod0")`,
		Opts().Module("mod0", `if 1 { } else { export true }`),
		tengo.UndefinedValue)
	expectRun(t, `out = import("mod0")`,
		Opts().Module("mod0", `for v:=0;;v++ { if v == 3 { export true } } }`),
		true)
	expectRun(t, `out = import("mod0")`,
		Opts().Module("mod0", `for v:=0;;v++ { if v == 3 { break } } }`),
		tengo.UndefinedValue)

	// duplicate compiled functions
	// NOTE: module "mod" has a function with some local variable, and it's
	//  imported twice by the main script. That causes the same CompiledFunction
	//  put in constants twice and the Bytecode optimization (removing duplicate
	//  constants) should still work correctly.
	expectRun(t, `
m1 := import("mod")
m2 := import("mod")
out = m1.x
	`,
		Opts().Module("mod", `
f1 := func(a, b) {
	c := a + b + 1
	return a + b + 1
}
export { x: 1 }
`),
		1)
}

func TestModuleBlockScopes(t *testing.T) {
	m := Opts().Module("rand",
		&tengo.BuiltinModule{
			Attrs: map[string]tengo.Object{
				"intn": &tengo.UserFunction{
					Name: "abs",
					Value: func(a ...tengo.Object) (tengo.Object, error) {
						v, _ := tengo.ToInt64(a[0])
						return &tengo.Int{Value: rand.Int63n(v)}, nil
					},
				},
			},
		})

	// block scopes in module
	expectRun(t, `out = import("mod1")()`, m.Module(
		"mod1", `
	rand := import("rand")
	foo := func() { return 1 }
	export func() {
		rand.intn(3)
		return foo()
	}`), 1)

	expectRun(t, `out = import("mod1")()`, m.Module(
		"mod1", `
rand := import("rand")
foo := func() { return 1 }
export func() {
	rand.intn(3)
	if foo() {}
	return 10
}
`), 10)

	expectRun(t, `out = import("mod1")()`, m.Module(
		"mod1", `
	rand := import("rand")
	foo := func() { return 1 }
	export func() {
		rand.intn(3)
		if true { foo() }
		return 10
	}
	`), 10)
}

func TestBangOperator(t *testing.T) {
	expectRun(t, `out = !true`, nil, false)
	expectRun(t, `out = !false`, nil, true)
	expectRun(t, `out = !0`, nil, true)
	expectRun(t, `out = !5`, nil, false)
	expectRun(t, `out = !!true`, nil, true)
	expectRun(t, `out = !!false`, nil, false)
	expectRun(t, `out = !!5`, nil, true)
}

func TestObjectsLimit(t *testing.T) {
	testAllocsLimit(t, `5`, 0)
	testAllocsLimit(t, `5 + 5`, 1)
	testAllocsLimit(t, `a := [1, 2, 3]`, 1)
	testAllocsLimit(t, `a := 1; b := 2; c := 3; d := [a, b, c]`, 1)
	testAllocsLimit(t, `a := {foo: 1, bar: 2}`, 1)
	testAllocsLimit(t, `a := 1; b := 2; c := {foo: a, bar: b}`, 1)
	testAllocsLimit(t, `
f := func() {
	return 5 + 5
}
a := f() + 5
`, 2)
	testAllocsLimit(t, `
f := func() {
	return 5 + 5
}
a := f()
`, 1)
	testAllocsLimit(t, `
a := []
f := func() {
	a = append(a, 5)
}
f()
f()
f()
`, 4)
}

func testAllocsLimit(t *testing.T, src string, limit int64) {
	expectRun(t, src,
		Opts().Skip2ndPass(), tengo.UndefinedValue) // no limit
	expectRun(t, src,
		Opts().MaxAllocs(limit).Skip2ndPass(), tengo.UndefinedValue)
	expectRun(t, src,
		Opts().MaxAllocs(limit+1).Skip2ndPass(), tengo.UndefinedValue)
	if limit > 1 {
		expectError(t, src,
			Opts().MaxAllocs(limit-1).Skip2ndPass(),
			"allocation limit exceeded")
	}
	if limit > 2 {
		expectError(t, src,
			Opts().MaxAllocs(limit-2).Skip2ndPass(),
			"allocation limit exceeded")
	}
}

func TestReturn(t *testing.T) {
	expectRun(t, `out = func() { return 10; }()`, nil, 10)
	expectRun(t, `out = func() { return 10; return 9; }()`, nil, 10)
	expectRun(t, `out = func() { return 2 * 5; return 9 }()`, nil, 10)
	expectRun(t, `out = func() { 9; return 2 * 5; return 9 }()`, nil, 10)
	expectRun(t, `
	out = func() { 
		if (10 > 1) {
			if (10 > 1) {
				return 10;
	  		}

	  		return 1;
		}
	}()`, nil, 10)

	expectRun(t, `f1 := func() { return 2 * 5; }; out = f1()`, nil, 10)
}

func TestVMScopes(t *testing.T) {
	// shadowed global variable
	expectRun(t, `
c := 5
if a := 3; a {
	c := 6
} else {
	c := 7
}
out = c
`, nil, 5)

	// shadowed local variable
	expectRun(t, `
func() {
	c := 5
	if a := 3; a {
		c := 6
	} else {
		c := 7
	}
	out = c
}()
`, nil, 5)

	// 'b' is declared in 2 separate blocks
	expectRun(t, `
c := 5
if a := 3; a {
	b := 8
	c = b
} else {
	b := 9
	c = b
}
out = c
`, nil, 8)

	// shadowing inside for statement
	expectRun(t, `
a := 4
b := 5
for i:=0;i<3;i++ {
	b := 6
	for j:=0;j<2;j++ {
		b := 7
		a = i*j
	}
}
out = a`, nil, 2)

	// shadowing variable declared in init statement
	expectRun(t, `
if a := 5; a {
	a := 6
	out = a
}`, nil, 6)
	expectRun(t, `
a := 4
if a := 5; a {
	a := 6
	out = a
}`, nil, 6)
	expectRun(t, `
a := 4
if a := 0; a {
	a := 6
	out = a
} else {
	a := 7
	out = a
}`, nil, 7)
	expectRun(t, `
a := 4
if a := 0; a {
	out = a
} else {
	out = a
}`, nil, 0)

	// shadowing function level
	expectRun(t, `
a := 5
func() {
	a := 6
	a = 7
}()
out = a
`, nil, 5)
	expectRun(t, `
a := 5
func() {
	if a := 7; true {
		a = 8
	}
}()
out = a
`, nil, 5)
}

func TestSelector(t *testing.T) {
	expectRun(t, `a := {k1: 5, k2: "foo"}; out = a.k1`,
		nil, 5)
	expectRun(t, `a := {k1: 5, k2: "foo"}; out = a.k2`,
		nil, "foo")
	expectRun(t, `a := {k1: 5, k2: "foo"}; out = a.k3`,
		nil, tengo.UndefinedValue)

	expectRun(t, `
a := {
	b: {
		c: 4,
		a: false
	},
	c: "foo bar"
}
out = a.b.c`, nil, 4)

	expectRun(t, `
a := {
	b: {
		c: 4,
		a: false
	},
	c: "foo bar"
}
b := a.x.c`, nil, tengo.UndefinedValue)

	expectRun(t, `
a := {
	b: {
		c: 4,
		a: false
	},
	c: "foo bar"
}
b := a.x.y`, nil, tengo.UndefinedValue)

	expectRun(t, `a := {b: 1, c: "foo"}; a.b = 2; out = a.b`,
		nil, 2)
	expectRun(t, `a := {b: 1, c: "foo"}; a.c = 2; out = a.c`,
		nil, 2) // type not checked on sub-field
	expectRun(t, `a := {b: {c: 1}}; a.b.c = 2; out = a.b.c`,
		nil, 2)
	expectRun(t, `a := {b: 1}; a.c = 2; out = a`,
		nil, MAP{"b": 1, "c": 2})
	expectRun(t, `a := {b: {c: 1}}; a.b.d = 2; out = a`,
		nil, MAP{"b": MAP{"c": 1, "d": 2}})

	expectRun(t, `func() { a := {b: 1, c: "foo"}; a.b = 2; out = a.b }()`,
		nil, 2)
	expectRun(t, `func() { a := {b: 1, c: "foo"}; a.c = 2; out = a.c }()`,
		nil, 2) // type not checked on sub-field
	expectRun(t, `func() { a := {b: {c: 1}}; a.b.c = 2; out = a.b.c }()`,
		nil, 2)
	expectRun(t, `func() { a := {b: 1}; a.c = 2; out = a }()`,
		nil, MAP{"b": 1, "c": 2})
	expectRun(t, `func() { a := {b: {c: 1}}; a.b.d = 2; out = a }()`,
		nil, MAP{"b": MAP{"c": 1, "d": 2}})

	expectRun(t, `func() { a := {b: 1, c: "foo"}; func() { a.b = 2 }(); out = a.b }()`,
		nil, 2)
	expectRun(t, `func() { a := {b: 1, c: "foo"}; func() { a.c = 2 }(); out = a.c }()`,
		nil, 2) // type not checked on sub-field
	expectRun(t, `func() { a := {b: {c: 1}}; func() { a.b.c = 2 }(); out = a.b.c }()`,
		nil, 2)
	expectRun(t, `func() { a := {b: 1}; func() { a.c = 2 }(); out = a }()`,
		nil, MAP{"b": 1, "c": 2})
	expectRun(t, `func() { a := {b: {c: 1}}; func() { a.b.d = 2 }(); out = a }()`,
		nil, MAP{"b": MAP{"c": 1, "d": 2}})

	expectRun(t, `
a := {
	b: [1, 2, 3],
	c: {
		d: 8,
		e: "foo",
		f: [9, 8]
	}
}
out = [a.b[2], a.c.d, a.c.e, a.c.f[1]]
`, nil, ARR{3, 8, "foo", 8})

	expectRun(t, `
func() {
	a := [1, 2, 3]
	b := 9
	a[1] = b
	b = 7     // make sure a[1] has a COPY of value of 'b'
	out = a[1]
}()
`, nil, 9)

	expectError(t, `a := {b: {c: 1}}; a.d.c = 2`,
		nil, "not index-assignable")
	expectError(t, `a := [1, 2, 3]; a.b = 2`,
		nil, "invalid index type")
	expectError(t, `a := "foo"; a.b = 2`,
		nil, "not index-assignable")
	expectError(t, `func() { a := {b: {c: 1}}; a.d.c = 2 }()`,
		nil, "not index-assignable")
	expectError(t, `func() { a := [1, 2, 3]; a.b = 2 }()`,
		nil, "invalid index type")
	expectError(t, `func() { a := "foo"; a.b = 2 }()`,
		nil, "not index-assignable")
}

func TestSourceModules(t *testing.T) {
	testEnumModule(t, `out = enum.key(0, 20)`, 0)
	testEnumModule(t, `out = enum.key(10, 20)`, 10)
	testEnumModule(t, `out = enum.value(0, 0)`, 0)
	testEnumModule(t, `out = enum.value(10, 20)`, 20)

	testEnumModule(t, `out = enum.all([], enum.value)`, true)
	testEnumModule(t, `out = enum.all([1], enum.value)`, true)
	testEnumModule(t, `out = enum.all([true, 1], enum.value)`, true)
	testEnumModule(t, `out = enum.all([true, 0], enum.value)`, false)
	testEnumModule(t, `out = enum.all([true, 0, 1], enum.value)`, false)
	testEnumModule(t, `out = enum.all(immutable([true, 0, 1]), enum.value)`,
		false) // immutable-array
	testEnumModule(t, `out = enum.all({}, enum.value)`, true)
	testEnumModule(t, `out = enum.all({a:1}, enum.value)`, true)
	testEnumModule(t, `out = enum.all({a:true, b:1}, enum.value)`, true)
	testEnumModule(t, `out = enum.all(immutable({a:true, b:1}), enum.value)`,
		true) // immutable-map
	testEnumModule(t, `out = enum.all({a:true, b:0}, enum.value)`, false)
	testEnumModule(t, `out = enum.all({a:true, b:0, c:1}, enum.value)`, false)
	testEnumModule(t, `out = enum.all(0, enum.value)`,
		tengo.UndefinedValue) // non-enumerable: undefined
	testEnumModule(t, `out = enum.all("123", enum.value)`,
		tengo.UndefinedValue) // non-enumerable: undefined

	testEnumModule(t, `out = enum.any([], enum.value)`, false)
	testEnumModule(t, `out = enum.any([1], enum.value)`, true)
	testEnumModule(t, `out = enum.any([true, 1], enum.value)`, true)
	testEnumModule(t, `out = enum.any([true, 0], enum.value)`, true)
	testEnumModule(t, `out = enum.any([true, 0, 1], enum.value)`, true)
	testEnumModule(t, `out = enum.any(immutable([true, 0, 1]), enum.value)`,
		true) // immutable-array
	testEnumModule(t, `out = enum.any([false], enum.value)`, false)
	testEnumModule(t, `out = enum.any([false, 0], enum.value)`, false)
	testEnumModule(t, `out = enum.any({}, enum.value)`, false)
	testEnumModule(t, `out = enum.any({a:1}, enum.value)`, true)
	testEnumModule(t, `out = enum.any({a:true, b:1}, enum.value)`, true)
	testEnumModule(t, `out = enum.any({a:true, b:0}, enum.value)`, true)
	testEnumModule(t, `out = enum.any({a:true, b:0, c:1}, enum.value)`, true)
	testEnumModule(t, `out = enum.any(immutable({a:true, b:0, c:1}), enum.value)`,
		true) // immutable-map
	testEnumModule(t, `out = enum.any({a:false}, enum.value)`, false)
	testEnumModule(t, `out = enum.any({a:false, b:0}, enum.value)`, false)
	testEnumModule(t, `out = enum.any(0, enum.value)`,
		tengo.UndefinedValue) // non-enumerable: undefined
	testEnumModule(t, `out = enum.any("123", enum.value)`,
		tengo.UndefinedValue) // non-enumerable: undefined

	testEnumModule(t, `out = enum.chunk([], 1)`, ARR{})
	testEnumModule(t, `out = enum.chunk([1], 1)`, ARR{ARR{1}})
	testEnumModule(t, `out = enum.chunk([1,2,3], 1)`,
		ARR{ARR{1}, ARR{2}, ARR{3}})
	testEnumModule(t, `out = enum.chunk([1,2,3], 2)`,
		ARR{ARR{1, 2}, ARR{3}})
	testEnumModule(t, `out = enum.chunk([1,2,3], 3)`,
		ARR{ARR{1, 2, 3}})
	testEnumModule(t, `out = enum.chunk([1,2,3], 4)`,
		ARR{ARR{1, 2, 3}})
	testEnumModule(t, `out = enum.chunk([1,2,3,4], 3)`,
		ARR{ARR{1, 2, 3}, ARR{4}})
	testEnumModule(t, `out = enum.chunk([], 0)`,
		tengo.UndefinedValue) // size=0: undefined
	testEnumModule(t, `out = enum.chunk([1], 0)`,
		tengo.UndefinedValue) // size=0: undefined
	testEnumModule(t, `out = enum.chunk([1,2,3], 0)`,
		tengo.UndefinedValue) // size=0: undefined
	testEnumModule(t, `out = enum.chunk({a:1,b:2,c:3}, 1)`,
		tengo.UndefinedValue) // map: undefined
	testEnumModule(t, `out = enum.chunk(0, 1)`,
		tengo.UndefinedValue) // non-enumerable: undefined
	testEnumModule(t, `out = enum.chunk("123", 1)`,
		tengo.UndefinedValue) // non-enumerable: undefined

	testEnumModule(t, `out = enum.at([], 0)`,
		tengo.UndefinedValue)
	testEnumModule(t, `out = enum.at([], 1)`,
		tengo.UndefinedValue)
	testEnumModule(t, `out = enum.at([], -1)`,
		tengo.UndefinedValue)
	testEnumModule(t, `out = enum.at(["one"], 0)`,
		"one")
	testEnumModule(t, `out = enum.at(["one"], 1)`,
		tengo.UndefinedValue)
	testEnumModule(t, `out = enum.at(["one"], -1)`,
		tengo.UndefinedValue)
	testEnumModule(t, `out = enum.at(["one","two","three"], 0)`,
		"one")
	testEnumModule(t, `out = enum.at(["one","two","three"], 1)`,
		"two")
	testEnumModule(t, `out = enum.at(["one","two","three"], 2)`,
		"three")
	testEnumModule(t, `out = enum.at(["one","two","three"], -1)`,
		tengo.UndefinedValue)
	testEnumModule(t, `out = enum.at(["one","two","three"], 3)`,
		tengo.UndefinedValue)
	testEnumModule(t, `out = enum.at(["one","two","three"], "1")`,
		tengo.UndefinedValue) // non-int index: undefined
	testEnumModule(t, `out = enum.at({}, "a")`,
		tengo.UndefinedValue)
	testEnumModule(t, `out = enum.at({a:"one"}, "a")`,
		"one")
	testEnumModule(t, `out = enum.at({a:"one"}, "b")`,
		tengo.UndefinedValue)
	testEnumModule(t, `out = enum.at({a:"one",b:"two",c:"three"}, "a")`,
		"one")
	testEnumModule(t, `out = enum.at({a:"one",b:"two",c:"three"}, "b")`,
		"two")
	testEnumModule(t, `out = enum.at({a:"one",b:"two",c:"three"}, "c")`,
		"three")
	testEnumModule(t, `out = enum.at({a:"one",b:"two",c:"three"}, "d")`,
		tengo.UndefinedValue)
	testEnumModule(t, `out = enum.at({a:"one",b:"two",c:"three"}, 'a')`,
		tengo.UndefinedValue) // non-string index: undefined
	testEnumModule(t, `out = enum.at(0, 1)`,
		tengo.UndefinedValue) // non-enumerable: undefined
	testEnumModule(t, `out = enum.at("abc", 1)`,
		tengo.UndefinedValue) // non-enumerable: undefined

	testEnumModule(t, `out=0; enum.each([],func(k,v){out+=v})`, 0)
	testEnumModule(t, `out=0; enum.each([1,2,3],func(k,v){out+=v})`, 6)
	testEnumModule(t, `out=0; enum.each([1,2,3],func(k,v){out+=k})`, 3)
	testEnumModule(t, `out=0; enum.each({a:1,b:2,c:3},func(k,v){out+=v})`, 6)
	testEnumModule(t, `out=""; enum.each({a:1,b:2,c:3},func(k,v){out+=k}); out=len(out)`,
		3)
	testEnumModule(t, `out=0; enum.each(5,func(k,v){out+=v})`, 0)     // non-enumerable: no iteration
	testEnumModule(t, `out=0; enum.each("123",func(k,v){out+=v})`, 0) // non-enumerable: no iteration

	testEnumModule(t, `out = enum.filter([], enum.value)`,
		ARR{})
	testEnumModule(t, `out = enum.filter([false,1,2], enum.value)`,
		ARR{1, 2})
	testEnumModule(t, `out = enum.filter([false,1,0,2], enum.value)`,
		ARR{1, 2})
	testEnumModule(t, `out = enum.filter({}, enum.value)`,
		tengo.UndefinedValue) // non-array: undefined
	testEnumModule(t, `out = enum.filter(0, enum.value)`,
		tengo.UndefinedValue) // non-array: undefined
	testEnumModule(t, `out = enum.filter("123", enum.value)`,
		tengo.UndefinedValue) // non-array: undefined

	testEnumModule(t, `out = enum.find([], enum.value)`,
		tengo.UndefinedValue)
	testEnumModule(t, `out = enum.find([0], enum.value)`,
		tengo.UndefinedValue)
	testEnumModule(t, `out = enum.find([1], enum.value)`, 1)
	testEnumModule(t, `out = enum.find([false,0,undefined,1], enum.value)`, 1)
	testEnumModule(t, `out = enum.find([1,2,3], enum.value)`, 1)
	testEnumModule(t, `out = enum.find({}, enum.value)`,
		tengo.UndefinedValue)
	testEnumModule(t, `out = enum.find({a:0}, enum.value)`,
		tengo.UndefinedValue)
	testEnumModule(t, `out = enum.find({a:1}, enum.value)`, 1)
	testEnumModule(t, `out = enum.find({a:false,b:0,c:undefined,d:1}, enum.value)`,
		1)
	//testEnumModule(t, `out = enum.find({a:1,b:2,c:3}, enum.value)`, 1)
	testEnumModule(t, `out = enum.find(0, enum.value)`,
		tengo.UndefinedValue) // non-enumerable: undefined
	testEnumModule(t, `out = enum.find("123", enum.value)`,
		tengo.UndefinedValue) // non-enumerable: undefined

	testEnumModule(t, `out = enum.find_key([], enum.value)`,
		tengo.UndefinedValue)
	testEnumModule(t, `out = enum.find_key([0], enum.value)`,
		tengo.UndefinedValue)
	testEnumModule(t, `out = enum.find_key([1], enum.value)`, 0)
	testEnumModule(t, `out = enum.find_key([false,0,undefined,1], enum.value)`,
		3)
	testEnumModule(t, `out = enum.find_key([1,2,3], enum.value)`, 0)
	testEnumModule(t, `out = enum.find_key({}, enum.value)`,
		tengo.UndefinedValue)
	testEnumModule(t, `out = enum.find_key({a:0}, enum.value)`,
		tengo.UndefinedValue)
	testEnumModule(t, `out = enum.find_key({a:1}, enum.value)`,
		"a")
	testEnumModule(t, `out = enum.find_key({a:false,b:0,c:undefined,d:1}, enum.value)`,
		"d")
	//testEnumModule(t, `out = enum.find_key({a:1,b:2,c:3}, enum.value)`, "a")
	testEnumModule(t, `out = enum.find_key(0, enum.value)`,
		tengo.UndefinedValue) // non-enumerable: undefined
	testEnumModule(t, `out = enum.find_key("123", enum.value)`,
		tengo.UndefinedValue) // non-enumerable: undefined

	testEnumModule(t, `out = enum.map([], enum.value)`,
		ARR{})
	testEnumModule(t, `out = enum.map([1,2,3], enum.value)`,
		ARR{1, 2, 3})
	testEnumModule(t, `out = enum.map([1,2,3], enum.key)`,
		ARR{0, 1, 2})
	testEnumModule(t, `out = enum.map([1,2,3], func(k,v) { return v*2 })`,
		ARR{2, 4, 6})
	testEnumModule(t, `out = enum.map({}, enum.value)`,
		ARR{})
	testEnumModule(t, `out = enum.map({a:1}, func(k,v) { return v*2 })`,
		ARR{2})
	testEnumModule(t, `out = enum.map(0, enum.value)`,
		tengo.UndefinedValue) // non-enumerable: undefined
	testEnumModule(t, `out = enum.map("123", enum.value)`,
		tengo.UndefinedValue) // non-enumerable: undefined
}

func testEnumModule(t *testing.T, input string, expected interface{}) {
	expectRun(t, `enum := import("enum"); `+input,
		Opts().Module("enum", stdlib.SourceModules["enum"]),
		expected)
}

func TestSrcModEnum(t *testing.T) {
	expectRun(t, `
x := import("enum")
out = x.all([1, 2, 3], func(_, v) { return v >= 1 }) 
`, Opts().Stdlib(), true)
	expectRun(t, `
x := import("enum")
out = x.all([1, 2, 3], func(_, v) { return v >= 2 }) 
`, Opts().Stdlib(), false)

	expectRun(t, `
x := import("enum")
out = x.any([1, 2, 3], func(_, v) { return v >= 1 }) 
`, Opts().Stdlib(), true)
	expectRun(t, `
x := import("enum")
out = x.any([1, 2, 3], func(_, v) { return v >= 2 }) 
`, Opts().Stdlib(), true)

	expectRun(t, `
x := import("enum")
out = x.chunk([1, 2, 3], 1) 
`, Opts().Stdlib(), ARR{ARR{1}, ARR{2}, ARR{3}})
	expectRun(t, `
x := import("enum")
out = x.chunk([1, 2, 3], 2) 
`, Opts().Stdlib(), ARR{ARR{1, 2}, ARR{3}})
	expectRun(t, `
x := import("enum")
out = x.chunk([1, 2, 3], 3) 
`, Opts().Stdlib(), ARR{ARR{1, 2, 3}})
	expectRun(t, `
x := import("enum")
out = x.chunk([1, 2, 3], 4) 
`, Opts().Stdlib(), ARR{ARR{1, 2, 3}})
	expectRun(t, `
x := import("enum")
out = x.chunk([1, 2, 3, 4, 5, 6], 2) 
`, Opts().Stdlib(), ARR{ARR{1, 2}, ARR{3, 4}, ARR{5, 6}})

	expectRun(t, `
x := import("enum")
out = x.at([1, 2, 3], 0) 
`, Opts().Stdlib(), 1)
}

func TestVMStackOverflow(t *testing.T) {
	expectError(t, `f := func() { return f() + 1 }; f()`,
		nil, "stack overflow")
}

func TestString(t *testing.T) {
	expectRun(t, `out = "Hello World!"`, nil, "Hello World!")
	expectRun(t, `out = "Hello" + " " + "World!"`, nil, "Hello World!")

	expectRun(t, `out = "Hello" == "Hello"`, nil, true)
	expectRun(t, `out = "Hello" == "World"`, nil, false)
	expectRun(t, `out = "Hello" != "Hello"`, nil, false)
	expectRun(t, `out = "Hello" != "World"`, nil, true)

	expectRun(t, `out = "Hello" > "World"`, nil, false)
	expectRun(t, `out = "World" < "Hello"`, nil, false)
	expectRun(t, `out = "Hello" < "World"`, nil, true)
	expectRun(t, `out = "World" > "Hello"`, nil, true)
	expectRun(t, `out = "Hello" >= "World"`, nil, false)
	expectRun(t, `out = "Hello" <= "World"`, nil, true)
	expectRun(t, `out = "Hello" >= "Hello"`, nil, true)
	expectRun(t, `out = "World" <= "World"`, nil, true)

	// index operator
	str := "abcdef"
	strStr := `"abcdef"`
	strLen := 6
	for idx := 0; idx < strLen; idx++ {
		expectRun(t, fmt.Sprintf("out = %s[%d]", strStr, idx),
			nil, str[idx])
		expectRun(t, fmt.Sprintf("out = %s[0 + %d]", strStr, idx),
			nil, str[idx])
		expectRun(t, fmt.Sprintf("out = %s[1 + %d - 1]", strStr, idx),
			nil, str[idx])
		expectRun(t, fmt.Sprintf("idx := %d; out = %s[idx]", idx, strStr),
			nil, str[idx])
	}

	expectRun(t, fmt.Sprintf("%s[%d]", strStr, -1),
		nil, tengo.UndefinedValue)
	expectRun(t, fmt.Sprintf("%s[%d]", strStr, strLen),
		nil, tengo.UndefinedValue)

	// slice operator
	for low := 0; low <= strLen; low++ {
		expectRun(t, fmt.Sprintf("out = %s[%d:%d]", strStr, low, low),
			nil, "")
		for high := low; high <= strLen; high++ {
			expectRun(t, fmt.Sprintf("out = %s[%d:%d]", strStr, low, high),
				nil, str[low:high])
			expectRun(t,
				fmt.Sprintf("out = %s[0 + %d : 0 + %d]", strStr, low, high),
				nil, str[low:high])
			expectRun(t,
				fmt.Sprintf("out = %s[1 + %d - 1 : 1 + %d - 1]",
					strStr, low, high),
				nil, str[low:high])
			expectRun(t,
				fmt.Sprintf("out = %s[:%d]", strStr, high),
				nil, str[:high])
			expectRun(t,
				fmt.Sprintf("out = %s[%d:]", strStr, low),
				nil, str[low:])
		}
	}

	expectRun(t, fmt.Sprintf("out = %s[:]", strStr),
		nil, str[:])
	expectRun(t, fmt.Sprintf("out = %s[:]", strStr),
		nil, str)
	expectRun(t, fmt.Sprintf("out = %s[%d:]", strStr, -1),
		nil, str)
	expectRun(t, fmt.Sprintf("out = %s[:%d]", strStr, strLen+1),
		nil, str)
	expectRun(t, fmt.Sprintf("out = %s[%d:%d]", strStr, 2, 2),
		nil, "")

	expectError(t, fmt.Sprintf("%s[:%d]", strStr, -1),
		nil, "invalid slice index")
	expectError(t, fmt.Sprintf("%s[%d:]", strStr, strLen+1),
		nil, "invalid slice index")
	expectError(t, fmt.Sprintf("%s[%d:%d]", strStr, 0, -1),
		nil, "invalid slice index")
	expectError(t, fmt.Sprintf("%s[%d:%d]", strStr, 2, 1),
		nil, "invalid slice index")

	// string concatenation with other types
	expectRun(t, `out = "foo" + 1`, nil, "foo1")
	// Float.String() returns the smallest number of digits
	// necessary such that ParseFloat will return f exactly.
	expectRun(t, `out = "foo" + 1.0`, nil, "foo1") // <- note '1' instead of '1.0'
	expectRun(t, `out = "foo" + 1.5`, nil, "foo1.5")
	expectRun(t, `out = "foo" + true`, nil, "footrue")
	expectRun(t, `out = "foo" + 'X'`, nil, "fooX")
	expectRun(t, `out = "foo" + error(5)`, nil, "fooerror: 5")
	expectRun(t, `out = "foo" + undefined`, nil, "foo<undefined>")
	expectRun(t, `out = "foo" + [1,2,3]`, nil, "foo[1, 2, 3]")
	// also works with "+=" operator
	expectRun(t, `out = "foo"; out += 1.5`, nil, "foo1.5")
	// string concats works only when string is LHS
	expectError(t, `1 + "foo"`, nil, "invalid operation")

	expectError(t, `"foo" - "bar"`, nil, "invalid operation")
}

func TestTailCall(t *testing.T) {
	expectRun(t, `
	fac := func(n, a) {
		if n == 1 {
			return a
		}
		return fac(n-1, n*a)
	}
	out = fac(5, 1)`, nil, 120)

	expectRun(t, `
	fac := func(n, a) {
		if n == 1 {
			return a
		}
		x := {foo: fac} // indirection for test
		return x.foo(n-1, n*a)
	}
	out = fac(5, 1)`, nil, 120)

	expectRun(t, `
	fib := func(x, s) {
		if x == 0 {
			return 0 + s
		} else if x == 1 {
			return 1 + s
		}
		return fib(x-1, fib(x-2, s))
	}
	out = fib(15, 0)`, nil, 610)

	expectRun(t, `
	fib := func(n, a, b) {
		if n == 0 {
			return a
		} else if n == 1 {
			return b
		}
		return fib(n-1, b, a + b)
	}
	out = fib(15, 0, 1)`, nil, 610)

	// global variable and no return value
	expectRun(t, `
			out = 0
			foo := func(a) {
			   if a == 0 {
			       return
			   }
			   out += a
			   foo(a-1)
			}
			foo(10)`, nil, 55)

	expectRun(t, `
	f1 := func() {
		f2 := 0    // TODO: this might be fixed in the future
		f2 = func(n, s) {
			if n == 0 { return s }
			return f2(n-1, n + s)
		}
		return f2(5, 0)
	}
	out = f1()`, nil, 15)

	// tail-call replacing loop
	// without tail-call optimization, this code will cause stack overflow
	expectRun(t, `
iter := func(n, max) {
	if n == max {
		return n
	}

	return iter(n+1, max)
}
out = iter(0, 9999)
`, nil, 9999)
	expectRun(t, `
c := 0
iter := func(n, max) {
	if n == max {
		return
	}

	c++
	iter(n+1, max)
}
iter(0, 9999)
out = c 
`, nil, 9999)
}

// tail call with free vars
func TestTailCallFreeVars(t *testing.T) {
	expectRun(t, `
func() {
	a := 10
	f2 := 0
	f2 = func(n, s) {
		if n == 0 {
			return s + a
		}
		return f2(n-1, n+s)
	}
	out = f2(5, 0)
}()`, nil, 25)
}

func TestSpread(t *testing.T) {
	expectRun(t, `
	f := func(...a) {
		return append(a, 3)
	}
	out = f([1, 2]...)
	`, nil, ARR{1, 2, 3})

	expectRun(t, `
	f := func(a, ...b) {
		return append([a], append(b, 3)...)
	}
	out = f([1, 2]...)
	`, nil, ARR{1, 2, 3})

	expectRun(t, `
	f := func(a, ...b) {
		return append(append([a], b), 3)
	}
	out = f(1, [2]...)
	`, nil, ARR{1, ARR{2}, 3})

	expectRun(t, `
	f1 := func(...a){
		return append([3], a...)
	}
	f2 := func(a, ...b) {
		return f1(append([a], b...)...)
	}
	out = f2([1, 2]...)
	`, nil, ARR{3, 1, 2})

	expectRun(t, `
	f := func(a, ...b) {
		return func(...a) {
			return append([3], append(a, 4)...)
		}(a, b...)
	}
	out = f([1, 2]...)
	`, nil, ARR{3, 1, 2, 4})

	expectRun(t, `
	f := func(a, ...b) {
		c := append(b, 4)
		return func(){
			return append(append([a], b...), c...)
		}()
	}
	out = f(1, immutable([2, 3])...)
	`, nil, ARR{1, 2, 3, 2, 3, 4})

	expectError(t, `func(a) {}([1, 2]...)`, nil,
		"Runtime Error: wrong number of arguments: want=1, got=2")
	expectError(t, `func(a, b, c) {}([1, 2]...)`, nil,
		"Runtime Error: wrong number of arguments: want=3, got=2")
}

func expectRun(
	t *testing.T,
	input string,
	opts *testopts,
	expected interface{},
) {
	if opts == nil {
		opts = Opts()
	}

	symbols := opts.symbols
	modules := opts.modules
	maxAllocs := opts.maxAllocs

	expectedObj := toObject(expected)

	if symbols == nil {
		symbols = make(map[string]tengo.Object)
	}
	symbols[testOut] = objectZeroCopy(expectedObj)

	// first pass: run the code normally
	{
		// parse
		file := parse(t, input)
		if file == nil {
			return
		}

		// compiler/VM
		res, trace, err := traceCompileRun(file, symbols, modules, maxAllocs)
		require.NoError(t, err, "\n"+strings.Join(trace, "\n"))
		require.Equal(t, expectedObj, res[testOut],
			"\n"+strings.Join(trace, "\n"))
	}

	// second pass: run the code as import module
	if !opts.skip2ndPass {
		file := parse(t, `out = import("__code__")`)
		if file == nil {
			return
		}

		expectedObj := toObject(expected)
		switch eo := expectedObj.(type) {
		case *tengo.Array:
			expectedObj = &tengo.ImmutableArray{Value: eo.Value}
		case *tengo.Map:
			expectedObj = &tengo.ImmutableMap{Value: eo.Value}
		}

		modules.AddSourceModule("__code__",
			[]byte(fmt.Sprintf("out := undefined; %s; export out", input)))

		res, trace, err := traceCompileRun(file, symbols, modules, maxAllocs)
		require.NoError(t, err, "\n"+strings.Join(trace, "\n"))
		require.Equal(t, expectedObj, res[testOut],
			"\n"+strings.Join(trace, "\n"))
	}
}

func expectError(
	t *testing.T,
	input string,
	opts *testopts,
	expected string,
) {
	if opts == nil {
		opts = Opts()
	}

	symbols := opts.symbols
	modules := opts.modules
	maxAllocs := opts.maxAllocs

	expected = strings.TrimSpace(expected)
	if expected == "" {
		panic("expected must not be empty")
	}

	// parse
	program := parse(t, input)
	if program == nil {
		return
	}

	// compiler/VM
	_, trace, err := traceCompileRun(program, symbols, modules, maxAllocs)
	require.Error(t, err, "\n"+strings.Join(trace, "\n"))
	require.True(t, strings.Contains(err.Error(), expected),
		"expected error string: %s, got: %s\n%s",
		expected, err.Error(), strings.Join(trace, "\n"))
}

func expectErrorIs(
	t *testing.T,
	input string,
	opts *testopts,
	expected error,
) {
	if opts == nil {
		opts = Opts()
	}
	symbols := opts.symbols
	modules := opts.modules
	maxAllocs := opts.maxAllocs

	// parse
	program := parse(t, input)
	if program == nil {
		return
	}

	// compiler/VM
	_, trace, err := traceCompileRun(program, symbols, modules, maxAllocs)
	require.Error(t, err, "\n"+strings.Join(trace, "\n"))
	require.True(t, errors.Is(err, expected),
		"expected error is: %s, got: %s\n%s",
		expected.Error(), err.Error(), strings.Join(trace, "\n"))
}

func expectErrorAs(
	t *testing.T,
	input string,
	opts *testopts,
	expected interface{},
) {
	if opts == nil {
		opts = Opts()
	}
	symbols := opts.symbols
	modules := opts.modules
	maxAllocs := opts.maxAllocs

	// parse
	program := parse(t, input)
	if program == nil {
		return
	}

	// compiler/VM
	_, trace, err := traceCompileRun(program, symbols, modules, maxAllocs)
	require.Error(t, err, "\n"+strings.Join(trace, "\n"))
	require.True(t, errors.As(err, expected),
		"expected error as: %v, got: %v\n%s",
		expected, err, strings.Join(trace, "\n"))
}

type vmTracer struct {
	Out []string
}

func (o *vmTracer) Write(p []byte) (n int, err error) {
	o.Out = append(o.Out, string(p))
	return len(p), nil
}

func traceCompileRun(
	file *parser.File,
	symbols map[string]tengo.Object,
	modules *tengo.ModuleMap,
	maxAllocs int64,
) (res map[string]tengo.Object, trace []string, err error) {
	var v *tengo.VM

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic: %v", e)

			// stack trace
			var stackTrace []string
			for i := 2; ; i += 1 {
				_, file, line, ok := _runtime.Caller(i)
				if !ok {
					break
				}
				stackTrace = append(stackTrace,
					fmt.Sprintf("  %s:%d", file, line))
			}

			trace = append(trace,
				fmt.Sprintf("[Error Trace]\n\n  %s\n",
					strings.Join(stackTrace, "\n  ")))
		}
	}()

	globals := make([]tengo.Object, tengo.GlobalsSize)

	symTable := tengo.NewSymbolTable()
	for name, value := range symbols {
		sym := symTable.Define(name)

		// should not store pointer to 'value' variable
		// which is re-used in each iteration.
		valueCopy := value
		globals[sym.Index] = valueCopy
	}
	for idx, fn := range tengo.GetAllBuiltinFunctions() {
		symTable.DefineBuiltin(idx, fn.Name)
	}

	tr := &vmTracer{}
	c := tengo.NewCompiler(file.InputFile, symTable, nil, modules, tr)
	err = c.Compile(file)
	trace = append(trace,
		fmt.Sprintf("\n[Compiler Trace]\n\n%s",
			strings.Join(tr.Out, "")))
	if err != nil {
		return
	}

	bytecode := c.Bytecode()
	bytecode.RemoveDuplicates()
	trace = append(trace, fmt.Sprintf("\n[Compiled Constants]\n\n%s",
		strings.Join(bytecode.FormatConstants(), "\n")))
	trace = append(trace, fmt.Sprintf("\n[Compiled Instructions]\n\n%s\n",
		strings.Join(bytecode.FormatInstructions(), "\n")))

	v = tengo.NewVM(bytecode, globals, maxAllocs)

	err = v.Run()
	{
		res = make(map[string]tengo.Object)
		for name := range symbols {
			sym, depth, ok := symTable.Resolve(name, false)
			if !ok || depth != 0 {
				err = fmt.Errorf("symbol not found: %s", name)
				return
			}

			res[name] = globals[sym.Index]
		}
		trace = append(trace, fmt.Sprintf("\n[Globals]\n\n%s",
			strings.Join(formatGlobals(globals), "\n")))
	}
	if err == nil && !v.IsStackEmpty() {
		err = errors.New("non empty stack after execution")
	}

	return
}

func formatGlobals(globals []tengo.Object) (formatted []string) {
	for idx, global := range globals {
		if global == nil {
			return
		}
		formatted = append(formatted, fmt.Sprintf("[% 3d] %s (%s|%p)",
			idx, global.String(), reflect.TypeOf(global).Elem().Name(), global))
	}
	return
}

func parse(t *testing.T, input string) *parser.File {
	testFileSet := parser.NewFileSet()
	testFile := testFileSet.AddFile("test", -1, len(input))

	p := parser.NewParser(testFile, []byte(input), nil)
	file, err := p.ParseFile()
	require.NoError(t, err)
	return file
}

func errorObject(v interface{}) *tengo.Error {
	return &tengo.Error{Value: toObject(v)}
}

func toObject(v interface{}) tengo.Object {
	switch v := v.(type) {
	case tengo.Object:
		return v
	case string:
		return &tengo.String{Value: v}
	case int64:
		return &tengo.Int{Value: v}
	case int: // for convenience
		return &tengo.Int{Value: int64(v)}
	case bool:
		if v {
			return tengo.TrueValue
		}
		return tengo.FalseValue
	case rune:
		return &tengo.Char{Value: v}
	case byte: // for convenience
		return &tengo.Char{Value: rune(v)}
	case float64:
		return &tengo.Float{Value: v}
	case []byte:
		return &tengo.Bytes{Value: v}
	case MAP:
		objs := make(map[string]tengo.Object)
		for k, v := range v {
			objs[k] = toObject(v)
		}

		return &tengo.Map{Value: objs}
	case ARR:
		var objs []tengo.Object
		for _, e := range v {
			objs = append(objs, toObject(e))
		}

		return &tengo.Array{Value: objs}
	case IMAP:
		objs := make(map[string]tengo.Object)
		for k, v := range v {
			objs[k] = toObject(v)
		}

		return &tengo.ImmutableMap{Value: objs}
	case IARR:
		var objs []tengo.Object
		for _, e := range v {
			objs = append(objs, toObject(e))
		}

		return &tengo.ImmutableArray{Value: objs}
	}

	panic(fmt.Errorf("unknown type: %T", v))
}

func objectZeroCopy(o tengo.Object) tengo.Object {
	switch o.(type) {
	case *tengo.Int:
		return &tengo.Int{}
	case *tengo.Float:
		return &tengo.Float{}
	case *tengo.Bool:
		return &tengo.Bool{}
	case *tengo.Char:
		return &tengo.Char{}
	case *tengo.String:
		return &tengo.String{}
	case *tengo.Array:
		return &tengo.Array{}
	case *tengo.Map:
		return &tengo.Map{}
	case *tengo.Undefined:
		return tengo.UndefinedValue
	case *tengo.Error:
		return &tengo.Error{}
	case *tengo.Bytes:
		return &tengo.Bytes{}
	case *tengo.ImmutableArray:
		return &tengo.ImmutableArray{}
	case *tengo.ImmutableMap:
		return &tengo.ImmutableMap{}
	case nil:
		panic("nil")
	default:
		panic(fmt.Errorf("unknown object type: %s", o.TypeName()))
	}
}

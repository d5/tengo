package runtime_test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/d5/tengo/objects"
)

func TestBuiltin(t *testing.T) {
	m := Opts().Module("math",
		&objects.BuiltinModule{
			Attrs: map[string]objects.Object{
				"abs": &objects.UserFunction{
					Name: "abs",
					Value: func(args ...objects.Object) (ret objects.Object, err error) {
						v, _ := objects.ToFloat64(args[0])
						return &objects.Float{Value: math.Abs(v)}, nil
					},
				},
			},
		})

	// builtin
	expectOpts(t, `math := import("math"); out = math.abs(1)`, m, 1.0)
	expectOpts(t, `math := import("math"); out = math.abs(-1)`, m, 1.0)
	expectOpts(t, `math := import("math"); out = math.abs(1.0)`, m, 1.0)
	expectOpts(t, `math := import("math"); out = math.abs(-1.0)`, m, 1.0)
}

func TestUserModules(t *testing.T) {
	// export none
	expectOpts(t, `out = import("mod1")`, Opts().Module("mod1", `fn := func() { return 5.0 }; a := 2`), objects.UndefinedValue)

	// export values
	expectOpts(t, `out = import("mod1")`, Opts().Module("mod1", `export 5`), 5)
	expectOpts(t, `out = import("mod1")`, Opts().Module("mod1", `export "foo"`), "foo")

	// export compound types
	expectOpts(t, `out = import("mod1")`, Opts().Module("mod1", `export [1, 2, 3]`), IARR{1, 2, 3})
	expectOpts(t, `out = import("mod1")`, Opts().Module("mod1", `export {a: 1, b: 2}`), IMAP{"a": 1, "b": 2})

	// export value is immutable
	expectErrorOpts(t, `m1 := import("mod1"); m1.a = 5`, Opts().Module("mod1", `export {a: 1, b: 2}`), "not index-assignable")
	expectErrorOpts(t, `m1 := import("mod1"); m1[1] = 5`, Opts().Module("mod1", `export [1, 2, 3]`), "not index-assignable")

	// code after export statement will not be executed
	expectOpts(t, `out = import("mod1")`, Opts().Module("mod1", `a := 10; export a; a = 20`), 10)
	expectOpts(t, `out = import("mod1")`, Opts().Module("mod1", `a := 10; export a; a = 20; export a`), 10)

	// export function
	expectOpts(t, `out = import("mod1")()`, Opts().Module("mod1", `export func() { return 5.0 }`), 5.0)
	// export function that reads module-global variable
	expectOpts(t, `out = import("mod1")()`, Opts().Module("mod1", `a := 1.5; export func() { return a + 5.0 }`), 6.5)
	// export function that read local variable
	expectOpts(t, `out = import("mod1")()`, Opts().Module("mod1", `export func() { a := 1.5; return a + 5.0 }`), 6.5)
	// export function that read free variables
	expectOpts(t, `out = import("mod1")()`, Opts().Module("mod1", `export func() { a := 1.5; return func() { return a + 5.0 }() }`), 6.5)

	// recursive function in module
	expectOpts(t, `out = import("mod1")`, Opts().Module(
		"mod1", `
a := func(x) {
	return x == 0 ? 0 : x + a(x-1)
}

export a(5)
`), 15)
	expectOpts(t, `out = import("mod1")`, Opts().Module(
		"mod1", `
export func() {
	a := func(x) {
		return x == 0 ? 0 : x + a(x-1)
	}

	return a(5)
}()
`), 15)

	// (main) -> mod1 -> mod2
	expectOpts(t, `out = import("mod1")()`,
		Opts().Module("mod1", `export import("mod2")`).
			Module("mod2", `export func() { return 5.0 }`),
		5.0)
	// (main) -> mod1 -> mod2
	//        -> mod2
	expectOpts(t, `import("mod1"); out = import("mod2")()`,
		Opts().Module("mod1", `export import("mod2")`).
			Module("mod2", `export func() { return 5.0 }`),
		5.0)
	// (main) -> mod1 -> mod2 -> mod3
	//        -> mod2 -> mod3
	expectOpts(t, `import("mod1"); out = import("mod2")()`,
		Opts().Module("mod1", `export import("mod2")`).
			Module("mod2", `export import("mod3")`).
			Module("mod3", `export func() { return 5.0 }`),
		5.0)

	// cyclic imports
	// (main) -> mod1 -> mod2 -> mod1
	expectErrorOpts(t, `import("mod1")`,
		Opts().Module("mod1", `import("mod2")`).
			Module("mod2", `import("mod1")`),
		"Compile Error: cyclic module import: mod1\n\tat mod2:1:1")
	// (main) -> mod1 -> mod2 -> mod3 -> mod1
	expectErrorOpts(t, `import("mod1")`,
		Opts().Module("mod1", `import("mod2")`).
			Module("mod2", `import("mod3")`).
			Module("mod3", `import("mod1")`),
		"Compile Error: cyclic module import: mod1\n\tat mod3:1:1")
	// (main) -> mod1 -> mod2 -> mod3 -> mod2
	expectErrorOpts(t, `import("mod1")`,
		Opts().Module("mod1", `import("mod2")`).
			Module("mod2", `import("mod3")`).
			Module("mod3", `import("mod2")`),
		"Compile Error: cyclic module import: mod2\n\tat mod3:1:1")

	// unknown modules
	expectErrorOpts(t, `import("mod0")`, Opts().Module("mod1", `a := 5`), "module 'mod0' not found")
	expectErrorOpts(t, `import("mod1")`, Opts().Module("mod1", `import("mod2")`), "module 'mod2' not found")

	// module is immutable but its variables is not necessarily immutable.
	expectOpts(t, `m1 := import("mod1"); m1.a.b = 5; out = m1.a.b`,
		Opts().Module("mod1", `export {a: {b: 3}}`),
		5)

	// make sure module has same builtin functions
	expectOpts(t, `out = import("mod1")`,
		Opts().Module("mod1", `export func() { return type_name(0) }()`),
		"int")

	// 'export' statement is ignored outside module
	expectOpts(t, `a := 5; export func() { a = 10 }(); out = a`, Opts().Skip2ndPass(), 5)

	// 'export' must be in the top-level
	expectErrorOpts(t, `import("mod1")`,
		Opts().Module("mod1", `func() { export 5 }()`),
		"Compile Error: export not allowed inside function\n\tat mod1:1:10")
	expectErrorOpts(t, `import("mod1")`,
		Opts().Module("mod1", `func() { func() { export 5 }() }()`),
		"Compile Error: export not allowed inside function\n\tat mod1:1:19")

	// module cannot access outer scope
	expectErrorOpts(t, `a := 5; import("mod1")`,
		Opts().Module("mod1", `export a`),
		"Compile Error: unresolved reference 'a'\n\tat mod1:1:8")

	// runtime error within modules
	expectErrorOpts(t, `
a := 1;
b := import("mod1");
b(a)`,
		Opts().Module("mod1", `
export func(a) {
   a()
}
`), "Runtime Error: not callable: int\n\tat mod1:3:4\n\tat test:4:1")
}

func TestModuleBlockScopes(t *testing.T) {
	m := Opts().Module("rand",
		&objects.BuiltinModule{
			Attrs: map[string]objects.Object{
				"intn": &objects.UserFunction{
					Name: "abs",
					Value: func(args ...objects.Object) (ret objects.Object, err error) {
						v, _ := objects.ToInt64(args[0])
						return &objects.Int{Value: rand.Int63n(v)}, nil
					},
				},
			},
		})

	// block scopes in module
	expectOpts(t, `out = import("mod1")()`, m.Module(
		"mod1", `
	rand := import("rand")
	foo := func() { return 1 }
	export func() {
		rand.intn(3)
		return foo()
	}`), 1)

	expectOpts(t, `out = import("mod1")()`, m.Module(
		"mod1", `
rand := import("rand")
foo := func() { return 1 }
export func() {
	rand.intn(3)
	if foo() {}
	return 10
}
`), 10)

	expectOpts(t, `out = import("mod1")()`, m.Module(
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

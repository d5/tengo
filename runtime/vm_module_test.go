package runtime_test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/d5/tengo/objects"
)

func TestBuiltin(t *testing.T) {
	mods := map[string]objects.Importable{
		"math": &objects.BuiltinModule{
			Attrs: map[string]objects.Object{
				"abs": &objects.UserFunction{
					Name: "abs",
					Value: func(args ...objects.Object) (ret objects.Object, err error) {
						v, _ := objects.ToFloat64(args[0])
						return &objects.Float{Value: math.Abs(v)}, nil
					},
				},
			},
		},
	}

	// builtin
	expectWithBuiltinModules(t, `math := import("math"); out = math.abs(1)`, 1.0, mods)
	expectWithBuiltinModules(t, `math := import("math"); out = math.abs(-1)`, 1.0, mods)
	expectWithBuiltinModules(t, `math := import("math"); out = math.abs(1.0)`, 1.0, mods)
	expectWithBuiltinModules(t, `math := import("math"); out = math.abs(-1.0)`, 1.0, mods)
}

func TestUserModules(t *testing.T) {
	// user modules

	// export none
	expectWithUserModules(t, `out = import("mod1")`, objects.UndefinedValue, map[string]string{
		"mod1": `fn := func() { return 5.0 }; a := 2`,
	})

	// export values
	expectWithUserModules(t, `out = import("mod1")`, 5, map[string]string{
		"mod1": `export 5`,
	})
	expectWithUserModules(t, `out = import("mod1")`, "foo", map[string]string{
		"mod1": `export "foo"`,
	})

	// export compound types
	expectWithUserModules(t, `out = import("mod1")`, IARR{1, 2, 3}, map[string]string{
		"mod1": `export [1, 2, 3]`,
	})
	expectWithUserModules(t, `out = import("mod1")`, IMAP{"a": 1, "b": 2}, map[string]string{
		"mod1": `export {a: 1, b: 2}`,
	})
	// export value is immutable
	expectErrorWithUserModules(t, `m1 := import("mod1"); m1.a = 5`, map[string]string{
		"mod1": `export {a: 1, b: 2}`,
	}, "not index-assignable")
	expectErrorWithUserModules(t, `m1 := import("mod1"); m1[1] = 5`, map[string]string{
		"mod1": `export [1, 2, 3]`,
	}, "not index-assignable")

	// code after export statement will not be executed
	expectWithUserModules(t, `out = import("mod1")`, 10, map[string]string{
		"mod1": `a := 10; export a; a = 20`,
	})
	expectWithUserModules(t, `out = import("mod1")`, 10, map[string]string{
		"mod1": `a := 10; export a; a = 20; export a`,
	})

	// export function
	expectWithUserModules(t, `out = import("mod1")()`, 5.0, map[string]string{
		"mod1": `export func() { return 5.0 }`,
	})
	// export function that reads module-global variable
	expectWithUserModules(t, `out = import("mod1")()`, 6.5, map[string]string{
		"mod1": `a := 1.5; export func() { return a + 5.0 }`,
	})
	// export function that read local variable
	expectWithUserModules(t, `out = import("mod1")()`, 6.5, map[string]string{
		"mod1": `export func() { a := 1.5; return a + 5.0 }`,
	})
	// export function that read free variables
	expectWithUserModules(t, `out = import("mod1")()`, 6.5, map[string]string{
		"mod1": `export func() { a := 1.5; return func() { return a + 5.0 }() }`,
	})

	// recursive function in module
	expectWithUserModules(t, `out = import("mod1")`, 15, map[string]string{
		"mod1": `
a := func(x) {
	return x == 0 ? 0 : x + a(x-1)
}

export a(5)
`})
	expectWithUserModules(t, `out = import("mod1")`, 15, map[string]string{
		"mod1": `
export func() {
	a := func(x) {
		return x == 0 ? 0 : x + a(x-1)
	}

	return a(5)
}()
`})

	// (main) -> mod1 -> mod2
	expectWithUserModules(t, `out = import("mod1")()`, 5.0, map[string]string{
		"mod1": `export import("mod2")`,
		"mod2": `export func() { return 5.0 }`,
	})
	// (main) -> mod1 -> mod2
	//        -> mod2
	expectWithUserModules(t, `import("mod1"); out = import("mod2")()`, 5.0, map[string]string{
		"mod1": `export import("mod2")`,
		"mod2": `export func() { return 5.0 }`,
	})
	// (main) -> mod1 -> mod2 -> mod3
	//        -> mod2 -> mod3
	expectWithUserModules(t, `import("mod1"); out = import("mod2")()`, 5.0, map[string]string{
		"mod1": `export import("mod2")`,
		"mod2": `export import("mod3")`,
		"mod3": `export func() { return 5.0 }`,
	})

	// cyclic imports
	// (main) -> mod1 -> mod2 -> mod1
	expectErrorWithUserModules(t, `import("mod1")`, map[string]string{
		"mod1": `import("mod2")`,
		"mod2": `import("mod1")`,
	}, "Compile Error: cyclic module import: mod1\n\tat mod2:1:1")
	// (main) -> mod1 -> mod2 -> mod3 -> mod1
	expectErrorWithUserModules(t, `import("mod1")`, map[string]string{
		"mod1": `import("mod2")`,
		"mod2": `import("mod3")`,
		"mod3": `import("mod1")`,
	}, "Compile Error: cyclic module import: mod1\n\tat mod3:1:1")
	// (main) -> mod1 -> mod2 -> mod3 -> mod2
	expectErrorWithUserModules(t, `import("mod1")`, map[string]string{
		"mod1": `import("mod2")`,
		"mod2": `import("mod3")`,
		"mod3": `import("mod2")`,
	}, "Compile Error: cyclic module import: mod2\n\tat mod3:1:1")

	// unknown modules
	expectErrorWithUserModules(t, `import("mod0")`, map[string]string{
		"mod1": `a := 5`,
	}, "module 'mod0' not found")
	expectErrorWithUserModules(t, `import("mod1")`, map[string]string{
		"mod1": `import("mod2")`,
	}, "module 'mod2' not found")

	// module is immutable but its variables is not necessarily immutable.
	expectWithUserModules(t, `m1 := import("mod1"); m1.a.b = 5; out = m1.a.b`, 5, map[string]string{
		"mod1": `export {a: {b: 3}}`,
	})

	// make sure module has same builtin functions
	expectWithUserModules(t, `out = import("mod1")`, "int", map[string]string{
		"mod1": `export func() { return type_name(0) }()`,
	})

	// 'export' statement is ignored outside module
	expectNoMod(t, `a := 5; export func() { a = 10 }(); out = a`, 5)

	// 'export' must be in the top-level
	expectErrorWithUserModules(t, `import("mod1")`, map[string]string{
		"mod1": `func() { export 5 }()`,
	}, "Compile Error: export not allowed inside function\n\tat mod1:1:10")
	expectErrorWithUserModules(t, `import("mod1")`, map[string]string{
		"mod1": `func() { func() { export 5 }() }()`,
	}, "Compile Error: export not allowed inside function\n\tat mod1:1:19")

	// module cannot access outer scope
	expectErrorWithUserModules(t, `a := 5; import("mod1")`, map[string]string{
		"mod1": `export a`,
	}, "Compile Error: unresolved reference 'a'\n\tat mod1:1:8")

	// runtime error within modules
	expectErrorWithUserModules(t, `
a := 1;
b := import("mod1");
b(a)`,
		map[string]string{"mod1": `
export func(a) {
   a()
}
`,
		}, "Runtime Error: not callable: int\n\tat mod1:3:4\n\tat test:4:1")
}

func TestModuleBlockScopes(t *testing.T) {
	mods := map[string]objects.Importable{
		"rand": &objects.BuiltinModule{
			Attrs: map[string]objects.Object{
				"intn": &objects.UserFunction{
					Name: "abs",
					Value: func(args ...objects.Object) (ret objects.Object, err error) {
						v, _ := objects.ToInt64(args[0])
						return &objects.Int{Value: rand.Int63n(v)}, nil
					},
				},
			},
		},
	}

	// block scopes in module
	expectWithUserAndBuiltinModules(t, `out = import("mod1")()`, 1, map[string]string{
		"mod1": `
	rand := import("rand")
	foo := func() { return 1 }
	export func() {
		rand.intn(3)
		return foo()
	}
	`,
	}, mods)

	expectWithUserAndBuiltinModules(t, `out = import("mod1")()`, 10, map[string]string{
		"mod1": `
rand := import("rand")
foo := func() { return 1 }
export func() {
	rand.intn(3)
	if foo() {}
	return 10
}
`,
	}, mods)

	expectWithUserAndBuiltinModules(t, `out = import("mod1")()`, 10, map[string]string{
		"mod1": `
	rand := import("rand")
	foo := func() { return 1 }
	export func() {
		rand.intn(3)
		if true { foo() }
		return 10
	}
	`,
	}, mods)
}

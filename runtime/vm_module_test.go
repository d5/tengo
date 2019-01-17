package runtime_test

import "testing"

func TestModule(t *testing.T) {
	// stdmods
	expect(t, `math := import("math"); out = math.abs(1)`, 1.0)
	expect(t, `math := import("math"); out = math.abs(-1)`, 1.0)
	expect(t, `math := import("math"); out = math.abs(1.0)`, 1.0)
	expect(t, `math := import("math"); out = math.abs(-1.0)`, 1.0)

	// user modules
	expectWithUserModules(t, `out = import("mod1").bar()`, 5.0, map[string]string{
		"mod1": `bar := func() { return 5.0 }`,
	})
	// (main) -> mod1 -> mod2
	expectWithUserModules(t, `out = import("mod1").mod2.bar()`, 5.0, map[string]string{
		"mod1": `mod2 := import("mod2")`,
		"mod2": `bar := func() { return 5.0 }`,
	})
	// (main) -> mod1 -> mod2
	//        -> mod2
	expectWithUserModules(t, `import("mod1"); out = import("mod2").bar()`, 5.0, map[string]string{
		"mod1": `mod2 := import("mod2")`,
		"mod2": `bar := func() { return 5.0 }`,
	})
	// (main) -> mod1 -> mod2 -> mod3
	//        -> mod2 -> mod3
	expectWithUserModules(t, `import("mod1"); out = import("mod2").mod3.bar()`, 5.0, map[string]string{
		"mod1": `mod2 := import("mod2")`,
		"mod2": `mod3 := import("mod3")`,
		"mod3": `bar := func() { return 5.0 }`,
	})

	// cyclic imports
	// (main) -> mod1 -> mod2 -> mod1
	expectErrorWithUserModules(t, `import("mod1")`, map[string]string{
		"mod1": `import("mod2")`,
		"mod2": `import("mod1")`,
	})
	// (main) -> mod1 -> mod2 -> mod3 -> mod1
	expectErrorWithUserModules(t, `import("mod1")`, map[string]string{
		"mod1": `import("mod2")`,
		"mod2": `import("mod3")`,
		"mod3": `import("mod1")`,
	})
	// (main) -> mod1 -> mod2 -> mod3 -> mod2
	expectErrorWithUserModules(t, `import("mod1")`, map[string]string{
		"mod1": `import("mod2")`,
		"mod2": `import("mod3")`,
		"mod3": `import("mod2")`,
	})

	// unknown modules
	expectErrorWithUserModules(t, `import("mod0")`, map[string]string{
		"mod1": `a := 5`,
	})
	expectErrorWithUserModules(t, `import("mod1")`, map[string]string{
		"mod1": `import("mod2")`,
	})

	// for-in
	expectWithUserModules(t, `for _, n in import("mod1") { out += n }`, 6, map[string]string{
		"mod1": `a := 1; b := 2; c := 3`,
	})
	expectWithUserModules(t, `for k, _ in import("mod1") { out += k }`, "a", map[string]string{
		"mod1": `a := 1`, // only 1 global variable because module map does not sort the keys
	})

	// mutating global variables inside the module does not affect exported values
	expectWithUserModules(t, `m1 := import("mod1"); m1.mutate(); out = m1.a`, 3, map[string]string{
		"mod1": `a := 3; mutate := func() { a = 10 }`,
	})

	// module map is immutable
	expectErrorWithUserModules(t, `m1 := import("mod1"); m1.a = 5`, map[string]string{
		"mod1": `a := 3`,
	})
	expectErrorWithUserModules(t, `m1 := import("mod1"); m1.a.b = 5`, map[string]string{
		"mod1": `a := {b: 3}`,
	})
}

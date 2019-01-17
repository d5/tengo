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
	// script -> mod1 -> mod2
	expectWithUserModules(t, `out = import("mod1").mod2.bar()`, 5.0, map[string]string{
		"mod1": `mod2 := import("mod2")`,
		"mod2": `bar := func() { return 5.0 }`,
	})
	// cyclic: script -> mod1 -> mod2 -> mod1
	//expectWithUserModules(t, `out = import("mod1").mod2.bar()`, 5.0, map[string]string{
	//	"mod1": `mod2 := import("mod2")`,
	//	"mod2": `mod1 := import("mod1")`,
	//})

	//

	// TODO: test for-in statement with module map

	// TODO: test sharing of global variables

	// TODO: test module function that mutate other variables in the same module

	// TODO: test for cyclic dependencies

	// TODO: VM should reuse moduleVMs whenever possible
}

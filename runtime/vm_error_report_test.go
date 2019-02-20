package runtime_test

import "testing"

func TestVMErrorInfo(t *testing.T) {
	expectError(t, `a := 5
a + "boo"`,
		"test:2:1: invalid operation: int + string")

	expectError(t, `a := 5
b := a(5)`,
		"test:2:6: not callable: int")

	expectError(t, `a := 5
b := {}
b.x.y = 10`,
		"test:3:1: not index-assignable: undefined")

	expectError(t, `
a := func() {
	b := 5
	b += "foo"
}
a()`,
		"test:4:2: invalid operation: int + string")

	expectErrorWithUserModules(t, `a := 5
	a + import("mod1")`, map[string]string{
		"mod1": `export "foo"`,
	}, "test:2:2: invalid operation: int + string")

	expectErrorWithUserModules(t, `a := import("mod1")()`, map[string]string{
		"mod1": `
export func() {
	b := 5
	return b + "foo"
}`,
	}, "mod1:4:9: invalid operation: int + string")

	expectErrorWithUserModules(t, `a := import("mod1")()`, map[string]string{
		"mod1": `export import("mod2")()`,
		"mod2": `
export func() {
	b := 5
	return b + "foo"
}`,
	}, "mod2:4:9: invalid operation: int + string")
}

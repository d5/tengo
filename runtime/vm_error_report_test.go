package runtime_test

import "testing"

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

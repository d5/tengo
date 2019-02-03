package runtime_test

import (
	"testing"

	"github.com/d5/tengo/objects"
)

func TestStdLib(t *testing.T) {
	// stdlib
	expect(t, `math := import("math"); out = math.abs(1)`, 1.0)
	expect(t, `math := import("math"); out = math.abs(-1)`, 1.0)
	expect(t, `math := import("math"); out = math.abs(1.0)`, 1.0)
	expect(t, `math := import("math"); out = math.abs(-1.0)`, 1.0)

	// os.File
	expect(t, `
os := import("os")

write_file := func(filename, data) {
	file := os.create(filename)
	if !file { return file }

	if res := file.write(bytes(data)); is_error(res) {
		return res
	}

	return file.close()
}

read_file := func(filename) {
	file := os.open(filename)
	if !file { return file }

	data := bytes(100)
	cnt := file.read(data)
	if  is_error(cnt) {
		return cnt
	}

	file.close()
	return data[:cnt]
}

if write_file("./temp", "foobar") {
	out = string(read_file("./temp"))
}

os.remove("./temp")
`, "foobar")

	// exec.command
	expect(t, `
os := import("os")
cmd := os.exec("echo", "foo", "bar")
if !is_error(cmd) { 
	out = cmd.output()
}
`, []byte("foo bar\n"))

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

	// export composite types
	expectWithUserModules(t, `out = import("mod1")`, IARR{1, 2, 3}, map[string]string{
		"mod1": `export [1, 2, 3]`,
	})
	expectWithUserModules(t, `out = import("mod1")`, IMAP{"a": 1, "b": 2}, map[string]string{
		"mod1": `export {a: 1, b: 2}`,
	})
	// export value is immutable
	expectErrorWithUserModules(t, `m1 := import("mod1"); m1.a = 5`, map[string]string{
		"mod1": `export {a: 1, b: 2}`,
	})
	expectErrorWithUserModules(t, `m1 := import("mod1"); m1[1] = 5`, map[string]string{
		"mod1": `export [1, 2, 3]`,
	})

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

	// module is immutable but its variables is not necessarily immutable.
	expectWithUserModules(t, `m1 := import("mod1"); m1.a.b = 5; out = m1.a.b`, 5, map[string]string{
		"mod1": `export {a: {b: 3}}`,
	})

	// make sure module has same builtin functions
	expectWithUserModules(t, `out = import("mod1")`, "int", map[string]string{
		"mod1": `export func() { return type_name(0) }()`,
	})

	// 'export' statement is ignored outside module
	expect(t, `a := 5; export func() { a = 10 }(); out = a`, 5)

	// 'export' must be in the top-level
	expectErrorWithUserModules(t, `import("mod1")`, map[string]string{
		"mod1": `func() { export 5 }()`,
	})
	expectErrorWithUserModules(t, `import("mod1")`, map[string]string{
		"mod1": `func() { func() { export 5 }() }()`,
	})

	// module cannot access outer scope
	expectErrorWithUserModules(t, `a := 5; import("mod")`, map[string]string{
		"mod1": `export a`,
	})
}

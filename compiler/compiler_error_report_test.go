package compiler_test

import "testing"

func TestCompilerErrorReport(t *testing.T) {
	expectError(t, `import("user1")`, "Compile Error: module 'user1' not found\n\tat test:1:1")

	expectError(t, `a = 1`, "Compile Error: unresolved reference 'a'\n\tat test:1:1")
	expectError(t, `a, b := 1, 2`, "Compile Error: tuple assignment not allowed\n\tat test:1:1")
	expectError(t, `a.b := 1`, "not allowed with selector")
	expectError(t, `a:=1; a:=3`, "Compile Error: 'a' redeclared in this block\n\tat test:1:7")

	expectError(t, `return 5`, "Compile Error: return not allowed outside function\n\tat test:1:1")
	expectError(t, `func() { break }`, "Compile Error: break not allowed outside loop\n\tat test:1:10")
	expectError(t, `func() { continue }`, "Compile Error: continue not allowed outside loop\n\tat test:1:10")
	expectError(t, `func() { export 5 }`, "Compile Error: export not allowed inside function\n\tat test:1:10")
}

package compiler_test

import "testing"

func TestCompilerErrorReport(t *testing.T) {
	expectError(t, `import("user1")`, "test:1:1: module file read error: open user1.tengo: no such file or directory")

	expectError(t, `a = 1`, "test:1:1: unresolved reference 'a'")
	expectError(t, `a, b := 1, 2`, "test:1:1: tuple assignment not allowed")
	expectError(t, `a.b := 1`, "not allowed with selector")
	expectError(t, `a:=1; a:=3`, "test:1:7: 'a' redeclared in this block")

	expectError(t, `return 5`, "test:1:1: return not allowed outside function")
	expectError(t, `func() { break }`, "test:1:10: break not allowed outside loop")
	expectError(t, `func() { continue }`, "test:1:10: continue not allowed outside loop")
	expectError(t, `func() { export 5 }`, "test:1:10: export not allowed inside function")
}

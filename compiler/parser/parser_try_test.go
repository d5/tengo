package parser_test

import (
	"testing"

	"github.com/d5/tengo/compiler/ast"
)

func TestTry(t *testing.T) {
	expectString(t, "try(someExpr)", "try(someExpr)")

	expect(t, `try(someExpr)`, func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				tryExpr(
					p(1, 1),
					p(1, 4),
					p(1, 13),
					ident("someExpr", p(1, 5)),
				)))
	})

	expectError(t, "try(a, b)")
	expectError(t, "try(a...)")
	expectError(t, "x := try")
	expectError(t, "try := func(){}")
	expectError(t, "try = func(){}")
}

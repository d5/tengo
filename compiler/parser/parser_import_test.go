package parser_test

import (
	"testing"

	"github.com/d5/tengo/compiler/ast"
	"github.com/d5/tengo/compiler/token"
)

func TestImport(t *testing.T) {
	expect(t, `a := import("mod1")`, func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(ident("a", p(1, 1))),
				exprs(importExpr("mod1", p(1, 6))),
				token.Define, p(1, 3)))
	})

	expect(t, `import("mod1").var1`, func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				selectorExpr(
					importExpr("mod1", p(1, 1)),
					stringLit("var1", p(1, 16)))))
	})

	expect(t, `import("mod1").func1()`, func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				callExpr(
					selectorExpr(
						importExpr("mod1", p(1, 1)),
						stringLit("func1", p(1, 16))),
					p(1, 21), p(1, 22))))
	})

	expect(t, `for x, y in import("mod1") {}`, func(p pfn) []ast.Stmt {
		return stmts(
			forInStmt(
				ident("x", p(1, 5)),
				ident("y", p(1, 8)),
				importExpr("mod1", p(1, 13)),
				blockStmt(p(1, 28), p(1, 29)),
				p(1, 1)))
	})
}

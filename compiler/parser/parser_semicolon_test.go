package parser_test

import (
	"testing"

	"github.com/d5/tengo/compiler/ast"
)

func TestSemicolon(t *testing.T) {
	expect(t, "1", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(intLit(1, p(1, 1))))
	})

	expect(t, "1;", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(intLit(1, p(1, 1))))
	})

	expect(t, "1;;", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(intLit(1, p(1, 1))),
			emptyStmt(false, p(1, 3)))
	})

	expect(t, `1
`, func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(intLit(1, p(1, 1))))
	})

	expect(t, `1
;`, func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(intLit(1, p(1, 1))),
			emptyStmt(false, p(2, 1)))
	})

	expect(t, `1;
;`, func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(intLit(1, p(1, 1))),
			emptyStmt(false, p(2, 1)))
	})
}

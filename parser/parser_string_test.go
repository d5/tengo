package parser_test

import (
	"testing"

	"github.com/d5/tengo/ast"
	"github.com/d5/tengo/token"
)

func TestString(t *testing.T) {
	expect(t, `a = "foo\nbar"`, func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(ident("a", p(1, 1))),
				exprs(stringLit("foo\nbar", p(1, 5))),
				token.Assign,
				p(1, 3)))
	})

	expect(t, "a = `raw string`", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(ident("a", p(1, 1))),
				exprs(stringLit("raw string", p(1, 5))),
				token.Assign,
				p(1, 3)))
	})
}

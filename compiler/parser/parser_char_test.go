package parser_test

import (
	"testing"

	"github.com/d5/tengo/compiler/ast"
)

func TestChar(t *testing.T) {
	expect(t, `'A'`, func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				charLit('A', 1)))
	})
	expect(t, `'九'`, func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				charLit('九', 1)))
	})

	expectError(t, `''`)
	expectError(t, `'AB'`)
	expectError(t, `'A九'`)
}

package parser_test

import (
	"testing"

	"github.com/d5/tengo/compiler/ast"
	"github.com/d5/tengo/compiler/token"
)

func TestIndex(t *testing.T) {
	expect(t, "[1, 2, 3][1]", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				indexExpr(
					arrayLit(p(1, 1), p(1, 9),
						intLit(1, p(1, 2)),
						intLit(2, p(1, 5)),
						intLit(3, p(1, 8))),
					intLit(1, p(1, 11)),
					p(1, 10), p(1, 12))))
	})

	expect(t, "[1, 2, 3][5 - a]", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				indexExpr(
					arrayLit(p(1, 1), p(1, 9),
						intLit(1, p(1, 2)),
						intLit(2, p(1, 5)),
						intLit(3, p(1, 8))),
					binaryExpr(
						intLit(5, p(1, 11)),
						ident("a", p(1, 15)),
						token.Sub,
						p(1, 13)),
					p(1, 10), p(1, 16))))
	})

	expect(t, "[1, 2, 3][5 : a]", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				sliceExpr(
					arrayLit(p(1, 1), p(1, 9),
						intLit(1, p(1, 2)),
						intLit(2, p(1, 5)),
						intLit(3, p(1, 8))),
					intLit(5, p(1, 11)),
					ident("a", p(1, 15)),
					p(1, 10), p(1, 16))))
	})

	expect(t, "[1, 2, 3][a + 3 : b - 8]", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				sliceExpr(
					arrayLit(p(1, 1), p(1, 9),
						intLit(1, p(1, 2)),
						intLit(2, p(1, 5)),
						intLit(3, p(1, 8))),
					binaryExpr(
						ident("a", p(1, 11)),
						intLit(3, p(1, 15)),
						token.Add,
						p(1, 13)),
					binaryExpr(
						ident("b", p(1, 19)),
						intLit(8, p(1, 23)),
						token.Sub,
						p(1, 21)),
					p(1, 10), p(1, 24))))
	})

	expect(t, `{a: 1, b: 2}["b"]`, func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				indexExpr(
					mapLit(p(1, 1), p(1, 12),
						mapElementLit("a", p(1, 2), p(1, 3), intLit(1, p(1, 5))),
						mapElementLit("b", p(1, 8), p(1, 9), intLit(2, p(1, 11)))),
					stringLit("b", p(1, 14)),
					p(1, 13), p(1, 17))))
	})

	expect(t, `{a: 1, b: 2}[a + b]`, func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				indexExpr(
					mapLit(p(1, 1), p(1, 12),
						mapElementLit("a", p(1, 2), p(1, 3), intLit(1, p(1, 5))),
						mapElementLit("b", p(1, 8), p(1, 9), intLit(2, p(1, 11)))),
					binaryExpr(
						ident("a", p(1, 14)),
						ident("b", p(1, 18)),
						token.Add,
						p(1, 16)),
					p(1, 13), p(1, 19))))
	})
}

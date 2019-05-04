package parser_test

import (
	"testing"

	"github.com/d5/tengo/compiler/ast"
)

func TestSpreadArray(t *testing.T) {
	expect(t, "[a..., b, c]", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(arrayLit(
				p(1, 1), p(1, 12),
				spreadExpr(ident("a", p(1, 2)), p(1, 3)),
				ident("b", p(1, 8)),
				ident("c", p(1, 11)),
			)))
	})

	expect(t, "[a..., b..., c]", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(arrayLit(
				p(1, 1), p(1, 15),
				spreadExpr(ident("a", p(1, 2)), p(1, 3)),
				spreadExpr(ident("b", p(1, 8)), p(1, 9)),
				ident("c", p(1, 14)),
			)))
	})

	expect(t, "[a..., b..., c...]", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(arrayLit(
				p(1, 1), p(1, 18),
				spreadExpr(ident("a", p(1, 2)), p(1, 3)),
				spreadExpr(ident("b", p(1, 8)), p(1, 9)),
				spreadExpr(ident("c", p(1, 14)), p(1, 15)),
			)))
	})

	expect(t, "[a.b.c...]", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(arrayLit(
				p(1, 1), p(1, 10),
				spreadExpr(
					selectorExpr(
						selectorExpr(ident("a", p(1, 2)), stringLit("b", p(1, 4))),
						stringLit("c", p(1, 6)),
					),
					p(1, 7)),
			)))
	})

	expect(t, "[[1,2,3]...]", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(arrayLit(
				p(1, 1), p(1, 12),
				spreadExpr(
					arrayLit(p(1, 2), p(1, 8),
						intLit(1, p(1, 3)),
						intLit(2, p(1, 5)),
						intLit(3, p(1, 7)),
					),
					p(1, 9)),
			)))
	})
}

func TestSpreadCall(t *testing.T) {
	expect(t, "fn(a..., b, c)", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(callExpr(
				ident("fn", p(1, 1)),
				p(1, 3), p(1, 14),
				spreadExpr(ident("a", p(1, 4)), p(1, 5)),
				ident("b", p(1, 10)),
				ident("c", p(1, 13)),
			)))
	})

	expect(t, "fn(a..., b..., c)", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(callExpr(
				ident("fn", p(1, 1)),
				p(1, 3), p(1, 17),
				spreadExpr(ident("a", p(1, 4)), p(1, 5)),
				spreadExpr(ident("b", p(1, 10)), p(1, 11)),
				ident("c", p(1, 16)),
			)))
	})

	expect(t, "fn(a..., b..., c...)", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(callExpr(
				ident("fn", p(1, 1)),
				p(1, 3), p(1, 20),
				spreadExpr(ident("a", p(1, 4)), p(1, 5)),
				spreadExpr(ident("b", p(1, 10)), p(1, 11)),
				spreadExpr(ident("c", p(1, 16)), p(1, 17)),
			)))
	})

	expect(t, "fn(a.b.c...)", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(callExpr(
				ident("fn", p(1, 1)),
				p(1, 3), p(1, 12),
				spreadExpr(
					selectorExpr(
						selectorExpr(ident("a", p(1, 4)), stringLit("b", p(1, 6))),
						stringLit("c", p(1, 8)),
					),
					p(1, 9)),
			)))
	})

	expect(t, "obj.fn(a.b.c...)", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(callExpr(
				selectorExpr(ident("obj", p(1, 1)), stringLit("fn", p(1, 5))),
				p(1, 7), p(1, 16),
				spreadExpr(
					selectorExpr(
						selectorExpr(ident("a", p(1, 8)), stringLit("b", p(1, 10))),
						stringLit("c", p(1, 12)),
					),
					p(1, 13)),
			)))
	})
}

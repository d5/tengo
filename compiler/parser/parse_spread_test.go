package parser_test

import (
	"testing"

	"github.com/d5/tengo/compiler/ast"
)

func TestSpreadArray(t *testing.T) {
	expect(t, "[a...]", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				arrayLit(p(1, 1), p(1, 6),
					SpreadExpr(
						ident("a", p(1, 2)),
					))))
	})

	expect(t, "[a..., b]", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				arrayLit(p(1, 1), p(1, 9),
					SpreadExpr(
						ident("a", p(1, 2)),
					),
					ident("b", p(1, 8)),
				)))
	})

	expect(t, "[a, b...]", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				arrayLit(p(1, 1), p(1, 9),
					ident("a", p(1, 2)),
					SpreadExpr(
						ident("b", p(1, 5)),
					),
				)))
	})

	expect(t, "[a..., b...]", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				arrayLit(p(1, 1), p(1, 12),
					SpreadExpr(
						ident("a", p(1, 2)),
					),
					SpreadExpr(
						ident("b", p(1, 8)),
					),
				)))
	})
}

func TestSpreadFunc(t *testing.T) {
	expect(t, "fn(a...)", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				callExpr(
					ident("fn", p(1, 1)),
					p(1, 3), p(1, 8),
					SpreadExpr(
						ident("a", p(1, 4)),
					),
				),
			))
	})

	expect(t, "fn(a..., b)", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				callExpr(
					ident("fn", p(1, 1)),
					p(1, 3), p(1, 11),
					SpreadExpr(
						ident("a", p(1, 4)),
					),
					ident("b", p(1, 10)),
				),
			))
	})

	expect(t, "fn(a, b...)", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				callExpr(
					ident("fn", p(1, 1)),
					p(1, 3), p(1, 11),
					ident("a", p(1, 4)),
					SpreadExpr(
						ident("b", p(1, 7)),
					),
				),
			))
	})

	expect(t, "fn(a..., b...)", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				callExpr(
					ident("fn", p(1, 1)),
					p(1, 3), p(1, 14),
					SpreadExpr(
						ident("a", p(1, 4)),
					),
					SpreadExpr(
						ident("b", p(1, 10)),
					),
				),
			))
	})
}

package parser_test

import (
	"testing"

	"github.com/d5/tengo/compiler/ast"
)

func TestExplodeArray(t *testing.T) {
	expect(t, "[a...]", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				arrayLit(p(1, 1), p(1, 6),
					explodeExpr(
						ident("a", p(1, 2)),
					))))
	})

	expect(t, "[a..., b]", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				arrayLit(p(1, 1), p(1, 9),
					explodeExpr(
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
					explodeExpr(
						ident("b", p(1, 5)),
					),
				)))
	})

	expect(t, "[a..., b...]", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				arrayLit(p(1, 1), p(1, 12),
					explodeExpr(
						ident("a", p(1, 2)),
					),
					explodeExpr(
						ident("b", p(1, 8)),
					),
				)))
	})
}

func TestExplodeFunc(t *testing.T) {
	expect(t, "fn(a...)", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				callExpr(
					ident("fn", p(1, 1)),
					p(1, 3), p(1, 8),
					explodeExpr(
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
					explodeExpr(
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
					explodeExpr(
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
					explodeExpr(
						ident("a", p(1, 4)),
					),
					explodeExpr(
						ident("b", p(1, 10)),
					),
				),
			))
	})
}

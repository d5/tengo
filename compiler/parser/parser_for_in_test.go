package parser_test

import (
	"testing"

	"github.com/d5/tengo/compiler/ast"
)

func TestForIn(t *testing.T) {
	expect(t, "for x in y {}", func(p pfn) []ast.Stmt {
		return stmts(
			forInStmt(
				ident("_", p(1, 5)),
				ident("x", p(1, 5)),
				ident("y", p(1, 10)),
				blockStmt(p(1, 12), p(1, 13)),
				p(1, 1)))
	})

	expect(t, "for _ in y {}", func(p pfn) []ast.Stmt {
		return stmts(
			forInStmt(
				ident("_", p(1, 5)),
				ident("_", p(1, 5)),
				ident("y", p(1, 10)),
				blockStmt(p(1, 12), p(1, 13)),
				p(1, 1)))
	})

	expect(t, "for x in [1, 2, 3] {}", func(p pfn) []ast.Stmt {
		return stmts(
			forInStmt(
				ident("_", p(1, 5)),
				ident("x", p(1, 5)),
				arrayLit(
					p(1, 10), p(1, 18),
					intLit(1, p(1, 11)),
					intLit(2, p(1, 14)),
					intLit(3, p(1, 17))),
				blockStmt(p(1, 20), p(1, 21)),
				p(1, 1)))
	})

	expect(t, "for x, y in z {}", func(p pfn) []ast.Stmt {
		return stmts(
			forInStmt(
				ident("x", p(1, 5)),
				ident("y", p(1, 8)),
				ident("z", p(1, 13)),
				blockStmt(p(1, 15), p(1, 16)),
				p(1, 1)))
	})

	expect(t, "for x, y in {k1: 1, k2: 2} {}", func(p pfn) []ast.Stmt {
		return stmts(
			forInStmt(
				ident("x", p(1, 5)),
				ident("y", p(1, 8)),
				mapLit(
					p(1, 13), p(1, 26),
					mapElementLit("k1", p(1, 14), p(1, 16), intLit(1, p(1, 18))),
					mapElementLit("k2", p(1, 21), p(1, 23), intLit(2, p(1, 25)))),
				blockStmt(p(1, 28), p(1, 29)),
				p(1, 1)))
	})
}

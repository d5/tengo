package parser_test

import (
	"testing"

	"github.com/d5/tengo/compiler/ast"
	"github.com/d5/tengo/compiler/token"
)

func TestAssignment(t *testing.T) {
	expect(t, "a = 5", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(ident("a", p(1, 1))),
				exprs(intLit(5, p(1, 5))),
				token.Assign,
				p(1, 3)))
	})

	expect(t, "a := 5", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(ident("a", p(1, 1))),
				exprs(intLit(5, p(1, 6))),
				token.Define,
				p(1, 3)))
	})

	expect(t, "a, b = 5, 10", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(
					ident("a", p(1, 1)),
					ident("b", p(1, 4))),
				exprs(
					intLit(5, p(1, 8)),
					intLit(10, p(1, 11))),
				token.Assign,
				p(1, 6)))
	})

	expect(t, "a, b := 5, 10", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(
					ident("a", p(1, 1)),
					ident("b", p(1, 4))),
				exprs(
					intLit(5, p(1, 9)),
					intLit(10, p(1, 12))),
				token.Define,
				p(1, 6)))
	})

	expect(t, "a, b = a + 2, b - 8", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(
					ident("a", p(1, 1)),
					ident("b", p(1, 4))),
				exprs(
					binaryExpr(
						ident("a", p(1, 8)),
						intLit(2, p(1, 12)),
						token.Add,
						p(1, 10)),
					binaryExpr(
						ident("b", p(1, 15)),
						intLit(8, p(1, 19)),
						token.Sub,
						p(1, 17))),
				token.Assign,
				p(1, 6)))
	})

	expect(t, "a = [1, 2, 3]", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(ident("a", p(1, 1))),
				exprs(arrayLit(p(1, 5), p(1, 13),
					intLit(1, p(1, 6)),
					intLit(2, p(1, 9)),
					intLit(3, p(1, 12)))),
				token.Assign,
				p(1, 3)))
	})

	expect(t, "a = [1 + 2, b * 4, [4, c]]", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(ident("a", p(1, 1))),
				exprs(arrayLit(p(1, 5), p(1, 26),
					binaryExpr(
						intLit(1, p(1, 6)),
						intLit(2, p(1, 10)),
						token.Add,
						p(1, 8)),
					binaryExpr(
						ident("b", p(1, 13)),
						intLit(4, p(1, 17)),
						token.Mul,
						p(1, 15)),
					arrayLit(p(1, 20), p(1, 25),
						intLit(4, p(1, 21)),
						ident("c", p(1, 24))))),
				token.Assign,
				p(1, 3)))
	})

	expect(t, "a += 5", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(ident("a", p(1, 1))),
				exprs(intLit(5, p(1, 6))),
				token.AddAssign,
				p(1, 3)))
	})

	expect(t, "a *= 5 + 10", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(ident("a", p(1, 1))),
				exprs(
					binaryExpr(
						intLit(5, p(1, 6)),
						intLit(10, p(1, 10)),
						token.Add,
						p(1, 8))),
				token.MulAssign,
				p(1, 3)))
	})
}

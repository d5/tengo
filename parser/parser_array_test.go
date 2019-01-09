package parser_test

import (
	"testing"

	"github.com/d5/tengo/ast"
	"github.com/d5/tengo/token"
)

func TestArray(t *testing.T) {
	expect(t, "[1, 2, 3]", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				arrayLit(p(1, 1), p(1, 9),
					intLit(1, p(1, 2)),
					intLit(2, p(1, 5)),
					intLit(3, p(1, 8)))))
	})

	expect(t, `
[
	1, 
	2, 
	3
]`, func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				arrayLit(p(2, 1), p(6, 1),
					intLit(1, p(3, 2)),
					intLit(2, p(4, 2)),
					intLit(3, p(5, 2)))))
	})
	expect(t, `
[
	1, 
	2, 
	3

]`, func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				arrayLit(p(2, 1), p(7, 1),
					intLit(1, p(3, 2)),
					intLit(2, p(4, 2)),
					intLit(3, p(5, 2)))))
	})

	expect(t, `[1, "foo", 12.34]`, func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				arrayLit(p(1, 1), p(1, 17),
					intLit(1, p(1, 2)),
					stringLit("foo", p(1, 5)),
					floatLit(12.34, p(1, 12)))))
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

	expectError(t, `[1, 2, 3,]`)
	expectError(t, `
[
	1, 
	2, 
	3,
]`)
	expectError(t, `
[
	1, 
	2, 
	3,

]`)
	expectError(t, `[1, 2, 3, ,]`)
}

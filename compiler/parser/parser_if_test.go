package parser_test

import (
	"testing"

	"github.com/d5/tengo/compiler/ast"
	"github.com/d5/tengo/compiler/token"
)

func TestIf(t *testing.T) {
	expect(t, "if a == 5 {}", func(p pfn) []ast.Stmt {
		return stmts(
			ifStmt(
				nil,
				binaryExpr(
					ident("a", p(1, 4)),
					intLit(5, p(1, 9)),
					token.Equal,
					p(1, 6)),
				blockStmt(
					p(1, 11), p(1, 12)),
				nil,
				p(1, 1)))
	})

	expect(t, "if a == 5 && b != 3 {}", func(p pfn) []ast.Stmt {
		return stmts(
			ifStmt(
				nil,
				binaryExpr(
					binaryExpr(
						ident("a", p(1, 4)),
						intLit(5, p(1, 9)),
						token.Equal,
						p(1, 6)),
					binaryExpr(
						ident("b", p(1, 14)),
						intLit(3, p(1, 19)),
						token.NotEqual,
						p(1, 16)),
					token.LAnd,
					p(1, 11)),
				blockStmt(
					p(1, 21), p(1, 22)),
				nil,
				p(1, 1)))
	})

	expect(t, "if a == 5 { a = 3; a = 1 }", func(p pfn) []ast.Stmt {
		return stmts(
			ifStmt(
				nil,
				binaryExpr(
					ident("a", p(1, 4)),
					intLit(5, p(1, 9)),
					token.Equal,
					p(1, 6)),
				blockStmt(
					p(1, 11), p(1, 26),
					assignStmt(
						exprs(ident("a", p(1, 13))),
						exprs(intLit(3, p(1, 17))),
						token.Assign,
						p(1, 15)),
					assignStmt(
						exprs(ident("a", p(1, 20))),
						exprs(intLit(1, p(1, 24))),
						token.Assign,
						p(1, 22))),
				nil,
				p(1, 1)))
	})

	expect(t, "if a == 5 { a = 3; a = 1 } else { a = 2; a = 4 }", func(p pfn) []ast.Stmt {
		return stmts(
			ifStmt(
				nil,
				binaryExpr(
					ident("a", p(1, 4)),
					intLit(5, p(1, 9)),
					token.Equal,
					p(1, 6)),
				blockStmt(
					p(1, 11), p(1, 26),
					assignStmt(
						exprs(ident("a", p(1, 13))),
						exprs(intLit(3, p(1, 17))),
						token.Assign,
						p(1, 15)),
					assignStmt(
						exprs(ident("a", p(1, 20))),
						exprs(intLit(1, p(1, 24))),
						token.Assign,
						p(1, 22))),
				blockStmt(
					p(1, 33), p(1, 48),
					assignStmt(
						exprs(ident("a", p(1, 35))),
						exprs(intLit(2, p(1, 39))),
						token.Assign,
						p(1, 37)),
					assignStmt(
						exprs(ident("a", p(1, 42))),
						exprs(intLit(4, p(1, 46))),
						token.Assign,
						p(1, 44))),
				p(1, 1)))
	})

	expect(t, `
if a == 5 { 
	b = 3 
	c = 1
} else if d == 3 { 
	e = 8
	f = 3
} else { 
	g = 2
	h = 4
}`, func(p pfn) []ast.Stmt {
		return stmts(
			ifStmt(
				nil,
				binaryExpr(
					ident("a", p(2, 4)),
					intLit(5, p(2, 9)),
					token.Equal,
					p(2, 6)),
				blockStmt(
					p(2, 11), p(5, 1),
					assignStmt(
						exprs(ident("b", p(3, 2))),
						exprs(intLit(3, p(3, 6))),
						token.Assign,
						p(3, 4)),
					assignStmt(
						exprs(ident("c", p(4, 2))),
						exprs(intLit(1, p(4, 6))),
						token.Assign,
						p(4, 4))),
				ifStmt(
					nil,
					binaryExpr(
						ident("d", p(5, 11)),
						intLit(3, p(5, 16)),
						token.Equal,
						p(5, 13)),
					blockStmt(
						p(5, 18), p(8, 1),
						assignStmt(
							exprs(ident("e", p(6, 2))),
							exprs(intLit(8, p(6, 6))),
							token.Assign,
							p(6, 4)),
						assignStmt(
							exprs(ident("f", p(7, 2))),
							exprs(intLit(3, p(7, 6))),
							token.Assign,
							p(7, 4))),
					blockStmt(
						p(8, 8), p(11, 1),
						assignStmt(
							exprs(ident("g", p(9, 2))),
							exprs(intLit(2, p(9, 6))),
							token.Assign,
							p(9, 4)),
						assignStmt(
							exprs(ident("h", p(10, 2))),
							exprs(intLit(4, p(10, 6))),
							token.Assign,
							p(10, 4))),
					p(5, 8)),
				p(2, 1)))
	})

	expect(t, "if a := 3; a < b {}", func(p pfn) []ast.Stmt {
		return stmts(
			ifStmt(
				assignStmt(
					exprs(ident("a", p(1, 4))),
					exprs(intLit(3, p(1, 9))),
					token.Define, p(1, 6)),
				binaryExpr(
					ident("a", p(1, 12)),
					ident("b", p(1, 16)),
					token.Less, p(1, 14)),
				blockStmt(
					p(1, 18), p(1, 19)),
				nil,
				p(1, 1)))
	})

	expect(t, "if a++; a < b {}", func(p pfn) []ast.Stmt {
		return stmts(
			ifStmt(
				incDecStmt(ident("a", p(1, 4)), token.Inc, p(1, 5)),
				binaryExpr(
					ident("a", p(1, 9)),
					ident("b", p(1, 13)),
					token.Less, p(1, 11)),
				blockStmt(
					p(1, 15), p(1, 16)),
				nil,
				p(1, 1)))
	})

	expectError(t, `if {}`)
	expectError(t, `if a == b { } else a != b { }`)
	expectError(t, `if a == b { } else if { }`)
	expectError(t, `else { }`)
	expectError(t, `if ; {}`)
	expectError(t, `if a := 3; {}`)
	expectError(t, `if ; a < 3 {}`)

}

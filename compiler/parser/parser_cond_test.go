package parser_test

import (
	"testing"

	"github.com/d5/tengo/compiler/ast"
)

func TestCondExpr(t *testing.T) {
	expect(t, "a ? b : c", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				condExpr(
					ident("a", p(1, 1)),
					ident("b", p(1, 5)),
					ident("c", p(1, 9)),
					p(1, 3),
					p(1, 7))))
	})
	expect(t, `a ?
b :
c`, func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				condExpr(
					ident("a", p(1, 1)),
					ident("b", p(1, 5)),
					ident("c", p(1, 9)),
					p(1, 3),
					p(1, 7))))
	})

	expectString(t, `a ? b : c`, "(a ? b : c)")
	expectString(t, `a + b ? c - d : e * f`, "((a + b) ? (c - d) : (e * f))")
	expectString(t, `a == b ? c + (d / e) : f ? g : h + i`, "((a == b) ? (c + ((d / e))) : (f ? g : (h + i)))")
	expectString(t, `(a + b) ? (c - d) : (e * f)`, "(((a + b)) ? ((c - d)) : ((e * f)))")
	expectString(t, `a + (b ? c : d) - e`, "((a + ((b ? c : d))) - e)")
	expectString(t, `a ? b ? c : d : e`, "(a ? (b ? c : d) : e)")
	expectString(t, `a := b ? c : d`, "a := (b ? c : d)")
	expectString(t, `x := a ? b ? c : d : e`, "x := (a ? (b ? c : d) : e)")

	// ? : should be at the end of each line if it's multi-line
	expectError(t, `a 
? b 
: c`)
	expectError(t, `a ? (b : e)`)
	expectError(t, `(a ? b) : e`)
}

package parser_test

import (
	"testing"

	"github.com/d5/tengo/ast"
	"github.com/d5/tengo/token"
)

func TestFunction(t *testing.T) {
	// TODO: function declaration currently not parsed.
	// All functions are parsed as function literal instead.
	// In Go, function declaration is parsed only at the top level.
	//expect(t, "func a(b, c, d) {}", func(p pfn) []ast.Stmt {
	//	return stmts(
	//		declStmt(
	//			funcDecl(
	//				ident("a", p(1, 6)),
	//				funcType(
	//					identList(p(1, 7), p(1, 15),
	//						ident("b", p(1, 8)),
	//						ident("c", p(1, 11)),
	//						ident("d", p(1, 14))),
	//					p(1, 12)),
	//				blockStmt(p(1, 17), p(1, 18)))))
	//})

	expect(t, "a = func(b, c, d) { return d }", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(
					ident("a", p(1, 1))),
				exprs(
					funcLit(
						funcType(
							identList(p(1, 9), p(1, 17),
								ident("b", p(1, 10)),
								ident("c", p(1, 13)),
								ident("d", p(1, 16))),
							p(1, 5)),
						blockStmt(p(1, 19), p(1, 30),
							returnStmt(p(1, 21), ident("d", p(1, 28)))))),
				token.Assign,
				p(1, 3)))
	})
}

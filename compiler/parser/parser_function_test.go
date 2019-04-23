package parser_test

import (
	"testing"

	"github.com/d5/tengo/compiler/ast"
	"github.com/d5/tengo/compiler/token"
)

func TestFunction(t *testing.T) {
	expect(t, "a = func(b, c, d) { return d }", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(
					ident("a", p(1, 1))),
				exprs(
					funcLit(
						funcType(
							identList(p(1, 9), p(1, 17), false,
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

func TestVariableFunction(t *testing.T) {
	expect(t, "a = func(...args) { return args[0] }", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(
					ident("a", p(1, 1))),
				exprs(
					funcLit(
						funcType(
							identList(
								p(1, 9), p(1, 17),
								true,
								ident("args", p(1, 13)),
							), p(1, 5)),
						blockStmt(p(1, 19), p(1, 36),
							returnStmt(p(1, 21),
								indexExpr(
									ident("args", p(1, 28)),
									intLit(0, p(1, 33)),
									p(1, 32), p(1, 34),
								),
							),
						),
					),
				),
				token.Assign,
				p(1, 3)))
	})
}

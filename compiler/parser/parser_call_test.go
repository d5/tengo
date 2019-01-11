package parser_test

import (
	"testing"

	"github.com/d5/tengo/compiler/ast"
	"github.com/d5/tengo/compiler/token"
)

func TestCall(t *testing.T) {
	expect(t, "add(1, 2, 3)", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				callExpr(
					ident("add", p(1, 1)),
					p(1, 4), p(1, 12),
					intLit(1, p(1, 5)),
					intLit(2, p(1, 8)),
					intLit(3, p(1, 11)))))
	})

	expect(t, "a = add(1, 2, 3)", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(
					ident("a", p(1, 1))),
				exprs(
					callExpr(
						ident("add", p(1, 5)),
						p(1, 8), p(1, 16),
						intLit(1, p(1, 9)),
						intLit(2, p(1, 12)),
						intLit(3, p(1, 15)))),
				token.Assign,
				p(1, 3)))
	})

	expect(t, "a, b = add(1, 2, 3)", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(
					ident("a", p(1, 1)),
					ident("b", p(1, 4))),
				exprs(
					callExpr(
						ident("add", p(1, 8)),
						p(1, 11), p(1, 19),
						intLit(1, p(1, 12)),
						intLit(2, p(1, 15)),
						intLit(3, p(1, 18)))),
				token.Assign,
				p(1, 6)))
	})

	expect(t, "add(a + 1, 2 * 1, (b + c))", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				callExpr(
					ident("add", p(1, 1)),
					p(1, 4), p(1, 26),
					binaryExpr(
						ident("a", p(1, 5)),
						intLit(1, p(1, 9)),
						token.Add,
						p(1, 7)),
					binaryExpr(
						intLit(2, p(1, 12)),
						intLit(1, p(1, 16)),
						token.Mul,
						p(1, 14)),
					parenExpr(
						binaryExpr(
							ident("b", p(1, 20)),
							ident("c", p(1, 24)),
							token.Add,
							p(1, 22)),
						p(1, 19), p(1, 25)))))
	})

	expectString(t, "a + add(b * c) + d", "((a + add((b * c))) + d)")
	expectString(t, "add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))")
	expectString(t, "f1(a) + f2(b) * f3(c)", "(f1(a) + (f2(b) * f3(c)))")
	expectString(t, "(f1(a) + f2(b)) * f3(c)", "(((f1(a) + f2(b))) * f3(c))")

	expect(t, "func(a, b) { a + b }(1, 2)", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				callExpr(
					funcLit(
						funcType(
							identList(
								p(1, 5), p(1, 10),
								ident("a", p(1, 6)),
								ident("b", p(1, 9))),
							p(1, 1)),
						blockStmt(
							p(1, 12), p(1, 20),
							exprStmt(
								binaryExpr(
									ident("a", p(1, 14)),
									ident("b", p(1, 18)),
									token.Add,
									p(1, 16))))),
					p(1, 21), p(1, 26),
					intLit(1, p(1, 22)),
					intLit(2, p(1, 25)))))
	})

	expect(t, `a.b()`, func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				callExpr(
					selectorExpr(
						ident("a", p(1, 1)),
						stringLit("b", p(1, 3))),
					p(1, 4), p(1, 5))))
	})

	expect(t, `a.b.c()`, func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				callExpr(
					selectorExpr(
						selectorExpr(
							ident("a", p(1, 1)),
							stringLit("b", p(1, 3))),
						stringLit("c", p(1, 5))),
					p(1, 6), p(1, 7))))
	})

	expect(t, `a["b"].c()`, func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				callExpr(
					selectorExpr(
						indexExpr(
							ident("a", p(1, 1)),
							stringLit("b", p(1, 3)),
							p(1, 2), p(1, 6)),
						stringLit("c", p(1, 8))),
					p(1, 9), p(1, 10))))
	})
}

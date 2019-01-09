package parser_test

import (
	"testing"

	"github.com/d5/tengo/ast"
	"github.com/d5/tengo/token"
)

func TestSelector(t *testing.T) {
	expect(t, "a.b", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				selectorExpr(
					ident("a", p(1, 1)),
					stringLit("b", p(1, 3)))))
	})

	expect(t, "a.b.c", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				selectorExpr(
					selectorExpr(
						ident("a", p(1, 1)),
						stringLit("b", p(1, 3))),
					stringLit("c", p(1, 5)))))
	})

	expect(t, "{k1:1}.k1", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				selectorExpr(
					mapLit(
						p(1, 1), p(1, 6),
						mapElementLit("k1", p(1, 2), p(1, 4), intLit(1, p(1, 5)))),
					stringLit("k1", p(1, 8)))))

	})
	expect(t, "{k1:{v1:1}}.k1.v1", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				selectorExpr(
					selectorExpr(
						mapLit(
							p(1, 1), p(1, 11),
							mapElementLit("k1", p(1, 2), p(1, 4),
								mapLit(p(1, 5), p(1, 10),
									mapElementLit("v1", p(1, 6), p(1, 8), intLit(1, p(1, 9)))))),
						stringLit("k1", p(1, 13))),
					stringLit("v1", p(1, 16)))))
	})

	expect(t, "a.b = 4", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(selectorExpr(ident("a", p(1, 1)), stringLit("b", p(1, 3)))),
				exprs(intLit(4, p(1, 7))),
				token.Assign, p(1, 5)))
	})

	expect(t, "a.b.c = 4", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(selectorExpr(selectorExpr(ident("a", p(1, 1)), stringLit("b", p(1, 3))), stringLit("c", p(1, 5)))),
				exprs(intLit(4, p(1, 9))),
				token.Assign, p(1, 7)))
	})

	expect(t, "a.b.c = 4 + 5", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(selectorExpr(selectorExpr(ident("a", p(1, 1)), stringLit("b", p(1, 3))), stringLit("c", p(1, 5)))),
				exprs(binaryExpr(intLit(4, p(1, 9)), intLit(5, p(1, 13)), token.Add, p(1, 11))),
				token.Assign, p(1, 7)))
	})

	expect(t, "a[0].c = 4", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(
					selectorExpr(
						indexExpr(
							ident("a", p(1, 1)),
							intLit(0, p(1, 3)),
							p(1, 2), p(1, 4)),
						stringLit("c", p(1, 6)))),
				exprs(intLit(4, p(1, 10))),
				token.Assign, p(1, 8)))
	})

	expect(t, "a.b[0].c = 4", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(
					selectorExpr(
						indexExpr(
							selectorExpr(
								ident("a", p(1, 1)),
								stringLit("b", p(1, 3))),
							intLit(0, p(1, 5)),
							p(1, 4), p(1, 6)),
						stringLit("c", p(1, 8)))),
				exprs(intLit(4, p(1, 12))),
				token.Assign, p(1, 10)))
	})

	expect(t, "a.b[0][2].c = 4", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(
					selectorExpr(
						indexExpr(
							indexExpr(
								selectorExpr(
									ident("a", p(1, 1)),
									stringLit("b", p(1, 3))),
								intLit(0, p(1, 5)),
								p(1, 4), p(1, 6)),
							intLit(2, p(1, 8)),
							p(1, 7), p(1, 9)),
						stringLit("c", p(1, 11)))),
				exprs(intLit(4, p(1, 15))),
				token.Assign, p(1, 13)))
	})

	expect(t, `a.b["key1"][2].c = 4`, func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(
					selectorExpr(
						indexExpr(
							indexExpr(
								selectorExpr(
									ident("a", p(1, 1)),
									stringLit("b", p(1, 3))),
								stringLit("key1", p(1, 5)),
								p(1, 4), p(1, 11)),
							intLit(2, p(1, 13)),
							p(1, 12), p(1, 14)),
						stringLit("c", p(1, 16)))),
				exprs(intLit(4, p(1, 20))),
				token.Assign, p(1, 18)))
	})

	expect(t, "a[0].b[2].c = 4", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(
					selectorExpr(
						indexExpr(
							selectorExpr(
								indexExpr(
									ident("a", p(1, 1)),
									intLit(0, p(1, 3)),
									p(1, 2), p(1, 4)),
								stringLit("b", p(1, 6))),
							intLit(2, p(1, 8)),
							p(1, 7), p(1, 9)),
						stringLit("c", p(1, 11)))),
				exprs(intLit(4, p(1, 15))),
				token.Assign, p(1, 13)))
	})

	expectError(t, `a.(b.c)`)
}

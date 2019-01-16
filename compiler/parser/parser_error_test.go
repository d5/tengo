package parser_test

import (
	"testing"

	"github.com/d5/tengo/compiler/ast"
	"github.com/d5/tengo/compiler/token"
)

func TestImport(t *testing.T) {
	expect(t, `error(1234)`, func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				errorExpr(p(1, 1), intLit(1234, p(1, 7)), p(1, 6), p(1, 11))))
	})

	expect(t, `err1 := error("some error")`, func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(ident("err1", p(1, 1))),
				exprs(errorExpr(p(1, 9), stringLit("some error", p(1, 15)), p(1, 14), p(1, 27))),
				token.Define, p(1, 6)))
	})

	expect(t, `return error("some error")`, func(p pfn) []ast.Stmt {
		return stmts(
			returnStmt(p(1, 1),
				errorExpr(p(1, 8), stringLit("some error", p(1, 14)), p(1, 13), p(1, 26))))
	})

	expect(t, `return error("some" + "error")`, func(p pfn) []ast.Stmt {
		return stmts(
			returnStmt(p(1, 1),
				errorExpr(p(1, 8),
					binaryExpr(
						stringLit("some", p(1, 14)),
						stringLit("error", p(1, 23)),
						token.Add, p(1, 21)),
					p(1, 13), p(1, 30))))
	})

	expectError(t, `error()`) // must have a value
}

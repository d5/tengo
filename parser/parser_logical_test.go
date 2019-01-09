package parser_test

import (
	"testing"

	"github.com/d5/tengo/ast"
	"github.com/d5/tengo/token"
)

func TestLogical(t *testing.T) {
	expect(t, "a && 5 || true", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				binaryExpr(
					binaryExpr(
						ident("a", p(1, 1)),
						intLit(5, p(1, 6)),
						token.LAnd,
						p(1, 3)),
					boolLit(true, p(1, 11)),
					token.LOr,
					p(1, 8))))
	})

	expect(t, "a || 5 && true", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				binaryExpr(
					ident("a", p(1, 1)),
					binaryExpr(
						intLit(5, p(1, 6)),
						boolLit(true, p(1, 11)),
						token.LAnd,
						p(1, 8)),
					token.LOr,
					p(1, 3))))
	})

	expect(t, "a && (5 || true)", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				binaryExpr(
					ident("a", p(1, 1)),
					parenExpr(
						binaryExpr(
							intLit(5, p(1, 7)),
							boolLit(true, p(1, 12)),
							token.LOr,
							p(1, 9)),
						p(1, 6), p(1, 16)),
					token.LAnd,
					p(1, 3))))
	})
}

package parser_test

import (
	"testing"

	"github.com/d5/tengo/ast"
	"github.com/d5/tengo/token"
)

func TestMap(t *testing.T) {
	expect(t, "{ key1: 1, key2: \"2\", key3: true }", func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				mapLit(p(1, 1), p(1, 34),
					mapElementLit("key1", p(1, 3), p(1, 7), intLit(1, p(1, 9))),
					mapElementLit("key2", p(1, 12), p(1, 16), stringLit("2", p(1, 18))),
					mapElementLit("key3", p(1, 23), p(1, 27), boolLit(true, p(1, 29))))))
	})

	expect(t, "a = { key1: 1, key2: \"2\", key3: true }", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(ident("a", p(1, 1))),
				exprs(mapLit(p(1, 5), p(1, 38),
					mapElementLit("key1", p(1, 7), p(1, 11), intLit(1, p(1, 13))),
					mapElementLit("key2", p(1, 16), p(1, 20), stringLit("2", p(1, 22))),
					mapElementLit("key3", p(1, 27), p(1, 31), boolLit(true, p(1, 33))))),
				token.Assign,
				p(1, 3)))
	})

	expect(t, "a = { key1: 1, key2: \"2\", key3: { k1: `bar`, k2: 4 } }", func(p pfn) []ast.Stmt {
		return stmts(
			assignStmt(
				exprs(ident("a", p(1, 1))),
				exprs(mapLit(p(1, 5), p(1, 54),
					mapElementLit("key1", p(1, 7), p(1, 11), intLit(1, p(1, 13))),
					mapElementLit("key2", p(1, 16), p(1, 20), stringLit("2", p(1, 22))),
					mapElementLit("key3", p(1, 27), p(1, 31),
						mapLit(p(1, 33), p(1, 52),
							mapElementLit("k1", p(1, 35), p(1, 37), stringLit("bar", p(1, 39))),
							mapElementLit("k2", p(1, 46), p(1, 48), intLit(4, p(1, 50))))))),
				token.Assign,
				p(1, 3)))
	})

	expect(t, `
{
	key1: 1, 
	key2: "2", 
	key3: true
}`, func(p pfn) []ast.Stmt {
		return stmts(
			exprStmt(
				mapLit(p(2, 1), p(6, 1),
					mapElementLit("key1", p(3, 2), p(3, 6), intLit(1, p(3, 8))),
					mapElementLit("key2", p(4, 2), p(4, 6), stringLit("2", p(4, 8))),
					mapElementLit("key3", p(5, 2), p(5, 6), boolLit(true, p(5, 8))))))
	})

	expectError(t, `{ key1: 1, }`)
	expectError(t, `{
key1: 1,
key2: 2,
}`)
}

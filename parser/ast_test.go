package parser_test

import (
	"strings"
	"testing"

	"github.com/d5/tengo/v2/parser"
	"github.com/d5/tengo/v2/require"
)

func TestIdentListString(t *testing.T) {
	identListVar := &parser.IdentList{
		List: []*parser.Ident{
			{Name: "a"},
			{Name: "b"},
			{Name: "c"},
		},
		VarArgs: true,
	}

	expectedVar := "(a, b, ...c)"
	if str := identListVar.String(); str != expectedVar {
		t.Fatalf("expected string of %#v to be %s, got %s",
			identListVar, expectedVar, str)
	}

	identList := &parser.IdentList{
		List: []*parser.Ident{
			{Name: "a"},
			{Name: "b"},
			{Name: "c"},
		},
		VarArgs: false,
	}

	expected := "(a, b, c)"
	if str := identList.String(); str != expected {
		t.Fatalf("expected string of %#v to be %s, got %s",
			identList, expected, str)
	}
}

func TestIdentListPosEnd(t *testing.T) {
	noParens := &parser.IdentList{
		List: []*parser.Ident{
			{Name: "a", NamePos: 6},
			{Name: "bc", NamePos: 9},
		},
	}
	require.Equal(t, parser.Pos(6), noParens.Pos())
	require.Equal(t, parser.Pos(11), noParens.End())

	empty := &parser.IdentList{}
	require.Equal(t, parser.NoPos, empty.Pos())
	require.Equal(t, parser.NoPos, empty.End())
}

func TestIdentListNumFields(t *testing.T) {
	var nilList *parser.IdentList
	require.Equal(t, 0, nilList.NumFields())

	for _, tc := range []struct {
		input    string
		expected int
	}{
		{"func() {}", 0},
		{"func(a) {}", 1},
		{"func(a, b, ...c) {}", 3},
	} {
		file, err := parseSource("test", []byte(tc.input), nil)
		require.NoError(t, err)
		fn := file.Stmts[0].(*parser.ExprStmt).Expr.(*parser.FuncLit)
		require.Equal(t, tc.expected, fn.Type.Params.NumFields())
	}
}

func TestNodePosEnd(t *testing.T) {
	for _, code := range []string{
		"[1, 2]",
		"a + b * c",
		"true",
		"false",
		"f(1, a...)",
		"'x'",
		"a ? b : c",
		"1.5",
		"func(a, b) { return a }",
		"foo",
		`import("fmt")`,
		"a[1]",
		"42",
		`{a: 1, "b": 2}`,
		"(a)",
		"a.bc",
		"a[1:2]",
		`"foo"`,
		"`foo`",
		"undefined",
		"a = 1",
		"a, b := 1, 2",
		"a += 1",
		";",
		"export a",
		"for {}",
		"for a < 1 { a++ }",
		"for i := 0; i < 1; i++ {}",
		"for k, v in m {}",
		"if a {}",
		"if a {} else { b }",
		"if a {} else if b {}",
		"a++",
		"a--",
		"return",
		"return a",
	} {
		t.Run(code, func(t *testing.T) {
			src := "x\n" + code + "\ny"
			file := parseNodes(t, src)
			stmt := file.Stmts[1]
			expectSpan(t, file, src, code, stmt)
			if es, ok := stmt.(*parser.ExprStmt); ok {
				expectSpan(t, file, src, code, es.Expr)
			}
		})
	}

	expr := func(s parser.Stmt) parser.Expr {
		return s.(*parser.ExprStmt).Expr
	}
	for _, tc := range []struct {
		code string
		text string
		node func(s parser.Stmt) parser.Node
	}{
		{"func(a, b) { return a }", "func(a, b)", func(s parser.Stmt) parser.Node {
			return expr(s).(*parser.FuncLit).Type
		}},
		{"func(a, b) { return a }", "(a, b)", func(s parser.Stmt) parser.Node {
			return expr(s).(*parser.FuncLit).Type.Params
		}},
		{"func(a, b) { return a }", "{ return a }", func(s parser.Stmt) parser.Node {
			return expr(s).(*parser.FuncLit).Body
		}},
		{"func(a, b) { return a }", "return a", func(s parser.Stmt) parser.Node {
			return expr(s).(*parser.FuncLit).Body.Stmts[0]
		}},
		{`{a: 1, "b": 2}`, "a: 1", func(s parser.Stmt) parser.Node {
			return expr(s).(*parser.MapLit).Elements[0]
		}},
		{`{a: 1, "b": 2}`, `"b": 2`, func(s parser.Stmt) parser.Node {
			return expr(s).(*parser.MapLit).Elements[1]
		}},
		{"a.bc", "bc", func(s parser.Stmt) parser.Node {
			return expr(s).(*parser.SelectorExpr).Sel
		}},
		{"for k, v in m {}", "k", func(s parser.Stmt) parser.Node {
			return s.(*parser.ForInStmt).Key
		}},
		{"for k, v in m {}", "v", func(s parser.Stmt) parser.Node {
			return s.(*parser.ForInStmt).Value
		}},
		{"for { break }", "break", func(s parser.Stmt) parser.Node {
			return s.(*parser.ForStmt).Body.Stmts[0]
		}},
		{"for { continue }", "continue", func(s parser.Stmt) parser.Node {
			return s.(*parser.ForStmt).Body.Stmts[0]
		}},
		{"if a {} else { b }", "{ b }", func(s parser.Stmt) parser.Node {
			return s.(*parser.IfStmt).Else
		}},
	} {
		t.Run(tc.code+" "+tc.text, func(t *testing.T) {
			src := "x\n" + tc.code + "\ny"
			file := parseNodes(t, src)
			expectSpan(t, file, src, tc.text, tc.node(file.Stmts[1]))
		})
	}

	src := "x\nerror(1)\nimmutable(1)\n-a\ny"
	file := parseNodes(t, src)
	require.Equal(t, file.InputFile.FileSetPos(strings.Index(src, "error")),
		expr(file.Stmts[1]).Pos())
	require.Equal(t, file.InputFile.FileSetPos(strings.Index(src, "immutable")),
		expr(file.Stmts[2]).Pos())
	require.Equal(t, file.InputFile.FileSetPos(strings.Index(src, "-a")+2),
		expr(file.Stmts[3]).End())
}

func TestFilePosEnd(t *testing.T) {
	fileSet := parser.NewFileSet()
	fileSet.AddFile("first", -1, 10)
	src := "a := 1\nb := 2"
	srcFile := fileSet.AddFile("second", -1, len(src))
	file, err := parser.NewParser(srcFile, []byte(src), nil).ParseFile()
	require.NoError(t, err)

	require.Equal(t, file.Stmts[0].Pos(), file.Pos())
	require.Equal(t, file.Stmts[1].End(), file.End())
	require.True(t, fileSet.File(file.Pos()) == srcFile)
	require.True(t, fileSet.File(file.End()) == srcFile)
}

func TestBadNodes(t *testing.T) {
	badExpr := &parser.BadExpr{From: 3, To: 7}
	require.Equal(t, parser.Pos(3), badExpr.Pos())
	require.Equal(t, parser.Pos(7), badExpr.End())
	require.Equal(t, "<bad expression>", badExpr.String())

	badStmt := &parser.BadStmt{From: 3, To: 7}
	require.Equal(t, parser.Pos(3), badStmt.Pos())
	require.Equal(t, parser.Pos(7), badStmt.End())
	require.Equal(t, "<bad statement>", badStmt.String())
}

func TestEmptyStmtEnd(t *testing.T) {
	explicit := &parser.EmptyStmt{Semicolon: 5}
	require.Equal(t, parser.Pos(6), explicit.End())

	implicit := &parser.EmptyStmt{Semicolon: 5, Implicit: true}
	require.Equal(t, parser.Pos(5), implicit.Pos())
	require.Equal(t, parser.Pos(5), implicit.End())
}

func TestNodeString(t *testing.T) {
	for _, tc := range []struct {
		input    string
		expected string
	}{
		{`[1, "a", 1.5]`, `[1, "a", 1.5]`},
		{"[]", "[]"},
		{"true", "true"},
		{"false", "false"},
		{"f()", "f()"},
		{"f(a, b...)", "f(a, b...)"},
		{"'x'", "'x'"},
		{`error("x")`, `error("x")`},
		{"1e3", "1e3"},
		{"func(a, ...b) { return a }", "func(a, ...b) {return a}"},
		{`import("fmt")`, `import("fmt")`},
		{"immutable([1, 2])", "immutable([1, 2])"},
		{"a[1]", "a[1]"},
		{"0x1F", "0x1F"},
		{"{a: 1, b: 2}", "{a: 1, b: 2}"},
		{"a.b.c", "a.b.c"},
		{"a[1:2]", "a[1:2]"},
		{"a[:2]", "a[:2]"},
		{"a[1:]", "a[1:]"},
		{"a[:]", "a[:]"},
		{`"foo"`, `"foo"`},
		{"`foo`", "`foo`"},
		{"-a", "(-a)"},
		{"!a", "(!a)"},
		{"^a", "(^a)"},
		{"+a", "(+a)"},
		{"undefined", "undefined"},
		{"a, b = 1, 2", "a, b = 1, 2"},
		{"a += 1", "a += 1"},
		{"1;;", "1; ;"},
		{"export a", "export a"},
		{"for {}", "for {}"},
		{"for a < 5 { a++ }", "for (a < 5) {a++}"},
		{"for { break; continue }", "for {break; continue}"},
		{"for k, v in m {}", "for k, v in m {}"},
		{"for v in m {}", "for _, v in m {}"},
		{"if a {}", "if a {}"},
		{"if a := f(); a { b-- }", "if a := f(); a {b--}"},
		{"if a {} else { b }", "if a {} else {b}"},
		{"if a {} else if b {}", "if a {} else if b {}"},
		{"return", "return"},
		{"return a", "return a"},
	} {
		expectParseString(t, tc.input, tc.expected)
	}

	file := parseNodes(t, "func(a, ...b) {}")
	fn := file.Stmts[0].(*parser.ExprStmt).Expr.(*parser.FuncLit)
	require.Equal(t, "func(a, ...b)", fn.Type.String())

	var nilIdent *parser.Ident
	require.Equal(t, "<null>", nilIdent.String())
}

func parseNodes(t *testing.T, src string) *parser.File {
	file, err := parseSource("test", []byte(src), nil)
	require.NoError(t, err)
	return file
}

func expectSpan(
	t *testing.T,
	file *parser.File,
	src, text string,
	node parser.Node,
) {
	start := strings.Index(src, text)
	require.True(t, start >= 0, "%q not found in %q", text, src)
	require.Equal(t, file.InputFile.FileSetPos(start), node.Pos())
	require.Equal(t, file.InputFile.FileSetPos(start+len(text)), node.End())
}

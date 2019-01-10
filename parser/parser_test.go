package parser_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/ast"
	"github.com/d5/tengo/parser"
	"github.com/d5/tengo/scanner"
	"github.com/d5/tengo/token"
)

type pfn func(int, int) scanner.Pos      // position conversion function
type expectedFn func(pos pfn) []ast.Stmt // callback function to return expected results

type tracer struct {
	out []string
}

func (o *tracer) Write(p []byte) (n int, err error) {
	o.out = append(o.out, string(p))
	return len(p), nil
}

//type slowPrinter struct {
//}
//
//func (o *slowPrinter) Write(p []byte) (n int, err error) {
//	fmt.Print(string(p))
//	time.Sleep(25 * time.Millisecond)
//	return len(p), nil
//}

func expect(t *testing.T, input string, fn expectedFn) (ok bool) {
	testFileSet := scanner.NewFileSet()
	testFile := testFileSet.AddFile("", -1, len(input))

	defer func() {
		if !ok {
			// print trace
			tr := &tracer{}
			actual, _ := parser.ParseFile(testFile, []byte(input), tr)
			if actual != nil {
				t.Logf("Parsed:\n%s", actual.String())
			}
			t.Logf("Trace:\n%s", strings.Join(tr.out, ""))
		}
	}()

	actual, err := parser.ParseFile(testFile, []byte(input), nil)
	if !assert.NoError(t, err) {
		return
	}

	expected := fn(func(line, column int) scanner.Pos {
		return scanner.Pos(int(testFile.LineStart(line)) + (column - 1))
	})

	if !assert.Equal(t, len(expected), len(actual.Stmts)) {
		return
	}

	for i := 0; i < len(expected); i++ {
		if !equalStmt(t, expected[i], actual.Stmts[i]) {
			return
		}
	}

	ok = true

	return
}

func expectError(t *testing.T, input string) (ok bool) {
	testFileSet := scanner.NewFileSet()
	testFile := testFileSet.AddFile("", -1, len(input))

	defer func() {
		if !ok {
			// print trace
			tr := &tracer{}
			_, _ = parser.ParseFile(testFile, []byte(input), tr)
			t.Logf("Trace:\n%s", strings.Join(tr.out, ""))
		}
	}()

	_, err := parser.ParseFile(testFile, []byte(input), nil)
	if !assert.Error(t, err) {
		return
	}

	ok = true

	return
}

func expectString(t *testing.T, input, expected string) (ok bool) {
	testFileSet := scanner.NewFileSet()
	testFile := testFileSet.AddFile("", -1, len(input))

	defer func() {
		if !ok {
			// print trace
			tr := &tracer{}
			_, _ = parser.ParseFile(testFile, []byte(input), tr)
			t.Logf("Trace:\n%s", strings.Join(tr.out, ""))
		}
	}()

	actual, err := parser.ParseFile(testFile, []byte(input), nil)
	if !assert.NoError(t, err) {
		return
	}

	if !assert.Equal(t, expected, actual.String()) {
		return
	}

	ok = true

	return
}

//func printTrace(input string) {
//	testFileSet := scanner.NewFileSet()
//	testFile := testFileSet.AddFile("", -1, len(input))
//
//	_, _ = parser.ParseFile(testFile, []byte(input), &slowPrinter{})
//}

func stmts(s ...ast.Stmt) []ast.Stmt {
	return s
}

func exprStmt(x ast.Expr) *ast.ExprStmt {
	return &ast.ExprStmt{Expr: x}
}

func assignStmt(lhs, rhs []ast.Expr, token token.Token, pos scanner.Pos) *ast.AssignStmt {
	return &ast.AssignStmt{Lhs: lhs, Rhs: rhs, Token: token, TokenPos: pos}
}

func emptyStmt(implicit bool, pos scanner.Pos) *ast.EmptyStmt {
	return &ast.EmptyStmt{Implicit: implicit, Semicolon: pos}
}

func returnStmt(pos scanner.Pos, results ...ast.Expr) *ast.ReturnStmt {
	return &ast.ReturnStmt{Results: results, ReturnPos: pos}
}

func forStmt(init ast.Stmt, cond ast.Expr, post ast.Stmt, body *ast.BlockStmt, pos scanner.Pos) *ast.ForStmt {
	return &ast.ForStmt{Cond: cond, Init: init, Post: post, Body: body, ForPos: pos}
}

func forInStmt(key, value *ast.Ident, seq ast.Expr, body *ast.BlockStmt, pos scanner.Pos) *ast.ForInStmt {
	return &ast.ForInStmt{Key: key, Value: value, Iterable: seq, Body: body, ForPos: pos}
}

func ifStmt(init ast.Stmt, cond ast.Expr, body *ast.BlockStmt, elseStmt ast.Stmt, pos scanner.Pos) *ast.IfStmt {
	return &ast.IfStmt{Init: init, Cond: cond, Body: body, Else: elseStmt, IfPos: pos}
}

func incDecStmt(expr ast.Expr, tok token.Token, pos scanner.Pos) *ast.IncDecStmt {
	return &ast.IncDecStmt{Expr: expr, Token: tok, TokenPos: pos}
}

func funcType(params *ast.IdentList, pos scanner.Pos) *ast.FuncType {
	return &ast.FuncType{Params: params, FuncPos: pos}
}

func blockStmt(lbrace, rbrace scanner.Pos, list ...ast.Stmt) *ast.BlockStmt {
	return &ast.BlockStmt{Stmts: list, LBrace: lbrace, RBrace: rbrace}
}

func ident(name string, pos scanner.Pos) *ast.Ident {
	return &ast.Ident{Name: name, NamePos: pos}
}

func identList(opening, closing scanner.Pos, list ...*ast.Ident) *ast.IdentList {
	return &ast.IdentList{List: list, LParen: opening, RParen: closing}
}

func binaryExpr(x, y ast.Expr, op token.Token, pos scanner.Pos) *ast.BinaryExpr {
	return &ast.BinaryExpr{Lhs: x, Rhs: y, Token: op, TokenPos: pos}
}

func unaryExpr(x ast.Expr, op token.Token, pos scanner.Pos) *ast.UnaryExpr {
	return &ast.UnaryExpr{Expr: x, Token: op, TokenPos: pos}
}

func exprs(list ...ast.Expr) []ast.Expr {
	return list
}

func intLit(value int64, pos scanner.Pos) *ast.IntLit {
	return &ast.IntLit{Value: value, ValuePos: pos}
}

func floatLit(value float64, pos scanner.Pos) *ast.FloatLit {
	return &ast.FloatLit{Value: value, ValuePos: pos}
}

func stringLit(value string, pos scanner.Pos) *ast.StringLit {
	return &ast.StringLit{Value: value, ValuePos: pos}
}

func charLit(value rune, pos scanner.Pos) *ast.CharLit {
	return &ast.CharLit{Value: value, ValuePos: pos, Literal: fmt.Sprintf("'%c'", value)}
}

func boolLit(value bool, pos scanner.Pos) *ast.BoolLit {
	return &ast.BoolLit{Value: value, ValuePos: pos}
}

func arrayLit(lbracket, rbracket scanner.Pos, list ...ast.Expr) *ast.ArrayLit {
	return &ast.ArrayLit{LBrack: lbracket, RBrack: rbracket, Elements: list}
}

func mapElementLit(key string, keyPos scanner.Pos, colonPos scanner.Pos, value ast.Expr) *ast.MapElementLit {
	return &ast.MapElementLit{Key: key, KeyPos: keyPos, ColonPos: colonPos, Value: value}
}

func mapLit(lbrace, rbrace scanner.Pos, list ...*ast.MapElementLit) *ast.MapLit {
	return &ast.MapLit{LBrace: lbrace, RBrace: rbrace, Elements: list}
}

func funcLit(funcType *ast.FuncType, body *ast.BlockStmt) *ast.FuncLit {
	return &ast.FuncLit{Type: funcType, Body: body}
}

func parenExpr(x ast.Expr, lparen, rparen scanner.Pos) *ast.ParenExpr {
	return &ast.ParenExpr{Expr: x, LParen: lparen, RParen: rparen}
}

func callExpr(f ast.Expr, lparen, rparen scanner.Pos, args ...ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{Func: f, LParen: lparen, RParen: rparen, Args: args}
}

func indexExpr(x, index ast.Expr, lbrack, rbrack scanner.Pos) *ast.IndexExpr {
	return &ast.IndexExpr{Expr: x, Index: index, LBrack: lbrack, RBrack: rbrack}
}

func sliceExpr(x, low, high ast.Expr, lbrack, rbrack scanner.Pos) *ast.SliceExpr {
	return &ast.SliceExpr{Expr: x, Low: low, High: high, LBrack: lbrack, RBrack: rbrack}
}

func selectorExpr(x, sel ast.Expr) *ast.SelectorExpr {
	return &ast.SelectorExpr{Expr: x, Sel: sel}
}

func equalStmt(t *testing.T, expected, actual ast.Stmt) bool {
	if expected == nil || reflect.ValueOf(expected).IsNil() {
		return assert.Nil(t, actual, "expected nil, but got not nil")
	}
	if !assert.NotNil(t, actual, "expected not nil, but got nil") {
		return false
	}
	if !assert.IsType(t, expected, actual) {
		return false
	}

	switch expected := expected.(type) {
	case *ast.ExprStmt:
		return equalExpr(t, expected.Expr, actual.(*ast.ExprStmt).Expr)
	case *ast.EmptyStmt:
		return assert.Equal(t, expected.Implicit, actual.(*ast.EmptyStmt).Implicit) &&
			assert.Equal(t, expected.Semicolon, actual.(*ast.EmptyStmt).Semicolon)
	case *ast.BlockStmt:
		return assert.Equal(t, expected.LBrace, actual.(*ast.BlockStmt).LBrace) &&
			assert.Equal(t, expected.RBrace, actual.(*ast.BlockStmt).RBrace) &&
			equalStmts(t, expected.Stmts, actual.(*ast.BlockStmt).Stmts)
	case *ast.AssignStmt:
		return equalExprs(t, expected.Lhs, actual.(*ast.AssignStmt).Lhs) &&
			equalExprs(t, expected.Rhs, actual.(*ast.AssignStmt).Rhs) &&
			assert.Equal(t, int(expected.Token), int(actual.(*ast.AssignStmt).Token)) &&
			assert.Equal(t, int(expected.TokenPos), int(actual.(*ast.AssignStmt).TokenPos))
	case *ast.IfStmt:
		return equalStmt(t, expected.Init, actual.(*ast.IfStmt).Init) &&
			equalExpr(t, expected.Cond, actual.(*ast.IfStmt).Cond) &&
			equalStmt(t, expected.Body, actual.(*ast.IfStmt).Body) &&
			equalStmt(t, expected.Else, actual.(*ast.IfStmt).Else) &&
			assert.Equal(t, expected.IfPos, actual.(*ast.IfStmt).IfPos)
	case *ast.IncDecStmt:
		return equalExpr(t, expected.Expr, actual.(*ast.IncDecStmt).Expr) &&
			assert.Equal(t, expected.Token, actual.(*ast.IncDecStmt).Token) &&
			assert.Equal(t, expected.TokenPos, actual.(*ast.IncDecStmt).TokenPos)
	case *ast.ForStmt:
		return equalStmt(t, expected.Init, actual.(*ast.ForStmt).Init) &&
			equalExpr(t, expected.Cond, actual.(*ast.ForStmt).Cond) &&
			equalStmt(t, expected.Post, actual.(*ast.ForStmt).Post) &&
			equalStmt(t, expected.Body, actual.(*ast.ForStmt).Body) &&
			assert.Equal(t, expected.ForPos, actual.(*ast.ForStmt).ForPos)
	case *ast.ForInStmt:
		return equalExpr(t, expected.Key, actual.(*ast.ForInStmt).Key) &&
			equalExpr(t, expected.Value, actual.(*ast.ForInStmt).Value) &&
			equalExpr(t, expected.Iterable, actual.(*ast.ForInStmt).Iterable) &&
			equalStmt(t, expected.Body, actual.(*ast.ForInStmt).Body) &&
			assert.Equal(t, expected.ForPos, actual.(*ast.ForInStmt).ForPos)
	case *ast.ReturnStmt:
		return equalExprs(t, expected.Results, actual.(*ast.ReturnStmt).Results) &&
			assert.Equal(t, expected.ReturnPos, actual.(*ast.ReturnStmt).ReturnPos)
	case *ast.BranchStmt:
		return equalExpr(t, expected.Label, actual.(*ast.BranchStmt).Label) &&
			assert.Equal(t, expected.Token, actual.(*ast.BranchStmt).Token) &&
			assert.Equal(t, expected.TokenPos, actual.(*ast.BranchStmt).TokenPos)
	default:
		panic(fmt.Errorf("unknown type: %T", expected))
	}
}

func equalExpr(t *testing.T, expected, actual ast.Expr) bool {
	if expected == nil || reflect.ValueOf(expected).IsNil() {
		return assert.Nil(t, actual, "expected nil, but got not nil")
	}
	if !assert.NotNil(t, actual, "expected not nil, but got nil") {
		return false
	}
	if !assert.IsType(t, expected, actual) {
		return false
	}

	switch expected := expected.(type) {
	case *ast.Ident:
		return assert.Equal(t, expected.Name, actual.(*ast.Ident).Name) &&
			assert.Equal(t, int(expected.NamePos), int(actual.(*ast.Ident).NamePos))
	case *ast.IntLit:
		return assert.Equal(t, expected.Value, actual.(*ast.IntLit).Value) &&
			assert.Equal(t, int(expected.ValuePos), int(actual.(*ast.IntLit).ValuePos))
	case *ast.FloatLit:
		return assert.Equal(t, expected.Value, actual.(*ast.FloatLit).Value) &&
			assert.Equal(t, int(expected.ValuePos), int(actual.(*ast.FloatLit).ValuePos))
	case *ast.BoolLit:
		return assert.Equal(t, expected.Value, actual.(*ast.BoolLit).Value) &&
			assert.Equal(t, int(expected.ValuePos), int(actual.(*ast.BoolLit).ValuePos))
	case *ast.CharLit:
		return assert.Equal(t, expected.Value, actual.(*ast.CharLit).Value) &&
			assert.Equal(t, int(expected.ValuePos), int(actual.(*ast.CharLit).ValuePos))
	case *ast.StringLit:
		return assert.Equal(t, expected.Value, actual.(*ast.StringLit).Value) &&
			assert.Equal(t, int(expected.ValuePos), int(actual.(*ast.StringLit).ValuePos))
	case *ast.ArrayLit:
		return assert.Equal(t, expected.LBrack, actual.(*ast.ArrayLit).LBrack) &&
			assert.Equal(t, expected.RBrack, actual.(*ast.ArrayLit).RBrack) &&
			equalExprs(t, expected.Elements, actual.(*ast.ArrayLit).Elements)
	case *ast.MapLit:
		return assert.Equal(t, expected.LBrace, actual.(*ast.MapLit).LBrace) &&
			assert.Equal(t, expected.RBrace, actual.(*ast.MapLit).RBrace) &&
			equalMapElements(t, expected.Elements, actual.(*ast.MapLit).Elements)
	case *ast.BinaryExpr:
		return equalExpr(t, expected.Lhs, actual.(*ast.BinaryExpr).Lhs) &&
			equalExpr(t, expected.Rhs, actual.(*ast.BinaryExpr).Rhs) &&
			assert.Equal(t, expected.Token, actual.(*ast.BinaryExpr).Token) &&
			assert.Equal(t, expected.TokenPos, actual.(*ast.BinaryExpr).TokenPos)
	case *ast.UnaryExpr:
		return equalExpr(t, expected.Expr, actual.(*ast.UnaryExpr).Expr) &&
			assert.Equal(t, expected.Token, actual.(*ast.UnaryExpr).Token) &&
			assert.Equal(t, expected.TokenPos, actual.(*ast.UnaryExpr).TokenPos)
	case *ast.FuncLit:
		return equalFuncType(t, expected.Type, actual.(*ast.FuncLit).Type) &&
			equalStmt(t, expected.Body, actual.(*ast.FuncLit).Body)
	case *ast.CallExpr:
		return equalExpr(t, expected.Func, actual.(*ast.CallExpr).Func) &&
			assert.Equal(t, expected.LParen, actual.(*ast.CallExpr).LParen) &&
			assert.Equal(t, expected.RParen, actual.(*ast.CallExpr).RParen) &&
			equalExprs(t, expected.Args, actual.(*ast.CallExpr).Args)
	case *ast.ParenExpr:
		return equalExpr(t, expected.Expr, actual.(*ast.ParenExpr).Expr) &&
			assert.Equal(t, expected.LParen, actual.(*ast.ParenExpr).LParen) &&
			assert.Equal(t, expected.RParen, actual.(*ast.ParenExpr).RParen)
	case *ast.IndexExpr:
		return equalExpr(t, expected.Expr, actual.(*ast.IndexExpr).Expr) &&
			equalExpr(t, expected.Index, actual.(*ast.IndexExpr).Index) &&
			assert.Equal(t, expected.LBrack, actual.(*ast.IndexExpr).LBrack) &&
			assert.Equal(t, expected.RBrack, actual.(*ast.IndexExpr).RBrack)
	case *ast.SliceExpr:
		return equalExpr(t, expected.Expr, actual.(*ast.SliceExpr).Expr) &&
			equalExpr(t, expected.Low, actual.(*ast.SliceExpr).Low) &&
			equalExpr(t, expected.High, actual.(*ast.SliceExpr).High) &&
			assert.Equal(t, expected.LBrack, actual.(*ast.SliceExpr).LBrack) &&
			assert.Equal(t, expected.RBrack, actual.(*ast.SliceExpr).RBrack)
	case *ast.SelectorExpr:
		return equalExpr(t, expected.Expr, actual.(*ast.SelectorExpr).Expr) &&
			equalExpr(t, expected.Sel, actual.(*ast.SelectorExpr).Sel)
	default:
		panic(fmt.Errorf("unknown type: %T", expected))
	}
}

func equalFuncType(t *testing.T, expected, actual *ast.FuncType) bool {
	return assert.Equal(t, expected.Params.LParen, actual.Params.LParen) &&
		assert.Equal(t, expected.Params.RParen, actual.Params.RParen) &&
		equalIdents(t, expected.Params.List, actual.Params.List)
}

func equalIdents(t *testing.T, expected, actual []*ast.Ident) bool {
	if !assert.Equal(t, len(expected), len(actual)) {
		return false
	}

	for i := 0; i < len(expected); i++ {
		if !equalExpr(t, expected[i], actual[i]) {
			return false
		}
	}

	return true
}

func equalExprs(t *testing.T, expected, actual []ast.Expr) bool {
	if !assert.Equal(t, len(expected), len(actual)) {
		return false
	}

	for i := 0; i < len(expected); i++ {
		if !equalExpr(t, expected[i], actual[i]) {
			return false
		}
	}

	return true
}

func equalStmts(t *testing.T, expected, actual []ast.Stmt) bool {
	if !assert.Equal(t, len(expected), len(actual)) {
		return false
	}

	for i := 0; i < len(expected); i++ {
		if !equalStmt(t, expected[i], actual[i]) {
			return false
		}
	}

	return true
}

func equalMapElements(t *testing.T, expected, actual []*ast.MapElementLit) bool {
	if !assert.Equal(t, len(expected), len(actual)) {
		return false
	}

	for i := 0; i < len(expected); i++ {
		if !assert.Equal(t, expected[i].Key, actual[i].Key) ||
			!assert.Equal(t, expected[i].KeyPos, actual[i].KeyPos) ||
			!assert.Equal(t, expected[i].ColonPos, actual[i].ColonPos) ||
			!equalExpr(t, expected[i].Value, actual[i].Value) {
			return false
		}
	}

	return true
}

package token_test

import (
	"testing"

	"github.com/d5/tengo/v2/require"
	"github.com/d5/tengo/v2/token"
)

var keywords = []struct {
	ident string
	tok   token.Token
}{
	{"break", token.Break},
	{"continue", token.Continue},
	{"else", token.Else},
	{"for", token.For},
	{"func", token.Func},
	{"error", token.Error},
	{"immutable", token.Immutable},
	{"if", token.If},
	{"return", token.Return},
	{"export", token.Export},
	{"true", token.True},
	{"false", token.False},
	{"in", token.In},
	{"undefined", token.Undefined},
	{"import", token.Import},
}

func TestTokenString(t *testing.T) {
	for _, tc := range []struct {
		tok      token.Token
		expected string
	}{
		{token.Illegal, "ILLEGAL"},
		{token.EOF, "EOF"},
		{token.Comment, "COMMENT"},
		{token.Ident, "IDENT"},
		{token.Int, "INT"},
		{token.Float, "FLOAT"},
		{token.Char, "CHAR"},
		{token.String, "STRING"},
		{token.Add, "+"},
		{token.AndNot, "&^"},
		{token.AndNotAssign, "&^="},
		{token.ShlAssign, "<<="},
		{token.LAnd, "&&"},
		{token.LOr, "||"},
		{token.Inc, "++"},
		{token.NotEqual, "!="},
		{token.Define, ":="},
		{token.Ellipsis, "..."},
		{token.Semicolon, ";"},
		{token.Question, "?"},
		{token.Token(-1), "token(-1)"},
		{token.Token(1000), "token(1000)"},
	} {
		require.Equal(t, tc.expected, tc.tok.String())
	}

	for _, kw := range keywords {
		require.Equal(t, kw.ident, kw.tok.String())
	}
}

func TestTokenPrecedence(t *testing.T) {
	for _, tc := range []struct {
		tok      token.Token
		expected int
	}{
		{token.LOr, 1},
		{token.LAnd, 2},
		{token.Equal, 3},
		{token.NotEqual, 3},
		{token.Less, 3},
		{token.LessEq, 3},
		{token.Greater, 3},
		{token.GreaterEq, 3},
		{token.Add, 4},
		{token.Sub, 4},
		{token.Or, 4},
		{token.Xor, 4},
		{token.Mul, 5},
		{token.Quo, 5},
		{token.Rem, 5},
		{token.Shl, 5},
		{token.Shr, 5},
		{token.And, 5},
		{token.AndNot, 5},
		{token.Assign, token.LowestPrec},
		{token.Define, token.LowestPrec},
		{token.AddAssign, token.LowestPrec},
		{token.Not, token.LowestPrec},
		{token.Inc, token.LowestPrec},
		{token.Question, token.LowestPrec},
		{token.Colon, token.LowestPrec},
		{token.LParen, token.LowestPrec},
		{token.Period, token.LowestPrec},
		{token.Ident, token.LowestPrec},
		{token.In, token.LowestPrec},
		{token.EOF, token.LowestPrec},
	} {
		require.Equal(t, tc.expected, tc.tok.Precedence(), "%s", tc.tok)
	}
}

func TestTokenClass(t *testing.T) {
	for _, tc := range []struct {
		tok      token.Token
		literal  bool
		operator bool
		keyword  bool
	}{
		{token.Illegal, false, false, false},
		{token.EOF, false, false, false},
		{token.Comment, false, false, false},
		{token.Ident, true, false, false},
		{token.Int, true, false, false},
		{token.Float, true, false, false},
		{token.Char, true, false, false},
		{token.String, true, false, false},
		{token.Add, false, true, false},
		{token.Assign, false, true, false},
		{token.LParen, false, true, false},
		{token.Semicolon, false, true, false},
		{token.Question, false, true, false},
		{token.Break, false, false, true},
		{token.True, false, false, true},
		{token.Import, false, false, true},
		{token.Token(-1), false, false, false},
		{token.Token(1000), false, false, false},
	} {
		require.Equal(t, tc.literal, tc.tok.IsLiteral(), "%s", tc.tok)
		require.Equal(t, tc.operator, tc.tok.IsOperator(), "%s", tc.tok)
		require.Equal(t, tc.keyword, tc.tok.IsKeyword(), "%s", tc.tok)
	}
}

func TestLookup(t *testing.T) {
	for _, kw := range keywords {
		require.Equal(t, kw.tok, token.Lookup(kw.ident))
		require.True(t, kw.tok.IsKeyword(), kw.ident)
	}

	for _, ident := range []string{
		"", "foo", "_", "Break", "IF", "fors", "nil", "EOF", "IDENT", "+",
	} {
		require.Equal(t, token.Ident, token.Lookup(ident), ident)
	}

	n := 0
	for tok := token.Token(0); tok < 1000; tok++ {
		if tok.IsKeyword() {
			n++
			require.Equal(t, tok, token.Lookup(tok.String()))
		}
	}
	require.Equal(t, len(keywords), n)
}

package parser_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/d5/tengo/v2/parser"
	"github.com/d5/tengo/v2/require"
	"github.com/d5/tengo/v2/token"
)

var testFileSet = parser.NewFileSet()

type scanResult struct {
	Token   token.Token
	Literal string
	Line    int
	Column  int
}

func TestScanner_Scan(t *testing.T) {
	var testCases = [...]struct {
		token   token.Token
		literal string
	}{
		{token.Comment, "/* a comment */"},
		{token.Comment, "// a comment \n"},
		{token.Comment, "/*\r*/"},
		{token.Comment, "/**\r/*/"},
		{token.Comment, "/**\r\r/*/"},
		{token.Comment, "//\r\n"},
		{token.Ident, "foobar"},
		{token.Ident, "a۰۱۸"},
		{token.Ident, "foo६४"},
		{token.Ident, "bar９８７６"},
		{token.Ident, "ŝ"},
		{token.Ident, "ŝfoo"},
		{token.Int, "0"},
		{token.Int, "1"},
		{token.Int, "123456789012345678890"},
		{token.Int, "01234567"},
		{token.Int, "0xcafebabe"},
		{token.Float, "0."},
		{token.Float, ".0"},
		{token.Float, "3.14159265"},
		{token.Float, "1e0"},
		{token.Float, "1e+100"},
		{token.Float, "1e-100"},
		{token.Float, "2.71828e-1000"},
		{token.Char, "'a'"},
		{token.Char, "'\\000'"},
		{token.Char, "'\\xFF'"},
		{token.Char, "'\\uff16'"},
		{token.Char, "'\\U0000ff16'"},
		{token.String, "`foobar`"},
		{token.String, "`" + `foo
	                        bar` +
			"`",
		},
		{token.String, "`\r`"},
		{token.String, "`foo\r\nbar`"},
		{token.Add, "+"},
		{token.Sub, "-"},
		{token.Mul, "*"},
		{token.Quo, "/"},
		{token.Rem, "%"},
		{token.And, "&"},
		{token.Or, "|"},
		{token.Xor, "^"},
		{token.Shl, "<<"},
		{token.Shr, ">>"},
		{token.AndNot, "&^"},
		{token.AddAssign, "+="},
		{token.SubAssign, "-="},
		{token.MulAssign, "*="},
		{token.QuoAssign, "/="},
		{token.RemAssign, "%="},
		{token.AndAssign, "&="},
		{token.OrAssign, "|="},
		{token.XorAssign, "^="},
		{token.ShlAssign, "<<="},
		{token.ShrAssign, ">>="},
		{token.AndNotAssign, "&^="},
		{token.LAnd, "&&"},
		{token.LOr, "||"},
		{token.Inc, "++"},
		{token.Dec, "--"},
		{token.Equal, "=="},
		{token.Less, "<"},
		{token.Greater, ">"},
		{token.Assign, "="},
		{token.Not, "!"},
		{token.NotEqual, "!="},
		{token.LessEq, "<="},
		{token.GreaterEq, ">="},
		{token.Define, ":="},
		{token.Ellipsis, "..."},
		{token.LParen, "("},
		{token.LBrack, "["},
		{token.LBrace, "{"},
		{token.Comma, ","},
		{token.Period, "."},
		{token.RParen, ")"},
		{token.RBrack, "]"},
		{token.RBrace, "}"},
		{token.Semicolon, ";"},
		{token.Colon, ":"},
		{token.Break, "break"},
		{token.Continue, "continue"},
		{token.Else, "else"},
		{token.For, "for"},
		{token.Func, "func"},
		{token.If, "if"},
		{token.Return, "return"},
		{token.Export, "export"},
	}

	// combine
	var lines []string
	var lineSum int
	lineNos := make([]int, len(testCases))
	columnNos := make([]int, len(testCases))
	for i, tc := range testCases {
		// add 0-2 lines before each test case
		emptyLines := rand.Intn(3)
		for j := 0; j < emptyLines; j++ {
			lines = append(lines, strings.Repeat(" ", rand.Intn(10)))
		}

		// add test case line with some whitespaces around it
		emptyColumns := rand.Intn(10)
		lines = append(lines, fmt.Sprintf("%s%s%s",
			strings.Repeat(" ", emptyColumns),
			tc.literal,
			strings.Repeat(" ", rand.Intn(10))))

		lineNos[i] = lineSum + emptyLines + 1
		lineSum += emptyLines + countLines(tc.literal)
		columnNos[i] = emptyColumns + 1
	}

	// expected results
	var expected []scanResult
	var expectedSkipComments []scanResult
	for i, tc := range testCases {
		// expected literal
		var expectedLiteral string
		switch tc.token {
		case token.Comment:
			// strip CRs in comments
			expectedLiteral = string(parser.StripCR([]byte(tc.literal),
				tc.literal[1] == '*'))

			//-style comment literal doesn't contain newline
			if expectedLiteral[1] == '/' {
				expectedLiteral = expectedLiteral[:len(expectedLiteral)-1]
			}
		case token.Ident:
			expectedLiteral = tc.literal
		case token.Semicolon:
			expectedLiteral = ";"
		default:
			if tc.token.IsLiteral() {
				// strip CRs in raw string
				expectedLiteral = tc.literal
				if expectedLiteral[0] == '`' {
					expectedLiteral = string(parser.StripCR(
						[]byte(expectedLiteral), false))
				}
			} else if tc.token.IsKeyword() {
				expectedLiteral = tc.literal
			}
		}

		res := scanResult{
			Token:   tc.token,
			Literal: expectedLiteral,
			Line:    lineNos[i],
			Column:  columnNos[i],
		}

		expected = append(expected, res)
		if tc.token != token.Comment {
			expectedSkipComments = append(expectedSkipComments, res)
		}
	}

	scanExpect(t, strings.Join(lines, "\n"),
		parser.ScanComments|parser.DontInsertSemis, expected...)
	scanExpect(t, strings.Join(lines, "\n"),
		parser.DontInsertSemis, expectedSkipComments...)
}

func TestScanner_SemicolonBeforeComment(t *testing.T) {
	scanExpect(t, "a // c\nb", 0,
		scanResult{token.Ident, "a", 1, 1},
		scanResult{token.Semicolon, "\n", 1, 3},
		scanResult{token.Ident, "b", 2, 1},
		scanResult{token.Semicolon, "\n", 2, 2})
	scanExpect(t, "a // c\nb", parser.ScanComments,
		scanResult{token.Ident, "a", 1, 1},
		scanResult{token.Semicolon, "\n", 1, 3},
		scanResult{token.Comment, "// c", 1, 3},
		scanResult{token.Ident, "b", 2, 1},
		scanResult{token.Semicolon, "\n", 2, 2})
	scanExpect(t, "a /* c */", 0,
		scanResult{token.Ident, "a", 1, 1},
		scanResult{token.Semicolon, "\n", 1, 3})
	scanExpect(t, "a /* c */\nb", 0,
		scanResult{token.Ident, "a", 1, 1},
		scanResult{token.Semicolon, "\n", 1, 3},
		scanResult{token.Ident, "b", 2, 1},
		scanResult{token.Semicolon, "\n", 2, 2})
	scanExpect(t, "a /* c\n */ b", 0,
		scanResult{token.Ident, "a", 1, 1},
		scanResult{token.Semicolon, "\n", 1, 3},
		scanResult{token.Ident, "b", 2, 5},
		scanResult{token.Semicolon, "\n", 2, 6})
	scanExpect(t, "a /* c */ // d\nb", 0,
		scanResult{token.Ident, "a", 1, 1},
		scanResult{token.Semicolon, "\n", 1, 3},
		scanResult{token.Ident, "b", 2, 1},
		scanResult{token.Semicolon, "\n", 2, 2})

	scanExpect(t, "a /* c */ b", 0,
		scanResult{token.Ident, "a", 1, 1},
		scanResult{token.Ident, "b", 1, 11},
		scanResult{token.Semicolon, "\n", 1, 12})
	scanExpect(t, "a /* c */ /* d */ b", 0,
		scanResult{token.Ident, "a", 1, 1},
		scanResult{token.Ident, "b", 1, 19},
		scanResult{token.Semicolon, "\n", 1, 20})
	scanExpect(t, "a /* c */ / b", 0,
		scanResult{token.Ident, "a", 1, 1},
		scanResult{token.Quo, "", 1, 11},
		scanResult{token.Ident, "b", 1, 13},
		scanResult{token.Semicolon, "\n", 1, 14})
}

func TestScanner_PeriodAtEOF(t *testing.T) {
	scanExpect(t, ".", 0,
		scanResult{token.Period, "", 1, 1})
	scanExpect(t, "a.", 0,
		scanResult{token.Ident, "a", 1, 1},
		scanResult{token.Period, "", 1, 2})
	scanExpect(t, "a..", 0,
		scanResult{token.Ident, "a", 1, 1},
		scanResult{token.Period, "", 1, 2},
		scanResult{token.Period, "", 1, 3})
}

func TestScanner_ScanEscape(t *testing.T) {
	for _, tc := range []struct {
		input    string
		expected []string
	}{
		{`"\a\b\f\n\r\t\v\\\""`, nil},
		{`'\''`, nil},
		{`'\\'`, nil},
		{`"\101\x41A\U00000041"`, nil},
		{`"￿\U0010FFFF\377"`, nil},
		{`"\z"`, []string{"test:1:3: unknown escape sequence"}},
		{`"\'"`, []string{"test:1:3: unknown escape sequence"}},
		{`'\"'`, []string{"test:1:3: unknown escape sequence"}},
		{`"\x4g"`, []string{
			"test:1:5: illegal character U+0067 'g' in escape sequence"}},
		{`"\08"`, []string{
			"test:1:4: illegal character U+0038 '8' in escape sequence"}},
		{`"\400"`, []string{
			"test:1:3: escape sequence is invalid Unicode code point"}},
		{`"\uD800"`, []string{
			"test:1:3: escape sequence is invalid Unicode code point"}},
		{`"\U00110000"`, []string{
			"test:1:3: escape sequence is invalid Unicode code point"}},
		{`"\x4`, []string{
			"test:1:5: escape sequence not terminated",
			"test:1:1: string literal not terminated"}},
		{`"\`, []string{
			"test:1:3: escape sequence not terminated",
			"test:1:1: string literal not terminated"}},
	} {
		require.Equal(t, tc.expected, scanErrors(tc.input), tc.input)
	}
}

func TestStripCR(t *testing.T) {
	for _, tc := range []struct {
		input  string
		expect string
	}{
		{"//\n", "//\n"},
		{"//\r\n", "//\n"},
		{"//\r\r\r\n", "//\n"},
		{"//\r*\r/\r\n", "//*/\n"},
		{"/**/", "/**/"},
		{"/*\r/*/", "/*/*/"},
		{"/*\r*/", "/**/"},
		{"/**\r/*/", "/**\r/*/"},
		{"/*\r/\r*\r/*/", "/*/*\r/*/"},
		{"/*\r\r\r\r*/", "/**/"},
	} {
		actual := string(parser.StripCR([]byte(tc.input),
			len(tc.input) >= 2 && tc.input[1] == '*'))
		require.Equal(t, tc.expect, actual)
	}
}

func scanExpect(
	t *testing.T,
	input string,
	mode parser.ScanMode,
	expected ...scanResult,
) {
	testFile := testFileSet.AddFile("test", -1, len(input))

	s := parser.NewScanner(
		testFile,
		[]byte(input),
		func(_ parser.SourceFilePos, msg string) { require.Fail(t, msg) },
		mode)

	for idx, e := range expected {
		tok, literal, pos := s.Scan()

		filePos := testFile.Position(pos)

		require.Equal(t, e.Token, tok, "[%d] expected: %s, actual: %s",
			idx, e.Token.String(), tok.String())
		require.Equal(t, e.Literal, literal)
		require.Equal(t, e.Line, filePos.Line)
		require.Equal(t, e.Column, filePos.Column)
	}

	tok, _, _ := s.Scan()
	require.Equal(t, token.EOF, tok, "more tokens left")
	require.Equal(t, 0, s.ErrorCount())
}

func scanErrors(input string) (errs []string) {
	testFile := testFileSet.AddFile("test", -1, len(input))
	s := parser.NewScanner(
		testFile,
		[]byte(input),
		func(pos parser.SourceFilePos, msg string) {
			errs = append(errs, pos.String()+": "+msg)
		},
		0)
	for {
		if tok, _, _ := s.Scan(); tok == token.EOF {
			return
		}
	}
}

func countLines(s string) int {
	if s == "" {
		return 0
	}
	n := 1
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			n++
		}
	}
	return n
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

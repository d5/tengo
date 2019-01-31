package scanner_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/compiler/scanner"
	"github.com/d5/tengo/compiler/source"
	"github.com/d5/tengo/compiler/token"
)

var testFileSet = source.NewFileSet()

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
			expectedLiteral = string(scanner.StripCR([]byte(tc.literal), tc.literal[1] == '*'))

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
					expectedLiteral = string(scanner.StripCR([]byte(expectedLiteral), false))
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

	scanExpect(t, strings.Join(lines, "\n"), scanner.ScanComments|scanner.DontInsertSemis, expected...)
	scanExpect(t, strings.Join(lines, "\n"), scanner.DontInsertSemis, expectedSkipComments...)
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
		actual := string(scanner.StripCR([]byte(tc.input), len(tc.input) >= 2 && tc.input[1] == '*'))
		assert.Equal(t, tc.expect, actual)
	}
}

func scanExpect(t *testing.T, input string, mode scanner.Mode, expected ...scanResult) bool {
	testFile := testFileSet.AddFile("", testFileSet.Base(), len(input))

	s := scanner.NewScanner(
		testFile,
		[]byte(input),
		func(_ source.FilePos, msg string) { assert.Fail(t, msg) },
		mode)

	for idx, e := range expected {
		tok, literal, pos := s.Scan()

		filePos := testFile.Position(pos)

		if !assert.Equal(t, e.Token, tok, "[%d] expected: %s, actual: %s", idx, e.Token.String(), tok.String()) ||
			!assert.Equal(t, e.Literal, literal) ||
			!assert.Equal(t, e.Line, filePos.Line) ||
			!assert.Equal(t, e.Column, filePos.Column) {
			return false
		}
	}

	tok, _, _ := s.Scan()
	assert.Equal(t, token.EOF, tok, "more tokens left")

	return assert.Equal(t, 0, s.ErrorCount())
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

package token

import "strconv"

// Token represents a token.
type Token int

// List of tokens
const (
	Illegal Token = iota
	EOF
	Comment
	_literalBeg
	Ident
	Int
	Float
	Char
	String
	_literalEnd
	_operatorBeg
	Add          // +
	Sub          // -
	Mul          // *
	Quo          // /
	Rem          // %
	And          // &
	Or           // |
	Xor          // ^
	Shl          // <<
	Shr          // >>
	AndNot       // &^
	AddAssign    // +=
	SubAssign    // -=
	MulAssign    // *=
	QuoAssign    // /=
	RemAssign    // %=
	AndAssign    // &=
	OrAssign     // |=
	XorAssign    // ^=
	ShlAssign    // <<=
	ShrAssign    // >>=
	AndNotAssign // &^=
	LAnd         // &&
	LOr          // ||
	Inc          // ++
	Dec          // --
	Equal        // ==
	Less         // <
	Greater      // >
	Assign       // =
	Not          // !
	NotEqual     // !=
	LessEq       // <=
	GreaterEq    // >=
	Define       // :=
	Ellipsis     // ...
	LParen       // (
	LBrack       // [
	LBrace       // {
	Comma        // ,
	Period       // .
	RParen       // )
	RBrack       // ]
	RBrace       // }
	Semicolon    // ;
	Colon        // :
	_operatorEnd
	_keywordBeg
	Break
	Case
	Continue
	Default
	Else
	For
	Func
	Error
	If
	Return
	Switch
	Var
	True
	False
	In
	Undefined
	_keywordEnd
)

var tokens = [...]string{
	Illegal:      "ILLEGAL",
	EOF:          "EOF",
	Comment:      "COMMENT",
	Ident:        "IDENT",
	Int:          "INT",
	Float:        "FLOAT",
	Char:         "CHAR",
	String:       "STRING",
	Add:          "+",
	Sub:          "-",
	Mul:          "*",
	Quo:          "/",
	Rem:          "%",
	And:          "&",
	Or:           "|",
	Xor:          "^",
	Shl:          "<<",
	Shr:          ">>",
	AndNot:       "&^",
	AddAssign:    "+=",
	SubAssign:    "-=",
	MulAssign:    "*=",
	QuoAssign:    "/=",
	RemAssign:    "%=",
	AndAssign:    "&=",
	OrAssign:     "|=",
	XorAssign:    "^=",
	ShlAssign:    "<<=",
	ShrAssign:    ">>=",
	AndNotAssign: "&^=",
	LAnd:         "&&",
	LOr:          "||",
	Inc:          "++",
	Dec:          "--",
	Equal:        "==",
	Less:         "<",
	Greater:      ">",
	Assign:       "=",
	Not:          "!",
	NotEqual:     "!=",
	LessEq:       "<=",
	GreaterEq:    ">=",
	Define:       ":=",
	Ellipsis:     "...",
	LParen:       "(",
	LBrack:       "[",
	LBrace:       "{",
	Comma:        ",",
	Period:       ".",
	RParen:       ")",
	RBrack:       "]",
	RBrace:       "}",
	Semicolon:    ";",
	Colon:        ":",
	Break:        "break",
	Case:         "case",
	Continue:     "continue",
	Default:      "default",
	Else:         "else",
	For:          "for",
	Func:         "func",
	Error:        "error",
	If:           "if",
	Return:       "return",
	Switch:       "switch",
	Var:          "var",
	True:         "true",
	False:        "false",
	In:           "in",
	Undefined:    "undefined",
}

func (tok Token) String() string {
	s := ""

	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}

	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}

	return s
}

// LowestPrec represents lowest operator precedence.
const LowestPrec = 0

// Precedence returns the precedence for the operator token.
func (tok Token) Precedence() int {
	switch tok {
	case LOr:
		return 1
	case LAnd:
		return 2
	case Equal, NotEqual, Less, LessEq, Greater, GreaterEq:
		return 3
	case Add, Sub, Or, Xor:
		return 4
	case Mul, Quo, Rem, Shl, Shr, And, AndNot:
		return 5
	}
	return LowestPrec
}

// IsLiteral returns true if the token is a literal.
func (tok Token) IsLiteral() bool {
	return _literalBeg < tok && tok < _literalEnd
}

// IsOperator returns true if the token is an operator.
func (tok Token) IsOperator() bool {
	return _operatorBeg < tok && tok < _operatorEnd
}

// IsKeyword returns true if the token is a keyword.
func (tok Token) IsKeyword() bool {
	return _keywordBeg < tok && tok < _keywordEnd
}

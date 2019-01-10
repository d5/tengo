package ast

import "github.com/d5/tengo/scanner"

type UndefinedLit struct {
	TokenPos scanner.Pos
}

func (e *UndefinedLit) exprNode() {}

func (e *UndefinedLit) Pos() scanner.Pos {
	return e.TokenPos
}

func (e *UndefinedLit) End() scanner.Pos {
	return e.TokenPos + 9 // len(undefined) == 9
}

func (e *UndefinedLit) String() string {
	return "undefined"
}

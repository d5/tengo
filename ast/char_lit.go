package ast

import "github.com/d5/tengo/scanner"

type CharLit struct {
	Value    rune
	ValuePos scanner.Pos
	Literal  string
}

func (e *CharLit) exprNode() {}

func (e *CharLit) Pos() scanner.Pos {
	return e.ValuePos
}

func (e *CharLit) End() scanner.Pos {
	return scanner.Pos(int(e.ValuePos) + len(e.Literal))
}

func (e *CharLit) String() string {
	return e.Literal
}

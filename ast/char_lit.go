package ast

import "github.com/d5/tengo/source"

type CharLit struct {
	Value    rune
	ValuePos source.Pos
	Literal  string
}

func (e *CharLit) exprNode() {}

func (e *CharLit) Pos() source.Pos {
	return e.ValuePos
}

func (e *CharLit) End() source.Pos {
	return source.Pos(int(e.ValuePos) + len(e.Literal))
}

func (e *CharLit) String() string {
	return e.Literal
}

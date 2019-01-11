package ast

import "github.com/d5/tengo/source"

type StringLit struct {
	Value    string
	ValuePos source.Pos
	Literal  string
}

func (e *StringLit) exprNode() {}

func (e *StringLit) Pos() source.Pos {
	return e.ValuePos
}

func (e *StringLit) End() source.Pos {
	return source.Pos(int(e.ValuePos) + len(e.Literal))
}

func (e *StringLit) String() string {
	return e.Literal
}

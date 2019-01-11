package ast

import "github.com/d5/tengo/source"

type IntLit struct {
	Value    int64
	ValuePos source.Pos
	Literal  string
}

func (e *IntLit) exprNode() {}

func (e *IntLit) Pos() source.Pos {
	return e.ValuePos
}

func (e *IntLit) End() source.Pos {
	return source.Pos(int(e.ValuePos) + len(e.Literal))
}

func (e *IntLit) String() string {
	return e.Literal
}

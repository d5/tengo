package ast

import "github.com/d5/tengo/compiler/source"

type BoolLit struct {
	Value    bool
	ValuePos source.Pos
	Literal  string
}

func (e *BoolLit) exprNode() {}

func (e *BoolLit) Pos() source.Pos {
	return e.ValuePos
}

func (e *BoolLit) End() source.Pos {
	return source.Pos(int(e.ValuePos) + len(e.Literal))
}

func (e *BoolLit) String() string {
	return e.Literal
}

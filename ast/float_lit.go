package ast

import "github.com/d5/tengo/source"

type FloatLit struct {
	Value    float64
	ValuePos source.Pos
	Literal  string
}

func (e *FloatLit) exprNode() {}

func (e *FloatLit) Pos() source.Pos {
	return e.ValuePos
}

func (e *FloatLit) End() source.Pos {
	return source.Pos(int(e.ValuePos) + len(e.Literal))
}

func (e *FloatLit) String() string {
	return e.Literal
}

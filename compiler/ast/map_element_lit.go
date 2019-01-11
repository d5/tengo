package ast

import "github.com/d5/tengo/compiler/source"

type MapElementLit struct {
	Key      string
	KeyPos   source.Pos
	ColonPos source.Pos
	Value    Expr
}

func (e *MapElementLit) exprNode() {}

func (e *MapElementLit) Pos() source.Pos {
	return e.KeyPos
}

func (e *MapElementLit) End() source.Pos {
	return e.Value.End()
}

func (e *MapElementLit) String() string {
	return e.Key + ": " + e.Value.String()
}

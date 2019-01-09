package ast

import "github.com/d5/tengo/scanner"

type MapElementLit struct {
	Key      string
	KeyPos   scanner.Pos
	ColonPos scanner.Pos
	Value    Expr
}

func (e *MapElementLit) exprNode() {}

func (e *MapElementLit) Pos() scanner.Pos {
	return e.KeyPos
}

func (e *MapElementLit) End() scanner.Pos {
	return e.Value.End()
}

func (e *MapElementLit) String() string {
	return e.Key + ": " + e.Value.String()
}

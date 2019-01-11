package ast

import "github.com/d5/tengo/compiler/source"

type BadExpr struct {
	From source.Pos
	To   source.Pos
}

func (e *BadExpr) exprNode() {}

func (e *BadExpr) Pos() source.Pos {
	return e.From
}

func (e *BadExpr) End() source.Pos {
	return e.To
}

func (e *BadExpr) String() string {
	return "<bad expression>"
}

package ast

import "github.com/d5/tengo/source"

type ParenExpr struct {
	Expr   Expr
	LParen source.Pos
	RParen source.Pos
}

func (e *ParenExpr) exprNode() {}

func (e *ParenExpr) Pos() source.Pos {
	return e.LParen
}

func (e *ParenExpr) End() source.Pos {
	return e.RParen + 1
}

func (e *ParenExpr) String() string {
	return "(" + e.Expr.String() + ")"
}

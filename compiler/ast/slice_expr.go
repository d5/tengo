package ast

import "github.com/d5/tengo/compiler/source"

type SliceExpr struct {
	Expr   Expr
	LBrack source.Pos
	Low    Expr
	High   Expr
	RBrack source.Pos
}

func (e *SliceExpr) exprNode() {}

func (e *SliceExpr) Pos() source.Pos {
	return e.Expr.Pos()
}

func (e *SliceExpr) End() source.Pos {
	return e.RBrack + 1
}

func (e *SliceExpr) String() string {
	var low, high string
	if e.Low != nil {
		low = e.Low.String()
	}
	if e.High != nil {
		high = e.High.String()
	}

	return e.Expr.String() + "[" + low + ":" + high + "]"
}

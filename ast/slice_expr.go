package ast

import "github.com/d5/tengo/scanner"

type SliceExpr struct {
	Expr   Expr
	LBrack scanner.Pos
	Low    Expr
	High   Expr
	RBrack scanner.Pos
}

func (e *SliceExpr) exprNode() {}

func (e *SliceExpr) Pos() scanner.Pos {
	return e.Expr.Pos()
}

func (e *SliceExpr) End() scanner.Pos {
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

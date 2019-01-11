package ast

import "github.com/d5/tengo/source"

type SelectorExpr struct {
	Expr Expr
	Sel  Expr
}

func (e *SelectorExpr) exprNode() {}

func (e *SelectorExpr) Pos() source.Pos {
	return e.Expr.Pos()
}

func (e *SelectorExpr) End() source.Pos {
	return e.Sel.End()
}

func (e *SelectorExpr) String() string {
	return e.Expr.String() + "." + e.Sel.String()
}

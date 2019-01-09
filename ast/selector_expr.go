package ast

import "github.com/d5/tengo/scanner"

type SelectorExpr struct {
	Expr Expr
	Sel  Expr
}

func (e *SelectorExpr) exprNode() {}

func (e *SelectorExpr) Pos() scanner.Pos {
	return e.Expr.Pos()
}

func (e *SelectorExpr) End() scanner.Pos {
	return e.Sel.End()
}

func (e *SelectorExpr) String() string {
	return e.Expr.String() + "." + e.Sel.String()
}

package ast

import "github.com/d5/tengo/scanner"

type IndexExpr struct {
	Expr   Expr
	LBrack scanner.Pos
	Index  Expr
	RBrack scanner.Pos
}

func (e *IndexExpr) exprNode() {}

func (e *IndexExpr) Pos() scanner.Pos {
	return e.Expr.Pos()
}

func (e *IndexExpr) End() scanner.Pos {
	return e.RBrack + 1
}

func (e *IndexExpr) String() string {
	var index string
	if e.Index != nil {
		index = e.Index.String()
	}

	return e.Expr.String() + "[" + index + "]"
}

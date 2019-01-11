package ast

import "github.com/d5/tengo/compiler/source"

type IndexExpr struct {
	Expr   Expr
	LBrack source.Pos
	Index  Expr
	RBrack source.Pos
}

func (e *IndexExpr) exprNode() {}

func (e *IndexExpr) Pos() source.Pos {
	return e.Expr.Pos()
}

func (e *IndexExpr) End() source.Pos {
	return e.RBrack + 1
}

func (e *IndexExpr) String() string {
	var index string
	if e.Index != nil {
		index = e.Index.String()
	}

	return e.Expr.String() + "[" + index + "]"
}

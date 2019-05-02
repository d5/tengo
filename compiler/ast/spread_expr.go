package ast

import "github.com/d5/tengo/compiler/source"

// SpreadExpr represents an explosion expression
type SpreadExpr struct {
	Expr     Expr
	Ellipsis source.Pos
}

func (e *SpreadExpr) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *SpreadExpr) Pos() source.Pos {
	return e.Expr.Pos()
}

// End returns the position of first character immediately after the node.
func (e *SpreadExpr) End() source.Pos {
	return e.Ellipsis + 3
}

func (e *SpreadExpr) String() string {
	return e.Expr.String() + "..."
}

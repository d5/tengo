package ast

import "github.com/d5/tengo/compiler/source"

// ExplodeExpr represents an explosion expression
type ExplodeExpr struct {
	Expr     Expr
	Ellipsis source.Pos
}

func (e *ExplodeExpr) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *ExplodeExpr) Pos() source.Pos {
	return e.Expr.Pos()
}

// End returns the position of first character immediately after the node.
func (e *ExplodeExpr) End() source.Pos {
	return e.Ellipsis + 3
}

func (e *ExplodeExpr) String() string {
	return e.Expr.String() + "..."
}

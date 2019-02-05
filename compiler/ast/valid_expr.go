package ast

import (
	"github.com/d5/tengo/compiler/source"
)

// ValidExpr represents a ternary conditional expression.
type ValidExpr struct {
	Cond        Expr
	True        Expr
	False       Expr
	QuestionPos source.Pos
}

func (e *ValidExpr) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *ValidExpr) Pos() source.Pos {
	return e.Cond.Pos()
}

// End returns the position of first character immediately after the node.
func (e *ValidExpr) End() source.Pos {
	return e.False.End()
}

func (e *ValidExpr) String() string {
	return "(" + e.Cond.String() + " ?? " + e.False.String() + ")"
}

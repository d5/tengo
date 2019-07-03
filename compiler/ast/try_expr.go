package ast

import (
	"fmt"

	"github.com/d5/tengo/compiler/source"
)

// TryExpr represents a try expression.
type TryExpr struct {
	TryPos source.Pos
	LParen source.Pos
	Expr   Expr
	RParen source.Pos
}

func (e *TryExpr) exprNode() {}

// Pos returns the position of first character belonging to the node.
func (e *TryExpr) Pos() source.Pos { return e.TryPos }

// End returns the position of first character immediately after the node.
func (e *TryExpr) End() source.Pos { return e.RParen + 1 }

func (e *TryExpr) String() string { return fmt.Sprintf("try(%s)", e.Expr.String()) }

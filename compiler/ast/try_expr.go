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

func (e *TryExpr) Pos() source.Pos { return e.TryPos }
func (e *TryExpr) End() source.Pos { return e.RParen }
func (e *TryExpr) String() string  { return fmt.Sprintf("try(%s)", e.Expr.String()) }

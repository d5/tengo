package ast

import (
	"github.com/d5/tengo/compiler/source"
	"github.com/d5/tengo/compiler/token"
)

type BinaryExpr struct {
	Lhs      Expr
	Rhs      Expr
	Token    token.Token
	TokenPos source.Pos
}

func (e *BinaryExpr) exprNode() {}

func (e *BinaryExpr) Pos() source.Pos {
	return e.Lhs.Pos()
}

func (e *BinaryExpr) End() source.Pos {
	return e.Rhs.End()
}

func (e *BinaryExpr) String() string {
	return "(" + e.Lhs.String() + " " + e.Token.String() + " " + e.Rhs.String() + ")"
}

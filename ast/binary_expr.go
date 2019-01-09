package ast

import (
	"github.com/d5/tengo/scanner"
	"github.com/d5/tengo/token"
)

type BinaryExpr struct {
	Lhs      Expr
	Rhs      Expr
	Token    token.Token
	TokenPos scanner.Pos
}

func (e *BinaryExpr) exprNode() {}

func (e *BinaryExpr) Pos() scanner.Pos {
	return e.Lhs.Pos()
}

func (e *BinaryExpr) End() scanner.Pos {
	return e.Rhs.End()
}

func (e *BinaryExpr) String() string {
	return "(" + e.Lhs.String() + " " + e.Token.String() + " " + e.Rhs.String() + ")"
}

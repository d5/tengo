package ast

import (
	"github.com/d5/tengo/scanner"
	"github.com/d5/tengo/token"
)

type UnaryExpr struct {
	Expr     Expr
	Token    token.Token
	TokenPos scanner.Pos
}

func (e *UnaryExpr) exprNode() {}

func (e *UnaryExpr) Pos() scanner.Pos {
	return e.Expr.Pos()
}

func (e *UnaryExpr) End() scanner.Pos {
	return e.Expr.End()
}

func (e *UnaryExpr) String() string {
	return "(" + e.Token.String() + e.Expr.String() + ")"
}

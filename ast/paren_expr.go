package ast

import "github.com/d5/tengo/scanner"

type ParenExpr struct {
	Expr   Expr
	LParen scanner.Pos
	RParen scanner.Pos
}

func (e *ParenExpr) exprNode() {}

func (e *ParenExpr) Pos() scanner.Pos {
	return e.LParen
}

func (e *ParenExpr) End() scanner.Pos {
	return e.RParen + 1
}

func (e *ParenExpr) String() string {
	return "(" + e.Expr.String() + ")"
}

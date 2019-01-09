package ast

import "github.com/d5/tengo/scanner"

type BadExpr struct {
	From scanner.Pos
	To   scanner.Pos
}

func (e *BadExpr) exprNode() {}

func (e *BadExpr) Pos() scanner.Pos {
	return e.From
}

func (e *BadExpr) End() scanner.Pos {
	return e.To
}

func (e *BadExpr) String() string {
	return "<bad expression>"
}

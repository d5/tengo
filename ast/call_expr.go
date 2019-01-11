package ast

import (
	"strings"

	"github.com/d5/tengo/source"
)

type CallExpr struct {
	Func   Expr
	LParen source.Pos
	Args   []Expr
	RParen source.Pos
}

func (e *CallExpr) exprNode() {}

func (e *CallExpr) Pos() source.Pos {
	return e.Func.Pos()
}

func (e *CallExpr) End() source.Pos {
	return e.RParen + 1
}

func (e *CallExpr) String() string {
	var args []string
	for _, e := range e.Args {
		args = append(args, e.String())
	}

	return e.Func.String() + "(" + strings.Join(args, ", ") + ")"
}

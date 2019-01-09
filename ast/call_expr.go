package ast

import (
	"strings"

	"github.com/d5/tengo/scanner"
)

type CallExpr struct {
	Func   Expr
	LParen scanner.Pos
	Args   []Expr
	RParen scanner.Pos
}

func (e *CallExpr) exprNode() {}

func (e *CallExpr) Pos() scanner.Pos {
	return e.Func.Pos()
}

func (e *CallExpr) End() scanner.Pos {
	return e.RParen + 1
}

func (e *CallExpr) String() string {
	var args []string
	for _, e := range e.Args {
		args = append(args, e.String())
	}

	return e.Func.String() + "(" + strings.Join(args, ", ") + ")"
}

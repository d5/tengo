package ast

import (
	"strings"

	"github.com/d5/tengo/compiler/source"
)

type ArrayLit struct {
	Elements []Expr
	LBrack   source.Pos
	RBrack   source.Pos
}

func (e *ArrayLit) exprNode() {}

func (e *ArrayLit) Pos() source.Pos {
	return e.LBrack
}

func (e *ArrayLit) End() source.Pos {
	return e.RBrack + 1
}

func (e *ArrayLit) String() string {
	var elts []string
	for _, m := range e.Elements {
		elts = append(elts, m.String())
	}

	return "[" + strings.Join(elts, ", ") + "]"
}

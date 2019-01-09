package ast

import (
	"strings"

	"github.com/d5/tengo/scanner"
)

type ArrayLit struct {
	Elements []Expr
	LBrack   scanner.Pos
	RBrack   scanner.Pos
}

func (e *ArrayLit) exprNode() {}

func (e *ArrayLit) Pos() scanner.Pos {
	return e.LBrack
}

func (e *ArrayLit) End() scanner.Pos {
	return e.RBrack + 1
}

func (e *ArrayLit) String() string {
	var elts []string
	for _, m := range e.Elements {
		elts = append(elts, m.String())
	}

	return "[" + strings.Join(elts, ", ") + "]"
}

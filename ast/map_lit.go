package ast

import (
	"strings"

	"github.com/d5/tengo/scanner"
)

type MapLit struct {
	LBrace   scanner.Pos
	Elements []*MapElementLit
	RBrace   scanner.Pos
}

func (e *MapLit) exprNode() {}

func (e *MapLit) Pos() scanner.Pos {
	return e.LBrace
}

func (e *MapLit) End() scanner.Pos {
	return e.RBrace + 1
}

func (e *MapLit) String() string {
	var elts []string
	for _, m := range e.Elements {
		elts = append(elts, m.String())
	}

	return "{" + strings.Join(elts, ", ") + "}"
}

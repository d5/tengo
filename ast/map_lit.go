package ast

import (
	"strings"

	"github.com/d5/tengo/source"
)

type MapLit struct {
	LBrace   source.Pos
	Elements []*MapElementLit
	RBrace   source.Pos
}

func (e *MapLit) exprNode() {}

func (e *MapLit) Pos() source.Pos {
	return e.LBrace
}

func (e *MapLit) End() source.Pos {
	return e.RBrace + 1
}

func (e *MapLit) String() string {
	var elts []string
	for _, m := range e.Elements {
		elts = append(elts, m.String())
	}

	return "{" + strings.Join(elts, ", ") + "}"
}

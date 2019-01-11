package ast

import (
	"strings"

	"github.com/d5/tengo/source"
)

type IdentList struct {
	LParen source.Pos
	List   []*Ident
	RParen source.Pos
}

func (n *IdentList) Pos() source.Pos {
	if n.LParen.IsValid() {
		return n.LParen
	}

	if len(n.List) > 0 {
		return n.List[0].Pos()
	}

	return source.NoPos
}

func (n *IdentList) End() source.Pos {
	if n.RParen.IsValid() {
		return n.RParen + 1
	}

	if l := len(n.List); l > 0 {
		return n.List[l-1].End()
	}

	return source.NoPos
}

func (n *IdentList) NumFields() int {
	if n == nil {
		return 0
	}

	return len(n.List)
}

func (n *IdentList) String() string {
	var list []string
	for _, e := range n.List {
		list = append(list, e.String())
	}

	return "(" + strings.Join(list, ", ") + ")"
}

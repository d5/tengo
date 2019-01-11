package ast

import "github.com/d5/tengo/source"

type Ident struct {
	Name    string
	NamePos source.Pos
}

func (e *Ident) exprNode() {}

func (e *Ident) Pos() source.Pos {
	return e.NamePos
}

func (e *Ident) End() source.Pos {
	return source.Pos(int(e.NamePos) + len(e.Name))
}

func (e *Ident) String() string {
	if e != nil {
		return e.Name
	}

	return nullRep
}

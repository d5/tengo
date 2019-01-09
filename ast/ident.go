package ast

import "github.com/d5/tengo/scanner"

type Ident struct {
	Name    string
	NamePos scanner.Pos
}

func (e *Ident) exprNode() {}

func (e *Ident) Pos() scanner.Pos {
	return e.NamePos
}

func (e *Ident) End() scanner.Pos {
	return scanner.Pos(int(e.NamePos) + len(e.Name))
}

func (e *Ident) String() string {
	if e != nil {
		return e.Name
	}

	return nullRep
}

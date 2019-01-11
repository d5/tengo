package ast

import "github.com/d5/tengo/compiler/source"

type ForInStmt struct {
	ForPos   source.Pos
	Key      *Ident
	Value    *Ident
	Iterable Expr
	Body     *BlockStmt
}

func (s *ForInStmt) stmtNode() {}

func (s *ForInStmt) Pos() source.Pos {
	return s.ForPos
}

func (s *ForInStmt) End() source.Pos {
	return s.Body.End()
}

func (s *ForInStmt) String() string {
	if s.Value != nil {
		return "for " + s.Key.String() + ", " + s.Value.String() + " in " + s.Iterable.String() + " " + s.Body.String()
	}

	return "for " + s.Key.String() + " in " + s.Iterable.String() + " " + s.Body.String()
}

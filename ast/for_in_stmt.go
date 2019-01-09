package ast

import "github.com/d5/tengo/scanner"

type ForInStmt struct {
	ForPos   scanner.Pos
	Key      *Ident
	Value    *Ident
	Iterable Expr
	Body     *BlockStmt
}

func (s *ForInStmt) stmtNode() {}

func (s *ForInStmt) Pos() scanner.Pos {
	return s.ForPos
}

func (s *ForInStmt) End() scanner.Pos {
	return s.Body.End()
}

func (s *ForInStmt) String() string {
	if s.Value != nil {
		return "for " + s.Key.String() + ", " + s.Value.String() + " in " + s.Iterable.String() + " " + s.Body.String()
	}

	return "for " + s.Key.String() + " in " + s.Iterable.String() + " " + s.Body.String()
}

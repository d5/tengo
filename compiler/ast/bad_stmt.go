package ast

import "github.com/d5/tengo/compiler/source"

type BadStmt struct {
	From source.Pos
	To   source.Pos
}

func (s *BadStmt) stmtNode() {}

func (s *BadStmt) Pos() source.Pos {
	return s.From
}

func (s *BadStmt) End() source.Pos {
	return s.To
}

func (s *BadStmt) String() string {
	return "<bad statement>"
}

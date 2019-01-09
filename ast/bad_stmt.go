package ast

import "github.com/d5/tengo/scanner"

type BadStmt struct {
	From scanner.Pos
	To   scanner.Pos
}

func (s *BadStmt) stmtNode() {}

func (s *BadStmt) Pos() scanner.Pos {
	return s.From
}

func (s *BadStmt) End() scanner.Pos {
	return s.To
}

func (s *BadStmt) String() string {
	return "<bad statement>"
}

package ast

import "github.com/d5/tengo/compiler/source"

type EmptyStmt struct {
	Semicolon source.Pos
	Implicit  bool
}

func (s *EmptyStmt) stmtNode() {}

func (s *EmptyStmt) Pos() source.Pos {
	return s.Semicolon
}

func (s *EmptyStmt) End() source.Pos {
	if s.Implicit {
		return s.Semicolon
	}

	return s.Semicolon + 1
}

func (s *EmptyStmt) String() string {
	return ";"
}

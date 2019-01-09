package ast

import "github.com/d5/tengo/scanner"

type EmptyStmt struct {
	Semicolon scanner.Pos
	Implicit  bool
}

func (s *EmptyStmt) stmtNode() {}

func (s *EmptyStmt) Pos() scanner.Pos {
	return s.Semicolon
}

func (s *EmptyStmt) End() scanner.Pos {
	if s.Implicit {
		return s.Semicolon
	}

	return s.Semicolon + 1
}

func (s *EmptyStmt) String() string {
	return ";"
}

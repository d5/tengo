package ast

import (
	"github.com/d5/tengo/compiler/source"
	"github.com/d5/tengo/compiler/token"
)

type BranchStmt struct {
	Token    token.Token
	TokenPos source.Pos
	Label    *Ident
}

func (s *BranchStmt) stmtNode() {}

func (s *BranchStmt) Pos() source.Pos {
	return s.TokenPos
}

func (s *BranchStmt) End() source.Pos {
	if s.Label != nil {
		return s.Label.End()
	}

	return source.Pos(int(s.TokenPos) + len(s.Token.String()))
}

func (s *BranchStmt) String() string {
	var label string
	if s.Label != nil {
		label = " " + s.Label.Name
	}

	return s.Token.String() + label
}

package ast

import (
	"github.com/d5/tengo/scanner"
	"github.com/d5/tengo/token"
)

type BranchStmt struct {
	Token    token.Token
	TokenPos scanner.Pos
	Label    *Ident
}

func (s *BranchStmt) stmtNode() {}

func (s *BranchStmt) Pos() scanner.Pos {
	return s.TokenPos
}

func (s *BranchStmt) End() scanner.Pos {
	if s.Label != nil {
		return s.Label.End()
	}

	return scanner.Pos(int(s.TokenPos) + len(s.Token.String()))
}

func (s *BranchStmt) String() string {
	var label string
	if s.Label != nil {
		label = " " + s.Label.Name
	}

	return s.Token.String() + label
}

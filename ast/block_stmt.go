package ast

import (
	"strings"

	"github.com/d5/tengo/source"
)

type BlockStmt struct {
	Stmts  []Stmt
	LBrace source.Pos
	RBrace source.Pos
}

func (s *BlockStmt) stmtNode() {}

func (s *BlockStmt) Pos() source.Pos {
	return s.LBrace
}

func (s *BlockStmt) End() source.Pos {
	return s.RBrace + 1
}

func (s *BlockStmt) String() string {
	var list []string
	for _, e := range s.Stmts {
		list = append(list, e.String())
	}

	return "{" + strings.Join(list, "; ") + "}"
}

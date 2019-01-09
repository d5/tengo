package ast

import (
	"strings"

	"github.com/d5/tengo/scanner"
)

type BlockStmt struct {
	Stmts  []Stmt
	LBrace scanner.Pos
	RBrace scanner.Pos
}

func (s *BlockStmt) stmtNode() {}

func (s *BlockStmt) Pos() scanner.Pos {
	return s.LBrace
}

func (s *BlockStmt) End() scanner.Pos {
	return s.RBrace + 1
}

func (s *BlockStmt) String() string {
	var list []string
	for _, e := range s.Stmts {
		list = append(list, e.String())
	}

	return "{" + strings.Join(list, "; ") + "}"
}

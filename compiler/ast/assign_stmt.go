package ast

import (
	"strings"

	"github.com/d5/tengo/compiler/source"
	"github.com/d5/tengo/compiler/token"
)

type AssignStmt struct {
	Lhs      []Expr
	Rhs      []Expr
	Token    token.Token
	TokenPos source.Pos
}

func (s *AssignStmt) stmtNode() {}

func (s *AssignStmt) Pos() source.Pos {
	return s.Lhs[0].Pos()
}

func (s *AssignStmt) End() source.Pos {
	return s.Rhs[len(s.Rhs)-1].End()
}

func (s *AssignStmt) String() string {
	var lhs, rhs []string
	for _, e := range s.Lhs {
		lhs = append(lhs, e.String())
	}
	for _, e := range s.Rhs {
		rhs = append(rhs, e.String())
	}

	return strings.Join(lhs, ", ") + " " + s.Token.String() + " " + strings.Join(rhs, ", ")
}

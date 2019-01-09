package ast

import (
	"strings"

	"github.com/d5/tengo/scanner"
	"github.com/d5/tengo/token"
)

type AssignStmt struct {
	Lhs      []Expr
	Rhs      []Expr
	Token    token.Token
	TokenPos scanner.Pos
}

func (s *AssignStmt) stmtNode() {}

func (s *AssignStmt) Pos() scanner.Pos {
	return s.Lhs[0].Pos()
}

func (s *AssignStmt) End() scanner.Pos {
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

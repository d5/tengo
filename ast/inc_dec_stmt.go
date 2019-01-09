package ast

import (
	"github.com/d5/tengo/scanner"
	"github.com/d5/tengo/token"
)

type IncDecStmt struct {
	Expr     Expr
	Token    token.Token
	TokenPos scanner.Pos
}

func (s *IncDecStmt) stmtNode() {}

func (s *IncDecStmt) Pos() scanner.Pos {
	return s.Expr.Pos()
}

func (s *IncDecStmt) End() scanner.Pos {
	return scanner.Pos(int(s.TokenPos) + 2)
}

func (s *IncDecStmt) String() string {
	return s.Expr.String() + s.Token.String()
}

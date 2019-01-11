package ast

import "github.com/d5/tengo/source"

type ExprStmt struct {
	Expr Expr
}

func (s *ExprStmt) stmtNode() {}

func (s *ExprStmt) Pos() source.Pos {
	return s.Expr.Pos()
}

func (s *ExprStmt) End() source.Pos {
	return s.Expr.End()
}

func (s *ExprStmt) String() string {
	return s.Expr.String()
}

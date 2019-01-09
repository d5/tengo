package ast

import "github.com/d5/tengo/scanner"

type ExprStmt struct {
	Expr Expr
}

func (s *ExprStmt) stmtNode() {}

func (s *ExprStmt) Pos() scanner.Pos {
	return s.Expr.Pos()
}

func (s *ExprStmt) End() scanner.Pos {
	return s.Expr.End()
}

func (s *ExprStmt) String() string {
	return s.Expr.String()
}

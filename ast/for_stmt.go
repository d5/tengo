package ast

import "github.com/d5/tengo/scanner"

type ForStmt struct {
	ForPos scanner.Pos
	Init   Stmt
	Cond   Expr
	Post   Stmt
	Body   *BlockStmt
}

func (s *ForStmt) stmtNode() {}

func (s *ForStmt) Pos() scanner.Pos {
	return s.ForPos
}

func (s *ForStmt) End() scanner.Pos {
	return s.Body.End()
}

func (s *ForStmt) String() string {
	var init, cond, post string
	if s.Init != nil {
		init = s.Init.String()
	}
	if s.Cond != nil {
		cond = s.Cond.String() + " "
	}
	if s.Post != nil {
		post = s.Post.String()
	}

	if init != "" || post != "" {
		return "for " + init + " ; " + cond + " ; " + post + s.Body.String()
	}

	return "for " + cond + s.Body.String()
}

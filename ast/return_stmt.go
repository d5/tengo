package ast

import (
	"strings"

	"github.com/d5/tengo/scanner"
)

type ReturnStmt struct {
	ReturnPos scanner.Pos
	Results   []Expr
}

func (s *ReturnStmt) stmtNode() {}

func (s *ReturnStmt) Pos() scanner.Pos {
	return s.ReturnPos
}

func (s *ReturnStmt) End() scanner.Pos {
	if n := len(s.Results); n > 0 {
		return s.Results[n-1].End()
	}

	return s.ReturnPos + 6
}

func (s *ReturnStmt) String() string {
	if len(s.Results) > 0 {
		var res []string
		for _, e := range s.Results {
			res = append(res, e.String())
		}

		return "return " + strings.Join(res, ", ")
	}

	return "return"
}

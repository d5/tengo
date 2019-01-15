package ast

import (
	"strings"

	"github.com/d5/tengo/compiler/source"
)

// ReturnStmt represents a return statement.
type ReturnStmt struct {
	ReturnPos source.Pos
	Results   []Expr
}

func (s *ReturnStmt) stmtNode() {}

// Pos returns the position of first character belonging to the node.
func (s *ReturnStmt) Pos() source.Pos {
	return s.ReturnPos
}

// End returns the position of first character immediately after the node.
func (s *ReturnStmt) End() source.Pos {
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

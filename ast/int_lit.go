package ast

import "github.com/d5/tengo/scanner"

type IntLit struct {
	Value    int64
	ValuePos scanner.Pos
	Literal  string
}

func (e *IntLit) exprNode() {}

func (e *IntLit) Pos() scanner.Pos {
	return e.ValuePos
}

func (e *IntLit) End() scanner.Pos {
	return scanner.Pos(int(e.ValuePos) + len(e.Literal))
}

func (e *IntLit) String() string {
	return e.Literal
}

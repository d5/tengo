package ast

import "github.com/d5/tengo/scanner"

type StringLit struct {
	Value    string
	ValuePos scanner.Pos
	Literal  string
}

func (e *StringLit) exprNode() {}

func (e *StringLit) Pos() scanner.Pos {
	return e.ValuePos
}

func (e *StringLit) End() scanner.Pos {
	return scanner.Pos(int(e.ValuePos) + len(e.Literal))
}

func (e *StringLit) String() string {
	return e.Literal
}

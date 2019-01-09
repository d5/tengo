package ast

import "github.com/d5/tengo/scanner"

type BoolLit struct {
	Value    bool
	ValuePos scanner.Pos
	Literal  string
}

func (e *BoolLit) exprNode() {}

func (e *BoolLit) Pos() scanner.Pos {
	return e.ValuePos
}

func (e *BoolLit) End() scanner.Pos {
	return scanner.Pos(int(e.ValuePos) + len(e.Literal))
}

func (e *BoolLit) String() string {
	return e.Literal
}

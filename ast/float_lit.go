package ast

import "github.com/d5/tengo/scanner"

type FloatLit struct {
	Value    float64
	ValuePos scanner.Pos
	Literal  string
}

func (e *FloatLit) exprNode() {}

func (e *FloatLit) Pos() scanner.Pos {
	return e.ValuePos
}

func (e *FloatLit) End() scanner.Pos {
	return scanner.Pos(int(e.ValuePos) + len(e.Literal))
}

func (e *FloatLit) String() string {
	return e.Literal
}

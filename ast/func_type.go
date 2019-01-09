package ast

import "github.com/d5/tengo/scanner"

type FuncType struct {
	FuncPos scanner.Pos
	Params  *IdentList
}

func (e *FuncType) exprNode() {}

func (e *FuncType) Pos() scanner.Pos {
	return e.FuncPos
}

func (e *FuncType) End() scanner.Pos {
	return e.Params.End()
}

func (e *FuncType) String() string {
	return "func" + e.Params.String()
}

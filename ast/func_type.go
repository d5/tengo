package ast

import "github.com/d5/tengo/source"

type FuncType struct {
	FuncPos source.Pos
	Params  *IdentList
}

func (e *FuncType) exprNode() {}

func (e *FuncType) Pos() source.Pos {
	return e.FuncPos
}

func (e *FuncType) End() source.Pos {
	return e.Params.End()
}

func (e *FuncType) String() string {
	return "func" + e.Params.String()
}

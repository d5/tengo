package ast

import "github.com/d5/tengo/compiler/source"

type FuncLit struct {
	Type *FuncType
	Body *BlockStmt
}

func (e *FuncLit) exprNode() {}

func (e *FuncLit) Pos() source.Pos {
	return e.Type.Pos()
}

func (e *FuncLit) End() source.Pos {
	return e.Body.End()
}

func (e *FuncLit) String() string {
	return "func" + e.Type.Params.String() + " " + e.Body.String()
}

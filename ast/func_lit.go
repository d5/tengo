package ast

import "github.com/d5/tengo/scanner"

type FuncLit struct {
	Type *FuncType
	Body *BlockStmt
}

func (e *FuncLit) exprNode() {}

func (e *FuncLit) Pos() scanner.Pos {
	return e.Type.Pos()
}

func (e *FuncLit) End() scanner.Pos {
	return e.Body.End()
}

func (e *FuncLit) String() string {
	return "func" + e.Type.Params.String() + " " + e.Body.String()
}

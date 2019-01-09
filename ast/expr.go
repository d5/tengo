package ast

type Expr interface {
	Node
	exprNode()
}

package ast

type Stmt interface {
	Node
	stmtNode()
}

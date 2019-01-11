package compiler

import (
	"github.com/d5/tengo/compiler/ast"
	"github.com/d5/tengo/compiler/token"
)

func (c *Compiler) compileLogical(node *ast.BinaryExpr) error {
	// left side term
	if err := c.Compile(node.Lhs); err != nil {
		return err
	}

	// jump position
	var jumpPos int
	if node.Token == token.LAnd {
		jumpPos = c.emit(OpAndJump, 0)
	} else {
		jumpPos = c.emit(OpOrJump, 0)
	}

	// right side term
	if err := c.Compile(node.Rhs); err != nil {
		return err
	}

	c.changeOperand(jumpPos, len(c.currentInstructions()))

	return nil
}

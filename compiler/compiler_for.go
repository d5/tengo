package compiler

import (
	"github.com/d5/tengo/compiler/ast"
)

func (c *Compiler) compileForStmt(stmt *ast.ForStmt) error {
	c.symbolTable = c.symbolTable.Fork(true)
	defer func() {
		c.symbolTable = c.symbolTable.Parent(false)
	}()

	// init statement
	if stmt.Init != nil {
		if err := c.Compile(stmt.Init); err != nil {
			return err
		}
	}

	// pre-condition position
	preCondPos := len(c.currentInstructions())

	// condition expression
	if err := c.Compile(stmt.Cond); err != nil {
		return err
	}
	// condition jump position
	postCondPos := c.emit(OpJumpFalsy, 0)

	// enter loop
	loop := c.enterLoop()

	// body statement
	if err := c.Compile(stmt.Body); err != nil {
		c.leaveLoop()
		return err
	}

	c.leaveLoop()

	// post-body position
	postBodyPos := len(c.currentInstructions())

	// post statement
	if stmt.Post != nil {
		if err := c.Compile(stmt.Post); err != nil {
			return err
		}
	}

	// back to condition
	c.emit(OpJump, preCondPos)

	// post-statement position
	postStmtPos := len(c.currentInstructions())
	c.changeOperand(postCondPos, postStmtPos)

	// update all break/continue jump positions
	for _, pos := range loop.Breaks {
		c.changeOperand(pos, postStmtPos)
	}
	for _, pos := range loop.Continues {
		c.changeOperand(pos, postBodyPos)
	}

	return nil
}

func (c *Compiler) compileForInStmt(stmt *ast.ForInStmt) error {
	c.symbolTable = c.symbolTable.Fork(true)
	defer func() {
		c.symbolTable = c.symbolTable.Parent(false)
	}()

	// for-in statement is compiled like following:
	//
	//   for :it := iterator(iterable); :it.next();  {
	//     k, v := :it.get()  // DEFINE operator
	//
	//     ... body ...
	//   }
	//
	// ":it" is a local variable but will be conflict with other user variables
	// because character ":" is not allowed.

	// init
	//   :it = iterator(iterable)
	itSymbol := c.symbolTable.Define(":it")
	if err := c.Compile(stmt.Iterable); err != nil {
		return err
	}
	c.emit(OpIteratorInit)
	if itSymbol.Scope == ScopeGlobal {
		c.emit(OpSetGlobal, itSymbol.Index)
	} else {
		c.emit(OpDefineLocal, itSymbol.Index)
	}

	// pre-condition position
	preCondPos := len(c.currentInstructions())

	// condition
	//  :it.HasMore()
	if itSymbol.Scope == ScopeGlobal {
		c.emit(OpGetGlobal, itSymbol.Index)
	} else {
		c.emit(OpGetLocal, itSymbol.Index)
	}
	c.emit(OpIteratorNext)

	// condition jump position
	postCondPos := c.emit(OpJumpFalsy, 0)

	// enter loop
	loop := c.enterLoop()

	// assign key variable
	if stmt.Key.Name != "_" {
		keySymbol := c.symbolTable.Define(stmt.Key.Name)
		if itSymbol.Scope == ScopeGlobal {
			c.emit(OpGetGlobal, itSymbol.Index)
		} else {
			c.emit(OpGetLocal, itSymbol.Index)
		}
		c.emit(OpIteratorKey)
		if keySymbol.Scope == ScopeGlobal {
			c.emit(OpSetGlobal, keySymbol.Index)
		} else {
			c.emit(OpDefineLocal, keySymbol.Index)
		}
	}

	// assign value variable
	if stmt.Value.Name != "_" {
		valueSymbol := c.symbolTable.Define(stmt.Value.Name)
		if itSymbol.Scope == ScopeGlobal {
			c.emit(OpGetGlobal, itSymbol.Index)
		} else {
			c.emit(OpGetLocal, itSymbol.Index)
		}
		c.emit(OpIteratorValue)
		if valueSymbol.Scope == ScopeGlobal {
			c.emit(OpSetGlobal, valueSymbol.Index)
		} else {
			c.emit(OpDefineLocal, valueSymbol.Index)
		}
	}

	// body statement
	if err := c.Compile(stmt.Body); err != nil {
		c.leaveLoop()
		return err
	}

	c.leaveLoop()

	// post-body position
	postBodyPos := len(c.currentInstructions())

	// back to condition
	c.emit(OpJump, preCondPos)

	// post-statement position
	postStmtPos := len(c.currentInstructions())
	c.changeOperand(postCondPos, postStmtPos)

	// update all break/continue jump positions
	for _, pos := range loop.Breaks {
		c.changeOperand(pos, postStmtPos)
	}
	for _, pos := range loop.Continues {
		c.changeOperand(pos, postBodyPos)
	}

	return nil
}

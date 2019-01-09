package compiler

func (c *Compiler) currentInstructions() []byte {
	return c.scopes[c.scopeIndex].instructions
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions: make([]byte, 0),
	}

	c.scopes = append(c.scopes, scope)
	c.scopeIndex++

	c.symbolTable = c.symbolTable.Fork(false)

	if c.trace != nil {
		c.printTrace("SCOPE", c.scopeIndex)
	}
}

func (c *Compiler) leaveScope() []byte {
	instructions := c.currentInstructions()

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--

	c.symbolTable = c.symbolTable.Parent(true)

	if c.trace != nil {
		c.printTrace("SCOPL", c.scopeIndex)
	}

	return instructions
}

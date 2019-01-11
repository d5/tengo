package script

import (
	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/runtime"
)

var undefined objects.Object = &objects.Undefined{}

// Compiled is a compiled instance of the user script.
// Use Script.Compile() to create Compiled object.
type Compiled struct {
	symbolTable *compiler.SymbolTable
	machine     *runtime.VM
}

// Run executes the compiled script in the virtual machine.
func (c *Compiled) Run() error {
	return c.machine.Run()
}

// IsDefined returns true if the variable name is defined (has value) before or after the execution.
func (c *Compiled) IsDefined(name string) bool {
	symbol, _, ok := c.symbolTable.Resolve(name)
	if !ok {
		return false
	}

	v := c.machine.Globals()[symbol.Index]
	if v == nil {
		return false
	}

	_, isUndefined := (*v).(*objects.Undefined)

	return !isUndefined
}

// Get returns a variable identified by the name.
func (c *Compiled) Get(name string) *Variable {
	value := &undefined

	symbol, _, ok := c.symbolTable.Resolve(name)
	if ok && symbol.Scope == compiler.ScopeGlobal {
		value = c.machine.Globals()[symbol.Index]
		if value == nil {
			value = &undefined
		}
	}

	return &Variable{
		name:  name,
		value: value,
	}
}

// GetAll returns all the variables that are defined by the compiled script.
func (c *Compiled) GetAll() []*Variable {
	var vars []*Variable
	for _, name := range c.symbolTable.Names() {
		symbol, _, ok := c.symbolTable.Resolve(name)
		if ok && symbol.Scope == compiler.ScopeGlobal {
			value := c.machine.Globals()[symbol.Index]
			if value == nil {
				value = &undefined
			}

			vars = append(vars, &Variable{
				name:  name,
				value: value,
			})
		}
	}

	return vars
}

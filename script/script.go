package script

import (
	"fmt"

	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/parser"
	"github.com/d5/tengo/scanner"
	"github.com/d5/tengo/vm"
)

// Script can simplify compilation and execution of embedded scripts.
type Script struct {
	variables map[string]*Variable
	input     []byte
}

// New creates a Script instance with an input script.
func New(input []byte) *Script {
	return &Script{
		variables: make(map[string]*Variable),
		input:     input,
	}
}

// Add adds a new variable or updates an existing variable to the script.
func (s *Script) Add(name string, value interface{}) error {
	obj, err := interfaceToObject(value)
	if err != nil {
		return err
	}

	s.variables[name] = &Variable{
		name:  name,
		value: &obj,
	}

	return nil
}

// Remove removes (undefines) an existing variable for the script.
// It returns false if the variable name is not defined.
func (s *Script) Remove(name string) bool {
	if _, ok := s.variables[name]; !ok {
		return false
	}

	delete(s.variables, name)

	return true
}

// Compile compiles the script with all the defined variables, and, returns Compiled object.
func (s *Script) Compile() (*Compiled, error) {
	symbolTable, globals := s.prepCompile()

	fileSet := scanner.NewFileSet()

	p := parser.NewParser(fileSet.AddFile("", -1, len(s.input)), s.input, nil)
	file, err := p.ParseFile()
	if err != nil {
		return nil, fmt.Errorf("parse error: %s", err.Error())
	}

	c := compiler.NewCompiler(symbolTable, nil)
	if err := c.Compile(file); err != nil {
		return nil, err
	}

	return &Compiled{
		symbolTable: symbolTable,
		machine:     vm.NewVM(c.Bytecode(), globals),
	}, nil
}

func (s *Script) prepCompile() (symbolTable *compiler.SymbolTable, globals []*objects.Object) {
	var names []string
	for name := range s.variables {
		names = append(names, name)
	}

	symbolTable = compiler.NewSymbolTable()
	globals = make([]*objects.Object, len(names), len(names))

	for idx, name := range names {
		symbol := symbolTable.Define(name)
		if symbol.Index != idx {
			panic(fmt.Errorf("wrong symbol index: %d != %d", idx, symbol.Index))
		}

		globals[symbol.Index] = s.variables[name].value
	}

	return
}

func (s *Script) copyVariables() map[string]*Variable {
	vars := make(map[string]*Variable)
	for n, v := range s.variables {
		vars[n] = v
	}

	return vars
}

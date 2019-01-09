package tengo

import (
	"fmt"

	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/parser"
	"github.com/d5/tengo/scanner"
	"github.com/d5/tengo/vm"
)

type Variable struct {
	name  string
	value *objects.Object
}

func (v *Variable) Name() string {
	return v.name
}

func (v *Variable) Value() interface{} {
	return nil
}

type Script struct {
	variables map[string]*objects.Object
	input     []byte
}

func NewScript(input []byte) *Script {
	return &Script{
		variables: make(map[string]*objects.Object),
		input:     input,
	}
}

func (s *Script) Add(name string, value interface{}) error {
	obj, err := objects.FromValue(value)
	if err != nil {
		return err
	}

	s.variables[name] = &obj

	return nil
}

func (s *Script) Remove(name string) bool {
	if _, ok := s.variables[name]; !ok {
		return false
	}

	delete(s.variables, name)

	return true
}

func (s *Script) Compile() (*CompiledScript, error) {
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

	return &CompiledScript{
		bytecode:    c.Bytecode(),
		symbolTable: symbolTable,
		globals:     globals,
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

		globals[symbol.Index] = s.variables[name]
	}

	return
}

type CompiledScript struct {
	bytecode    *compiler.Bytecode
	symbolTable *compiler.SymbolTable
	globals     []*objects.Object
}

func (c *CompiledScript) Run() error {
	v := vm.NewVM(c.bytecode, c.globals)

	return v.Run()
}

func (c *CompiledScript) Update(name string, value interface{}) error {
	symbol, _, ok := c.symbolTable.Resolve(name)
	if !ok {
		return fmt.Errorf("name not found: %s", name)
	}

	updated, err := objects.FromValue(value)
	if err != nil {
		return err
	}

	c.globals[symbol.Index] = &updated

	return nil
}

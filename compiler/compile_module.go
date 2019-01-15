package compiler

import (
	"github.com/d5/tengo/compiler/parser"
	"github.com/d5/tengo/compiler/source"
	"github.com/d5/tengo/objects"
)

var (
	fileSet = source.NewFileSet()
)

func compileModule(fileName string, src []byte, precompileGlobals map[string]objects.Object) (*objects.CompiledModule, error) {
	p := parser.NewParser(fileSet.AddFile(fileName, -1, len(src)), src, nil)
	file, err := p.ParseFile()
	if err != nil {
		return nil, err
	}

	var names []string
	for name := range precompileGlobals {
		names = append(names, name)
	}

	symbolTable := NewSymbolTable()
	globals := make(map[string]objects.CompiledModuleGlobal)

	for _, name := range names {
		_ = symbolTable.Define(name)
	}

	c := NewCompiler(symbolTable, nil)
	if err := c.Compile(file); err != nil {
		return nil, err
	}

	for _, name := range symbolTable.Names() {
		symbol, _, _ := symbolTable.Resolve(name)
		if symbol.Scope == ScopeGlobal {
			value := precompileGlobals[name]

			globals[name] = objects.CompiledModuleGlobal{
				Index: symbol.Index,
				Value: &value,
			}
		}
	}

	bytecode := c.Bytecode()

	return &objects.CompiledModule{
		Instructions: bytecode.Instructions,
		Constants:    bytecode.Constants,
		Globals:      globals,
	}, nil
}

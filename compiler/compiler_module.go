package compiler

import (
	"io/ioutil"
	"strings"

	"github.com/d5/tengo/compiler/parser"
	"github.com/d5/tengo/compiler/source"
	"github.com/d5/tengo/objects"
)

var (
	fileSet = source.NewFileSet()
)

func (c *Compiler) compileUserModule(moduleName string) (*objects.CompiledModule, error) {
	compiledModule, exists := c.loadCompiledModule(moduleName)
	if exists {
		return compiledModule, nil
	}

	// read module source from loader
	var moduleSrc []byte
	if c.moduleLoader == nil {
		// default loader: read from local file
		if !strings.HasSuffix(moduleName, ".tengo") {
			moduleName += ".tengo"
		}

		var err error
		moduleSrc, err = ioutil.ReadFile(moduleName)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		moduleSrc, err = c.moduleLoader(moduleName)
		if err != nil {
			return nil, err
		}
	}

	compiledModule, err := c.compileModule(moduleName, moduleSrc)
	if err != nil {
		return nil, err
	}

	c.storeCompiledModule(moduleName, compiledModule)

	return compiledModule, nil
}

func (c *Compiler) compileModule(fileName string, src []byte) (*objects.CompiledModule, error) {
	p := parser.NewParser(fileSet.AddFile(fileName, -1, len(src)), src, nil)
	file, err := p.ParseFile()
	if err != nil {
		return nil, err
	}

	symbolTable := NewSymbolTable()
	globals := make(map[string]int)

	moduleCompiler := c.fork(symbolTable)
	if err := moduleCompiler.Compile(file); err != nil {
		return nil, err
	}

	for _, name := range symbolTable.Names() {
		symbol, _, _ := symbolTable.Resolve(name)
		if symbol.Scope == ScopeGlobal {
			globals[name] = symbol.Index
		}
	}

	return &objects.CompiledModule{
		Instructions: moduleCompiler.Bytecode().Instructions,
		Globals:      globals,
	}, nil
}

func (c *Compiler) loadCompiledModule(moduleName string) (mod *objects.CompiledModule, ok bool) {
	if c.parent != nil {
		return c.parent.loadCompiledModule(moduleName)
	}

	mod, ok = c.compiledModules[moduleName]

	return
}

func (c *Compiler) storeCompiledModule(moduleName string, module *objects.CompiledModule) {
	if c.parent != nil {
		c.parent.storeCompiledModule(moduleName, module)
	}

	c.compiledModules[moduleName] = module
}

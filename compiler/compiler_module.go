package compiler

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/d5/tengo/compiler/parser"
	"github.com/d5/tengo/compiler/source"
	"github.com/d5/tengo/objects"
)

var (
	fileSet = source.NewFileSet()
)

func (c *Compiler) compileModule(moduleName string) (*objects.CompiledFunction, error) {
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

		if err := c.checkCyclicImports(moduleName); err != nil {
			return nil, err
		}

		var err error
		moduleSrc, err = ioutil.ReadFile(moduleName)
		if err != nil {
			return nil, err
		}
	} else {
		if err := c.checkCyclicImports(moduleName); err != nil {
			return nil, err
		}

		var err error
		moduleSrc, err = c.moduleLoader(moduleName)
		if err != nil {
			return nil, err
		}
	}

	compiledModule, err := c.doCompileModule(moduleName, moduleSrc)
	if err != nil {
		return nil, err
	}

	c.storeCompiledModule(moduleName, compiledModule)

	return compiledModule, nil
}

func (c *Compiler) checkCyclicImports(moduleName string) error {
	if c.moduleName == moduleName {
		return fmt.Errorf("cyclic module import: %s", moduleName)
	} else if c.parent != nil {
		return c.parent.checkCyclicImports(moduleName)
	}

	return nil
}

func (c *Compiler) doCompileModule(moduleName string, src []byte) (*objects.CompiledFunction, error) {
	p := parser.NewParser(fileSet.AddFile(moduleName, -1, len(src)), src, nil)
	file, err := p.ParseFile()
	if err != nil {
		return nil, err
	}

	symbolTable := NewSymbolTable()

	// inherit builtin functions
	for idx, fn := range objects.Builtins {
		s, _, ok := c.symbolTable.Resolve(fn.Name)
		if ok && s.Scope == ScopeBuiltin {
			symbolTable.DefineBuiltin(idx, fn.Name)
		}
	}

	// no global scope for the module
	symbolTable = symbolTable.Fork(false)

	// compile module
	moduleCompiler := c.fork(moduleName, symbolTable)
	if err := moduleCompiler.Compile(file); err != nil {
		return nil, err
	}

	// add OpReturn (== export undefined) if export is missing
	if !moduleCompiler.lastInstructionIs(OpReturnValue) {
		moduleCompiler.emit(OpReturn)
	}

	return &objects.CompiledFunction{
		Instructions: moduleCompiler.Bytecode().Instructions,
		NumLocals:    symbolTable.MaxSymbols(),
	}, nil
}

func (c *Compiler) loadCompiledModule(moduleName string) (mod *objects.CompiledFunction, ok bool) {
	if c.parent != nil {
		return c.parent.loadCompiledModule(moduleName)
	}

	mod, ok = c.compiledModules[moduleName]

	return
}

func (c *Compiler) storeCompiledModule(moduleName string, module *objects.CompiledFunction) {
	if c.parent != nil {
		c.parent.storeCompiledModule(moduleName, module)
	}

	c.compiledModules[moduleName] = module
}

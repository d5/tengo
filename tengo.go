package tengo

import (
	"fmt"

	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/parser"
	"github.com/d5/tengo/scanner"
	"github.com/d5/tengo/vm"
)

func Compile(input []byte) (*compiler.Bytecode, error) {
	fileSet := scanner.NewFileSet()

	p := parser.NewParser(fileSet.AddFile("", -1, len(input)), input, nil)
	file, err := p.ParseFile()
	if err != nil {
		return nil, fmt.Errorf("parse error: %s", err.Error())
	}

	c := compiler.NewCompiler(nil, nil)
	if err := c.Compile(file); err != nil {
		return nil, err
	}

	return c.Bytecode(), nil
}

func Run(b *compiler.Bytecode) error {
	v := vm.NewVM(b, nil)

	return v.Run()
}

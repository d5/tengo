package tengo

import (
	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/compiler/parser"
	"github.com/d5/tengo/compiler/source"
	"github.com/d5/tengo/runtime"
)

func Compile(input []byte, filename string) (*compiler.Bytecode, error) {
	fileSet := source.NewFileSet()

	p := parser.NewParser(fileSet.AddFile(filename, -1, len(input)), input, nil)
	file, err := p.ParseFile()
	if err != nil {
		return nil, err
	}

	c := compiler.NewCompiler(nil, nil)
	if err := c.Compile(file); err != nil {
		return nil, err
	}

	return c.Bytecode(), nil
}

func Run(b *compiler.Bytecode) error {
	v := runtime.NewVM(b, nil)

	return v.Run()
}

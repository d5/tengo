package parser

import (
	"io"

	"github.com/d5/tengo/ast"
	"github.com/d5/tengo/scanner"
)

func ParseFile(file *scanner.File, src []byte, trace io.Writer) (res *ast.File, err error) {
	p := NewParser(file, src, trace)

	defer func() {
		if e := recover(); e != nil {
			if _, ok := e.(bailout); !ok {
				panic(e)
			}
		}

		p.errors.Sort()
		err = p.errors.Err()
	}()

	res, err = p.ParseFile()

	return
}

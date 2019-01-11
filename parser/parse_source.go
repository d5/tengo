package parser

import (
	"io"

	"github.com/d5/tengo/ast"
	"github.com/d5/tengo/source"
)

// ParseSource parses source code 'src' and builds an AST.
func ParseSource(src []byte, trace io.Writer) (res *ast.File, err error) {
	fileSet := source.NewFileSet()
	file := fileSet.AddFile("", -1, len(src))

	return ParseFile(file, src, trace)
}

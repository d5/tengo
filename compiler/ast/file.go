package ast

import (
	"strings"

	"github.com/d5/tengo/compiler/source"
)

type File struct {
	InputFile *source.File
	Stmts     []Stmt
}

func (n *File) Pos() source.Pos {
	return source.Pos(n.InputFile.Base())
}

func (n *File) End() source.Pos {
	return source.Pos(n.InputFile.Base() + n.InputFile.Size())
}

func (n *File) String() string {
	var stmts []string
	for _, e := range n.Stmts {
		stmts = append(stmts, e.String())
	}

	return strings.Join(stmts, "; ")
}

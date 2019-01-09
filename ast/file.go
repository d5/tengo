package ast

import (
	"strings"

	"github.com/d5/tengo/scanner"
)

type File struct {
	InputFile *scanner.File
	Stmts     []Stmt
}

func (n *File) Pos() scanner.Pos {
	return scanner.Pos(n.InputFile.Base())
}

func (n *File) End() scanner.Pos {
	return scanner.Pos(n.InputFile.Base() + n.InputFile.Size())
}

func (n *File) String() string {
	var stmts []string
	for _, e := range n.Stmts {
		stmts = append(stmts, e.String())
	}

	return strings.Join(stmts, "; ")
}

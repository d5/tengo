package ast

import (
	"github.com/d5/tengo/scanner"
)

type Node interface {
	Pos() scanner.Pos
	End() scanner.Pos
	String() string
}

package ast

import "github.com/d5/tengo/compiler/source"

type Node interface {
	Pos() source.Pos
	End() source.Pos
	String() string
}

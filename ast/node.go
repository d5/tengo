package ast

import "github.com/d5/tengo/source"

type Node interface {
	Pos() source.Pos
	End() source.Pos
	String() string
}

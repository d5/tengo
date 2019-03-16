package compiler

import "github.com/d5/tengo/objects"

type Module interface {
	Compile() (objects.Object, error)
}

package objects

import (
	"github.com/d5/tengo/compiler/token"
)

type Closure struct {
	Fn   *CompiledFunction
	Free []*Object
}

func (o *Closure) TypeName() string {
	return "closure"
}

func (o *Closure) String() string {
	return "<closure>"
}

func (o *Closure) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

func (o *Closure) Copy() Object {
	return &Closure{
		Fn:   o.Fn.Copy().(*CompiledFunction),
		Free: append([]*Object{}, o.Free...), // DO NOT Copy() of elements; these are variable pointers
	}
}

func (o *Closure) IsFalsy() bool {
	return false
}

func (o *Closure) Equals(x Object) bool {
	return false
}

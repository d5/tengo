package objects

import (
	"github.com/d5/tengo/token"
)

type CompiledFunction struct {
	Instructions  []byte
	NumLocals     int
	NumParameters int
}

func (o *CompiledFunction) TypeName() string {
	return "compiled-function"
}

func (o *CompiledFunction) String() string {
	return "<compiled-function>"
}

func (o *CompiledFunction) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

func (o *CompiledFunction) Copy() Object {
	return &CompiledFunction{
		Instructions:  append([]byte{}, o.Instructions...),
		NumLocals:     o.NumLocals,
		NumParameters: o.NumParameters,
	}
}

func (o *CompiledFunction) IsFalsy() bool {
	return false
}

func (o *CompiledFunction) Equals(x Object) bool {
	return false
}

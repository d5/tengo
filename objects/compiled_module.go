package objects

import (
	"github.com/d5/tengo/compiler/token"
)

// CompiledModule represents a compiled module.
type CompiledModule struct {
	Instructions []byte // compiled instructions
	NumGlobals   int
}

// TypeName returns the name of the type.
func (o *CompiledModule) TypeName() string {
	return "compiled-module"
}

func (o *CompiledModule) String() string {
	return "<compiled-module>"
}

// BinaryOp returns another object that is the result of
// a given binary operator and a right-hand side object.
func (o *CompiledModule) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *CompiledModule) Copy() Object {
	return &CompiledModule{
		Instructions: append([]byte{}, o.Instructions...),
		NumGlobals:   o.NumGlobals,
	}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *CompiledModule) IsFalsy() bool {
	return false
}

// Equals returns true if the value of the type
// is equal to the value of another object.
func (o *CompiledModule) Equals(x Object) bool {
	return false
}

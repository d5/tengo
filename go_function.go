package tengo

import (
	"github.com/d5/tengo/compiler/token"
)

// GoFunction represents a Go function.
type GoFunction struct {
	ObjectImpl
	Name       string
	Value      CallableFunc
	EncodingID string
}

// TypeName returns the name of the type.
func (o *GoFunction) TypeName() string {
	if o.Name == "" {
		return "go-function"
	}
	return "go-function:" + o.Name
}

func (o *GoFunction) String() string {
	return "<go-function>"
}

// BinaryOp returns another object that is the result of
// a given binary operator and a right-hand side object.
func (o *GoFunction) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *GoFunction) Copy() Object {
	return &GoFunction{Value: o.Value}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *GoFunction) IsFalsy() bool {
	return false
}

// Equals returns true if the value of the type
// is equal to the value of another object.
func (o *GoFunction) Equals(x Object) bool {
	return false
}

// Call invokes a user function.
func (o *GoFunction) Call(rt Interop, args ...Object) (Object, error) {
	return o.Value(rt, args...)
}

// CanCall returns whether the Object can be Called.
func (o *GoFunction) CanCall() bool {
	return true
}

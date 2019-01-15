package objects

import "github.com/d5/tengo/compiler/token"

// Undefined represents an undefined value.
type Undefined struct{}

// TypeName returns the name of the type.
func (o Undefined) TypeName() string {
	return "undefined"
}

func (o Undefined) String() string {
	return "<undefined>"
}

// BinaryOp returns another object that is the result of
// a given binary operator and a right-hand side object.
func (o Undefined) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o Undefined) Copy() Object {
	return Undefined{}
}

// IsFalsy returns true if the value of the type is falsy.
func (o Undefined) IsFalsy() bool {
	return true
}

// Equals returns true if the value of the type
// is equal to the value of another object.
func (o Undefined) Equals(x Object) bool {
	_, ok := x.(*Undefined)

	return ok
}

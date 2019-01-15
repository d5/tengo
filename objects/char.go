package objects

import (
	"github.com/d5/tengo/compiler/token"
)

// Char represents a character value.
type Char struct {
	Value rune
}

func (o *Char) String() string {
	return string(o.Value)
}

// TypeName returns the name of the type.
func (o *Char) TypeName() string {
	return "char"
}

// BinaryOp returns another object that is the result of
// a given binary operator and a right-hand side object.
func (o *Char) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *Char) Copy() Object {
	return &Char{Value: o.Value}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Char) IsFalsy() bool {
	return o.Value == 0
}

// Equals returns true if the value of the type
// is equal to the value of another object.
func (o *Char) Equals(x Object) bool {
	t, ok := x.(*Char)
	if !ok {
		return false
	}

	return o.Value == t.Value
}

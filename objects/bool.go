package objects

import (
	"github.com/d5/tengo/compiler/token"
)

// Bool represents a boolean value.
type Bool struct {
	Value bool
}

func (o *Bool) String() string {
	if o.Value {
		return "true"
	}

	return "false"
}

// TypeName returns the name of the type.
func (o *Bool) TypeName() string {
	return "bool"
}

// BinaryOp returns another object that is the result of
// a given binary operator and a right-hand side object.
func (o *Bool) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *Bool) Copy() Object {
	v := Bool{Value: o.Value}
	return &v
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Bool) IsFalsy() bool {
	return !o.Value
}

// Equals returns true if the value of the type
// is equal to the value of another object.
func (o *Bool) Equals(x Object) bool {
	t, ok := x.(*Bool)
	if !ok {
		return false
	}

	return o.Value == t.Value
}

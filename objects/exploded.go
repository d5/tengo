package objects

import (
	"fmt"

	"github.com/d5/tengo/compiler/token"
)

// Exploded represents the result of exploding an object.
type Exploded struct {
	ObjectImpl
	Value []Object
}

// TypeName returns the name of the type.
func (o *Exploded) TypeName() string { return "exploded" }

func (o *Exploded) String() string { return fmt.Sprintf("<%s>", o.TypeName()) }

// IsFalsy returns true if the value of the type is falsy.
func (o *Exploded) IsFalsy() bool { return true }

// Equals returns true if the value of the type
// is equal to the value of another object.
func (o *Exploded) Equals(Object) bool { return false }

// BinaryOp returns another object that is the result of
// a given binary operator and a right-hand side object.
func (o *Exploded) BinaryOp(token.Token, Object) (Object, error) { return nil, ErrInvalidOperator }

// Copy returns a deep copy of the Exploded value
func (o *Exploded) Copy() Object {
	copied := make([]Object, len(o.Value))

	for i, v := range o.Value {
		copied[i] = v.Copy()
	}

	return &Exploded{Value: copied}
}

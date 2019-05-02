package objects

import (
	"fmt"

	"github.com/d5/tengo/compiler/token"
)

// SpreadValues represents the result of spreading an object.
type SpreadValues struct {
	ObjectImpl
	Value []Object
}

// TypeName returns the name of the type.
func (o *SpreadValues) TypeName() string { return "spread-values" }

func (o *SpreadValues) String() string { return fmt.Sprintf("<%s>", o.TypeName()) }

// IsFalsy returns true if the value of the type is falsy.
func (o *SpreadValues) IsFalsy() bool { return true }

// Equals returns true if the value of the type
// is equal to the value of another object.
func (o *SpreadValues) Equals(Object) bool { return false }

// BinaryOp returns another object that is the result of
// a given binary operator and a right-hand side object.
func (o *SpreadValues) BinaryOp(token.Token, Object) (Object, error) { return nil, ErrInvalidOperator }

// Copy returns a copy of the SpreadValues value
func (o *SpreadValues) Copy() Object { return &SpreadValues{Value: o.Value} }

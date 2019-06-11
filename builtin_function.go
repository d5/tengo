package tengo

import (
	"github.com/d5/tengo/compiler/token"
)

// BuiltinFunction represents a builtin function.
type BuiltinFunction struct {
	ObjectImpl
	Name  string
	Value CallableFunc
}

// TypeName returns the name of the type.
func (o *BuiltinFunction) TypeName() string {
	return "builtin-function:" + o.Name
}

func (o *BuiltinFunction) String() string {
	return "<builtin-function>"
}

// BinaryOp returns another object that is the result of
// a given binary operator and a right-hand side object.
func (o *BuiltinFunction) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *BuiltinFunction) Copy() Object {
	return &BuiltinFunction{Value: o.Value}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *BuiltinFunction) IsFalsy() bool {
	return false
}

// Equals returns true if the value of the type
// is equal to the value of another object.
func (o *BuiltinFunction) Equals(x Object) bool {
	return false
}

// Call executes a builtin function.
func (o *BuiltinFunction) Call(rt Interop, args ...Object) (Object, error) {
	return o.Value(rt, args...)
}

// CanCall returns whether the Object can be Called.
func (o *BuiltinFunction) CanCall() bool {
	return true
}
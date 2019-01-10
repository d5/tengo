package objects

import (
	"github.com/d5/tengo/token"
)

type BuiltinFunction struct {
	Value BuiltinFunc
}

func (o *BuiltinFunction) TypeName() string {
	return "builtin-function"
}

func (o *BuiltinFunction) String() string {
	return "<builtin-function>"
}

func (o *BuiltinFunction) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

func (o *BuiltinFunction) Copy() Object {
	return &BuiltinFunction{Value: o.Value}
}

func (o *BuiltinFunction) IsFalsy() bool {
	return false
}

func (o *BuiltinFunction) Equals(x Object) bool {
	return false
}

func (o *BuiltinFunction) Call(args ...Object) (Object, error) {
	return o.Value(args...)
}

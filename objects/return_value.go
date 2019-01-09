package objects

import "github.com/d5/tengo/token"

type ReturnValue struct {
	Value Object
}

func (o *ReturnValue) TypeName() string {
	return "return-value"
}

func (o *ReturnValue) String() string {
	return "<return>"
}

func (o *ReturnValue) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

func (o *ReturnValue) Copy() Object {
	return &ReturnValue{Value: o.Copy()}
}

func (o *ReturnValue) IsFalsy() bool {
	return false
}

func (o *ReturnValue) Equals(x Object) bool {
	return false
}

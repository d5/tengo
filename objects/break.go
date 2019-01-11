package objects

import "github.com/d5/tengo/compiler/token"

type Break struct{}

func (o *Break) TypeName() string {
	return "break"
}

func (o *Break) String() string {
	return "<break>"
}

func (o *Break) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

func (o *Break) Copy() Object {
	return &Break{}
}

func (o *Break) IsFalsy() bool {
	return false
}

func (o *Break) Equals(x Object) bool {
	return false
}

package objects

import "github.com/d5/tengo/token"

var undefined = &Undefined{}

type Undefined struct{}

func (o Undefined) TypeName() string {
	return "undefined"
}

func (o Undefined) String() string {
	return "<undefined>"
}

func (o Undefined) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

func (o Undefined) Copy() Object {
	return Undefined{}
}

func (o Undefined) IsFalsy() bool {
	return true
}

func (o Undefined) Equals(x Object) bool {
	_, ok := x.(Undefined)

	return ok
}

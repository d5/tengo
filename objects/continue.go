package objects

import "github.com/d5/tengo/token"

type Continue struct {
}

func (o *Continue) TypeName() string {
	return "continue"
}

func (o *Continue) String() string {
	return "<continue>"
}

func (o *Continue) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

func (o *Continue) Copy() Object {
	return &Continue{}
}

func (o *Continue) IsFalsy() bool {
	return false
}

func (o *Continue) Equals(x Object) bool {
	return false
}

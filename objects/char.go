package objects

import (
	"github.com/d5/tengo/compiler/token"
)

type Char struct {
	Value rune
}

func (o *Char) String() string {
	return string(o.Value)
}

func (o *Char) TypeName() string {
	return "char"
}

func (o *Char) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

func (o *Char) Copy() Object {
	return &Char{Value: o.Value}
}

func (o *Char) IsFalsy() bool {
	return o.Value == 0
}

func (o *Char) Equals(x Object) bool {
	t, ok := x.(*Char)
	if !ok {
		return false
	}

	return o.Value == t.Value
}

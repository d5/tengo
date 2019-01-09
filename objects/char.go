package objects

import (
	"fmt"

	"github.com/d5/tengo/token"
)

type Char struct {
	Value rune
}

func (o *Char) String() string {
	return fmt.Sprintf("%q", string(o.Value))
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
	return false
}

func (o *Char) Equals(x Object) bool {
	t, ok := x.(*Char)
	if !ok {
		return false
	}

	return o.Value == t.Value
}

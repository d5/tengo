package objects

import (
	"github.com/d5/tengo/token"
)

type Bool struct {
	Value bool
}

func (o *Bool) String() string {
	if o.Value {
		return "true"
	}

	return "false"
}

func (o *Bool) TypeName() string {
	return "bool"
}

func (o *Bool) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

func (o *Bool) Copy() Object {
	v := Bool(*o)
	return &v
}

func (o *Bool) IsFalsy() bool {
	return !o.Value
}

func (o *Bool) Equals(x Object) bool {
	t, ok := x.(*Bool)
	if !ok {
		return false
	}

	return o.Value == t.Value
}

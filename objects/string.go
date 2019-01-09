package objects

import (
	"fmt"

	"github.com/d5/tengo/token"
)

type String struct {
	Value string
}

func (o *String) TypeName() string {
	return "string"
}

func (o *String) String() string {
	return fmt.Sprintf("%q", o.Value)
}

func (o *String) BinaryOp(op token.Token, rhs Object) (Object, error) {
	switch rhs := rhs.(type) {
	case *String:
		switch op {
		case token.Add:
			return &String{Value: o.Value + rhs.Value}, nil
		}
	case *Char:
		switch op {
		case token.Add:
			return &String{Value: o.Value + string(rhs.Value)}, nil
		}
	}

	return nil, ErrInvalidOperator
}

func (o *String) IsFalsy() bool {
	return len(o.Value) == 0
}

func (o *String) Copy() Object {
	return &String{Value: o.Value}
}

func (o *String) Equals(x Object) bool {
	t, ok := x.(*String)
	if !ok {
		return false
	}

	return o.Value == t.Value
}

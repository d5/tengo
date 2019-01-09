package objects

import (
	"strconv"

	"github.com/d5/tengo/token"
)

type Int struct {
	Value int64
}

func (o *Int) String() string {
	return strconv.FormatInt(o.Value, 10)
}

func (o *Int) TypeName() string {
	return "int"
}

func (o *Int) BinaryOp(op token.Token, rhs Object) (Object, error) {
	switch rhs := rhs.(type) {
	case *Int:
		switch op {
		case token.Add:
			return &Int{o.Value + rhs.Value}, nil
		case token.Sub:
			return &Int{o.Value - rhs.Value}, nil
		case token.Mul:
			return &Int{o.Value * rhs.Value}, nil
		case token.Quo:
			return &Int{o.Value / rhs.Value}, nil
		case token.Rem:
			return &Int{o.Value % rhs.Value}, nil
		case token.And:
			return &Int{o.Value & rhs.Value}, nil
		case token.Or:
			return &Int{o.Value | rhs.Value}, nil
		case token.Xor:
			return &Int{o.Value ^ rhs.Value}, nil
		case token.AndNot:
			return &Int{o.Value &^ rhs.Value}, nil
		case token.Shl:
			return &Int{o.Value << uint(rhs.Value)}, nil
		case token.Shr:
			return &Int{o.Value >> uint(rhs.Value)}, nil
		case token.Less:
			return &Bool{o.Value < rhs.Value}, nil
		case token.Greater:
			return &Bool{o.Value > rhs.Value}, nil
		case token.LessEq:
			return &Bool{o.Value <= rhs.Value}, nil
		case token.GreaterEq:
			return &Bool{o.Value >= rhs.Value}, nil
		}
	case *Float:
		switch op {
		case token.Add:
			return &Float{float64(o.Value) + rhs.Value}, nil
		case token.Sub:
			return &Float{float64(o.Value) - rhs.Value}, nil
		case token.Mul:
			return &Float{float64(o.Value) * rhs.Value}, nil
		case token.Quo:
			return &Float{float64(o.Value) / rhs.Value}, nil
		case token.Less:
			return &Bool{float64(o.Value) < rhs.Value}, nil
		case token.Greater:
			return &Bool{float64(o.Value) > rhs.Value}, nil
		case token.LessEq:
			return &Bool{float64(o.Value) <= rhs.Value}, nil
		case token.GreaterEq:
			return &Bool{float64(o.Value) >= rhs.Value}, nil
		}
	}

	return nil, ErrInvalidOperator
}

func (o *Int) Copy() Object {
	return &Int{o.Value}
}

func (o *Int) IsFalsy() bool {
	return o.Value == 0
}

func (o *Int) Equals(x Object) bool {
	t, ok := x.(*Int)
	if !ok {
		return false
	}

	return o.Value == t.Value
}

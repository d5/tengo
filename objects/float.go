package objects

import (
	"math"
	"strconv"

	"github.com/d5/tengo/token"
)

type Float struct {
	Value float64
}

func (o *Float) String() string {
	return strconv.FormatFloat(o.Value, 'f', -1, 64)
}

func (o *Float) TypeName() string {
	return "float"
}

func (o *Float) BinaryOp(op token.Token, rhs Object) (Object, error) {
	switch rhs := rhs.(type) {
	case *Float:
		switch op {
		case token.Add:
			return &Float{o.Value + rhs.Value}, nil
		case token.Sub:
			return &Float{o.Value - rhs.Value}, nil
		case token.Mul:
			return &Float{o.Value * rhs.Value}, nil
		case token.Quo:
			return &Float{o.Value / rhs.Value}, nil
		case token.Less:
			return &Bool{o.Value < rhs.Value}, nil
		case token.Greater:
			return &Bool{o.Value > rhs.Value}, nil
		case token.LessEq:
			return &Bool{o.Value <= rhs.Value}, nil
		case token.GreaterEq:
			return &Bool{o.Value >= rhs.Value}, nil
		}
	case *Int:
		switch op {
		case token.Add:
			return &Float{o.Value + float64(rhs.Value)}, nil
		case token.Sub:
			return &Float{o.Value - float64(rhs.Value)}, nil
		case token.Mul:
			return &Float{o.Value * float64(rhs.Value)}, nil
		case token.Quo:
			return &Float{o.Value / float64(rhs.Value)}, nil
		case token.Less:
			return &Bool{o.Value < float64(rhs.Value)}, nil
		case token.Greater:
			return &Bool{o.Value > float64(rhs.Value)}, nil
		case token.LessEq:
			return &Bool{o.Value <= float64(rhs.Value)}, nil
		case token.GreaterEq:
			return &Bool{o.Value >= float64(rhs.Value)}, nil
		}
	}

	return nil, ErrInvalidOperator
}

func (o *Float) Copy() Object {
	return &Float{Value: o.Value}
}

func (o *Float) IsFalsy() bool {
	return math.IsNaN(o.Value)
}

func (o *Float) Equals(x Object) bool {
	t, ok := x.(*Float)
	if !ok {
		return false
	}

	return o.Value == t.Value
}

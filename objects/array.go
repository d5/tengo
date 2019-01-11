package objects

import (
	"errors"
	"fmt"
	"strings"

	"github.com/d5/tengo/compiler/token"
)

type Array struct {
	Value []Object
}

func (o *Array) TypeName() string {
	return "array"
}

func (o *Array) String() string {
	var elements []string
	for _, e := range o.Value {
		elements = append(elements, e.String())
	}

	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

func (o *Array) BinaryOp(op token.Token, rhs Object) (Object, error) {
	if rhs, ok := rhs.(*Array); ok {
		switch op {
		case token.Add:
			return &Array{Value: append(o.Value, rhs.Value...)}, nil
		}
	}

	return nil, ErrInvalidOperator
}

func (o *Array) Copy() Object {
	var c []Object
	for _, elem := range o.Value {
		c = append(c, elem.Copy())
	}

	return &Array{Value: c}
}

func (o *Array) IsFalsy() bool {
	return len(o.Value) == 0
}

func (o *Array) Equals(x Object) bool {
	t, ok := x.(*Array)
	if !ok {
		return false
	}

	if len(o.Value) != len(t.Value) {
		return false
	}

	for i, e := range o.Value {
		if !e.Equals(t.Value[i]) {
			return false
		}
	}

	return true
}

func (o *Array) Get(index int) (Object, error) {
	if index < 0 || index >= len(o.Value) {
		return nil, errors.New("array index out of bounds")
	}

	return o.Value[index], nil
}

func (o *Array) Set(index int, value Object) error {
	if index < 0 || index >= len(o.Value) {
		return errors.New("array index out of bounds")
	}

	o.Value[index] = value

	return nil
}

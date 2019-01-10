package objects

import (
	"fmt"
	"strings"

	"github.com/d5/tengo/token"
)

type Map struct {
	Value map[string]Object
}

func (o *Map) TypeName() string {
	return "map"
}

func (o *Map) String() string {
	var pairs []string
	for k, v := range o.Value {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.String()))
	}

	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

func (o *Map) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

func (o *Map) Copy() Object {
	c := make(map[string]Object)
	for k, v := range o.Value {
		c[k] = v.Copy()
	}

	return &Map{Value: c}
}

func (o *Map) IsFalsy() bool {
	return len(o.Value) == 0
}

func (o *Map) Get(key string) (Object, bool) {
	val, ok := o.Value[key]

	return val, ok
}

func (o *Map) Set(key string, value Object) {
	o.Value[key] = value
}

func (o *Map) Equals(x Object) bool {
	t, ok := x.(*Map)
	if !ok {
		return false
	}

	if len(o.Value) != len(t.Value) {
		return false
	}

	for k, v := range o.Value {
		tv := t.Value[k]
		if !v.Equals(tv) {
			return false
		}
	}

	return true
}

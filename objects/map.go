package objects

import (
	"fmt"
	"strings"

	"github.com/d5/tengo/compiler/token"
)

// Map represents a map of objects.
type Map struct {
	Value map[string]Object
}

// TypeName returns the name of the type.
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

// BinaryOp returns another object that is the result of
// a given binary operator and a right-hand side object.
func (o *Map) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *Map) Copy() Object {
	c := make(map[string]Object)
	for k, v := range o.Value {
		c[k] = v.Copy()
	}

	return &Map{Value: c}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Map) IsFalsy() bool {
	return len(o.Value) == 0
}

// Get returns the value for the given key.
func (o *Map) Get(key string) (Object, bool) {
	val, ok := o.Value[key]

	return val, ok
}

// Set sets the value for the given key.
func (o *Map) Set(key string, value Object) {
	o.Value[key] = value
}

// Equals returns true if the value of the type
// is equal to the value of another object.
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

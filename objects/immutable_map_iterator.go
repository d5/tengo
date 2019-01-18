package objects

import "github.com/d5/tengo/compiler/token"

// ImmutableMapIterator represents an iterator for the immutable map.
type ImmutableMapIterator struct {
	v map[string]Object
	k []string
	i int
	l int
}

// NewModuleMapIterator creates a module iterator.
func NewModuleMapIterator(v *ImmutableMap) Iterator {
	var keys []string
	for k := range v.Value {
		keys = append(keys, k)
	}

	return &ImmutableMapIterator{
		v: v.Value,
		k: keys,
		l: len(keys),
	}
}

// TypeName returns the name of the type.
func (i *ImmutableMapIterator) TypeName() string {
	return "module-iterator"
}

func (i *ImmutableMapIterator) String() string {
	return "<module-iterator>"
}

// BinaryOp returns another object that is the result of
// a given binary operator and a right-hand side object.
func (i *ImmutableMapIterator) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

// IsFalsy returns true if the value of the type is falsy.
func (i *ImmutableMapIterator) IsFalsy() bool {
	return true
}

// Equals returns true if the value of the type
// is equal to the value of another object.
func (i *ImmutableMapIterator) Equals(Object) bool {
	return false
}

// Copy returns a copy of the type.
func (i *ImmutableMapIterator) Copy() Object {
	return &ImmutableMapIterator{v: i.v, k: i.k, i: i.i, l: i.l}
}

// Next returns true if there are more elements to iterate.
func (i *ImmutableMapIterator) Next() bool {
	i.i++
	return i.i <= i.l
}

// Key returns the key or index value of the current element.
func (i *ImmutableMapIterator) Key() Object {
	k := i.k[i.i-1]

	return &String{Value: k}
}

// Value returns the value of the current element.
func (i *ImmutableMapIterator) Value() Object {
	k := i.k[i.i-1]

	return i.v[k]
}

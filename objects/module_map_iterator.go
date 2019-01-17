package objects

import "github.com/d5/tengo/compiler/token"

// ModuleMapIterator represents an iterator for the module map.
type ModuleMapIterator struct {
	v map[string]Object
	k []string
	i int
	l int
}

// NewModuleMapIterator creates a module iterator.
func NewModuleMapIterator(v *ModuleMap) Iterator {
	var keys []string
	for k := range v.Value {
		keys = append(keys, k)
	}

	return &ModuleMapIterator{
		v: v.Value,
		k: keys,
		l: len(keys),
	}
}

// TypeName returns the name of the type.
func (i *ModuleMapIterator) TypeName() string {
	return "module-iterator"
}

func (i *ModuleMapIterator) String() string {
	return "<module-iterator>"
}

// BinaryOp returns another object that is the result of
// a given binary operator and a right-hand side object.
func (i *ModuleMapIterator) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

// IsFalsy returns true if the value of the type is falsy.
func (i *ModuleMapIterator) IsFalsy() bool {
	return true
}

// Equals returns true if the value of the type
// is equal to the value of another object.
func (i *ModuleMapIterator) Equals(Object) bool {
	return false
}

// Copy returns a copy of the type.
func (i *ModuleMapIterator) Copy() Object {
	return &ModuleMapIterator{v: i.v, k: i.k, i: i.i, l: i.l}
}

// Next returns true if there are more elements to iterate.
func (i *ModuleMapIterator) Next() bool {
	i.i++
	return i.i <= i.l
}

// Key returns the key or index value of the current element.
func (i *ModuleMapIterator) Key() Object {
	k := i.k[i.i-1]

	return &String{Value: k}
}

// Value returns the value of the current element.
func (i *ModuleMapIterator) Value() Object {
	k := i.k[i.i-1]

	return i.v[k]
}

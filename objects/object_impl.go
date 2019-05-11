package objects

import "github.com/d5/tengo/compiler/token"

// ObjectImpl represents a default Object Implementation. To defined a new value type,
// one can embed ObjectImpl in their type declarations to avoid implementing all non-significant
// methods. TypeName() and String() methods still need to be implemented.
type ObjectImpl struct {
}

// TypeName returns the name of the type.
func (o *ObjectImpl) TypeName() string {
	panic(ErrNotImplemented)
}

func (o *ObjectImpl) String() string {
	panic(ErrNotImplemented)
}

// BinaryOp returns another object that is the result of
// a given binary operator and a right-hand side object.
func (o *ObjectImpl) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *ObjectImpl) Copy() Object {
	return nil
}

// IsFalsy returns true if the value of the type is falsy.
func (o *ObjectImpl) IsFalsy() bool {
	return false
}

// Equals returns true if the value of the type
// is equal to the value of another object.
func (o *ObjectImpl) Equals(x Object) bool {
	return o == x
}

// IndexGet returns an element at a given index.
func (o *ObjectImpl) IndexGet(index Object) (res Object, err error) {
	return nil, ErrNotIndexable
}

// IndexSet sets an element at a given index.
func (o *ObjectImpl) IndexSet(index, value Object) (err error) {
	return ErrNotIndexAssignable
}

// Iterate returns an iterator.
func (o *ObjectImpl) Iterate() Iterator {
	return nil
}

// CanIterate returns whether the Object can be Iterated.
func (o *ObjectImpl) CanIterate() bool {
	return false
}

// Call takes an arbitrary number of arguments
// and returns a return value and/or an error.
func (o *ObjectImpl) Call(rt Runtime, args ...Object) (ret Object, err error) {
	return nil, nil
}

// CanCall returns whether the Object can be Called.
func (o *ObjectImpl) CanCall() bool {
	return false
}

// Spread returns a list of Objects.
func (o *ObjectImpl) Spread() []Object {
	return nil
}

// CanSpread returns whether the Object can be Spread.
func (o *ObjectImpl) CanSpread() bool {
	return false
}

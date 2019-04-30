package objects

import "github.com/d5/tengo/compiler/token"

// Object represents an object in the VM.
type Object interface {
	// TypeName should return the name of the type.
	TypeName() string

	// String should return a string representation of the type's value.
	String() string

	// BinaryOp should return another object that is the result of
	// a given binary operator and a right-hand side object.
	// If BinaryOp returns an error, the VM will treat it as a run-time error.
	BinaryOp(op token.Token, rhs Object) (Object, error)

	// IsFalsy should return true if the value of the type
	// should be considered as falsy.
	IsFalsy() bool

	// Equals should return true if the value of the type
	// should be considered as equal to the value of another object.
	Equals(another Object) bool

	// Copy should return a copy of the type (and its value).
	// Copy function will be used for copy() builtin function
	// which is expected to deep-copy the values generally.
	Copy() Object

	// IndexGet should take an index Object and return a result Object or an error for indexable objects.
	// Indexable is an object that can take an index and return an object.
	// If error is returned, the runtime will treat it as a run-time error and ignore returned value.
	// If Object is not indexable, ErrNotIndexable should be returned as error.
	// If nil is returned as value, it will be converted to Undefined value by the runtime.
	IndexGet(index Object) (value Object, err error)

	// IndexSet should take an index Object and a value Object for index assignable objects.
	// Index assignable is an object that can take an index and a value
	// on the left-hand side of the assignment statement.
	// If Object is not index assignable, ErrNotIndexAssignable should be returned as error.
	// If an error is returned, it will be treated as a run-time error.
	IndexSet(index, value Object) error

	// Iterate should return an Iterator for the type.
	Iterate() Iterator

	// CanIterate should return whether the Object can be Iterated.
	CanIterate() bool

	// Call should take an arbitrary number of arguments
	// and returns a return value and/or an error,
	// which the VM will consider as a run-time error.
	Call(args ...Object) (ret Object, err error)

	// CanCall should return whether the Object can be Called.
	CanCall() bool
}

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
	return true
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
func (o *ObjectImpl) Call(args ...Object) (ret Object, err error) {
	return nil, nil
}

// CanCall returns whether the Object can be Called.
func (o *ObjectImpl) CanCall() bool {
	return false
}

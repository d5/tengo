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

	// Call should take a reference to the currently
	// executing Runtime, an arbitrary number of arguments,
	// and should return a return value and/or an error,
	// which the VM will consider as a run-time error.
	Call(rt Runtime, args ...Object) (ret Object, err error)

	// CanCall should return whether the Object can be Called.
	CanCall() bool

	// Spread should return a list of Objects.
	Spread() []Object

	// CanSpread should return whether the Object can be Spread.
	CanSpread() bool
}

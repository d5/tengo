package objects

import "errors"

// ErrNotIndexable means the type is not indexable.
var ErrNotIndexable = errors.New("type is not indexable")

// ErrNotIndexAssignable means the type is not index-assignable.
var ErrNotIndexAssignable = errors.New("type is not index-assignable")

// ErrIndexOutOfBounds is an error where a given index is out of the bounds.
var ErrIndexOutOfBounds = errors.New("index out of bounds")

// ErrInvalidIndexType means the type is not supported as an index.
var ErrInvalidIndexType = errors.New("invalid index type")

// ErrInvalidOperator represents an error for invalid operator usage.
var ErrInvalidOperator = errors.New("invalid operator")

// ErrWrongNumArguments represents a wrong number of arguments error.
var ErrWrongNumArguments = errors.New("wrong number of arguments")

// ErrInvalidTypeConversion  represents an invalid type conversion error.
var ErrInvalidTypeConversion = errors.New("invalid type conversion")

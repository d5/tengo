package objects

import "errors"

// ErrNotCallable represents an error for calling on non-function objects.
var ErrNotCallable = errors.New("not a callable object")

// ErrNotIndexable represents an error for indexing on non-indexable objects.
var ErrNotIndexable = errors.New("non-indexable object")

// ErrInvalidOperator represents an error for invalid operator usage.
var ErrInvalidOperator = errors.New("invalid operator")

// ErrWrongNumArguments represents a wrong number of arguments error.
var ErrWrongNumArguments = errors.New("wrong number of arguments")

// ErrInvalidTypeConversion  represents an invalid type conversion error.
var ErrInvalidTypeConversion = errors.New("invalid type conversion")

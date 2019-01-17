package objects

import "errors"

// ErrInvalidOperator represents an error for invalid operator usage.
var ErrInvalidOperator = errors.New("invalid operator")

// ErrWrongNumArguments represents a wrong number of arguments error.
var ErrWrongNumArguments = errors.New("wrong number of arguments")

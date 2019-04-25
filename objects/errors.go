package objects

import (
	"errors"
	"fmt"
)

// ErrIndexOutOfBounds is an error where a given index is out of the bounds.
var ErrIndexOutOfBounds = errors.New("index out of bounds")

// ErrInvalidIndexType represents an invalid index type.
var ErrInvalidIndexType = errors.New("invalid index type")

// ErrInvalidIndexValueType represents an invalid index value type.
var ErrInvalidIndexValueType = errors.New("invalid index value type")

// ErrInvalidIndexOnError represents an invalid index on error.
var ErrInvalidIndexOnError = errors.New("invalid index on error")

// ErrInvalidOperator represents an error for invalid operator usage.
var ErrInvalidOperator = errors.New("invalid operator")

// ErrWrongNumArguments represents a wrong number of arguments error.
var ErrWrongNumArguments = errors.New("wrong number of arguments")

// ErrBytesLimit represents an error where the size of bytes value exceeds the limit.
var ErrBytesLimit = errors.New("exceeding bytes size limit")

// ErrStringLimit represents an error where the size of string value exceeds the limit.
var ErrStringLimit = errors.New("exceeding string size limit")

// ErrNotIndexable is an error where a given index is out of the bounds.
var ErrNotIndexable = errors.New("not indexable")

// ErrNotIndexable is an error where a given index is out of the bounds.
var ErrNotIndexAssignable = errors.New("not index-assignable")

// ErrInvalidArgumentType represents an invalid argument value type error.
type ErrInvalidArgumentType struct {
	Name     string
	Expected string
	Found    string
}

func (e ErrInvalidArgumentType) Error() string {
	return fmt.Sprintf("invalid type for argument '%s': expected %s, found %s", e.Name, e.Expected, e.Found)
}

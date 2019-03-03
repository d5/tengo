package runtime

import (
	"errors"
)

// ErrStackOverflow is a stack overflow error.
var ErrStackOverflow = errors.New("stack overflow")

// ErrObjectsLimit is an objects limit error.
var ErrObjectsLimit = errors.New("objects limit exceeded")

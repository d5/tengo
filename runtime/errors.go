package runtime

import (
	"errors"
)

var ErrStackOverflow = errors.New("stack overflow")

package objects

import (
	"fmt"
)

func builtinIsError(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments (got=%d, want=1)", len(args))
	}

	switch args[0].(type) {
	case *Error:
		return TrueValue, nil
	}

	return FalseValue, nil
}

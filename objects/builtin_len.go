package objects

import (
	"fmt"
)

func builtinLen(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	switch arg := args[0].(type) {
	case *Array:
		return &Int{Value: int64(len(arg.Value))}, nil
	case *String:
		return &Int{Value: int64(len(arg.Value))}, nil
	default:
		return nil, fmt.Errorf("unsupported type for 'len' function: %s", arg.TypeName())
	}
}

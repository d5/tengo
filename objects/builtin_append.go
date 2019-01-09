package objects

import (
	"fmt"
)

// append(src, items...)
func builtinAppend(args ...Object) (Object, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("not enough arguments in call to append")
	}

	switch arg := args[0].(type) {
	case *Array:
		return &Array{Value: append(arg.Value, args[1:]...)}, nil
	default:
		return nil, fmt.Errorf("unsupported type for 'append' function: %s", arg.TypeName())
	}
}

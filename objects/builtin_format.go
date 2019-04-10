package objects

import (
	"fmt"

	"github.com/d5/tengo"
)

func builtinFormat(args ...Object) (Object, error) {
	numArgs := len(args)
	if numArgs == 0 {
		return nil, ErrWrongNumArguments
	}

	format, ok := args[0].(*String)
	if !ok {
		return nil, ErrInvalidArgumentType{
			Name:     "format",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}
	if numArgs == 1 {
		return format, nil // okay to return 'format' directly as String is immutable
	}

	formatArgs := make([]interface{}, numArgs-1)
	for idx, arg := range args[1:] {
		switch arg := arg.(type) {
		case *Int, *Float, *Bool, *Char, *String, *Bytes:
			formatArgs[idx] = ToInterface(arg)
		default:
			formatArgs[idx] = arg
		}
	}

	s := fmt.Sprintf(format.Value, formatArgs...)

	if len(s) > tengo.MaxStringLen {
		return nil, ErrStringLimit
	}

	return &String{Value: s}, nil
}

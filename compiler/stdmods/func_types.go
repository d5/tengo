package stdmods

import (
	"fmt"

	"github.com/d5/tengo/objects"
)

func funcAFRF(fn func(float64) float64) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}

			switch arg := args[0].(type) {
			case *objects.Int:
				return &objects.Float{Value: fn(float64(arg.Value))}, nil
			case *objects.Float:
				return &objects.Float{Value: fn(arg.Value)}, nil
			default:
				return nil, fmt.Errorf("invalid argument type: %s", arg.TypeName())
			}
		},
	}
}

func funcAFRFI(fn func(float64) float64) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}

			switch arg := args[0].(type) {
			case *objects.Int:
				return &objects.Float{Value: fn(float64(arg.Value))}, nil
			case *objects.Float:
				return &objects.Float{Value: fn(arg.Value)}, nil
			default:
				return nil, fmt.Errorf("invalid argument type: %s", arg.TypeName())
			}
		},
	}
}

func funcAFFRF(fn func(float64, float64) float64) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 2 {
				return nil, objects.ErrWrongNumArguments
			}

			var arg0, arg1 float64

			switch arg := args[0].(type) {
			case *objects.Int:
				arg0 = float64(arg.Value)
			case *objects.Float:
				arg0 = arg.Value
			default:
				return nil, fmt.Errorf("invalid argument type: %s", arg.TypeName())
			}
			switch arg := args[1].(type) {
			case *objects.Int:
				arg1 = float64(arg.Value)
			case *objects.Float:
				arg1 = arg.Value
			default:
				return nil, fmt.Errorf("invalid argument type: %s", arg.TypeName())
			}

			return &objects.Float{Value: fn(arg0, arg1)}, nil
		},
	}
}

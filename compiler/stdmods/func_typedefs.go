package stdmods

import (
	"fmt"

	"github.com/d5/tengo/objects"
)

// FuncAFRF transform a function of 'func(float64) float64' signature
// into a user function object.
func FuncAFRF(fn func(float64) float64) *objects.UserFunction {
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

// FuncARF transform a function of 'func() float64' signature
// into a user function object.
func FuncARF(fn func() float64) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 0 {
				return nil, objects.ErrWrongNumArguments
			}

			return &objects.Float{Value: fn()}, nil
		},
	}
}

// FuncAIRF transform a function of 'func(int) float64' signature
// into a user function object.
func FuncAIRF(fn func(int) float64) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}

			switch arg := args[0].(type) {
			case *objects.Int:
				return &objects.Float{Value: fn(int(arg.Value))}, nil
			case *objects.Float:
				return &objects.Float{Value: fn(int(arg.Value))}, nil
			default:
				return nil, fmt.Errorf("invalid argument type: %s", arg.TypeName())
			}
		},
	}
}

// FuncAFRI transform a function of 'func(float64) int' signature
// into a user function object.
func FuncAFRI(fn func(float64) int) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}

			switch arg := args[0].(type) {
			case *objects.Int:
				return &objects.Int{Value: int64(fn(float64(arg.Value)))}, nil
			case *objects.Float:
				return &objects.Int{Value: int64(fn(arg.Value))}, nil
			default:
				return nil, fmt.Errorf("invalid argument type: %s", arg.TypeName())
			}
		},
	}
}

// FuncAFFRF transform a function of 'func(float64, float64) float64' signature
// into a user function object.
func FuncAFFRF(fn func(float64, float64) float64) *objects.UserFunction {
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

// FuncAIFRF transform a function of 'func(int, float64) float64' signature
// into a user function object.
func FuncAIFRF(fn func(int, float64) float64) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 2 {
				return nil, objects.ErrWrongNumArguments
			}

			var arg0 int
			var arg1 float64

			switch arg := args[0].(type) {
			case *objects.Int:
				arg0 = int(arg.Value)
			case *objects.Float:
				arg0 = int(arg.Value)
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

// FuncAFIRF transform a function of 'func(float64, int) float64' signature
// into a user function object.
func FuncAFIRF(fn func(float64, int) float64) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 2 {
				return nil, objects.ErrWrongNumArguments
			}

			var arg0 float64
			var arg1 int

			switch arg := args[0].(type) {
			case *objects.Int:
				arg0 = float64(arg.Value)
			case *objects.Float:
				arg0 = float64(arg.Value)
			default:
				return nil, fmt.Errorf("invalid argument type: %s", arg.TypeName())
			}
			switch arg := args[1].(type) {
			case *objects.Int:
				arg1 = int(arg.Value)
			case *objects.Float:
				arg1 = int(arg.Value)
			default:
				return nil, fmt.Errorf("invalid argument type: %s", arg.TypeName())
			}

			return &objects.Float{Value: fn(arg0, arg1)}, nil
		},
	}
}

// FuncAFIRB transform a function of 'func(float64, int) bool' signature
// into a user function object.
func FuncAFIRB(fn func(float64, int) bool) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 2 {
				return nil, objects.ErrWrongNumArguments
			}

			var arg0 float64
			var arg1 int

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
				arg1 = int(arg.Value)
			case *objects.Float:
				arg1 = int(arg.Value)
			default:
				return nil, fmt.Errorf("invalid argument type: %s", arg.TypeName())
			}

			return &objects.Bool{Value: fn(arg0, arg1)}, nil
		},
	}
}

// FuncAFRB transform a function of 'func(float64) bool' signature
// into a user function object.
func FuncAFRB(fn func(float64) bool) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}

			var arg0 float64
			switch arg := args[0].(type) {
			case *objects.Int:
				arg0 = float64(arg.Value)
			case *objects.Float:
				arg0 = arg.Value
			default:
				return nil, fmt.Errorf("invalid argument type: %s", arg.TypeName())
			}

			return &objects.Bool{Value: fn(arg0)}, nil
		},
	}
}

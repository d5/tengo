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

func funcARF(fn func() float64) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 0 {
				return nil, objects.ErrWrongNumArguments
			}

			return &objects.Float{Value: fn()}, nil
		},
	}
}

func funcAIRF(fn func(int) float64) *objects.UserFunction {
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

func funcAFRI(fn func(float64) int) *objects.UserFunction {
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

func funcAIFRF(fn func(int, float64) float64) *objects.UserFunction {
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

func funcAFIRF(fn func(float64, int) float64) *objects.UserFunction {
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

func funcAFIRB(fn func(float64, int) bool) *objects.UserFunction {
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

func funcAFRB(fn func(float64) bool) *objects.UserFunction {
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

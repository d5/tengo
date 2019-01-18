package stdlib

import (
	"github.com/d5/tengo/objects"
)

// FuncAR transform a function of 'func()' signature
// into a user function object.
func FuncAR(fn func()) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 0 {
				return nil, objects.ErrWrongNumArguments
			}

			fn()

			return objects.UndefinedValue, nil
		},
	}
}

// FuncARI transform a function of 'func() int' signature
// into a user function object.
func FuncARI(fn func() int) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 0 {
				return nil, objects.ErrWrongNumArguments
			}

			return &objects.Int{Value: int64(fn())}, nil
		},
	}
}

// FuncARB transform a function of 'func() bool' signature
// into a user function object.
func FuncARB(fn func() bool) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 0 {
				return nil, objects.ErrWrongNumArguments
			}

			return &objects.Bool{Value: fn()}, nil
		},
	}
}

// FuncARE transform a function of 'func() error' signature
// into a user function object.
func FuncARE(fn func() error) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 0 {
				return nil, objects.ErrWrongNumArguments
			}

			return wrapError(fn()), nil
		},
	}
}

// FuncARS transform a function of 'func() string' signature
// into a user function object.
func FuncARS(fn func() string) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 0 {
				return nil, objects.ErrWrongNumArguments
			}

			return &objects.String{Value: fn()}, nil
		},
	}
}

// FuncARSE transform a function of 'func() (string, error)' signature
// into a user function object.
func FuncARSE(fn func() (string, error)) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 0 {
				return nil, objects.ErrWrongNumArguments
			}

			res, err := fn()
			if err != nil {
				return wrapError(err), nil
			}

			return &objects.String{Value: res}, nil
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

// FuncARSs transform a function of 'func() []string' signature
// into a user function object.
func FuncARSs(fn func() []string) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 0 {
				return nil, objects.ErrWrongNumArguments
			}

			arr := &objects.Array{}
			for _, osArg := range fn() {
				arr.Value = append(arr.Value, &objects.String{Value: osArg})
			}

			return arr, nil
		},
	}
}

// FuncARIsE transform a function of 'func() ([]int, error)' signature
// into a user function object.
func FuncARIsE(fn func() ([]int, error)) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 0 {
				return nil, objects.ErrWrongNumArguments
			}

			res, err := fn()
			if err != nil {
				return wrapError(err), nil
			}

			arr := &objects.Array{}
			for _, v := range res {
				arr.Value = append(arr.Value, &objects.Int{Value: int64(v)})
			}

			return arr, nil
		},
	}
}

// FuncAFRF transform a function of 'func(float64) float64' signature
// into a user function object.
func FuncAFRF(fn func(float64) float64) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}

			f1, ok := objects.ToFloat64(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			return &objects.Float{Value: fn(f1)}, nil
		},
	}
}

// FuncAIR transform a function of 'func(int)' signature
// into a user function object.
func FuncAIR(fn func(int)) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}

			i1, ok := objects.ToInt(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			fn(i1)

			return objects.UndefinedValue, nil
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

			i1, ok := objects.ToInt(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			return &objects.Float{Value: fn(i1)}, nil
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

			f1, ok := objects.ToFloat64(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			return &objects.Int{Value: int64(fn(f1))}, nil
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

			f1, ok := objects.ToFloat64(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			f2, ok := objects.ToFloat64(args[1])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			return &objects.Float{Value: fn(f1, f2)}, nil
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

			i1, ok := objects.ToInt(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			f2, ok := objects.ToFloat64(args[1])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			return &objects.Float{Value: fn(i1, f2)}, nil
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

			f1, ok := objects.ToFloat64(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			i2, ok := objects.ToInt(args[1])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			return &objects.Float{Value: fn(f1, i2)}, nil
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

			f1, ok := objects.ToFloat64(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			i2, ok := objects.ToInt(args[1])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			return &objects.Bool{Value: fn(f1, i2)}, nil
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

			f1, ok := objects.ToFloat64(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			return &objects.Bool{Value: fn(f1)}, nil
		},
	}
}

// FuncASRS transform a function of 'func(string) string' signature into a user function object.
// User function will return 'true' if underlying native function returns nil.
func FuncASRS(fn func(string) string) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (objects.Object, error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}

			s1, ok := objects.ToString(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			return &objects.String{Value: fn(s1)}, nil
		},
	}
}

// FuncASRSE transform a function of 'func(string) (string, error)' signature into a user function object.
// User function will return 'true' if underlying native function returns nil.
func FuncASRSE(fn func(string) (string, error)) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (objects.Object, error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}

			s1, ok := objects.ToString(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			res, err := fn(s1)
			if err != nil {
				return wrapError(err), nil
			}

			return &objects.String{Value: res}, nil
		},
	}
}

// FuncASRE transform a function of 'func(string) error' signature into a user function object.
// User function will return 'true' if underlying native function returns nil.
func FuncASRE(fn func(string) error) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (objects.Object, error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}

			s1, ok := objects.ToString(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			return wrapError(fn(s1)), nil
		},
	}
}

// FuncASSRE transform a function of 'func(string, string) error' signature into a user function object.
// User function will return 'true' if underlying native function returns nil.
func FuncASSRE(fn func(string, string) error) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (objects.Object, error) {
			if len(args) != 2 {
				return nil, objects.ErrWrongNumArguments
			}

			s1, ok := objects.ToString(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			s2, ok := objects.ToString(args[1])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			return wrapError(fn(s1, s2)), nil
		},
	}
}

// FuncASI64RE transform a function of 'func(string, int64) error' signature
// into a user function object.
func FuncASI64RE(fn func(string, int64) error) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 2 {
				return nil, objects.ErrWrongNumArguments
			}

			s1, ok := objects.ToString(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			i2, ok := objects.ToInt64(args[1])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			return wrapError(fn(s1, i2)), nil
		},
	}
}

// FuncAIIRE transform a function of 'func(int, int) error' signature
// into a user function object.
func FuncAIIRE(fn func(int, int) error) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 2 {
				return nil, objects.ErrWrongNumArguments
			}

			i1, ok := objects.ToInt(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			i2, ok := objects.ToInt(args[1])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			return wrapError(fn(i1, i2)), nil
		},
	}
}

// FuncASIIRE transform a function of 'func(string, int, int) error' signature
// into a user function object.
func FuncASIIRE(fn func(string, int, int) error) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 3 {
				return nil, objects.ErrWrongNumArguments
			}

			s1, ok := objects.ToString(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			i2, ok := objects.ToInt(args[1])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			i3, ok := objects.ToInt(args[2])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			return wrapError(fn(s1, i2, i3)), nil
		},
	}
}

// FuncAYRIE transform a function of 'func([]byte) (int, error)' signature
// into a user function object.
func FuncAYRIE(fn func([]byte) (int, error)) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}

			y1, ok := objects.ToByteSlice(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			res, err := fn(y1)
			if err != nil {
				return wrapError(err), nil
			}

			return &objects.Int{Value: int64(res)}, nil
		},
	}
}

// FuncASRIE transform a function of 'func(string) (int, error)' signature
// into a user function object.
func FuncASRIE(fn func(string) (int, error)) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}

			s1, ok := objects.ToString(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			res, err := fn(s1)
			if err != nil {
				return wrapError(err), nil
			}

			return &objects.Int{Value: int64(res)}, nil
		},
	}
}

// FuncAIRSsE transform a function of 'func(int) ([]string, error)' signature
// into a user function object.
func FuncAIRSsE(fn func(int) ([]string, error)) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}

			i1, ok := objects.ToInt(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			res, err := fn(i1)
			if err != nil {
				return wrapError(err), nil
			}

			arr := &objects.Array{}
			for _, osArg := range res {
				arr.Value = append(arr.Value, &objects.String{Value: osArg})
			}

			return arr, nil
		},
	}
}

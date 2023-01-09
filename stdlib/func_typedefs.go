package stdlib

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/d5/tengo/v2"
)

var typeCodeSplitter = regexp.MustCompile(`[A-Z][^A-Z]*`)

type converter struct {
	arg func(i int, args ...tengo.Object) (interface{}, error)
	ret func(interface{}) (tengo.Object, error)
}

var typeCodeConverters = map[string]*converter{
	// string
	"S": {
		arg: func(i int, args ...tengo.Object) (interface{}, error) {
			return tengo.ToString(i, args...)
		},
		ret: func(i interface{}) (tengo.Object, error) {
			if len(i.(string)) > tengo.MaxStringLen {
				return nil, tengo.ErrStringLimit
			}
			return &tengo.String{Value: i.(string)}, nil
		},
	},
	"Ss": {
		arg: func(i int, args ...tengo.Object) (interface{}, error) {
			return tengo.ToStringSlice(i, args...)
		},
		ret: func(i interface{}) (tengo.Object, error) {
			arr := &tengo.Array{}
			for _, elem := range i.([]string) {
				if len(elem) > tengo.MaxStringLen {
					return nil, tengo.ErrStringLimit
				}
				arr.Value = append(arr.Value, &tengo.String{Value: elem})
			}
			return arr, nil
		},
	},
	// int
	"I": {
		arg: func(i int, args ...tengo.Object) (interface{}, error) {
			return tengo.ToInt(i, args...)
		},
		ret: func(i interface{}) (tengo.Object, error) {
			return &tengo.Int{Value: int64(i.(int))}, nil
		},
	},
	"Is": {
		arg: func(i int, args ...tengo.Object) (interface{}, error) {
			return tengo.ToIntSlice(i, args...)
		},
		ret: func(i interface{}) (tengo.Object, error) {
			arr := &tengo.Array{}
			for _, elem := range i.([]int) {
				arr.Value = append(arr.Value, &tengo.Int{Value: int64(elem)})
			}
			return arr, nil
		},
	},
	// int64
	"I64": {
		arg: func(i int, args ...tengo.Object) (interface{}, error) {
			return tengo.ToInt64(i, args...)
		},
		ret: func(i interface{}) (tengo.Object, error) {
			return &tengo.Int{Value: i.(int64)}, nil
		},
	},
	"I64s": {
		arg: func(i int, args ...tengo.Object) (interface{}, error) {
			return tengo.ToInt64Slice(i, args...)
		},
		ret: func(i interface{}) (tengo.Object, error) {
			arr := &tengo.Array{}
			for _, elem := range i.([]int64) {
				arr.Value = append(arr.Value, &tengo.Int{Value: elem})
			}
			return arr, nil
		},
	},
	// float64
	"F": {
		arg: func(i int, args ...tengo.Object) (interface{}, error) {
			return tengo.ToFloat64(i, args...)
		},
		ret: func(i interface{}) (tengo.Object, error) {
			return &tengo.Float{Value: i.(float64)}, nil
		},
	},
	"Fs": {
		arg: func(i int, args ...tengo.Object) (interface{}, error) {
			return tengo.ToFloat64Slice(i, args...)
		},
		ret: func(i interface{}) (tengo.Object, error) {
			arr := &tengo.Array{}
			for _, elem := range i.([]float64) {
				arr.Value = append(arr.Value, &tengo.Float{Value: elem})
			}
			return arr, nil
		},
	},
	// bool
	"B": {
		arg: func(i int, args ...tengo.Object) (interface{}, error) {
			return tengo.ToBool(i, args...), nil
		},
		ret: func(i interface{}) (tengo.Object, error) {
			if i.(bool) {
				return tengo.TrueValue, nil
			}
			return tengo.FalseValue, nil
		},
	},
	// byte
	"Y": {
		arg: func(i int, args ...tengo.Object) (interface{}, error) {
			return tengo.ToByteSlice(i, args...)
		},
		ret: func(i interface{}) (tengo.Object, error) {
			if len(i.([]byte)) > tengo.MaxBytesLen {
				return nil, tengo.ErrBytesLimit
			}
			return &tengo.Bytes{Value: i.([]byte)}, nil
		},
	},
	// error
	"E": {
		arg: func(i int, args ...tengo.Object) (interface{}, error) {
			panic("not implemented")
		},
		ret: func(i interface{}) (tengo.Object, error) {
			if i == nil {
				return wrapError(nil), nil
			}
			return wrapError(i.(error)), nil
		},
	},
}

func base(
	fnApplier func(args ...interface{}) []interface{}, fnSig string) tengo.CallableFunc {

	sigParts := strings.Split(fnSig, "-")
	if len(sigParts) != 2 {
		panic("invalid function signature, missing '-' arg values. return values splitter token")
	}

	argTypeCodes := typeCodeSplitter.FindAllString(sigParts[0], -1)
	argCount := len(argTypeCodes)
	for _, tc := range argTypeCodes {
		if _, exists := typeCodeConverters[tc]; !exists {
			panic(fmt.Errorf("invalid function signature, unexpected arg type code %s", tc))
		}
	}

	retTypeCodes := typeCodeSplitter.FindAllString(sigParts[1], -1)
	retCount := len(retTypeCodes)
	if retCount > 2 {
		panic("invalid function signature")
	}
	for _, tc := range retTypeCodes {
		if _, exists := typeCodeConverters[tc]; !exists {
			panic(fmt.Errorf("invalid function signature, unexpected return type code %s", tc))
		}
	}

	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != argCount {
			return nil, tengo.ErrInvalidArgumentCount{Min: argCount, Max: argCount, Actual: len(args)}
		}
		// convert tengo objects to raw type values
		rawArgs := make([]interface{}, 0, len(argTypeCodes))
		for i, tc := range argTypeCodes {
			v, err := typeCodeConverters[tc].arg(i, args...)
			if err != nil {
				return nil, err
			}
			rawArgs = append(rawArgs, v)
		}
		rets := fnApplier(rawArgs...)
		if len(rets) != len(retTypeCodes) {
			return nil, tengo.ErrInvalidReturnValueCount{Expected: retCount, Actual: len(rets)}
		}
		// convert raw type value to tengo objects
		if len(retTypeCodes) == 0 {
			// special case for no return value, return undefined
			return tengo.UndefinedValue, nil
		}
		if retTypeCodes[retCount-1] == "E" && rets[retCount-1] != nil {
			// special case for returning an error
			return wrapError(rets[retCount-1].(error)), nil
		}
		return typeCodeConverters[retTypeCodes[0]].ret(rets[0])
	}
}

// FuncAR transform a function of 'func()' signature into CallableFunc type.
func FuncAR(fn func()) tengo.CallableFunc {
	return base(func(_ ...interface{}) []interface{} {
		fn()
		return nil
	}, "-")
}

// FuncARI transform a function of 'func() int' signature into CallableFunc
// type.
func FuncARI(fn func() int) tengo.CallableFunc {
	return base(func(_ ...interface{}) []interface{} {
		return []interface{}{fn()}
	}, "-I")
}

// FuncARI64 transform a function of 'func() int64' signature into CallableFunc
// type.
func FuncARI64(fn func() int64) tengo.CallableFunc {
	return base(func(_ ...interface{}) []interface{} {
		return []interface{}{fn()}
	}, "-I64")
}

// FuncAI64RI64 transform a function of 'func(int64) int64' signature into
// CallableFunc type.
func FuncAI64RI64(fn func(int64) int64) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(int64))}
	}, "I64-I64")
}

// FuncAI64R transform a function of 'func(int64)' signature into CallableFunc
// type.
func FuncAI64R(fn func(int64)) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		fn(args[0].(int64))
		return nil
	}, "I64-")
}

// FuncARB transform a function of 'func() bool' signature into CallableFunc
// type.
func FuncARB(fn func() bool) tengo.CallableFunc {
	return base(func(_ ...interface{}) []interface{} {
		return []interface{}{fn()}
	}, "-B")
}

// FuncARE transform a function of 'func() error' signature into CallableFunc
// type.
func FuncARE(fn func() error) tengo.CallableFunc {
	return base(func(_ ...interface{}) []interface{} {
		return []interface{}{fn()}
	}, "-E")
}

// FuncARS transform a function of 'func() string' signature into CallableFunc
// type.
func FuncARS(fn func() string) tengo.CallableFunc {
	return base(func(_ ...interface{}) []interface{} {
		return []interface{}{fn()}
	}, "-S")
}

// FuncARSE transform a function of 'func() (string, error)' signature into
// CallableFunc type.
func FuncARSE(fn func() (string, error)) tengo.CallableFunc {
	return base(func(_ ...interface{}) []interface{} {
		v, e := fn()
		return []interface{}{v, e}
	}, "-SE")
}

// FuncARYE transform a function of 'func() ([]byte, error)' signature into
// CallableFunc type.
func FuncARYE(fn func() ([]byte, error)) tengo.CallableFunc {
	return base(func(_ ...interface{}) []interface{} {
		v, e := fn()
		return []interface{}{v, e}
	}, "-YE")
}

// FuncARF transform a function of 'func() float64' signature into CallableFunc
// type.
func FuncARF(fn func() float64) tengo.CallableFunc {
	return base(func(_ ...interface{}) []interface{} {
		return []interface{}{fn()}
	}, "-F")
}

// FuncARSs transform a function of 'func() []string' signature into
// CallableFunc type.
func FuncARSs(fn func() []string) tengo.CallableFunc {
	return base(func(_ ...interface{}) []interface{} {
		return []interface{}{fn()}
	}, "-Ss")
}

// FuncARIsE transform a function of 'func() ([]int, error)' signature into
// CallableFunc type.
func FuncARIsE(fn func() ([]int, error)) tengo.CallableFunc {
	return base(func(_ ...interface{}) []interface{} {
		v, e := fn()
		return []interface{}{v, e}
	}, "-IsE")
}

// FuncAIRIs transform a function of 'func(int) []int' signature into
// CallableFunc type.
func FuncAIRIs(fn func(int) []int) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(int))}
	}, "I-Is")
}

// FuncAFRF transform a function of 'func(float64) float64' signature into
// CallableFunc type.
func FuncAFRF(fn func(float64) float64) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(float64))}
	}, "F-F")
}

// FuncAIR transform a function of 'func(int)' signature into CallableFunc type.
func FuncAIR(fn func(int)) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		fn(args[0].(int))
		return nil
	}, "I-")
}

// FuncAIRF transform a function of 'func(int) float64' signature into
// CallableFunc type.
func FuncAIRF(fn func(int) float64) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(int))}
	}, "I-F")
}

// FuncAFRI transform a function of 'func(float64) int' signature into
// CallableFunc type.
func FuncAFRI(fn func(float64) int) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(float64))}
	}, "F-I")
}

// FuncAFFRF transform a function of 'func(float64, float64) float64' signature
// into CallableFunc type.
func FuncAFFRF(fn func(float64, float64) float64) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(float64), args[1].(float64))}
	}, "FF-F")
}

// FuncAIFRF transform a function of 'func(int, float64) float64' signature
// into CallableFunc type.
func FuncAIFRF(fn func(int, float64) float64) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(int), args[1].(float64))}
	}, "IF-F")
}

// FuncAFIRF transform a function of 'func(float64, int) float64' signature
// into CallableFunc type.
func FuncAFIRF(fn func(float64, int) float64) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(float64), args[1].(int))}
	}, "FI-F")
}

// FuncAFIRB transform a function of 'func(float64, int) bool' signature
// into CallableFunc type.
func FuncAFIRB(fn func(float64, int) bool) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(float64), args[1].(int))}
	}, "FI-B")
}

// FuncAFRB transform a function of 'func(float64) bool' signature
// into CallableFunc type.
func FuncAFRB(fn func(float64) bool) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(float64))}
	}, "F-B")
}

// FuncASRS transform a function of 'func(string) string' signature into
// CallableFunc type. User function will return 'true' if underlying native
// function returns nil.
func FuncASRS(fn func(string) string) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(string))}
	}, "S-S")
}

// FuncASRSs transform a function of 'func(string) []string' signature into
// CallableFunc type.
func FuncASRSs(fn func(string) []string) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(string))}
	}, "S-Ss")
}

// FuncASRSE transform a function of 'func(string) (string, error)' signature
// into CallableFunc type. User function will return 'true' if underlying
// native function returns nil.
func FuncASRSE(fn func(string) (string, error)) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		v, e := fn(args[0].(string))
		return []interface{}{v, e}
	}, "S-SE")
}

// FuncASRE transform a function of 'func(string) error' signature into
// CallableFunc type. User function will return 'true' if underlying native
// function returns nil.
func FuncASRE(fn func(string) error) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(string))}
	}, "S-E")
}

// FuncASSRE transform a function of 'func(string, string) error' signature
// into CallableFunc type. User function will return 'true' if underlying
// native function returns nil.
func FuncASSRE(fn func(string, string) error) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(string), args[1].(string))}
	}, "SS-E")
}

// FuncASSRSs transform a function of 'func(string, string) []string'
// signature into CallableFunc type.
func FuncASSRSs(fn func(string, string) []string) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(string), args[1].(string))}
	}, "SS-Ss")
}

// FuncASSIRSs transform a function of 'func(string, string, int) []string'
// signature into CallableFunc type.
func FuncASSIRSs(fn func(string, string, int) []string) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(string), args[1].(string), args[2].(int))}
	}, "SSI-Ss")
}

// FuncASSRI transform a function of 'func(string, string) int' signature into
// CallableFunc type.
func FuncASSRI(fn func(string, string) int) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(string), args[1].(string))}
	}, "SS-I")
}

// FuncASSRS transform a function of 'func(string, string) string' signature
// into CallableFunc type.
func FuncASSRS(fn func(string, string) string) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(string), args[1].(string))}
	}, "SS-S")
}

// FuncASSRB transform a function of 'func(string, string) bool' signature
// into CallableFunc type.
func FuncASSRB(fn func(string, string) bool) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(string), args[1].(string))}
	}, "SS-B")
}

// FuncASsSRS transform a function of 'func([]string, string) string' signature
// into CallableFunc type.
func FuncASsSRS(fn func([]string, string) string) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].([]string), args[1].(string))}
	}, "SsS-S")
}

// FuncASI64RE transform a function of 'func(string, int64) error' signature
// into CallableFunc type.
func FuncASI64RE(fn func(string, int64) error) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(string), args[1].(int64))}
	}, "SI64-E")
}

// FuncAIIRE transform a function of 'func(int, int) error' signature
// into CallableFunc type.
func FuncAIIRE(fn func(int, int) error) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(int), args[1].(int))}
	}, "II-E")
}

// FuncASIRS transform a function of 'func(string, int) string' signature
// into CallableFunc type.
func FuncASIRS(fn func(string, int) string) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(string), args[1].(int))}
	}, "SI-S")
}

// FuncASIIRE transform a function of 'func(string, int, int) error' signature
// into CallableFunc type.
func FuncASIIRE(fn func(string, int, int) error) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(string), args[1].(int), args[2].(int))}
	}, "SII-E")
}

// FuncAYRIE transform a function of 'func([]byte) (int, error)' signature
// into CallableFunc type.
func FuncAYRIE(fn func([]byte) (int, error)) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		v, e := fn(args[0].([]byte))
		return []interface{}{v, e}
	}, "Y-IE")
}

// FuncAYRS transform a function of 'func([]byte) string' signature into
// CallableFunc type.
func FuncAYRS(fn func([]byte) string) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].([]byte))}
	}, "Y-S")
}

// FuncASRIE transform a function of 'func(string) (int, error)' signature
// into CallableFunc type.
func FuncASRIE(fn func(string) (int, error)) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		v, e := fn(args[0].(string))
		return []interface{}{v, e}
	}, "S-IE")
}

// FuncASRYE transform a function of 'func(string) ([]byte, error)' signature
// into CallableFunc type.
func FuncASRYE(fn func(string) ([]byte, error)) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		v, e := fn(args[0].(string))
		return []interface{}{v, e}
	}, "S-YE")
}

// FuncAIRSsE transform a function of 'func(int) ([]string, error)' signature
// into CallableFunc type.
func FuncAIRSsE(fn func(int) ([]string, error)) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		v, e := fn(args[0].(int))
		return []interface{}{v, e}
	}, "I-SsE")
}

// FuncAIRS transform a function of 'func(int) string' signature into
// CallableFunc type.
func FuncAIRS(fn func(int) string) tengo.CallableFunc {
	return base(func(args ...interface{}) []interface{} {
		return []interface{}{fn(args[0].(int))}
	}, "I-S")
}

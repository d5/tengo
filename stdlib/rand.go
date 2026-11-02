package stdlib

import (
	"errors"
	"math/rand"

	"github.com/d5/tengo/v2"
)

var randModule = map[string]tengo.Object{
	"int": &tengo.UserFunction{
		Name:  "int",
		Value: FuncARI64(rand.Int63),
	},
	"float": &tengo.UserFunction{
		Name:  "float",
		Value: FuncARF(rand.Float64),
	},
	"intn": &tengo.UserFunction{
		Name:  "intn",
		Value: randIntn(rand.Int63n),
	},
	"exp_float": &tengo.UserFunction{
		Name:  "exp_float",
		Value: FuncARF(rand.ExpFloat64),
	},
	"norm_float": &tengo.UserFunction{
		Name:  "norm_float",
		Value: FuncARF(rand.NormFloat64),
	},
	"perm": &tengo.UserFunction{
		Name:  "perm",
		Value: randPerm(rand.Perm),
	},
	"seed": &tengo.UserFunction{
		Name:  "seed",
		Value: FuncAI64R(rand.Seed),
	},
	"read": &tengo.UserFunction{
		Name: "read",
		Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
			if len(args) != 1 {
				return nil, tengo.ErrWrongNumArguments
			}
			y1, ok := args[0].(*tengo.Bytes)
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     "first",
					Expected: "bytes",
					Found:    args[0].TypeName(),
				}
			}
			res, err := rand.Read(y1.Value)
			if err != nil {
				ret = wrapError(err)
				return
			}
			return &tengo.Int{Value: int64(res)}, nil
		},
	},
	"rand": &tengo.UserFunction{
		Name: "rand",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if len(args) != 1 {
				return nil, tengo.ErrWrongNumArguments
			}
			i1, ok := tengo.ToInt64(args[0])
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     "first",
					Expected: "int(compatible)",
					Found:    args[0].TypeName(),
				}
			}
			src := rand.NewSource(i1)
			return randRand(rand.New(src)), nil
		},
	},
}

func randRand(r *rand.Rand) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"int": &tengo.UserFunction{
				Name:  "int",
				Value: FuncARI64(r.Int63),
			},
			"float": &tengo.UserFunction{
				Name:  "float",
				Value: FuncARF(r.Float64),
			},
			"intn": &tengo.UserFunction{
				Name:  "intn",
				Value: randIntn(r.Int63n),
			},
			"exp_float": &tengo.UserFunction{
				Name:  "exp_float",
				Value: FuncARF(r.ExpFloat64),
			},
			"norm_float": &tengo.UserFunction{
				Name:  "norm_float",
				Value: FuncARF(r.NormFloat64),
			},
			"perm": &tengo.UserFunction{
				Name:  "perm",
				Value: randPerm(r.Perm),
			},
			"seed": &tengo.UserFunction{
				Name:  "seed",
				Value: FuncAI64R(r.Seed),
			},
			"read": &tengo.UserFunction{
				Name: "read",
				Value: func(args ...tengo.Object) (
					ret tengo.Object,
					err error,
				) {
					if len(args) != 1 {
						return nil, tengo.ErrWrongNumArguments
					}
					y1, ok := args[0].(*tengo.Bytes)
					if !ok {
						return nil, tengo.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "bytes",
							Found:    args[0].TypeName(),
						}
					}
					res, err := r.Read(y1.Value)
					if err != nil {
						ret = wrapError(err)
						return
					}
					return &tengo.Int{Value: int64(res)}, nil
				},
			},
		},
	}
}

// randIntn rejects the n that makes Int63n panic.
func randIntn(intn func(int64) int64) tengo.CallableFunc {
	call := FuncAI64RI64(intn)
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) == 1 {
			if n, ok := tengo.ToInt64(args[0]); ok && n <= 0 {
				return nil, errors.New("n must be positive")
			}
		}
		return call(args...)
	}
}

// randPerm rejects the n that makes Perm panic, and caps n as range caps
// the arrays it builds.
func randPerm(perm func(int) []int) tengo.CallableFunc {
	call := FuncAIRIs(perm)
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) == 1 {
			if n, ok := tengo.ToInt(args[0]); ok {
				if n < 0 {
					return nil, errors.New("n must not be negative")
				}
				if n > tengo.MaxRangeLen {
					return nil, tengo.ErrRangeLimit
				}
			}
		}
		return call(args...)
	}
}

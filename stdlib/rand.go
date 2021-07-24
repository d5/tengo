package stdlib

import (
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
		Value: FuncAI64RI64(rand.Int63n),
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
		Value: FuncAIRIs(rand.Perm),
	},
	"seed": &tengo.UserFunction{
		Name:  "seed",
		Value: FuncAI64R(rand.Seed),
	},
	"read": &tengo.UserFunction{
		Name: "read",
		Value: tengo.CheckStrictArgs(func(args ...tengo.Object) (tengo.Object, error) {
			y1, _ := args[0].(*tengo.Bytes)
			res, err := rand.Read(y1.Value)
			if err != nil {
				return wrapError(err), nil
			}
			return &tengo.Int{Value: int64(res)}, nil
		}, tengo.BytesTN),
	},
	"rand": &tengo.UserFunction{
		Name: "rand",
		Value: tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
			i1, err := tengo.ToInt64(0, args...)
			if err != nil {
				return nil, err
			}
			src := rand.NewSource(i1)
			return randRand(rand.New(src)), nil
		}, 1),
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
				Value: FuncAI64RI64(r.Int63n),
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
				Value: FuncAIRIs(r.Perm),
			},
			"seed": &tengo.UserFunction{
				Name:  "seed",
				Value: FuncAI64R(r.Seed),
			},
			"read": &tengo.UserFunction{
				Name: "read",
				Value: tengo.CheckStrictArgs(func(args ...tengo.Object) (
					tengo.Object,
					error,
				) {
					y1, _ := args[0].(*tengo.Bytes)
					res, err := r.Read(y1.Value)
					if err != nil {
						return wrapError(err), nil
					}
					return &tengo.Int{Value: int64(res)}, nil
				}, tengo.BytesTN),
			},
		},
	}
}

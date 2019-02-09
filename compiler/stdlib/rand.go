package stdlib

import (
	"math/rand"

	"github.com/d5/tengo/objects"
)

var randModule = map[string]objects.Object{
	"int":        FuncARI64(rand.Int63),
	"float":      FuncARF(rand.Float64),
	"intn":       FuncAI64RI64(rand.Int63n),
	"exp_float":  FuncARF(rand.ExpFloat64),
	"norm_float": FuncARF(rand.NormFloat64),
	"perm":       FuncAIRIs(rand.Perm),
	"seed":       FuncAI64R(rand.Seed),
	"read": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}

			y1, ok := args[0].(*objects.Bytes)
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			res, err := rand.Read(y1.Value)
			if err != nil {
				ret = wrapError(err)
				return
			}

			return &objects.Int{Value: int64(res)}, nil
		},
	},
	"rand": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}

			i1, ok := objects.ToInt64(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			src := rand.NewSource(i1)

			return randRand(rand.New(src)), nil
		},
	},
}

func randRand(r *rand.Rand) *objects.ImmutableMap {
	return &objects.ImmutableMap{
		Value: map[string]objects.Object{
			"int":        FuncARI64(r.Int63),
			"float":      FuncARF(r.Float64),
			"intn":       FuncAI64RI64(r.Int63n),
			"exp_float":  FuncARF(r.ExpFloat64),
			"norm_float": FuncARF(r.NormFloat64),
			"perm":       FuncAIRIs(r.Perm),
			"seed":       FuncAI64R(r.Seed),
			"read": &objects.UserFunction{
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					if len(args) != 1 {
						return nil, objects.ErrWrongNumArguments
					}

					y1, ok := args[0].(*objects.Bytes)
					if !ok {
						return nil, objects.ErrInvalidTypeConversion
					}

					res, err := r.Read(y1.Value)
					if err != nil {
						ret = wrapError(err)
						return
					}

					return &objects.Int{Value: int64(res)}, nil
				},
			},
		},
	}
}

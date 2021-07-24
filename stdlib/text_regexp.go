package stdlib

import (
	"regexp"

	"github.com/d5/tengo/v2"
)

func makeTextRegexp(re *regexp.Regexp) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			// match(text) => bool
			"match": &tengo.UserFunction{
				Value: tengo.CheckAnyArgs(func(args ...tengo.Object) (
					tengo.Object,
					error,
				) {
					s1, err := tengo.ToString(0, args...)
					if err != nil {
						return nil, err
					}

					if re.MatchString(s1) {
						return tengo.TrueValue, nil
					}
					return tengo.FalseValue, nil
				}, 1),
			},

			// find(text) 			=> array(array({text:,begin:,end:}))/undefined
			// find(text, maxCount) => array(array({text:,begin:,end:}))/undefined
			"find": &tengo.UserFunction{
				Value: tengo.CheckAnyArgs(func(args ...tengo.Object) (
					tengo.Object,
					error,
				) {
					numArgs := len(args)
					s1, err := tengo.ToString(0, args...)
					if err != nil {
						return nil, err
					}

					if numArgs == 1 {
						m := re.FindStringSubmatchIndex(s1)
						if m == nil {
							return tengo.UndefinedValue, nil
						}

						arr := &tengo.Array{}
						for i := 0; i < len(m); i += 2 {
							arr.Value = append(arr.Value,
								&tengo.ImmutableMap{
									Value: map[string]tengo.Object{
										"text": &tengo.String{
											Value: s1[m[i]:m[i+1]],
										},
										"begin": &tengo.Int{
											Value: int64(m[i]),
										},
										"end": &tengo.Int{
											Value: int64(m[i+1]),
										},
									}})
						}

						return &tengo.Array{Value: []tengo.Object{arr}}, nil
					}

					i2, err := tengo.ToInt(1, args...)
					if err != nil {
						return nil, err
					}
					m := re.FindAllStringSubmatchIndex(s1, i2)
					if m == nil {
						return tengo.UndefinedValue, nil
					}

					arr := &tengo.Array{}
					for _, m := range m {
						subMatch := &tengo.Array{}
						for i := 0; i < len(m); i += 2 {
							subMatch.Value = append(subMatch.Value,
								&tengo.ImmutableMap{
									Value: map[string]tengo.Object{
										"text": &tengo.String{
											Value: s1[m[i]:m[i+1]],
										},
										"begin": &tengo.Int{
											Value: int64(m[i]),
										},
										"end": &tengo.Int{
											Value: int64(m[i+1]),
										},
									}})
						}

						arr.Value = append(arr.Value, subMatch)
					}

					return arr, nil
				}, 1, 2),
			},

			// replace(src, repl) => string
			"replace": &tengo.UserFunction{
				Value: tengo.CheckAnyArgs(func(args ...tengo.Object) (
					tengo.Object,
					error,
				) {

					s1, err := tengo.ToString(0, args...)
					if err != nil {
						return nil, err
					}

					s2, err := tengo.ToString(1, args...)
					if err != nil {
						return nil, err
					}

					s, ok := doTextRegexpReplace(re, s1, s2)
					if !ok {
						return nil, tengo.ErrStringLimit
					}

					return &tengo.String{Value: s}, nil
				}, 2),
			},

			// split(text) 			 => array(string)
			// split(text, maxCount) => array(string)
			"split": &tengo.UserFunction{
				Value: tengo.CheckAnyArgs(func(args ...tengo.Object) (
					tengo.Object,
					error,
				) {
					numArgs := len(args)
					s1, err := tengo.ToString(0, args...)
					if err != nil {
						return nil, err
					}

					var i2 = -1
					if numArgs > 1 {
						i2, err = tengo.ToInt(1, args...)
						if err != nil {
							return nil, err
						}
					}

					arr := &tengo.Array{}
					for _, s := range re.Split(s1, i2) {
						arr.Value = append(arr.Value,
							&tengo.String{Value: s})
					}

					return arr, nil
				}, 1, 2),
			},
		},
	}
}

// Size-limit checking implementation of regexp.ReplaceAllString.
func doTextRegexpReplace(re *regexp.Regexp, src, repl string) (string, bool) {
	idx := 0
	out := ""
	for _, m := range re.FindAllStringSubmatchIndex(src, -1) {
		var exp []byte
		exp = re.ExpandString(exp, repl, src, m)
		if len(out)+m[0]-idx+len(exp) > tengo.MaxStringLen {
			return "", false
		}
		out += src[idx:m[0]] + string(exp)
		idx = m[1]
	}
	if idx < len(src) {
		if len(out)+len(src)-idx > tengo.MaxStringLen {
			return "", false
		}
		out += src[idx:]
	}
	return out, true
}

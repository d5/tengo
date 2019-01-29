package stdlib

import (
	"regexp"

	"github.com/d5/tengo/objects"
)

var textModule = map[string]objects.Object{
	// re_match(pattern, text) => bool/error
	"re_match": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 2 {
				err = objects.ErrWrongNumArguments
				return
			}

			s1, ok := objects.ToString(args[0])
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			s2, ok := objects.ToString(args[1])
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			matched, err := regexp.MatchString(s1, s2)
			if err != nil {
				ret = wrapError(err)
				return
			}

			if matched {
				ret = objects.TrueValue
			} else {
				ret = objects.FalseValue
			}

			return
		},
	},

	// re_find(pattern, text) 			 => array(array({text:,begin:,end:}))/undefined
	// re_find(pattern, text, maxCount)  => array(array({text:,begin:,end:}))/undefined
	"re_find": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			numArgs := len(args)
			if numArgs != 2 && numArgs != 3 {
				err = objects.ErrWrongNumArguments
				return
			}

			s1, ok := objects.ToString(args[0])
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			re, err := regexp.Compile(s1)
			if err != nil {
				ret = wrapError(err)
				return
			}

			s2, ok := objects.ToString(args[1])
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			if numArgs < 3 {
				m := re.FindStringSubmatchIndex(s2)
				if m == nil {
					ret = objects.UndefinedValue
					return
				}

				arr := &objects.Array{}
				for i := 0; i < len(m); i += 2 {
					arr.Value = append(arr.Value, &objects.ImmutableMap{Value: map[string]objects.Object{
						"text":  &objects.String{Value: s2[m[i]:m[i+1]]},
						"begin": &objects.Int{Value: int64(m[i])},
						"end":   &objects.Int{Value: int64(m[i+1])},
					}})
				}

				ret = &objects.Array{Value: []objects.Object{arr}}

				return
			}

			i3, ok := objects.ToInt(args[2])
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}
			m := re.FindAllStringSubmatchIndex(s2, i3)
			if m == nil {
				ret = objects.UndefinedValue
				return
			}

			arr := &objects.Array{}
			for _, m := range m {
				subMatch := &objects.Array{}
				for i := 0; i < len(m); i += 2 {
					subMatch.Value = append(subMatch.Value, &objects.ImmutableMap{Value: map[string]objects.Object{
						"text":  &objects.String{Value: s2[m[i]:m[i+1]]},
						"begin": &objects.Int{Value: int64(m[i])},
						"end":   &objects.Int{Value: int64(m[i+1])},
					}})
				}

				arr.Value = append(arr.Value, subMatch)
			}

			ret = arr

			return
		},
	},

	// re_replace(pattern, text, repl) => string/error
	"re_replace": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 3 {
				err = objects.ErrWrongNumArguments
				return
			}

			s1, ok := objects.ToString(args[0])
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			s2, ok := objects.ToString(args[1])
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			s3, ok := objects.ToString(args[2])
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			re, err := regexp.Compile(s1)
			if err != nil {
				ret = wrapError(err)
			} else {
				ret = &objects.String{Value: re.ReplaceAllString(s2, s3)}
			}

			return
		},
	},

	// re_split(pattern, text) => array(string)/error
	// re_split(pattern, text, maxCount) => array(string)/error
	"re_split": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			numArgs := len(args)
			if numArgs != 2 && numArgs != 3 {
				err = objects.ErrWrongNumArguments
				return
			}

			s1, ok := objects.ToString(args[0])
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			s2, ok := objects.ToString(args[1])
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			var i3 = -1
			if numArgs > 2 {
				i3, ok = objects.ToInt(args[2])
				if !ok {
					err = objects.ErrInvalidTypeConversion
					return
				}
			}

			re, err := regexp.Compile(s1)
			if err != nil {
				ret = wrapError(err)
				return
			}

			arr := &objects.Array{}
			for _, s := range re.Split(s2, i3) {
				arr.Value = append(arr.Value, &objects.String{Value: s})
			}

			ret = arr

			return
		},
	},

	// re_compile(pattern) => Regexp/error
	"re_compile": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 1 {
				err = objects.ErrWrongNumArguments
				return
			}

			s1, ok := objects.ToString(args[0])
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			re, err := regexp.Compile(s1)
			if err != nil {
				ret = wrapError(err)
			} else {
				ret = stringsRegexpImmutableMap(re)
			}

			return
		},
	},
}

func stringsRegexpImmutableMap(re *regexp.Regexp) *objects.ImmutableMap {
	return &objects.ImmutableMap{
		Value: map[string]objects.Object{
			// match(text) => bool
			"match": &objects.UserFunction{
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					if len(args) != 1 {
						err = objects.ErrWrongNumArguments
						return
					}

					s1, ok := objects.ToString(args[0])
					if !ok {
						err = objects.ErrInvalidTypeConversion
						return
					}

					if re.MatchString(s1) {
						ret = objects.TrueValue
					} else {
						ret = objects.FalseValue
					}

					return
				},
			},

			// find(text) 			=> array(array({text:,begin:,end:}))/undefined
			// find(text, maxCount) => array(array({text:,begin:,end:}))/undefined
			"find": &objects.UserFunction{
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					numArgs := len(args)
					if numArgs != 1 && numArgs != 2 {
						err = objects.ErrWrongNumArguments
						return
					}

					s1, ok := objects.ToString(args[0])
					if !ok {
						err = objects.ErrInvalidTypeConversion
						return
					}

					if numArgs == 1 {
						m := re.FindStringSubmatchIndex(s1)
						if m == nil {
							ret = objects.UndefinedValue
							return
						}

						arr := &objects.Array{}
						for i := 0; i < len(m); i += 2 {
							arr.Value = append(arr.Value, &objects.ImmutableMap{Value: map[string]objects.Object{
								"text":  &objects.String{Value: s1[m[i]:m[i+1]]},
								"begin": &objects.Int{Value: int64(m[i])},
								"end":   &objects.Int{Value: int64(m[i+1])},
							}})
						}

						ret = &objects.Array{Value: []objects.Object{arr}}

						return
					}

					i2, ok := objects.ToInt(args[1])
					if !ok {
						err = objects.ErrInvalidTypeConversion
						return
					}
					m := re.FindAllStringSubmatchIndex(s1, i2)
					if m == nil {
						ret = objects.UndefinedValue
						return
					}

					arr := &objects.Array{}
					for _, m := range m {
						subMatch := &objects.Array{}
						for i := 0; i < len(m); i += 2 {
							subMatch.Value = append(subMatch.Value, &objects.ImmutableMap{Value: map[string]objects.Object{
								"text":  &objects.String{Value: s1[m[i]:m[i+1]]},
								"begin": &objects.Int{Value: int64(m[i])},
								"end":   &objects.Int{Value: int64(m[i+1])},
							}})
						}

						arr.Value = append(arr.Value, subMatch)
					}

					ret = arr

					return
				},
			},

			// replace(src, repl) => string
			"replace": &objects.UserFunction{
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					if len(args) != 2 {
						err = objects.ErrWrongNumArguments
						return
					}

					s1, ok := objects.ToString(args[0])
					if !ok {
						err = objects.ErrInvalidTypeConversion
						return
					}

					s2, ok := objects.ToString(args[1])
					if !ok {
						err = objects.ErrInvalidTypeConversion
						return
					}

					ret = &objects.String{Value: re.ReplaceAllString(s1, s2)}

					return
				},
			},

			// split(text) 			 => array(string)
			// split(text, maxCount) => array(string)
			"split": &objects.UserFunction{
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					numArgs := len(args)
					if numArgs != 1 && numArgs != 2 {
						err = objects.ErrWrongNumArguments
						return
					}

					s1, ok := objects.ToString(args[0])
					if !ok {
						err = objects.ErrInvalidTypeConversion
						return
					}

					var i2 = -1
					if numArgs > 1 {
						i2, ok = objects.ToInt(args[1])
						if !ok {
							err = objects.ErrInvalidTypeConversion
							return
						}
					}

					arr := &objects.Array{}
					for _, s := range re.Split(s1, i2) {
						arr.Value = append(arr.Value, &objects.String{Value: s})
					}

					ret = arr

					return
				},
			},
		},
	}
}

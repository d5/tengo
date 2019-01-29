package stdlib

import (
	"regexp"
	"strconv"
	"strings"

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

	// compare(a, b) => int
	"compare": FuncASSRI(strings.Compare),
	// contains(s, substr) => bool
	"contains": FuncASSRB(strings.Contains),
	// contains_any(s, chars) => bool
	"contains_any": FuncASSRB(strings.ContainsAny),
	// count(s, substr) => int
	"count": FuncASSRI(strings.Count),
	// "equal_fold(s, t) => bool
	"equal_fold": FuncASSRB(strings.EqualFold),
	// fields(s) => array(string)
	"fields": FuncASRSs(strings.Fields),
	// has_prefix(s, prefix) => bool
	"has_prefix": FuncASSRB(strings.HasPrefix),
	// has_suffix(s, suffix) => bool
	"has_suffix": FuncASSRB(strings.HasSuffix),
	// index(s, substr) => int
	"index": FuncASSRI(strings.Index),
	// index_any(s, chars) => int
	"index_any": FuncASSRI(strings.IndexAny),
	// join(arr, sep) => string
	"join": FuncASsSRS(strings.Join),
	// last_index(s, substr) => int
	"last_index": FuncASSRI(strings.LastIndex),
	// last_index_any(s, chars) => int
	"last_index_any": FuncASSRI(strings.LastIndexAny),
	// repeat(s, count) => string
	"repeat": FuncASIRS(strings.Repeat),
	// replace(s, old, new, n) => string
	"replace": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 4 {
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

			i4, ok := objects.ToInt(args[3])
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			ret = &objects.String{Value: strings.Replace(s1, s2, s3, i4)}

			return
		},
	},
	// split(s, sep) => []string
	"split": FuncASSRSs(strings.Split),
	// split_after(s, sep) => []string
	"split_after": FuncASSRSs(strings.SplitAfter),
	// split_after_n(s, sep, n) => []string
	"split_after_n": FuncASSIRSs(strings.SplitAfterN),
	// split_n(s, sep, n) => []string
	"split_n": FuncASSIRSs(strings.SplitN),
	// title(s) => string
	"title": FuncASRS(strings.Title),
	// to_lower(s) => string
	"to_lower": FuncASRS(strings.ToLower),
	// to_title(s) => string
	"to_title": FuncASRS(strings.ToTitle),
	// to_upper(s) => string
	"to_upper": FuncASRS(strings.ToUpper),
	// trim_left(s, cutset) => string
	"trim_left": FuncASSRS(strings.TrimLeft),
	// trim_prefix(s, prefix) => string
	"trim_prefix": FuncASSRS(strings.TrimPrefix),
	// trim_right(s, cutset) => string
	"trim_right": FuncASSRS(strings.TrimRight),
	// trim_space(s) => string
	"trim_space": FuncASRS(strings.TrimSpace),
	// trim_suffix(s, suffix) => string
	"trim_suffix": FuncASSRS(strings.TrimSuffix),
	// atoi(str) => int/error
	"atoi": FuncASRIE(strconv.Atoi),
	// format_bool(b) => string
	"format_bool": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 1 {
				err = objects.ErrWrongNumArguments
				return
			}

			b1, ok := args[0].(*objects.Bool)
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			if b1 == objects.TrueValue {
				ret = &objects.String{Value: "true"}
			} else {
				ret = &objects.String{Value: "false"}
			}

			return
		},
	},
	// format_float(f, fmt, prec, bits) => string
	"format_float": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 4 {
				err = objects.ErrWrongNumArguments
				return
			}

			f1, ok := args[0].(*objects.Float)
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			s2, ok := objects.ToString(args[1])
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			i3, ok := objects.ToInt(args[2])
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			i4, ok := objects.ToInt(args[3])
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			ret = &objects.String{Value: strconv.FormatFloat(f1.Value, s2[0], i3, i4)}

			return
		},
	},
	// format_int(i, base) => string
	"format_int": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 2 {
				err = objects.ErrWrongNumArguments
				return
			}

			i1, ok := args[0].(*objects.Int)
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			i2, ok := objects.ToInt(args[1])
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			ret = &objects.String{Value: strconv.FormatInt(i1.Value, i2)}

			return
		},
	},
	// itoa(i) => string
	"itoa": FuncAIRS(strconv.Itoa),
	// parse_bool(str) => bool/error
	"parse_bool": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 1 {
				err = objects.ErrWrongNumArguments
				return
			}

			s1, ok := args[0].(*objects.String)
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			parsed, err := strconv.ParseBool(s1.Value)
			if err != nil {
				ret = wrapError(err)
				return
			}

			if parsed {
				ret = objects.TrueValue
			} else {
				ret = objects.FalseValue
			}

			return
		},
	},
	// parse_float(str, bits) => float/error
	"parse_float": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 2 {
				err = objects.ErrWrongNumArguments
				return
			}

			s1, ok := args[0].(*objects.String)
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			i2, ok := objects.ToInt(args[1])
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			parsed, err := strconv.ParseFloat(s1.Value, i2)
			if err != nil {
				ret = wrapError(err)
				return
			}

			ret = &objects.Float{Value: parsed}

			return
		},
	},
	// parse_int(str, base, bits) => int/error
	"parse_int": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) != 3 {
				err = objects.ErrWrongNumArguments
				return
			}

			s1, ok := args[0].(*objects.String)
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			i2, ok := objects.ToInt(args[1])
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			i3, ok := objects.ToInt(args[2])
			if !ok {
				err = objects.ErrInvalidTypeConversion
				return
			}

			parsed, err := strconv.ParseInt(s1.Value, i2, i3)
			if err != nil {
				ret = wrapError(err)
				return
			}

			ret = &objects.Int{Value: parsed}

			return
		},
	},
	// quote(str) => string
	"quote": FuncASRS(strconv.Quote),
	// unquote(str) => string/error
	"unquote": FuncASRSE(strconv.Unquote),
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

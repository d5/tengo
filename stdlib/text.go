package stdlib

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/d5/tengo/v2"
)

var textModule = map[string]tengo.Object{
	"re_match": &tengo.UserFunction{
		Name:  "re_match",
		Value: textREMatch,
	}, // re_match(pattern, text) => bool/error
	"re_find": &tengo.UserFunction{
		Name:  "re_find",
		Value: textREFind,
	}, // re_find(pattern, text, count) => [[{text:,begin:,end:}]]/undefined
	"re_replace": &tengo.UserFunction{
		Name:  "re_replace",
		Value: textREReplace,
	}, // re_replace(pattern, text, repl) => string/error
	"re_split": &tengo.UserFunction{
		Name:  "re_split",
		Value: textRESplit,
	}, // re_split(pattern, text, count) => [string]/error
	"re_compile": &tengo.UserFunction{
		Name:  "re_compile",
		Value: textRECompile,
	}, // re_compile(pattern) => Regexp/error
	"compare": &tengo.UserFunction{
		Name:  "compare",
		Value: FuncASSRI(strings.Compare),
	}, // compare(a, b) => int
	"contains": &tengo.UserFunction{
		Name:  "contains",
		Value: FuncASSRB(strings.Contains),
	}, // contains(s, substr) => bool
	"contains_any": &tengo.UserFunction{
		Name:  "contains_any",
		Value: FuncASSRB(strings.ContainsAny),
	}, // contains_any(s, chars) => bool
	"count": &tengo.UserFunction{
		Name:  "count",
		Value: FuncASSRI(strings.Count),
	}, // count(s, substr) => int
	"equal_fold": &tengo.UserFunction{
		Name:  "equal_fold",
		Value: FuncASSRB(strings.EqualFold),
	}, // "equal_fold(s, t) => bool
	"fields": &tengo.UserFunction{
		Name:  "fields",
		Value: FuncASRSs(strings.Fields),
	}, // fields(s) => [string]
	"has_prefix": &tengo.UserFunction{
		Name:  "has_prefix",
		Value: FuncASSRB(strings.HasPrefix),
	}, // has_prefix(s, prefix) => bool
	"has_suffix": &tengo.UserFunction{
		Name:  "has_suffix",
		Value: FuncASSRB(strings.HasSuffix),
	}, // has_suffix(s, suffix) => bool
	"index": &tengo.UserFunction{
		Name:  "index",
		Value: FuncASSRI(strings.Index),
	}, // index(s, substr) => int
	"index_any": &tengo.UserFunction{
		Name:  "index_any",
		Value: FuncASSRI(strings.IndexAny),
	}, // index_any(s, chars) => int
	"join": &tengo.UserFunction{
		Name:  "join",
		Value: textJoin,
	}, // join(arr, sep) => string
	"last_index": &tengo.UserFunction{
		Name:  "last_index",
		Value: FuncASSRI(strings.LastIndex),
	}, // last_index(s, substr) => int
	"last_index_any": &tengo.UserFunction{
		Name:  "last_index_any",
		Value: FuncASSRI(strings.LastIndexAny),
	}, // last_index_any(s, chars) => int
	"repeat": &tengo.UserFunction{
		Name:  "repeat",
		Value: textRepeat,
	}, // repeat(s, count) => string
	"replace": &tengo.UserFunction{
		Name:  "replace",
		Value: textReplace,
	}, // replace(s, old, new, n) => string
	"substr": &tengo.UserFunction{
		Name:  "substr",
		Value: textSubstring,
	}, // substr(s, lower, upper) => string
	"split": &tengo.UserFunction{
		Name:  "split",
		Value: FuncASSRSs(strings.Split),
	}, // split(s, sep) => [string]
	"split_after": &tengo.UserFunction{
		Name:  "split_after",
		Value: FuncASSRSs(strings.SplitAfter),
	}, // split_after(s, sep) => [string]
	"split_after_n": &tengo.UserFunction{
		Name:  "split_after_n",
		Value: FuncASSIRSs(strings.SplitAfterN),
	}, // split_after_n(s, sep, n) => [string]
	"split_n": &tengo.UserFunction{
		Name:  "split_n",
		Value: FuncASSIRSs(strings.SplitN),
	}, // split_n(s, sep, n) => [string]
	"title": &tengo.UserFunction{
		Name:  "title",
		Value: FuncASRS(strings.Title),
	}, // title(s) => string
	"to_lower": &tengo.UserFunction{
		Name:  "to_lower",
		Value: FuncASRS(strings.ToLower),
	}, // to_lower(s) => string
	"to_title": &tengo.UserFunction{
		Name:  "to_title",
		Value: FuncASRS(strings.ToTitle),
	}, // to_title(s) => string
	"to_upper": &tengo.UserFunction{
		Name:  "to_upper",
		Value: FuncASRS(strings.ToUpper),
	}, // to_upper(s) => string
	"pad_left": &tengo.UserFunction{
		Name:  "pad_left",
		Value: textPadLeft,
	}, // pad_left(s, pad_len, pad_with) => string
	"pad_right": &tengo.UserFunction{
		Name:  "pad_right",
		Value: textPadRight,
	}, // pad_right(s, pad_len, pad_with) => string
	"trim": &tengo.UserFunction{
		Name:  "trim",
		Value: FuncASSRS(strings.Trim),
	}, // trim(s, cutset) => string
	"trim_left": &tengo.UserFunction{
		Name:  "trim_left",
		Value: FuncASSRS(strings.TrimLeft),
	}, // trim_left(s, cutset) => string
	"trim_prefix": &tengo.UserFunction{
		Name:  "trim_prefix",
		Value: FuncASSRS(strings.TrimPrefix),
	}, // trim_prefix(s, prefix) => string
	"trim_right": &tengo.UserFunction{
		Name:  "trim_right",
		Value: FuncASSRS(strings.TrimRight),
	}, // trim_right(s, cutset) => string
	"trim_space": &tengo.UserFunction{
		Name:  "trim_space",
		Value: FuncASRS(strings.TrimSpace),
	}, // trim_space(s) => string
	"trim_suffix": &tengo.UserFunction{
		Name:  "trim_suffix",
		Value: FuncASSRS(strings.TrimSuffix),
	}, // trim_suffix(s, suffix) => string
	"atoi": &tengo.UserFunction{
		Name:  "atoi",
		Value: FuncASRIE(strconv.Atoi),
	}, // atoi(str) => int/error
	"format_bool": &tengo.UserFunction{
		Name:  "format_bool",
		Value: textFormatBool,
	}, // format_bool(b) => string
	"format_float": &tengo.UserFunction{
		Name:  "format_float",
		Value: textFormatFloat,
	}, // format_float(f, fmt, prec, bits) => string
	"format_int": &tengo.UserFunction{
		Name:  "format_int",
		Value: textFormatInt,
	}, // format_int(i, base) => string
	"itoa": &tengo.UserFunction{
		Name:  "itoa",
		Value: FuncAIRS(strconv.Itoa),
	}, // itoa(i) => string
	"parse_bool": &tengo.UserFunction{
		Name:  "parse_bool",
		Value: textParseBool,
	}, // parse_bool(str) => bool/error
	"parse_float": &tengo.UserFunction{
		Name:  "parse_float",
		Value: textParseFloat,
	}, // parse_float(str, bits) => float/error
	"parse_int": &tengo.UserFunction{
		Name:  "parse_int",
		Value: textParseInt,
	}, // parse_int(str, base, bits) => int/error
	"quote": &tengo.UserFunction{
		Name:  "quote",
		Value: FuncASRS(strconv.Quote),
	}, // quote(str) => string
	"unquote": &tengo.UserFunction{
		Name:  "unquote",
		Value: FuncASRSE(strconv.Unquote),
	}, // unquote(str) => string/error
}

var textREMatch = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	s1, err := tengo.ToString(0, args...)
	if err != nil {
		return nil, err
	}

	s2, err := tengo.ToString(1, args...)
	if err != nil {
		return nil, err
	}

	matched, err := regexp.MatchString(s1, s2)
	if err != nil {
		return wrapError(err), nil
	}

	if matched {
		return tengo.TrueValue, nil
	}
	return tengo.FalseValue, nil
}, 2)

var textREFind = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	numArgs := len(args)

	s1, err := tengo.ToString(0, args...)
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(s1)
	if err != nil {
		return wrapError(err), nil
	}

	s2, err := tengo.ToString(1, args...)
	if err != nil {
		return nil, err
	}

	if numArgs < 3 {
		m := re.FindStringSubmatchIndex(s2)
		if m == nil {
			return tengo.UndefinedValue, nil
		}

		arr := &tengo.Array{}
		for i := 0; i < len(m); i += 2 {
			arr.Value = append(arr.Value,
				&tengo.ImmutableMap{Value: map[string]tengo.Object{
					"text":  &tengo.String{Value: s2[m[i]:m[i+1]]},
					"begin": &tengo.Int{Value: int64(m[i])},
					"end":   &tengo.Int{Value: int64(m[i+1])},
				}})
		}

		return &tengo.Array{Value: []tengo.Object{arr}}, nil
	}

	i3, err := tengo.ToInt(2, args...)
	if err != nil {
		return nil, err
	}
	m := re.FindAllStringSubmatchIndex(s2, i3)
	if m == nil {
		return tengo.UndefinedValue, nil
	}

	arr := &tengo.Array{}
	for _, m := range m {
		subMatch := &tengo.Array{}
		for i := 0; i < len(m); i += 2 {
			subMatch.Value = append(subMatch.Value,
				&tengo.ImmutableMap{Value: map[string]tengo.Object{
					"text":  &tengo.String{Value: s2[m[i]:m[i+1]]},
					"begin": &tengo.Int{Value: int64(m[i])},
					"end":   &tengo.Int{Value: int64(m[i+1])},
				}})
		}

		arr.Value = append(arr.Value, subMatch)
	}

	return arr, nil
}, 2, 3)

var textREReplace = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	s1, err := tengo.ToString(0, args...)
	if err != nil {
		return nil, err
	}

	s2, err := tengo.ToString(1, args...)
	if err != nil {
		return nil, err
	}

	s3, err := tengo.ToString(2, args...)
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(s1)
	if err != nil {
		return wrapError(err), nil
	}
	s, ok := doTextRegexpReplace(re, s2, s3)
	if !ok {
		return nil, tengo.ErrStringLimit
	}

	return &tengo.String{Value: s}, nil
}, 3)

var textRESplit = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	numArgs := len(args)

	s1, err := tengo.ToString(0, args...)
	if err != nil {
		return nil, err
	}

	s2, err := tengo.ToString(1, args...)
	if err != nil {
		return nil, err
	}

	var i3 = -1
	if numArgs > 2 {
		i3, err = tengo.ToInt(2, args...)
		if err != nil {
			return nil, err
		}
	}

	re, err := regexp.Compile(s1)
	if err != nil {
		return wrapError(err), nil
	}

	arr := &tengo.Array{}
	for _, s := range re.Split(s2, i3) {
		arr.Value = append(arr.Value, &tengo.String{Value: s})
	}

	return arr, nil
}, 2, 3)

var textRECompile = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	s1, err := tengo.ToString(0, args...)
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(s1)
	if err != nil {
		return wrapError(err), nil
	}
	return makeTextRegexp(re), nil
}, 1)

var textReplace = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {

	s1, err := tengo.ToString(0, args...)
	if err != nil {
		return nil, err
	}

	s2, err := tengo.ToString(1, args...)
	if err != nil {
		return nil, err
	}

	s3, err := tengo.ToString(2, args...)
	if err != nil {
		return nil, err
	}

	i4, err := tengo.ToInt(3, args...)
	if err != nil {
		return nil, err
	}

	s, ok := doTextReplace(s1, s2, s3, i4)
	if !ok {
		return nil, tengo.ErrStringLimit
	}

	return &tengo.String{Value: s}, nil
}, 4)

var textSubstring = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	argslen := len(args)

	s1, err := tengo.ToString(0, args...)
	if err != nil {
		return nil, err
	}

	i2, err := tengo.ToInt(1, args...)
	if err != nil {
		return nil, err
	}

	strlen := len(s1)
	i3 := strlen
	if argslen == 3 {
		i3, err = tengo.ToInt(2, args...)
		if err != nil {
			return nil, err
		}
	}

	if i2 > i3 {
		return nil, tengo.ErrInvalidIndexType
	}

	if i2 < 0 {
		i2 = 0
	} else if i2 > strlen {
		i2 = strlen
	}

	if i3 < 0 {
		i3 = 0
	} else if i3 > strlen {
		i3 = strlen
	}

	return &tengo.String{Value: s1[i2:i3]}, nil
}, 2, 3)

var textPadLeft = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	argslen := len(args)

	s1, err := tengo.ToString(0, args...)
	if err != nil {
		return nil, err
	}

	i2, err := tengo.ToInt(1, args...)
	if err != nil {
		return nil, err
	}

	if i2 > tengo.MaxStringLen {
		return nil, tengo.ErrStringLimit
	}

	sLen := len(s1)
	if sLen >= i2 {
		return &tengo.String{Value: s1}, nil
	}

	s3 := " "
	if argslen == 3 {
		s3, err = tengo.ToString(2, args...)
		if err != nil {
			return nil, err
		}
	}

	padStrLen := len(s3)
	if padStrLen == 0 {
		return &tengo.String{Value: s1}, nil
	}

	padCount := ((i2 - padStrLen) / padStrLen) + 1
	retStr := strings.Repeat(s3, padCount) + s1
	return &tengo.String{Value: retStr[len(retStr)-i2:]}, nil
}, 2, 3)

var textPadRight = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	argslen := len(args)

	s1, err := tengo.ToString(0, args...)
	if err != nil {
		return nil, err
	}

	i2, err := tengo.ToInt(1, args...)
	if err != nil {
		return nil, err
	}

	if i2 > tengo.MaxStringLen {
		return nil, tengo.ErrStringLimit
	}

	sLen := len(s1)
	if sLen >= i2 {
		return &tengo.String{Value: s1}, nil
	}

	s3 := " "
	if argslen == 3 {
		s3, err = tengo.ToString(2, args...)
		if err != nil {
			return nil, err
		}
	}

	padStrLen := len(s3)
	if padStrLen == 0 {
		return &tengo.String{Value: s1}, nil
	}

	padCount := ((i2 - padStrLen) / padStrLen) + 1
	retStr := s1 + strings.Repeat(s3, padCount)
	return &tengo.String{Value: retStr[:i2]}, nil
}, 2, 3)

var textRepeat = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {

	s1, err := tengo.ToString(0, args...)
	if err != nil {
		return nil, err
	}

	i2, err := tengo.ToInt(1, args...)
	if err != nil {
		return nil, err
	}

	if len(s1)*i2 > tengo.MaxStringLen {
		return nil, tengo.ErrStringLimit
	}

	return &tengo.String{Value: strings.Repeat(s1, i2)}, nil
}, 2)

var textJoin = tengo.CheckArgs(func(args ...tengo.Object) (tengo.Object, error) {
	var slen int
	var ss1 []string
	switch arg0 := args[0].(type) {
	case *tengo.Array:
		for idx := range arg0.Value {
			as, err := tengo.ToString(idx, arg0.Value...)
			if err != nil {
				return nil, err
			}
			slen += len(as)
			ss1 = append(ss1, as)
		}
	case *tengo.ImmutableArray:
		for idx := range arg0.Value {
			as, err := tengo.ToString(idx, arg0.Value...)
			if err != nil {
				return nil, err
			}
			slen += len(as)
			ss1 = append(ss1, as)
		}
	default:
		panic("impossible")
	}

	s2, err := tengo.ToString(1, args...)
	if err != nil {
		return nil, err
	}

	// make sure output length does not exceed the limit
	if slen+len(s2)*(len(ss1)-1) > tengo.MaxStringLen {
		return nil, tengo.ErrStringLimit
	}

	return &tengo.String{Value: strings.Join(ss1, s2)}, nil
}, 2, 2, tengo.TNs{tengo.ArrayTN, tengo.ImmutableArrayTN}, nil)

var textFormatBool = tengo.CheckStrictArgs(func(args ...tengo.Object) (tengo.Object, error) {
	b1 := args[0].(*tengo.Bool)

	if b1 == tengo.TrueValue {
		return &tengo.String{Value: "true"}, nil
	}
	return &tengo.String{Value: "false"}, nil
}, tengo.BoolTN)

var textFormatFloat = tengo.CheckArgs(func(args ...tengo.Object) (tengo.Object, error) {
	f1 := args[0].(*tengo.Float)

	s2, err := tengo.ToString(1, args...)
	if err != nil {
		return nil, err
	}

	i3, err := tengo.ToInt(2, args...)
	if err != nil {
		return nil, err

	}

	i4, err := tengo.ToInt(3, args...)
	if err != nil {
		return nil, err

	}

	return &tengo.String{Value: strconv.FormatFloat(f1.Value, s2[0], i3, i4)}, nil

},
	4,
	4,
	tengo.TNs{tengo.FloatTN},
	nil,
	nil,
	nil)

var textFormatInt = tengo.CheckArgs(func(args ...tengo.Object) (tengo.Object, error) {
	i1 := args[0].(*tengo.Int)

	i2, err := tengo.ToInt(1, args...)
	if err != nil {
		return nil, err
	}

	return &tengo.String{Value: strconv.FormatInt(i1.Value, i2)}, nil

}, 2, 2, tengo.TNs{tengo.IntTN}, nil)

var textParseBool = tengo.CheckStrictArgs(func(args ...tengo.Object) (tengo.Object, error) {
	s1 := args[0].(*tengo.String)

	parsed, err := strconv.ParseBool(s1.Value)
	if err != nil {
		return wrapError(err), nil

	}

	if parsed {
		return tengo.TrueValue, nil
	}
	return tengo.FalseValue, nil
}, tengo.StringTN)

var textParseFloat = tengo.CheckArgs(func(args ...tengo.Object) (tengo.Object, error) {
	s1 := args[0].(*tengo.String)

	i2, err := tengo.ToInt(1, args...)
	if err != nil {
		return nil, err
	}

	parsed, err := strconv.ParseFloat(s1.Value, i2)
	if err != nil {
		return wrapError(err), nil

	}

	return &tengo.Float{Value: parsed}, nil

}, 2, 2, tengo.TNs{tengo.StringTN}, nil)

var textParseInt = tengo.CheckArgs(func(args ...tengo.Object) (tengo.Object, error) {
	s1 := args[0].(*tengo.String)

	i2, err := tengo.ToInt(1, args...)
	if err != nil {
		return nil, err
	}

	i3, err := tengo.ToInt(2, args...)
	if err != nil {
		return nil, err
	}

	parsed, err := strconv.ParseInt(s1.Value, i2, i3)
	if err != nil {
		return wrapError(err), nil

	}

	return &tengo.Int{Value: parsed}, nil

}, 3, 3, tengo.TNs{tengo.StringTN}, nil, nil)

// Modified implementation of strings.Replace
// to limit the maximum length of output string.
func doTextReplace(s, old, new string, n int) (string, bool) {
	if old == new || n == 0 {
		return s, true // avoid allocation
	}

	// Compute number of replacements.
	if m := strings.Count(s, old); m == 0 {
		return s, true // avoid allocation
	} else if n < 0 || m < n {
		n = m
	}

	// Apply replacements to buffer.
	t := make([]byte, len(s)+n*(len(new)-len(old)))
	w := 0
	start := 0
	for i := 0; i < n; i++ {
		j := start
		if len(old) == 0 {
			if i > 0 {
				_, wid := utf8.DecodeRuneInString(s[start:])
				j += wid
			}
		} else {
			j += strings.Index(s[start:], old)
		}

		ssj := s[start:j]
		if w+len(ssj)+len(new) > tengo.MaxStringLen {
			return "", false
		}

		w += copy(t[w:], ssj)
		w += copy(t[w:], new)
		start = j + len(old)
	}

	ss := s[start:]
	if w+len(ss) > tengo.MaxStringLen {
		return "", false
	}

	w += copy(t[w:], ss)

	return string(t[0:w]), true
}

package stdlib_test

import (
	"regexp"
	"testing"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/require"
)

func TestTextRE(t *testing.T) {
	// re_match(pattern, text)
	for _, d := range []struct {
		pattern string
		text    string
	}{
		{"abc", ""},
		{"abc", "abc"},
		{"a", "abc"},
		{"b", "abc"},
		{"^a", "abc"},
		{"^b", "abc"},
	} {
		expected := regexp.MustCompile(d.pattern).MatchString(d.text)
		module(t, "text").call("re_match", d.pattern, d.text).
			expect(expected, "pattern: %q, src: %q", d.pattern, d.text)
		module(t, "text").call("re_compile", d.pattern).call("match", d.text).
			expect(expected, "patter: %q, src: %q", d.pattern, d.text)
	}

	// re_find(pattern, text)
	for _, d := range []struct {
		pattern  string
		text     string
		expected interface{}
	}{
		{"a(b)", "", tengo.UndefinedValue},
		{"a(b)", "ab", ARR{
			ARR{
				IMAP{"text": "ab", "begin": 0, "end": 2},
				IMAP{"text": "b", "begin": 1, "end": 2},
			},
		}},
		{"a(bc)d", "abcdefgabcd", ARR{
			ARR{
				IMAP{"text": "abcd", "begin": 0, "end": 4},
				IMAP{"text": "bc", "begin": 1, "end": 3},
			},
		}},
		{"(a)b(c)d", "abcdefgabcd", ARR{
			ARR{
				IMAP{"text": "abcd", "begin": 0, "end": 4},
				IMAP{"text": "a", "begin": 0, "end": 1},
				IMAP{"text": "c", "begin": 2, "end": 3},
			},
		}},
	} {
		module(t, "text").call("re_find", d.pattern, d.text).
			expect(d.expected, "pattern: %q, text: %q", d.pattern, d.text)
		module(t, "text").call("re_compile", d.pattern).call("find", d.text).
			expect(d.expected, "pattern: %q, text: %q", d.pattern, d.text)
	}

	// re_find(pattern, text, count))
	for _, d := range []struct {
		pattern  string
		text     string
		count    int
		expected interface{}
	}{
		{"a(b)", "", -1, tengo.UndefinedValue},
		{"a(b)", "ab", -1, ARR{
			ARR{
				IMAP{"text": "ab", "begin": 0, "end": 2},
				IMAP{"text": "b", "begin": 1, "end": 2},
			},
		}},
		{"a(bc)d", "abcdefgabcd", -1, ARR{
			ARR{
				IMAP{"text": "abcd", "begin": 0, "end": 4},
				IMAP{"text": "bc", "begin": 1, "end": 3},
			},
			ARR{
				IMAP{"text": "abcd", "begin": 7, "end": 11},
				IMAP{"text": "bc", "begin": 8, "end": 10},
			},
		}},
		{"(a)b(c)d", "abcdefgabcd", -1, ARR{
			ARR{
				IMAP{"text": "abcd", "begin": 0, "end": 4},
				IMAP{"text": "a", "begin": 0, "end": 1},
				IMAP{"text": "c", "begin": 2, "end": 3},
			},
			ARR{
				IMAP{"text": "abcd", "begin": 7, "end": 11},
				IMAP{"text": "a", "begin": 7, "end": 8},
				IMAP{"text": "c", "begin": 9, "end": 10},
			},
		}},
		{"(a)b(c)d", "abcdefgabcd", 0, tengo.UndefinedValue},
		{"(a)b(c)d", "abcdefgabcd", 1, ARR{
			ARR{
				IMAP{"text": "abcd", "begin": 0, "end": 4},
				IMAP{"text": "a", "begin": 0, "end": 1},
				IMAP{"text": "c", "begin": 2, "end": 3},
			},
		}},
	} {
		module(t, "text").call("re_find", d.pattern, d.text, d.count).
			expect(d.expected, "pattern: %q, text: %q", d.pattern, d.text)
		module(t, "text").call("re_compile", d.pattern).
			call("find", d.text, d.count).
			expect(d.expected, "pattern: %q, text: %q", d.pattern, d.text)
	}

	// re_replace(pattern, text, repl)
	for _, d := range []struct {
		pattern string
		text    string
		repl    string
	}{
		{"a", "", "b"},
		{"a", "a", "b"},
		{"a", "acac", "b"},
		{"b", "acac", "x"},
		{"a", "acac", "123"},
		{"ac", "acac", "99"},
		{"ac$", "acac", "foo"},
		{"a(b)", "ababab", "$1"},
		{"a(b)(c)", "abcabcabc", "$2$1"},
		{"(a(b)c)", "abcabcabc", "$1$2"},
		{"(일(2)삼)", "일2삼12삼일23", "$1$2"},
		{"((일)(2)3)", "일23\n일이3\n일23", "$1$2$3"},
		{"(a(b)c)", "abc\nabc\nabc", "$1$2"},
	} {
		expected := regexp.MustCompile(d.pattern).
			ReplaceAllString(d.text, d.repl)
		module(t, "text").call("re_replace", d.pattern, d.text, d.repl).
			expect(expected, "pattern: %q, text: %q, repl: %q",
				d.pattern, d.text, d.repl)
		module(t, "text").call("re_compile", d.pattern).
			call("replace", d.text, d.repl).
			expect(expected, "pattern: %q, text: %q, repl: %q",
				d.pattern, d.text, d.repl)
	}

	// re_split(pattern, text)
	for _, d := range []struct {
		pattern string
		text    string
	}{
		{"a", ""},
		{"a", "abcabc"},
		{"ab", "abcabc"},
		{"^a", "abcabc"},
	} {
		var expected []interface{}
		for _, ex := range regexp.MustCompile(d.pattern).Split(d.text, -1) {
			expected = append(expected, ex)
		}
		module(t, "text").call("re_split", d.pattern, d.text).
			expect(expected, "pattern: %q, text: %q", d.pattern, d.text)
		module(t, "text").call("re_compile", d.pattern).call("split", d.text).
			expect(expected, "pattern: %q, text: %q", d.pattern, d.text)
	}

	// re_split(pattern, text, count))
	for _, d := range []struct {
		pattern string
		text    string
		count   int
	}{
		{"a", "", -1},
		{"a", "abcabc", -1},
		{"ab", "abcabc", -1},
		{"^a", "abcabc", -1},
		{"a", "abcabc", 0},
		{"a", "abcabc", 1},
		{"a", "abcabc", 2},
		{"a", "abcabc", 3},
		{"b", "abcabc", 1},
		{"b", "abcabc", 2},
		{"b", "abcabc", 3},
	} {
		var expected []interface{}
		for _, ex := range regexp.MustCompile(d.pattern).Split(d.text, d.count) {
			expected = append(expected, ex)
		}
		module(t, "text").call("re_split", d.pattern, d.text, d.count).
			expect(expected, "pattern: %q, text: %q", d.pattern, d.text)
		module(t, "text").call("re_compile", d.pattern).
			call("split", d.text, d.count).
			expect(expected, "pattern: %q, text: %q", d.pattern, d.text)
	}
}

func TestText(t *testing.T) {
	module(t, "text").call("compare", "", "").expect(0)
	module(t, "text").call("compare", "", "a").expect(-1)
	module(t, "text").call("compare", "a", "").expect(1)
	module(t, "text").call("compare", "a", "a").expect(0)
	module(t, "text").call("compare", "a", "b").expect(-1)
	module(t, "text").call("compare", "b", "a").expect(1)
	module(t, "text").call("compare", "abcde", "abcde").expect(0)
	module(t, "text").call("compare", "abcde", "abcdf").expect(-1)
	module(t, "text").call("compare", "abcdf", "abcde").expect(1)

	module(t, "text").call("contains", "", "").expect(true)
	module(t, "text").call("contains", "", "a").expect(false)
	module(t, "text").call("contains", "a", "").expect(true)
	module(t, "text").call("contains", "a", "a").expect(true)
	module(t, "text").call("contains", "abcde", "a").expect(true)
	module(t, "text").call("contains", "abcde", "abcde").expect(true)
	module(t, "text").call("contains", "abc", "abcde").expect(false)
	module(t, "text").call("contains", "ab cd", "bc").expect(false)

	module(t, "text").call("replace", "", "", "", -1).expect("")
	module(t, "text").call("replace", "abcd", "a", "x", -1).expect("xbcd")
	module(t, "text").call("replace", "aaaa", "a", "x", -1).expect("xxxx")
	module(t, "text").call("replace", "aaaa", "a", "x", 0).expect("aaaa")
	module(t, "text").call("replace", "aaaa", "a", "x", 2).expect("xxaa")
	module(t, "text").call("replace", "abcd", "bc", "x", -1).expect("axd")

	module(t, "text").call("format_bool", true).expect("true")
	module(t, "text").call("format_bool", false).expect("false")
	module(t, "text").call("format_float", -19.84, 'f', -1, 64).expect("-19.84")
	module(t, "text").call("format_int", -1984, 10).expect("-1984")
	module(t, "text").call("format_int", 1984, 8).expect("3700")
	module(t, "text").call("parse_bool", "true").expect(true)
	module(t, "text").call("parse_bool", "0").expect(false)
	module(t, "text").call("parse_float", "-19.84", 64).expect(-19.84)
	module(t, "text").call("parse_int", "-1984", 10, 64).expect(-1984)
}

func TestReplaceLimit(t *testing.T) {
	curMaxStringLen := tengo.MaxStringLen
	defer func() { tengo.MaxStringLen = curMaxStringLen }()
	tengo.MaxStringLen = 12

	module(t, "text").call("replace", "123456789012", "1", "x", -1).
		expect("x234567890x2")
	module(t, "text").call("replace", "123456789012", "12", "x", -1).
		expect("x34567890x")
	module(t, "text").call("replace", "123456789012", "1", "xy", -1).
		expectError()
	module(t, "text").call("replace", "123456789012", "0", "xy", -1).
		expectError()
	module(t, "text").call("replace", "123456789012", "012", "xyz", -1).
		expect("123456789xyz")
	module(t, "text").call("replace", "123456789012", "012", "xyzz", -1).
		expectError()

	module(t, "text").call("re_replace", "1", "123456789012", "x").
		expect("x234567890x2")
	module(t, "text").call("re_replace", "12", "123456789012", "x").
		expect("x34567890x")
	module(t, "text").call("re_replace", "1", "123456789012", "xy").
		expectError()
	module(t, "text").call("re_replace", "1(2)", "123456789012", "x$1").
		expect("x234567890x2")
	module(t, "text").call("re_replace", "(1)(2)", "123456789012", "$2$1").
		expect("213456789021")
	module(t, "text").call("re_replace", "(1)(2)", "123456789012", "${2}${1}x").
		expectError()
}

func TestTextArgumentNames(t *testing.T) {
	// a type error names the bad argument and its own type
	for _, fn := range []string{"compare", "split"} {
		res := module(t, "text").call(fn, "a", tengo.UndefinedValue)
		e, ok := res.e.(tengo.ErrInvalidArgumentType)
		require.True(t, ok, res.e)
		require.Equal(t, "second", e.Name, fn)
		require.Equal(t, "undefined", e.Found, fn)
	}
}

func TestTextRepeat(t *testing.T) {
	curMaxStringLen := tengo.MaxStringLen
	defer func() { tengo.MaxStringLen = curMaxStringLen }()
	tengo.MaxStringLen = 12

	module(t, "text").call("repeat", "1234", "3").
		expect("123412341234")
	module(t, "text").call("repeat", "1234", "4").
		expectError()
	module(t, "text").call("repeat", "1", "12").
		expect("111111111111")
	module(t, "text").call("repeat", "1", "13").
		expectError()
}

func TestSubstr(t *testing.T) {
	module(t, "text").call("substr", "", 0, 0).expect("")
	module(t, "text").call("substr", "abcdef", 0, 3).expect("abc")
	module(t, "text").call("substr", "abcdef", 0, 6).expect("abcdef")
	module(t, "text").call("substr", "abcdef", 0, 10).expect("abcdef")
	module(t, "text").call("substr", "abcdef", -10, 10).expect("abcdef")
	module(t, "text").call("substr", "abcdef", 0).expect("abcdef")
	module(t, "text").call("substr", "abcdef", 3).expect("def")

	module(t, "text").call("substr", "", 10, 0).expectError()
	module(t, "text").call("substr", "", "10", 0).expectError()
	module(t, "text").call("substr", "", 10, "0").expectError()
	module(t, "text").call("substr", "", "10", "0").expectError()

	module(t, "text").call("substr", 0, 0, 1).expect("0")
	module(t, "text").call("substr", 123, 0, 1).expect("1")
	module(t, "text").call("substr", 123.456, 4, 7).expect("456")
}

func TestPadLeft(t *testing.T) {
	module(t, "text").call("pad_left", "ab", 7, 0).expect("00000ab")
	module(t, "text").call("pad_right", "ab", 7, 0).expect("ab00000")
	module(t, "text").call("pad_left", "ab", 7, "+-").expect("-+-+-ab")
	module(t, "text").call("pad_right", "ab", 7, "+-").expect("ab+-+-+")
}

func TestTextJoin(t *testing.T) {
	module(t, "text").call("join", ARR{"a", "b", "c"}, ", ").expect("a, b, c")
	module(t, "text").call("join", IARR{"a", "b"}, "-").expect("a-b")
	module(t, "text").call("join", ARR{}, ",").expect("")
	module(t, "text").call("join", ARR{"a"}, ",").expect("a")
	module(t, "text").call("join", ARR{1, "b", true, 'c'}, "").
		expect("1btruec")

	curMaxStringLen := tengo.MaxStringLen
	defer func() { tengo.MaxStringLen = curMaxStringLen }()
	tengo.MaxStringLen = 5
	module(t, "text").call("join", ARR{"ab", "cd"}, "-").expect("ab-cd")
	module(t, "text").call("join", ARR{"ab", "cd"}, "--").
		expectErr(tengo.ErrStringLimit)
	module(t, "text").call("join", ARR{"abc", "def"}, "").
		expectErr(tengo.ErrStringLimit)
}

func TestTextPad(t *testing.T) {
	module(t, "text").call("pad_left", "ab", 4).expect("  ab")
	module(t, "text").call("pad_right", "ab", 4).expect("ab  ")
	module(t, "text").call("pad_left", "ab", 5, "").expect("ab")
	module(t, "text").call("pad_right", "ab", 5, "").expect("ab")
	module(t, "text").call("pad_left", "abc", 2, "x").expect("abc")
	module(t, "text").call("pad_right", "abc", 3, "x").expect("abc")

	curMaxStringLen := tengo.MaxStringLen
	defer func() { tengo.MaxStringLen = curMaxStringLen }()
	tengo.MaxStringLen = 10
	module(t, "text").call("pad_left", "a", 10, "x").expect("xxxxxxxxxa")
	module(t, "text").call("pad_right", "a", 10, "x").expect("axxxxxxxxx")
	module(t, "text").call("pad_left", "a", 11, "x").
		expectErr(tengo.ErrStringLimit)
	module(t, "text").call("pad_right", "a", 11, "x").
		expectErr(tengo.ErrStringLimit)
}

func TestTextFormatParse(t *testing.T) {
	module(t, "text").call("format_float", 1.5, "e", 2, 64).expect("1.50e+00")
	module(t, "text").call("format_float", 0.1, "g", -1, 32).expect("0.1")
	module(t, "text").call("format_int", 255, 16).expect("ff")
	module(t, "text").call("format_int", 5, 2).expect("101")
	module(t, "text").call("format_int", 35, 36).expect("z")
	module(t, "text").call("parse_bool", "T").expect(true)
	module(t, "text").call("parse_bool", "FALSE").expect(false)
	module(t, "text").call("parse_float", "1e3", 64).expect(1000.0)
	module(t, "text").call("parse_float", "0.5", 32).expect(0.5)
	module(t, "text").call("parse_int", "ff", 16, 64).expect(255)
	module(t, "text").call("parse_int", "0x1f", 0, 64).expect(31)
	module(t, "text").call("parse_int", "-101", 2, 8).expect(-5)
}

func TestTextArgErrors(t *testing.T) {
	wrongNum := tengo.ErrWrongNumArguments
	undefined := tengo.UndefinedValue
	invalidString := func(name string) error {
		return invalidArg(name, "string(compatible)", "undefined")
	}
	invalidInt := func(name string) error {
		return invalidArg(name, "int(compatible)", "string")
	}

	cases := []argErrorCase{
		{"re_match", ARR{"a"}, wrongNum},
		{"re_match", ARR{undefined, "a"}, invalidString("first")},
		{"re_match", ARR{"a", undefined}, invalidString("second")},
		{"re_find", ARR{"a"}, wrongNum},
		{"re_find", ARR{"a", "a", 1, 1}, wrongNum},
		{"re_find", ARR{undefined, "a"}, invalidString("first")},
		{"re_find", ARR{"a", undefined}, invalidString("second")},
		{"re_find", ARR{"a", "a", "x"}, invalidInt("third")},
		{"re_replace", ARR{"a", "a"}, wrongNum},
		{"re_replace", ARR{undefined, "a", "b"}, invalidString("first")},
		{"re_replace", ARR{"a", undefined, "b"}, invalidString("second")},
		{"re_replace", ARR{"a", "a", undefined}, invalidString("third")},
		{"re_split", ARR{"a"}, wrongNum},
		{"re_split", ARR{"a", "a", 1, 1}, wrongNum},
		{"re_split", ARR{undefined, "a"}, invalidString("first")},
		{"re_split", ARR{"a", undefined}, invalidString("second")},
		{"re_split", ARR{"a", "a", "x"}, invalidInt("third")},
		{"re_compile", ARR{}, wrongNum},
		{"re_compile", ARR{"a", "a"}, wrongNum},
		{"re_compile", ARR{undefined}, invalidString("first")},
		{"join", ARR{ARR{"a"}}, wrongNum},
		{"join", ARR{ARR{"a"}, ",", ","}, wrongNum},
		{"join", ARR{"a", ","}, invalidArg("first", "array", "string")},
		{"join", ARR{ARR{"a", undefined}, ","}, invalidString("first[1]")},
		{"join", ARR{IARR{undefined}, ","}, invalidString("first[0]")},
		{"join", ARR{ARR{"a"}, undefined}, invalidString("second")},
		{"replace", ARR{"a", "a", "b"}, wrongNum},
		{"replace", ARR{undefined, "a", "b", 1}, invalidString("first")},
		{"replace", ARR{"a", undefined, "b", 1}, invalidString("second")},
		{"replace", ARR{"a", "a", undefined, 1}, invalidString("third")},
		{"replace", ARR{"a", "a", "b", "x"}, invalidInt("fourth")},
		{"substr", ARR{"a"}, wrongNum},
		{"substr", ARR{"a", 0, 1, 2}, wrongNum},
		{"substr", ARR{undefined, 0}, invalidString("first")},
		{"substr", ARR{"a", "x"}, invalidInt("second")},
		{"substr", ARR{"a", 0, "x"}, invalidInt("third")},
		{"format_bool", ARR{}, wrongNum},
		{"format_bool", ARR{true, true}, wrongNum},
		{"format_bool", ARR{1}, invalidArg("first", "bool", "int")},
		{"format_bool", ARR{"true"}, invalidArg("first", "bool", "string")},
		{"format_float", ARR{1.5, "f", 2}, wrongNum},
		{"format_float", ARR{1, "f", 2, 64}, invalidArg("first", "float", "int")},
		{"format_float", ARR{1.5, undefined, 2, 64}, invalidString("second")},
		{"format_float", ARR{1.5, "f", "x", 64}, invalidInt("third")},
		{"format_float", ARR{1.5, "f", 2, "x"}, invalidInt("fourth")},
		{"format_int", ARR{1}, wrongNum},
		{"format_int", ARR{1.5, 10}, invalidArg("first", "int", "float")},
		{"format_int", ARR{1, "x"}, invalidInt("second")},
		{"parse_bool", ARR{}, wrongNum},
		{"parse_bool", ARR{"true", "true"}, wrongNum},
		{"parse_bool", ARR{1}, invalidArg("first", "string", "int")},
		{"parse_float", ARR{"1.5"}, wrongNum},
		{"parse_float", ARR{1.5, 64}, invalidArg("first", "string", "float")},
		{"parse_float", ARR{"1.5", "x"}, invalidInt("second")},
		{"parse_int", ARR{"1", 10}, wrongNum},
		{"parse_int", ARR{1, 10, 64}, invalidArg("first", "string", "int")},
		{"parse_int", ARR{"1", "x", 64}, invalidInt("second")},
		{"parse_int", ARR{"1", 10, "x"}, invalidInt("third")},
	}
	for _, fn := range []string{"pad_left", "pad_right"} {
		cases = append(cases,
			argErrorCase{fn, ARR{"a"}, wrongNum},
			argErrorCase{fn, ARR{"a", 3, "x", "x"}, wrongNum},
			argErrorCase{fn, ARR{undefined, 3}, invalidString("first")},
			argErrorCase{fn, ARR{"a", "x"}, invalidInt("second")},
			argErrorCase{fn, ARR{"a", 3, undefined}, invalidString("third")},
		)
	}
	expectArgErrors(t, "text", cases)
}

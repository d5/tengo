package stdlib_test

import (
	"regexp"
	"testing"

	"github.com/d5/tengo"
	"github.com/d5/tengo/objects"
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
		module(t, "text").call("re_match", mockRuntime{},  d.pattern, d.text).expect(expected, "pattern: %q, src: %q", d.pattern, d.text)
		module(t, "text").call("re_compile", mockRuntime{},  d.pattern).call("match", mockRuntime{}, d.text).expect(expected, "patter: %q, src: %q", d.pattern, d.text)
	}

	// re_find(pattern, text)
	for _, d := range []struct {
		pattern  string
		text     string
		expected interface{}
	}{
		{"a(b)", "", objects.UndefinedValue},
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
		module(t, "text").call("re_find", mockRuntime{},  d.pattern, d.text).expect(d.expected, "pattern: %q, text: %q", d.pattern, d.text)
		module(t, "text").call("re_compile", mockRuntime{},  d.pattern).call("find", mockRuntime{}, d.text).expect(d.expected, "pattern: %q, text: %q", d.pattern, d.text)
	}

	// re_find(pattern, text, count))
	for _, d := range []struct {
		pattern  string
		text     string
		count    int
		expected interface{}
	}{
		{"a(b)", "", -1, objects.UndefinedValue},
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
		{"(a)b(c)d", "abcdefgabcd", 0, objects.UndefinedValue},
		{"(a)b(c)d", "abcdefgabcd", 1, ARR{
			ARR{
				IMAP{"text": "abcd", "begin": 0, "end": 4},
				IMAP{"text": "a", "begin": 0, "end": 1},
				IMAP{"text": "c", "begin": 2, "end": 3},
			},
		}},
	} {
		module(t, "text").call("re_find", mockRuntime{},  d.pattern, d.text, d.count).expect(d.expected, "pattern: %q, text: %q", d.pattern, d.text)
		module(t, "text").call("re_compile", mockRuntime{},  d.pattern).call("find", mockRuntime{}, d.text, d.count).expect(d.expected, "pattern: %q, text: %q", d.pattern, d.text)
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
		expected := regexp.MustCompile(d.pattern).ReplaceAllString(d.text, d.repl)
		module(t, "text").call("re_replace", mockRuntime{},  d.pattern, d.text, d.repl).expect(expected, "pattern: %q, text: %q, repl: %q", d.pattern, d.text, d.repl)
		module(t, "text").call("re_compile", mockRuntime{},  d.pattern).call("replace", mockRuntime{}, d.text, d.repl).expect(expected, "pattern: %q, text: %q, repl: %q", d.pattern, d.text, d.repl)
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
		module(t, "text").call("re_split", mockRuntime{},  d.pattern, d.text).expect(expected, "pattern: %q, text: %q", d.pattern, d.text)
		module(t, "text").call("re_compile", mockRuntime{},  d.pattern).call("split", mockRuntime{}, d.text).expect(expected, "pattern: %q, text: %q", d.pattern, d.text)
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
		module(t, "text").call("re_split", mockRuntime{},  d.pattern, d.text, d.count).expect(expected, "pattern: %q, text: %q", d.pattern, d.text)
		module(t, "text").call("re_compile", mockRuntime{},  d.pattern).call("split", mockRuntime{}, d.text, d.count).expect(expected, "pattern: %q, text: %q", d.pattern, d.text)
	}
}

func TestText(t *testing.T) {
	module(t, "text").call("compare", mockRuntime{},  "", "").expect(0)
	module(t, "text").call("compare", mockRuntime{},  "", "a").expect(-1)
	module(t, "text").call("compare", mockRuntime{},  "a", "").expect(1)
	module(t, "text").call("compare", mockRuntime{},  "a", "a").expect(0)
	module(t, "text").call("compare", mockRuntime{},  "a", "b").expect(-1)
	module(t, "text").call("compare", mockRuntime{},  "b", "a").expect(1)
	module(t, "text").call("compare", mockRuntime{},  "abcde", "abcde").expect(0)
	module(t, "text").call("compare", mockRuntime{},  "abcde", "abcdf").expect(-1)
	module(t, "text").call("compare", mockRuntime{},  "abcdf", "abcde").expect(1)

	module(t, "text").call("contains", mockRuntime{},  "", "").expect(true)
	module(t, "text").call("contains", mockRuntime{},  "", "a").expect(false)
	module(t, "text").call("contains", mockRuntime{},  "a", "").expect(true)
	module(t, "text").call("contains", mockRuntime{},  "a", "a").expect(true)
	module(t, "text").call("contains", mockRuntime{},  "abcde", "a").expect(true)
	module(t, "text").call("contains", mockRuntime{},  "abcde", "abcde").expect(true)
	module(t, "text").call("contains", mockRuntime{},  "abc", "abcde").expect(false)
	module(t, "text").call("contains", mockRuntime{},  "ab cd", "bc").expect(false)

	module(t, "text").call("replace", mockRuntime{},  "", "", "", -1).expect("")
	module(t, "text").call("replace", mockRuntime{},  "abcd", "a", "x", -1).expect("xbcd")
	module(t, "text").call("replace", mockRuntime{},  "aaaa", "a", "x", -1).expect("xxxx")
	module(t, "text").call("replace", mockRuntime{},  "aaaa", "a", "x", 0).expect("aaaa")
	module(t, "text").call("replace", mockRuntime{},  "aaaa", "a", "x", 2).expect("xxaa")
	module(t, "text").call("replace", mockRuntime{},  "abcd", "bc", "x", -1).expect("axd")

	module(t, "text").call("format_bool", mockRuntime{},  true).expect("true")
	module(t, "text").call("format_bool", mockRuntime{},  false).expect("false")
	module(t, "text").call("format_float", mockRuntime{},  -19.84, 'f', -1, 64).expect("-19.84")
	module(t, "text").call("format_int", mockRuntime{},  -1984, 10).expect("-1984")
	module(t, "text").call("format_int", mockRuntime{},  1984, 8).expect("3700")
	module(t, "text").call("parse_bool", mockRuntime{},  "true").expect(true)
	module(t, "text").call("parse_bool", mockRuntime{},  "0").expect(false)
	module(t, "text").call("parse_float", mockRuntime{},  "-19.84", 64).expect(-19.84)
	module(t, "text").call("parse_int", mockRuntime{},  "-1984", 10, 64).expect(-1984)
}

func TestReplaceLimit(t *testing.T) {
	curMaxStringLen := tengo.MaxStringLen
	defer func() { tengo.MaxStringLen = curMaxStringLen }()
	tengo.MaxStringLen = 12

	module(t, "text").call("replace", mockRuntime{},  "123456789012", "1", "x", -1).expect("x234567890x2")
	module(t, "text").call("replace", mockRuntime{},  "123456789012", "12", "x", -1).expect("x34567890x")
	module(t, "text").call("replace", mockRuntime{},  "123456789012", "1", "xy", -1).expectError()
	module(t, "text").call("replace", mockRuntime{},  "123456789012", "0", "xy", -1).expectError()
	module(t, "text").call("replace", mockRuntime{},  "123456789012", "012", "xyz", -1).expect("123456789xyz")
	module(t, "text").call("replace", mockRuntime{},  "123456789012", "012", "xyzz", -1).expectError()

	module(t, "text").call("re_replace", mockRuntime{},  "1", "123456789012", "x").expect("x234567890x2")
	module(t, "text").call("re_replace", mockRuntime{},  "12", "123456789012", "x").expect("x34567890x")
	module(t, "text").call("re_replace", mockRuntime{},  "1", "123456789012", "xy").expectError()
	module(t, "text").call("re_replace", mockRuntime{},  "1(2)", "123456789012", "x$1").expect("x234567890x2")
	module(t, "text").call("re_replace", mockRuntime{},  "(1)(2)", "123456789012", "$2$1").expect("213456789021")
	module(t, "text").call("re_replace", mockRuntime{},  "(1)(2)", "123456789012", "${2}${1}x").expectError()
}

func TestTextRepeat(t *testing.T) {
	curMaxStringLen := tengo.MaxStringLen
	defer func() { tengo.MaxStringLen = curMaxStringLen }()
	tengo.MaxStringLen = 12

	module(t, "text").call("repeat", mockRuntime{},  "1234", "3").expect("123412341234")
	module(t, "text").call("repeat", mockRuntime{},  "1234", "4").expectError()
	module(t, "text").call("repeat", mockRuntime{},  "1", "12").expect("111111111111")
	module(t, "text").call("repeat", mockRuntime{},  "1", "13").expectError()
}

func TestSubstr(t *testing.T) {
	module(t, "text").call("substr", mockRuntime{},  "", 0, 0).expect("")
	module(t, "text").call("substr", mockRuntime{},  "abcdef", 0, 3).expect("abc")
	module(t, "text").call("substr", mockRuntime{},  "abcdef", 0, 6).expect("abcdef")
	module(t, "text").call("substr", mockRuntime{},  "abcdef", 0, 10).expect("abcdef")
	module(t, "text").call("substr", mockRuntime{},  "abcdef", -10, 10).expect("abcdef")
	module(t, "text").call("substr", mockRuntime{},  "abcdef", 0).expect("abcdef")
	module(t, "text").call("substr", mockRuntime{},  "abcdef", 3).expect("def")

	module(t, "text").call("substr", mockRuntime{},  "", 10, 0).expectError()
	module(t, "text").call("substr", mockRuntime{},  "", "10", 0).expectError()
	module(t, "text").call("substr", mockRuntime{},  "", 10, "0").expectError()
	module(t, "text").call("substr", mockRuntime{},  "", "10", "0").expectError()

	module(t, "text").call("substr", mockRuntime{},  0, 0, 1).expect("0")
	module(t, "text").call("substr", mockRuntime{},  123, 0, 1).expect("1")
	module(t, "text").call("substr", mockRuntime{},  123.456, 4, 7).expect("456")
}

func TestPadLeft(t *testing.T) {
	module(t, "text").call("pad_left", mockRuntime{},  "ab", 7, 0).expect("00000ab")
	module(t, "text").call("pad_right", mockRuntime{},  "ab", 7, 0).expect("ab00000")
	module(t, "text").call("pad_left", mockRuntime{},  "ab", 7, "+-").expect("-+-+-ab")
	module(t, "text").call("pad_right", mockRuntime{},  "ab", 7, "+-").expect("ab+-+-+")
}

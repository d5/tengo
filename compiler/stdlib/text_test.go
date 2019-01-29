package stdlib_test

import (
	"testing"

	"github.com/d5/tengo/objects"
)

func TestTextRE(t *testing.T) {
	// re_match(pattern, text)
	for _, d := range []struct {
		pattern  string
		text     string
		expected interface{}
	}{
		{"abc", "", false},
		{"abc", "abc", true},
		{"a", "abc", true},
		{"b", "abc", true},
		{"^a", "abc", true},
		{"^b", "abc", false},
	} {
		module(t, "text").call("re_match", d.pattern, d.text).expect(d.expected, "pattern: %q, src: %q", d.pattern, d.text)
		module(t, "text").call("re_compile", d.pattern).call("match", d.text).expect(d.expected, "patter: %q, src: %q", d.pattern, d.text)
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
		module(t, "text").call("re_find", d.pattern, d.text).expect(d.expected, "pattern: %q, text: %q", d.pattern, d.text)
		module(t, "text").call("re_compile", d.pattern).call("find", d.text).expect(d.expected, "pattern: %q, text: %q", d.pattern, d.text)
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
		module(t, "text").call("re_find", d.pattern, d.text, d.count).expect(d.expected, "pattern: %q, text: %q", d.pattern, d.text)
		module(t, "text").call("re_compile", d.pattern).call("find", d.text, d.count).expect(d.expected, "pattern: %q, text: %q", d.pattern, d.text)
	}

	// re_replace(pattern, text, repl)
	for _, d := range []struct {
		pattern  string
		text     string
		repl     string
		expected interface{}
	}{
		{"a", "", "b", ""},
		{"a", "a", "b", "b"},
		{"a", "acac", "b", "bcbc"},
		{"a", "acac", "123", "123c123c"},
		{"ac", "acac", "99", "9999"},
		{"ac$", "acac", "foo", "acfoo"},
	} {
		module(t, "text").call("re_replace", d.pattern, d.text, d.repl).expect(d.expected, "pattern: %q, text: %q, repl: %q", d.pattern, d.text, d.repl)
		module(t, "text").call("re_compile", d.pattern).call("replace", d.text, d.repl).expect(d.expected, "pattern: %q, text: %q, repl: %q", d.pattern, d.text, d.repl)
	}

	// re_split(pattern, text)
	for _, d := range []struct {
		pattern  string
		text     string
		expected interface{}
	}{
		{"a", "", ARR{""}},
		{"a", "abcabc", ARR{"", "bc", "bc"}},
		{"ab", "abcabc", ARR{"", "c", "c"}},
		{"^a", "abcabc", ARR{"", "bcabc"}},
	} {
		module(t, "text").call("re_split", d.pattern, d.text).expect(d.expected, "pattern: %q, text: %q", d.pattern, d.text)
		module(t, "text").call("re_compile", d.pattern).call("split", d.text).expect(d.expected, "pattern: %q, text: %q", d.pattern, d.text)
	}

	// re_split(pattern, text, count))
	for _, d := range []struct {
		pattern  string
		text     string
		count    int
		expected interface{}
	}{
		{"a", "", -1, ARR{""}},
		{"a", "abcabc", -1, ARR{"", "bc", "bc"}},
		{"ab", "abcabc", -1, ARR{"", "c", "c"}},
		{"^a", "abcabc", -1, ARR{"", "bcabc"}},
		{"a", "abcabc", 0, ARR{}},
		{"a", "abcabc", 1, ARR{"abcabc"}},
		{"a", "abcabc", 2, ARR{"", "bcabc"}},
		{"a", "abcabc", 3, ARR{"", "bc", "bc"}},
		{"b", "abcabc", 1, ARR{"abcabc"}},
		{"b", "abcabc", 2, ARR{"a", "cabc"}},
		{"b", "abcabc", 3, ARR{"a", "ca", "c"}},
	} {
		module(t, "text").call("re_split", d.pattern, d.text, d.count).expect(d.expected, "pattern: %q, text: %q", d.pattern, d.text)
		module(t, "text").call("re_compile", d.pattern).call("split", d.text, d.count).expect(d.expected, "pattern: %q, text: %q", d.pattern, d.text)
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
}

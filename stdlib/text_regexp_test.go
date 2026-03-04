package stdlib_test

import (
	"testing"

	"github.com/d5/tengo/v2"
)

func TestTextREAlternation(t *testing.T) {
	module(t, "text").call("re_find", "([a-zA-Z])|([0-9])", "a").expect(ARR{
		ARR{
			IMAP{"text": "a", "begin": 0, "end": 1},
			IMAP{"text": "a", "begin": 0, "end": 1},
		},
	}, "alternation with letter")

	module(t, "text").call("re_find", "([a-zA-Z])|([0-9])", "5").expect(ARR{
		ARR{
			IMAP{"text": "5", "begin": 0, "end": 1},
			IMAP{"text": "5", "begin": 0, "end": 1},
		},
	}, "alternation with number")

	module(t, "text").call("re_find", "([a-zA-Z])|([0-9])", "").expect(tengo.UndefinedValue, "empty input")

	module(t, "text").call("re_find", "([a-zA-Z])|([0-9])", "!").expect(tengo.UndefinedValue, "non-matching input")

	module(t, "text").call("re_find", "(?:([a-zA-Z])|([0-9]))+", "a5b").expect(ARR{
		ARR{
			IMAP{"text": "a5b", "begin": 0, "end": 3},
			IMAP{"text": "b", "begin": 2, "end": 3},
			IMAP{"text": "5", "begin": 1, "end": 2},
		},
	}, "multiple alternations")

	module(t, "text").call("re_find", "(foo)|(bar)|(baz)", "foo").expect(ARR{
		ARR{
			IMAP{"text": "foo", "begin": 0, "end": 3},
			IMAP{"text": "foo", "begin": 0, "end": 3},
		},
	}, "multiple groups with non-matches")

	module(t, "text").call("re_find", "((cat)|(dog))((run)|(walk))", "catrun").expect(ARR{
		ARR{
			IMAP{"text": "catrun", "begin": 0, "end": 6},
			IMAP{"text": "cat", "begin": 0, "end": 3},
			IMAP{"text": "cat", "begin": 0, "end": 3},
			IMAP{"text": "run", "begin": 3, "end": 6},
			IMAP{"text": "run", "begin": 3, "end": 6},
		},
	}, "nested groups with alternation")
}

package stdlib_test

import "testing"

func TestFmtSprintf(t *testing.T) {
	module(t, `fmt`).call("sprintf", mockRuntime{}, "").expect("")
	module(t, `fmt`).call("sprintf", mockRuntime{}, "foo").expect("foo")
	module(t, `fmt`).call("sprintf", mockRuntime{}, `foo %d %v %s`, 1, 2, "bar").expect("foo 1 2 bar")
	module(t, `fmt`).call("sprintf", mockRuntime{}, "foo %v", ARR{1, "bar", true}).expect(`foo [1, "bar", true]`)
	module(t, `fmt`).call("sprintf", mockRuntime{}, "foo %v %d", ARR{1, "bar", true}, 19).expect(`foo [1, "bar", true] 19`)
	module(t, `fmt`).call("sprintf", mockRuntime{}, "foo %v", MAP{"a": IMAP{"b": IMAP{"c": ARR{1, 2, 3}}}}).expect(`foo {a: {b: {c: [1, 2, 3]}}}`)
	module(t, `fmt`).call("sprintf", mockRuntime{}, "%v", IARR{1, IARR{2, IARR{3, 4}}}).expect(`[1, [2, [3, 4]]]`)
}

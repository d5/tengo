package stdlib_test

import "testing"

func TestFmtSprintf(t *testing.T) {
	module(t, `fmt`).call("sprintf", "").expect("")
	module(t, `fmt`).call("sprintf", "foo").expect("foo")
	module(t, `fmt`).call("sprintf", `foo %d %v %s`, 1, 2, "bar").expect("foo 1 2 bar")
	module(t, `fmt`).call("sprintf", "foo %v", `[1, "bar", true]`).expect(`foo [1, "bar", true]`)
	module(t, `fmt`).call("sprintf", "foo %v %d", `[1, "bar", true]`, 19).expect(`foo [1, "bar", true] 19`)
	module(t, `fmt`).call("sprintf", "foo %v", `{a: {b: {c: [1, 2, 3]}}}`).expect(`foo {a: {b: {c: [1, 2, 3]}}}`)
}

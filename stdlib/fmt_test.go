package stdlib_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/require"
)

func TestFmtSprintf(t *testing.T) {
	module(t, `fmt`).call("sprintf", "").expect("")
	module(t, `fmt`).call("sprintf", "foo").expect("foo")
	module(t, `fmt`).call("sprintf", `foo %d %v %s`, 1, 2, "bar").
		expect("foo 1 2 bar")
	module(t, `fmt`).call("sprintf", "foo %v", ARR{1, "bar", true}).
		expect(`foo [1, "bar", true]`)
	module(t, `fmt`).call("sprintf", "foo %v %d", ARR{1, "bar", true}, 19).
		expect(`foo [1, "bar", true] 19`)
	module(t, `fmt`).
		call("sprintf", "foo %v", MAP{"a": IMAP{"b": IMAP{"c": ARR{1, 2, 3}}}}).
		expect(`foo {a: {b: {c: [1, 2, 3]}}}`)
	module(t, `fmt`).call("sprintf", "%v", IARR{1, IARR{2, IARR{3, 4}}}).
		expect(`[1, [2, [3, 4]]]`)
}

func captureStdout(t *testing.T, fn func()) string {
	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer func() { _ = r.Close() }()

	stdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = stdout }()
	fn()
	os.Stdout = stdout

	require.NoError(t, w.Close())
	out, err := ioutil.ReadAll(r)
	require.NoError(t, err)
	return string(out)
}

func TestFmtPrint(t *testing.T) {
	var res callres
	out := captureStdout(t, func() {
		res = module(t, "fmt").call("print", "a", 1, "b", true)
	})
	require.NoError(t, res.e)
	require.Equal(t, "a1btrue", out)

	out = captureStdout(t, func() {
		res = module(t, "fmt").call("println", "a", 1)
	})
	require.NoError(t, res.e)
	require.Equal(t, "a1\n", out)

	out = captureStdout(t, func() {
		res = module(t, "fmt").call("println")
	})
	require.NoError(t, res.e)
	require.Equal(t, "\n", out)

	out = captureStdout(t, func() {
		res = module(t, "fmt").call("printf", "%d-%s\n", 1, "x")
	})
	require.NoError(t, res.e)
	require.Equal(t, "1-x\n", out)
}

func TestFmtPrintErrors(t *testing.T) {
	cases := []argErrorCase{
		{"printf", ARR{}, tengo.ErrWrongNumArguments},
		{"printf", ARR{1, 2}, invalidArg("format", "string", "int")},
		{"sprintf", ARR{}, tengo.ErrWrongNumArguments},
		{"sprintf", ARR{1, 2}, invalidArg("format", "string", "int")},
	}
	out := captureStdout(t, func() { expectArgErrors(t, "fmt", cases) })
	require.Equal(t, "", out)

	curMaxStringLen := tengo.MaxStringLen
	defer func() { tengo.MaxStringLen = curMaxStringLen }()
	tengo.MaxStringLen = 5

	var res callres
	for fn, want := range map[string]string{
		"print":   "abcde",
		"println": "abcde\n",
	} {
		out = captureStdout(t, func() {
			res = module(t, "fmt").call(fn, "abc", "de")
		})
		require.NoError(t, res.e, fn)
		require.Equal(t, want, out, fn)

		out = captureStdout(t, func() {
			res = module(t, "fmt").call(fn, "abc", "def")
		})
		require.Equal(t, tengo.ErrStringLimit, res.e, fn)
		require.Equal(t, "", out, fn)
	}
}

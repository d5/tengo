package runtime_test

import (
	"testing"
)

func TestError(t *testing.T) {
	expect(t, `out = error(1)`, nil, errorObject(1))
	expect(t, `out = error(1).value`, nil, 1)
	expect(t, `out = error("some error")`, nil, errorObject("some error"))
	expect(t, `out = error("some" + " error")`, nil, errorObject("some error"))
	expect(t, `out = func() { return error(5) }()`, nil, errorObject(5))
	expect(t, `out = error(error("foo"))`, nil, errorObject(errorObject("foo")))
	expect(t, `out = error("some error")`, nil, errorObject("some error"))
	expect(t, `out = error("some error").value`, nil, "some error")
	expect(t, `out = error("some error")["value"]`, nil, "some error")

	expectError(t, `error("error").err`, nil, "invalid index on error")
	expectError(t, `error("error").value_`, nil, "invalid index on error")
	expectError(t, `error([1,2,3])[1]`, nil, "invalid index on error")
}

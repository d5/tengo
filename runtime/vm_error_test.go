package runtime_test

import (
	"testing"
)

func TestError(t *testing.T) {
	expect(t, `out = error(1)`, errorObject(1))
	expect(t, `out = error(1).value`, 1)
	expect(t, `out = error("some error")`, errorObject("some error"))
	expect(t, `out = error("some" + " error")`, errorObject("some error"))
	expect(t, `out = func() { return error(5) }()`, errorObject(5))
	expect(t, `out = error(error("foo"))`, errorObject(errorObject("foo")))
	expect(t, `out = error("some error")`, errorObject("some error"))
	expect(t, `out = error("some error").value`, "some error")
	expect(t, `out = error("some error")["value"]`, "some error")

	expectError(t, `error("error").err`)
	expectError(t, `error("error").value_`)
	expectError(t, `error([1,2,3])[1]`)
}

package runtime_test

import (
	"testing"

	"github.com/d5/tengo"
)

func TestTryExpr(t *testing.T) {
	// top-level try failure considered a runtime error
	expectError(t, `try(error("oops"))`, nil, "oops")
	// if try succeeds, value is passed through
	expect(t, `out = try(10)`, nil, 10)
	// also works with undefined values
	expect(t, `out = try(undefined)`, nil, tengo.UndefinedValue)

	oopsErr := &tengo.Error{Value: &tengo.String{Value: "oops"}}

	expect(t, `
out = func() {
	try(error("oops"))
	return 10;
}()`, nil, oopsErr)

	expect(t, `
out = func() {
	if is_error(func() {
		x := try(error("oops"))
		return x + 15;
	}()) {
		return 10;
	}
	return 11;
}()`, nil, 10)

	expect(t, `
f1 := func() {
	return error("oops");
};

f2 := func() {
	try(f1());
	return 16;
}

out = f2()
`, nil, oopsErr)

	expect(t, `
out = func() {
	return try(12);
}()`, nil, 12)

	expect(t, `
out = func() {
	return try(error("oops"));
}()`, nil, oopsErr)

	expect(t, `
x := func() {
	return try(error("oops"));
}();

if is_error(x) {
	out = 12;
} else {
	out = 0;	
}`, nil, 12)
}

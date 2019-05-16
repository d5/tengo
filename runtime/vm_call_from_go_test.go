package runtime_test

import (
	"testing"

	"github.com/d5/tengo/objects"
)

var interopOpts = &testopts{
	skip2ndPass: true,
	maxAllocs:   999999999,
	symbols: map[string]objects.Object{
		"call_from_go": &objects.UserFunction{
			Name: "call_from_go",
			Value: func(rt objects.Runtime, args ...objects.Object) (objects.Object, error) {
				if len(args) < 1 {
					return nil, objects.ErrWrongNumArguments
				}

				return rt.Call(args[0], args[1:]...)
			},
		},

		"call_from_go_noerr": &objects.UserFunction{
			Name: "call_from_go_noerr",
			Value: func(rt objects.Runtime, args ...objects.Object) (objects.Object, error) {
				if len(args) < 1 {
					return nil, objects.ErrWrongNumArguments
				}

				ret, _ := rt.Call(args[0], args[1:]...)
				if ret == nil {
					ret = objects.UndefinedValue
				}
				return ret, nil
			},
		},
	},
}

func interopExpect(t *testing.T, code string, want interface{}) {
	expect(t, code, interopOpts, want)
}

func interopExpectError(t *testing.T, code string, want string) {
	expectError(t, code, interopOpts, want)
}

func TestCallFromGo(t *testing.T) {
	interopExpect(t, `
		fn := func(x) {
			return x;
		}

		out = [
			call_from_go(fn, 101),
			call_from_go(fn, 202),
			call_from_go(fn, 303),
			call_from_go(fn, 404),
			call_from_go(fn, 505)
		];
	`, ARR{101, 202, 303, 404, 505})

	interopExpect(t, `
		fn := func(a, b, ...c) {
			return [a, b, c];
		}

		out = call_from_go(fn, 1,2,3,4,5);
	`, ARR{1, 2, ARR{3, 4, 5}})

	interopExpect(t, `
		obj := {}
		fn := func(k, v) {
			obj[k] = v
		}

		call_from_go(fn, "a", 1);
		call_from_go(fn, "b", 2);
		call_from_go(fn, "c", 3);
		out = [obj.a, obj.b, obj.c];
	`, ARR{1, 2, 3})

	interopExpect(t, `
		indirect := func(fn, ...arg) {
			return call_from_go(fn, arg[0]);
		}

		fn := func(x) {
			return x;
		}

		out = call_from_go(indirect, fn, 12)
	`, 12)

	interopExpect(t, `
		indirect := func(fn, ...arg) {
			return call_from_go(fn, arg[0], arg[1]);
		}
		obj := {}
		fn := func(k, v) {
			obj[k] = v
		}

		call_from_go(indirect, fn, "a",1)
		call_from_go(indirect, fn, "b",2)
		call_from_go(indirect, fn, "c",3)
		out = [obj.a, obj.b, obj.c]
	`, ARR{1, 2, 3})

	interopExpect(t, `
		indirect := func(fn, ...arg) {
			return call_from_go(fn, arg[0], arg[1]);
		}

		fn := func(x) {
			return x;
		}

		out = call_from_go(call_from_go, indirect, call_from_go, fn, 12)
	`, 12)

	interopExpect(t, `
		fib := func(x, s) {
			if x == 0 {
				return 0 + s
			} else if x == 1 {
				return 1 + s
			}
		
			return fib(x-1, fib(x-2, s))
		}

		out = call_from_go(fib, 25, 0)
	`, 75025)

	interopExpect(t, `
		fib := func(x, s) {
			if x == 0 {
				return 0 + s
			} else if x == 1 {
				return 1 + s
			}
		
			return fib(x-1, fib(x-2, s))
		}
		
		do_fib := func(nth) {
			return call_from_go(fib, nth, 0)
		}

		out = call_from_go(do_fib, 25)
	`, 75025)

	interopExpect(t, `
		out = call_from_go(format, "%d %d %d", 1, 2, 3)
	`, "1 2 3")

	interopExpect(t, `
		fn := func() {
			return 5/[]
		}

		out = call_from_go_noerr(fn)
	`, objects.UndefinedValue)

	interopExpectError(t, `
		fn := func(a, b, c) { return [a,b,c]; }

		out = call_from_go(fn, 1, 2)
	`, "Runtime Error: wrong number of arguments")
}

package runtime_test

import (
	"fmt"
	"testing"

	"github.com/d5/tengo/compiler/token"
	"github.com/d5/tengo/objects"
)

func TestInterop(t *testing.T) {
	// interop functions
	//
	//  'invoke'
	//    - invokes Callable object and propagate a runtime error
	//
	//  'invoke_ignore_err'
	//    - invokes Callable object and swallow a runtime error
	//
	//  'invoke_ignore_panic'
	//    - invokes Callable object and swallow panic
	//
	// other utility functions
	//
	//  'bind'
	//    - a binding function
	//
	//  'identity'
	//    - return the first argument as-is
	//
	//  'sum'
	//    - return args[0] + args[1]
	//
	//  'rt_err'
	//    - return a run-time error
	//
	//  'throw'
	//    - cause a panic
	//
	opts := Opts().Skip2ndPass().
		Symbol("invoke", &objects.UserFunction{
			Name: "invoke",
			Value: func(rt objects.Interop, args ...objects.Object) (objects.Object, error) {
				if len(args) < 1 {
					return nil, objects.ErrWrongNumArguments
				}

				return rt.InteropCall(args[0], args[1:]...)
			},
		}).
		Symbol("invoke_ignore_err", &objects.UserFunction{
			Name: "invoke_ignore_err",
			Value: func(rt objects.Interop, args ...objects.Object) (objects.Object, error) {
				if len(args) < 1 {
					return nil, objects.ErrWrongNumArguments
				}

				ret, _ := rt.InteropCall(args[0], args[1:]...)
				if ret == nil {
					ret = objects.UndefinedValue
				}
				return ret, nil
			},
		}).
		Symbol("invoke_ignore_panic", &objects.UserFunction{
			Name: "invoke_ignore_panic",
			Value: func(rt objects.Interop, args ...objects.Object) (objects.Object, error) {
				if len(args) < 1 {
					return nil, objects.ErrWrongNumArguments
				}

				defer func() {
					_ = recover() // swallow panic
				}()
				return rt.InteropCall(args[0], args[1:]...)
			},
		}).
		Symbol("bind", &objects.UserFunction{
			Name: "bind",
			Value: func(rt objects.Interop, args ...objects.Object) (objects.Object, error) {
				if len(args) < 1 {
					return nil, objects.ErrWrongNumArguments
				}
				fn := args[0]
				boundArgs := args[1:]
				return &objects.UserFunction{
					Value: func(rt objects.Interop, args ...objects.Object) (ret objects.Object, err error) {
						return rt.InteropCall(fn, append(boundArgs, args...)...)
					},
				}, nil
			},
		}).
		Symbol("identity", &objects.UserFunction{
			Name: "identity",
			Value: func(rt objects.Interop, args ...objects.Object) (objects.Object, error) {
				if len(args) < 1 {
					return nil, objects.ErrWrongNumArguments
				}
				return args[0], nil
			},
		}).
		Symbol("sum", &objects.UserFunction{
			Name: "sum",
			Value: func(rt objects.Interop, args ...objects.Object) (objects.Object, error) {
				if len(args) != 2 {
					return nil, objects.ErrWrongNumArguments
				}
				return args[0].BinaryOp(token.Add, args[1])
			},
		}).
		Symbol("rt_err", &objects.UserFunction{
			Name: "rt_err",
			Value: func(rt objects.Interop, args ...objects.Object) (objects.Object, error) {
				return nil, fmt.Errorf("rt_err: %s", args[0].String())
			},
		}).
		Symbol("throw", &objects.UserFunction{
			Name: "throw",
			Value: func(rt objects.Interop, args ...objects.Object) (objects.Object, error) {
				panic(fmt.Errorf("throw: %s", args[0].String()))
			},
		})

	// simple interop
	expect(t, `
out = invoke(identity, 10)
`, opts, 10)
	expect(t, `
sum4 := bind(sum, 4)
out = invoke(sum4, 5)
`, opts, 9)

	// runtime error propagated
	expectError(t, `
invoke(rt_err, 3)
`, opts, "rt_err: 3")

	// propagated runtime error should halt the runtime
	expectError(t, `
invoke(rt_err, 4)	// runtime error here
10 / 0				// this line must not be executed (integer divide by zero)
`, opts, "rt_err: 4")

	// runtime error swallowed
	expect(t, `
out = invoke_ignore_err(rt_err, 5)
`, opts, objects.UndefinedValue)

	// swallowed runtime error should not halt the runtime
	expect(t, `
invoke_ignore_err(rt_err, 6)
out = 7
`, opts, 7)

	// panic propagated
	expectPanic(t, `
invoke(throw, 8)
`, opts, "throw: 8")

	// propagated panic should halt the runtime
	expectPanic(t, `
invoke(throw, 9)	// panic here
10 / 0				// this line must not be executed (integer divide by zero)
`, opts, "throw: 9")

	// panic swallowed but runtime should still throw the same panic
	expectPanic(t, `
out = invoke_ignore_panic(throw, 10)
`, opts, "throw: 10")

	// panic swallowed but runtime should still throw the same panic
	expectPanic(t, `
invoke_ignore_panic(throw, 11)
out = 12
`, opts, "throw: 11")
}

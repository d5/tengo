package runtime_test

import (
	"testing"

	"github.com/d5/tengo/objects"
)

func TestGoInterop(t *testing.T) {
	opts := &testopts{
		symbols: map[string]objects.Object{
			"go_fn": &objects.UserFunction{
				Name: "go_fn",
				Value: func(hooks objects.RuntimeHooks, args ...objects.Object) (objects.Object, error) {
					fn := args[0]
					rest := args[1:]
					return hooks.Call(fn, rest...)
				},
			},
			"go_fn_apply": &objects.UserFunction{
				Name: "go_fn_apply",
				Value: func(hooks objects.RuntimeHooks, args ...objects.Object) (objects.Object, error) {
					fn := args[0]
					rest := args[1].(*objects.Array).Value
					return hooks.Call(fn, rest...)
				},
			},
		},
		maxAllocs:   2048,
		skip2ndPass: true,
	}

	const code = `
	out = {}

	tengo_fn_1 := func(a, b, c) {
		return [a, b, c]
	}

	tengo_fn_2 := func(a, b, ...c) {
		return [a, b, c]
	}

	tengo_fn_3 := func(fn, ...args) {
		go_fn_apply(fn, args)
	}

	out.with_tengo_fn_1 = go_fn(tengo_fn_1, 1, 2, 3)
	out.with_tengo_fn_2 = go_fn(tengo_fn_2, 1, 2, 3, 4)
	out.with_tengo_fn_3 = go_fn(tengo_fn_2, 1, 2, 3, 4)
	out.with_callable   = go_fn(go_fn_apply, tengo_fn_1, [1,2,3])
	`

	expect(t, code, opts, MAP{
		"with_tengo_fn_1": ARR{1, 2, 3},
		"with_tengo_fn_2": ARR{1, 2, ARR{3, 4}},
		"with_tengo_fn_3": ARR{1, 2, ARR{3, 4}},
		"with_callable":   ARR{1, 2, 3},
	})
}

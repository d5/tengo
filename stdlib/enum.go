package stdlib

import (
	"github.com/d5/tengo"
)

var enumModule = map[string]tengo.Object{
	"all":      &tengo.UserFunction{Name: "all", Value: enumAll},
	"any":      &tengo.UserFunction{Name: "any", Value: enumAny},
	"chunk":    &tengo.UserFunction{Name: "chunk", Value: enumChunk},
	"at":       &tengo.UserFunction{Name: "at", Value: enumAt},
	"each":     &tengo.UserFunction{Name: "each", Value: enumEach},
	"filter":   &tengo.UserFunction{Name: "filter", Value: enumFilter},
	"find":     &tengo.UserFunction{Name: "find", Value: enumFind},
	"find_key": &tengo.UserFunction{Name: "find_key", Value: enumFindKey},
	"map":      &tengo.UserFunction{Name: "map", Value: enumMap},
	"key":      &tengo.UserFunction{Name: "key", Value: enumKey},
	"value":    &tengo.UserFunction{Name: "value", Value: enumValue},
}

func enumAll(rt tengo.Interop, args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	if !args[0].CanIterate() {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "x",
			Expected: "iterable",
			Found:    args[0].TypeName(),
		}
	}
	if !args[1].CanCall() {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "fn",
			Expected: "callable",
			Found:    args[1].TypeName(),
		}
	}

	it := args[0].Iterate()
	for it.Next() {
		k, v := it.Key(), it.Value()
		res, err := rt.InteropCall(args[1], k, v)
		if err != nil {
			return nil, err
		}
		if res.IsFalsy() {
			return tengo.FalseValue, nil
		}
	}
	return tengo.TrueValue, nil
}

func enumAny(rt tengo.Interop, args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	if !args[0].CanIterate() {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "x",
			Expected: "iterable",
			Found:    args[0].TypeName(),
		}
	}
	if !args[1].CanCall() {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "fn",
			Expected: "callable",
			Found:    args[1].TypeName(),
		}
	}

	it := args[0].Iterate()
	for it.Next() {
		k, v := it.Key(), it.Value()
		res, err := rt.InteropCall(args[1], k, v)
		if err != nil {
			return nil, err
		}
		if !res.IsFalsy() {
			return tengo.TrueValue, nil
		}
	}
	return tengo.FalseValue, nil
}

func enumChunk(_ tengo.Interop, args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	var arr []tengo.Object
	switch v := args[0].(type) {
	case *tengo.Array:
		arr = v.Value
	case *tengo.ImmutableArray:
		arr = v.Value
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "x",
			Expected: "array",
			Found:    args[0].TypeName(),
		}
	}

	size, ok := tengo.ToInt(args[1])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "size",
			Expected: "int(compatible)",
			Found:    args[1].TypeName(),
		}
	}
	if size == 0 {
		return tengo.UndefinedValue, nil
	}

	var res []tengo.Object
	numArr := len(arr)
	for idx := 0; idx < numArr; idx += size {
		res = append(res, &tengo.Array{Value: arr[idx:minI(idx+size, numArr)]})
	}
	return &tengo.Array{Value: res}, nil
}

func enumAt(_ tengo.Interop, args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	return args[0].IndexGet(args[1])
}

func enumEach(rt tengo.Interop, args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	if !args[0].CanIterate() {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "x",
			Expected: "iterable",
			Found:    args[0].TypeName(),
		}
	}
	if !args[1].CanCall() {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "fn",
			Expected: "callable",
			Found:    args[1].TypeName(),
		}
	}

	it := args[0].Iterate()
	for it.Next() {
		k, v := it.Key(), it.Value()
		_, err := rt.InteropCall(args[1], k, v)
		if err != nil {
			return nil, err
		}
	}
	return tengo.UndefinedValue, nil
}

func enumFilter(rt tengo.Interop, args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	if !args[0].CanIterate() {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "x",
			Expected: "iterable",
			Found:    args[0].TypeName(),
		}
	}
	if !args[1].CanCall() {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "fn",
			Expected: "callable",
			Found:    args[1].TypeName(),
		}
	}

	var res []tengo.Object
	it := args[0].Iterate()
	for it.Next() {
		k, v := it.Key(), it.Value()
		b, err := rt.InteropCall(args[1], k, v)
		if err != nil {
			return nil, err
		}
		if !b.IsFalsy() {
			res = append(res, v)
		}
	}
	return &tengo.Array{Value: res}, nil
}

func enumFind(rt tengo.Interop, args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	if !args[0].CanIterate() {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "x",
			Expected: "iterable",
			Found:    args[0].TypeName(),
		}
	}
	if !args[1].CanCall() {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "fn",
			Expected: "callable",
			Found:    args[1].TypeName(),
		}
	}

	it := args[0].Iterate()
	for it.Next() {
		k, v := it.Key(), it.Value()
		b, err := rt.InteropCall(args[1], k, v)
		if err != nil {
			return nil, err
		}
		if !b.IsFalsy() {
			return v, nil
		}
	}
	return tengo.UndefinedValue, nil
}

func enumFindKey(rt tengo.Interop, args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	if !args[0].CanIterate() {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "x",
			Expected: "iterable",
			Found:    args[0].TypeName(),
		}
	}
	if !args[1].CanCall() {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "fn",
			Expected: "callable",
			Found:    args[1].TypeName(),
		}
	}

	it := args[0].Iterate()
	for it.Next() {
		k, v := it.Key(), it.Value()
		b, err := rt.InteropCall(args[1], k, v)
		if err != nil {
			return nil, err
		}
		if !b.IsFalsy() {
			return k, nil
		}
	}
	return tengo.UndefinedValue, nil
}

func enumMap(rt tengo.Interop, args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	if !args[0].CanIterate() {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "x",
			Expected: "iterable",
			Found:    args[0].TypeName(),
		}
	}
	if !args[1].CanCall() {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "fn",
			Expected: "callable",
			Found:    args[1].TypeName(),
		}
	}

	var res []tengo.Object
	it := args[0].Iterate()
	for it.Next() {
		k, v := it.Key(), it.Value()
		r, err := rt.InteropCall(args[1], k, v)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return &tengo.Array{Value: res}, nil
}

func enumKey(_ tengo.Interop, args ...tengo.Object) (tengo.Object, error) {
	if len(args) < 1 {
		return nil, tengo.ErrWrongNumArguments
	}

	return args[0], nil
}

func enumValue(_ tengo.Interop, args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	return args[1], nil
}

func minI(a, b int) int {
	if a < b {
		return a
	}
	return b
}

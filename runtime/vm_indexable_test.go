package runtime_test

import (
	"strings"
	"testing"

	"github.com/d5/tengo"
	"github.com/d5/tengo/compiler/token"
)

type StringDict struct {
	tengo.ObjectImpl
	Value map[string]string
}

func (o *StringDict) String() string { return "" }

func (o *StringDict) TypeName() string {
	return "string-dict"
}

func (o *StringDict) IndexGet(index tengo.Object) (tengo.Object, error) {
	strIdx, ok := index.(*tengo.String)
	if !ok {
		return nil, tengo.ErrInvalidIndexType
	}

	for k, v := range o.Value {
		if strings.ToLower(strIdx.Value) == strings.ToLower(k) {
			return &tengo.String{Value: v}, nil
		}
	}

	return tengo.UndefinedValue, nil
}

func (o *StringDict) IndexSet(index, value tengo.Object) error {
	strIdx, ok := index.(*tengo.String)
	if !ok {
		return tengo.ErrInvalidIndexType
	}

	strVal, ok := tengo.ToString(value)
	if !ok {
		return tengo.ErrInvalidIndexValueType
	}

	o.Value[strings.ToLower(strIdx.Value)] = strVal

	return nil
}

type StringCircle struct {
	tengo.ObjectImpl
	Value []string
}

func (o *StringCircle) TypeName() string {
	return "string-circle"
}

func (o *StringCircle) String() string {
	return ""
}

func (o *StringCircle) IndexGet(index tengo.Object) (tengo.Object, error) {
	intIdx, ok := index.(*tengo.Int)
	if !ok {
		return nil, tengo.ErrInvalidIndexType
	}

	r := int(intIdx.Value) % len(o.Value)
	if r < 0 {
		r = len(o.Value) + r
	}

	return &tengo.String{Value: o.Value[r]}, nil
}

func (o *StringCircle) IndexSet(index, value tengo.Object) error {
	intIdx, ok := index.(*tengo.Int)
	if !ok {
		return tengo.ErrInvalidIndexType
	}

	r := int(intIdx.Value) % len(o.Value)
	if r < 0 {
		r = len(o.Value) + r
	}

	strVal, ok := tengo.ToString(value)
	if !ok {
		return tengo.ErrInvalidIndexValueType
	}

	o.Value[r] = strVal

	return nil
}

type StringArray struct {
	tengo.ObjectImpl
	Value []string
}

func (o *StringArray) String() string {
	return strings.Join(o.Value, ", ")
}

func (o *StringArray) BinaryOp(op token.Token, rhs tengo.Object) (tengo.Object, error) {
	if rhs, ok := rhs.(*StringArray); ok {
		switch op {
		case token.Add:
			if len(rhs.Value) == 0 {
				return o, nil
			}
			return &StringArray{Value: append(o.Value, rhs.Value...)}, nil
		}
	}

	return nil, tengo.ErrInvalidOperator
}

func (o *StringArray) IsFalsy() bool {
	return len(o.Value) == 0
}

func (o *StringArray) Equals(x tengo.Object) bool {
	if x, ok := x.(*StringArray); ok {
		if len(o.Value) != len(x.Value) {
			return false
		}

		for i, v := range o.Value {
			if v != x.Value[i] {
				return false
			}
		}

		return true
	}

	return false
}

func (o *StringArray) Copy() tengo.Object {
	return &StringArray{
		Value: append([]string{}, o.Value...),
	}
}

func (o *StringArray) TypeName() string {
	return "string-array"
}

func (o *StringArray) IndexGet(index tengo.Object) (tengo.Object, error) {
	intIdx, ok := index.(*tengo.Int)
	if ok {
		if intIdx.Value >= 0 && intIdx.Value < int64(len(o.Value)) {
			return &tengo.String{Value: o.Value[intIdx.Value]}, nil
		}

		return nil, tengo.ErrIndexOutOfBounds
	}

	strIdx, ok := index.(*tengo.String)
	if ok {
		for vidx, str := range o.Value {
			if strIdx.Value == str {
				return &tengo.Int{Value: int64(vidx)}, nil
			}
		}

		return tengo.UndefinedValue, nil
	}

	return nil, tengo.ErrInvalidIndexType
}

func (o *StringArray) IndexSet(index, value tengo.Object) error {
	strVal, ok := tengo.ToString(value)
	if !ok {
		return tengo.ErrInvalidIndexValueType
	}

	intIdx, ok := index.(*tengo.Int)
	if ok {
		if intIdx.Value >= 0 && intIdx.Value < int64(len(o.Value)) {
			o.Value[intIdx.Value] = strVal
			return nil
		}

		return tengo.ErrIndexOutOfBounds
	}

	return tengo.ErrInvalidIndexType
}

func (o *StringArray) Call(_ tengo.Interop, args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}

	s1, ok := tengo.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	for i, v := range o.Value {
		if v == s1 {
			return &tengo.Int{Value: int64(i)}, nil
		}
	}

	return tengo.UndefinedValue, nil
}

func (o *StringArray) CanCall() bool {
	return true
}

func TestIndexable(t *testing.T) {
	dict := func() *StringDict { return &StringDict{Value: map[string]string{"a": "foo", "b": "bar"}} }
	expect(t, `out = dict["a"]`, Opts().Symbol("dict", dict()).Skip2ndPass(), "foo")
	expect(t, `out = dict["B"]`, Opts().Symbol("dict", dict()).Skip2ndPass(), "bar")
	expect(t, `out = dict["x"]`, Opts().Symbol("dict", dict()).Skip2ndPass(), tengo.UndefinedValue)
	expectError(t, `dict[0]`, Opts().Symbol("dict", dict()).Skip2ndPass(), "invalid index type")

	strCir := func() *StringCircle { return &StringCircle{Value: []string{"one", "two", "three"}} }
	expect(t, `out = cir[0]`, Opts().Symbol("cir", strCir()).Skip2ndPass(), "one")
	expect(t, `out = cir[1]`, Opts().Symbol("cir", strCir()).Skip2ndPass(), "two")
	expect(t, `out = cir[-1]`, Opts().Symbol("cir", strCir()).Skip2ndPass(), "three")
	expect(t, `out = cir[-2]`, Opts().Symbol("cir", strCir()).Skip2ndPass(), "two")
	expect(t, `out = cir[3]`, Opts().Symbol("cir", strCir()).Skip2ndPass(), "one")
	expectError(t, `cir["a"]`, Opts().Symbol("cir", strCir()).Skip2ndPass(), "invalid index type")

	strArr := func() *StringArray { return &StringArray{Value: []string{"one", "two", "three"}} }
	expect(t, `out = arr["one"]`, Opts().Symbol("arr", strArr()).Skip2ndPass(), 0)
	expect(t, `out = arr["three"]`, Opts().Symbol("arr", strArr()).Skip2ndPass(), 2)
	expect(t, `out = arr["four"]`, Opts().Symbol("arr", strArr()).Skip2ndPass(), tengo.UndefinedValue)
	expect(t, `out = arr[0]`, Opts().Symbol("arr", strArr()).Skip2ndPass(), "one")
	expect(t, `out = arr[1]`, Opts().Symbol("arr", strArr()).Skip2ndPass(), "two")
	expectError(t, `arr[-1]`, Opts().Symbol("arr", strArr()).Skip2ndPass(), "index out of bounds")
}

func TestIndexAssignable(t *testing.T) {
	dict := func() *StringDict { return &StringDict{Value: map[string]string{"a": "foo", "b": "bar"}} }
	expect(t, `dict["a"] = "1984"; out = dict["a"]`, Opts().Symbol("dict", dict()).Skip2ndPass(), "1984")
	expect(t, `dict["c"] = "1984"; out = dict["c"]`, Opts().Symbol("dict", dict()).Skip2ndPass(), "1984")
	expect(t, `dict["c"] = 1984; out = dict["C"]`, Opts().Symbol("dict", dict()).Skip2ndPass(), "1984")
	expectError(t, `dict[0] = "1984"`, Opts().Symbol("dict", dict()).Skip2ndPass(), "invalid index type")

	strCir := func() *StringCircle { return &StringCircle{Value: []string{"one", "two", "three"}} }
	expect(t, `cir[0] = "ONE"; out = cir[0]`, Opts().Symbol("cir", strCir()).Skip2ndPass(), "ONE")
	expect(t, `cir[1] = "TWO"; out = cir[1]`, Opts().Symbol("cir", strCir()).Skip2ndPass(), "TWO")
	expect(t, `cir[-1] = "THREE"; out = cir[2]`, Opts().Symbol("cir", strCir()).Skip2ndPass(), "THREE")
	expect(t, `cir[0] = "ONE"; out = cir[3]`, Opts().Symbol("cir", strCir()).Skip2ndPass(), "ONE")
	expectError(t, `cir["a"] = "ONE"`, Opts().Symbol("cir", strCir()).Skip2ndPass(), "invalid index type")

	strArr := func() *StringArray { return &StringArray{Value: []string{"one", "two", "three"}} }
	expect(t, `arr[0] = "ONE"; out = arr[0]`, Opts().Symbol("arr", strArr()).Skip2ndPass(), "ONE")
	expect(t, `arr[1] = "TWO"; out = arr[1]`, Opts().Symbol("arr", strArr()).Skip2ndPass(), "TWO")
	expectError(t, `arr["one"] = "ONE"`, Opts().Symbol("arr", strArr()).Skip2ndPass(), "invalid index type")
}

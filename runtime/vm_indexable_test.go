package runtime_test

import (
	"strings"
	"testing"

	"github.com/d5/tengo/compiler/token"
	"github.com/d5/tengo/objects"
)

type objectImpl struct{}

func (objectImpl) TypeName() string                   { return "" }
func (objectImpl) String() string                     { return "" }
func (objectImpl) IsFalsy() bool                      { return false }
func (objectImpl) Equals(another objects.Object) bool { return false }
func (objectImpl) Copy() objects.Object               { return nil }
func (objectImpl) BinaryOp(token.Token, objects.Object) (objects.Object, error) {
	return nil, objects.ErrInvalidOperator
}

type StringDict struct {
	objectImpl
	Value map[string]string
}

func (o *StringDict) TypeName() string {
	return "string-dict"
}

func (o *StringDict) IndexGet(index objects.Object) (objects.Object, error) {
	strIdx, ok := index.(*objects.String)
	if !ok {
		return nil, objects.ErrInvalidIndexType
	}

	for k, v := range o.Value {
		if strings.ToLower(strIdx.Value) == strings.ToLower(k) {
			return &objects.String{Value: v}, nil
		}
	}

	return objects.UndefinedValue, nil
}

func (o *StringDict) IndexSet(index, value objects.Object) error {
	strIdx, ok := index.(*objects.String)
	if !ok {
		return objects.ErrInvalidIndexType
	}

	strVal, ok := objects.ToString(value)
	if !ok {
		return objects.ErrInvalidIndexValueType
	}

	o.Value[strings.ToLower(strIdx.Value)] = strVal

	return nil
}

type StringCircle struct {
	objectImpl
	Value []string
}

func (o *StringCircle) TypeName() string {
	return "string-circle"
}

func (o *StringCircle) IndexGet(index objects.Object) (objects.Object, error) {
	intIdx, ok := index.(*objects.Int)
	if !ok {
		return nil, objects.ErrInvalidIndexType
	}

	r := int(intIdx.Value) % len(o.Value)
	if r < 0 {
		r = len(o.Value) + r
	}

	return &objects.String{Value: o.Value[r]}, nil
}

func (o *StringCircle) IndexSet(index, value objects.Object) error {
	intIdx, ok := index.(*objects.Int)
	if !ok {
		return objects.ErrInvalidIndexType
	}

	r := int(intIdx.Value) % len(o.Value)
	if r < 0 {
		r = len(o.Value) + r
	}

	strVal, ok := objects.ToString(value)
	if !ok {
		return objects.ErrInvalidIndexValueType
	}

	o.Value[r] = strVal

	return nil
}

type StringArray struct {
	Value []string
}

func (o *StringArray) String() string {
	return strings.Join(o.Value, ", ")
}

func (o *StringArray) BinaryOp(op token.Token, rhs objects.Object) (objects.Object, error) {
	if rhs, ok := rhs.(*StringArray); ok {
		switch op {
		case token.Add:
			if len(rhs.Value) == 0 {
				return o, nil
			}
			return &StringArray{Value: append(o.Value, rhs.Value...)}, nil
		}
	}

	return nil, objects.ErrInvalidOperator
}

func (o *StringArray) IsFalsy() bool {
	return len(o.Value) == 0
}

func (o *StringArray) Equals(x objects.Object) bool {
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

func (o *StringArray) Copy() objects.Object {
	return &StringArray{
		Value: append([]string{}, o.Value...),
	}
}

func (o *StringArray) TypeName() string {
	return "string-array"
}

func (o *StringArray) IndexGet(index objects.Object) (objects.Object, error) {
	intIdx, ok := index.(*objects.Int)
	if ok {
		if intIdx.Value >= 0 && intIdx.Value < int64(len(o.Value)) {
			return &objects.String{Value: o.Value[intIdx.Value]}, nil
		}

		return nil, objects.ErrIndexOutOfBounds
	}

	strIdx, ok := index.(*objects.String)
	if ok {
		for vidx, str := range o.Value {
			if strIdx.Value == str {
				return &objects.Int{Value: int64(vidx)}, nil
			}
		}

		return objects.UndefinedValue, nil
	}

	return nil, objects.ErrInvalidIndexType
}

func (o *StringArray) IndexSet(index, value objects.Object) error {
	strVal, ok := objects.ToString(value)
	if !ok {
		return objects.ErrInvalidIndexValueType
	}

	intIdx, ok := index.(*objects.Int)
	if ok {
		if intIdx.Value >= 0 && intIdx.Value < int64(len(o.Value)) {
			o.Value[intIdx.Value] = strVal
			return nil
		}

		return objects.ErrIndexOutOfBounds
	}

	return objects.ErrInvalidIndexType
}

func (o *StringArray) Call(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}

	s1, ok := objects.ToString(args[0])
	if !ok {
		return nil, objects.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	for i, v := range o.Value {
		if v == s1 {
			return &objects.Int{Value: int64(i)}, nil
		}
	}

	return objects.UndefinedValue, nil
}

func TestIndexable(t *testing.T) {
	dict := func() *StringDict { return &StringDict{Value: map[string]string{"a": "foo", "b": "bar"}} }
	expect(t, `out = dict["a"]`, Opts().Symbol("dict", dict()).Skip2ndPass(), "foo")
	expect(t, `out = dict["B"]`, Opts().Symbol("dict", dict()).Skip2ndPass(), "bar")
	expect(t, `out = dict["x"]`, Opts().Symbol("dict", dict()).Skip2ndPass(), objects.UndefinedValue)
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
	expect(t, `out = arr["four"]`, Opts().Symbol("arr", strArr()).Skip2ndPass(), objects.UndefinedValue)
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

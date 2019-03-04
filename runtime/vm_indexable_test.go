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
func (objectImpl) NumObjects() int64                  { return 1 }
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

func (o *StringArray) NumObjects() int64 {
	return 1
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
	expectWithSymbols(t, `out = dict["a"]`, "foo", SYM{"dict": dict()})
	expectWithSymbols(t, `out = dict["B"]`, "bar", SYM{"dict": dict()})
	expectWithSymbols(t, `out = dict["x"]`, objects.UndefinedValue, SYM{"dict": dict()})
	expectErrorWithSymbols(t, `dict[0]`, SYM{"dict": dict()}, "invalid index type")

	strCir := func() *StringCircle { return &StringCircle{Value: []string{"one", "two", "three"}} }
	expectWithSymbols(t, `out = cir[0]`, "one", SYM{"cir": strCir()})
	expectWithSymbols(t, `out = cir[1]`, "two", SYM{"cir": strCir()})
	expectWithSymbols(t, `out = cir[-1]`, "three", SYM{"cir": strCir()})
	expectWithSymbols(t, `out = cir[-2]`, "two", SYM{"cir": strCir()})
	expectWithSymbols(t, `out = cir[3]`, "one", SYM{"cir": strCir()})
	expectErrorWithSymbols(t, `cir["a"]`, SYM{"cir": strCir()}, "invalid index type")

	strArr := func() *StringArray { return &StringArray{Value: []string{"one", "two", "three"}} }
	expectWithSymbols(t, `out = arr["one"]`, 0, SYM{"arr": strArr()})
	expectWithSymbols(t, `out = arr["three"]`, 2, SYM{"arr": strArr()})
	expectWithSymbols(t, `out = arr["four"]`, objects.UndefinedValue, SYM{"arr": strArr()})
	expectWithSymbols(t, `out = arr[0]`, "one", SYM{"arr": strArr()})
	expectWithSymbols(t, `out = arr[1]`, "two", SYM{"arr": strArr()})
	expectErrorWithSymbols(t, `arr[-1]`, SYM{"arr": strArr()}, "index out of bounds")
}

func TestIndexAssignable(t *testing.T) {
	dict := func() *StringDict { return &StringDict{Value: map[string]string{"a": "foo", "b": "bar"}} }
	expectWithSymbols(t, `dict["a"] = "1984"; out = dict["a"]`, "1984", SYM{"dict": dict()})
	expectWithSymbols(t, `dict["c"] = "1984"; out = dict["c"]`, "1984", SYM{"dict": dict()})
	expectWithSymbols(t, `dict["c"] = 1984; out = dict["C"]`, "1984", SYM{"dict": dict()})
	expectErrorWithSymbols(t, `dict[0] = "1984"`, SYM{"dict": dict()}, "invalid index type")

	strCir := func() *StringCircle { return &StringCircle{Value: []string{"one", "two", "three"}} }
	expectWithSymbols(t, `cir[0] = "ONE"; out = cir[0]`, "ONE", SYM{"cir": strCir()})
	expectWithSymbols(t, `cir[1] = "TWO"; out = cir[1]`, "TWO", SYM{"cir": strCir()})
	expectWithSymbols(t, `cir[-1] = "THREE"; out = cir[2]`, "THREE", SYM{"cir": strCir()})
	expectWithSymbols(t, `cir[0] = "ONE"; out = cir[3]`, "ONE", SYM{"cir": strCir()})
	expectErrorWithSymbols(t, `cir["a"] = "ONE"`, SYM{"cir": strCir()}, "invalid index type")

	strArr := func() *StringArray { return &StringArray{Value: []string{"one", "two", "three"}} }
	expectWithSymbols(t, `arr[0] = "ONE"; out = arr[0]`, "ONE", SYM{"arr": strArr()})
	expectWithSymbols(t, `arr[1] = "TWO"; out = arr[1]`, "TWO", SYM{"arr": strArr()})
	expectErrorWithSymbols(t, `arr["one"] = "ONE"`, SYM{"arr": strArr()}, "invalid index type")
}

package tengo_test

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/parser"
	"github.com/d5/tengo/v2/require"
	"github.com/d5/tengo/v2/token"
)

func TestObject_TypeName(t *testing.T) {
	var o tengo.Object = &tengo.Int{}
	require.Equal(t, "int", o.TypeName())
	o = &tengo.Float{}
	require.Equal(t, "float", o.TypeName())
	o = &tengo.Char{}
	require.Equal(t, "char", o.TypeName())
	o = &tengo.String{}
	require.Equal(t, "string", o.TypeName())
	o = &tengo.Bool{}
	require.Equal(t, "bool", o.TypeName())
	o = &tengo.Array{}
	require.Equal(t, "array", o.TypeName())
	o = &tengo.Map{}
	require.Equal(t, "map", o.TypeName())
	o = &tengo.ArrayIterator{}
	require.Equal(t, "array-iterator", o.TypeName())
	o = &tengo.StringIterator{}
	require.Equal(t, "string-iterator", o.TypeName())
	o = &tengo.MapIterator{}
	require.Equal(t, "map-iterator", o.TypeName())
	o = &tengo.BuiltinFunction{Name: "fn"}
	require.Equal(t, "builtin-function:fn", o.TypeName())
	o = &tengo.UserFunction{Name: "fn"}
	require.Equal(t, "user-function:fn", o.TypeName())
	o = &tengo.CompiledFunction{}
	require.Equal(t, "compiled-function", o.TypeName())
	o = &tengo.Undefined{}
	require.Equal(t, "undefined", o.TypeName())
	o = &tengo.Error{}
	require.Equal(t, "error", o.TypeName())
	o = &tengo.Bytes{}
	require.Equal(t, "bytes", o.TypeName())
	o = &tengo.BytesIterator{}
	require.Equal(t, "bytes-iterator", o.TypeName())
	o = &tengo.ImmutableArray{}
	require.Equal(t, "immutable-array", o.TypeName())
	o = &tengo.ImmutableMap{}
	require.Equal(t, "immutable-map", o.TypeName())
	o = &tengo.Time{}
	require.Equal(t, "time", o.TypeName())
}

func TestObject_IsFalsy(t *testing.T) {
	var o tengo.Object = &tengo.Int{Value: 0}
	require.True(t, o.IsFalsy())
	o = &tengo.Int{Value: 1}
	require.False(t, o.IsFalsy())
	o = &tengo.Float{Value: 0}
	require.False(t, o.IsFalsy())
	o = &tengo.Float{Value: 1}
	require.False(t, o.IsFalsy())
	o = &tengo.Char{Value: ' '}
	require.False(t, o.IsFalsy())
	o = &tengo.Char{Value: 'T'}
	require.False(t, o.IsFalsy())
	o = &tengo.String{Value: ""}
	require.True(t, o.IsFalsy())
	o = &tengo.String{Value: " "}
	require.False(t, o.IsFalsy())
	o = &tengo.Array{Value: nil}
	require.True(t, o.IsFalsy())
	o = &tengo.Array{Value: []tengo.Object{nil}} // nil is not valid but still count as 1 element
	require.False(t, o.IsFalsy())
	o = &tengo.Map{Value: nil}
	require.True(t, o.IsFalsy())
	o = &tengo.Map{Value: map[string]tengo.Object{"a": nil}} // nil is not valid but still count as 1 element
	require.False(t, o.IsFalsy())
	o = &tengo.StringIterator{}
	require.True(t, o.IsFalsy())
	o = &tengo.ArrayIterator{}
	require.True(t, o.IsFalsy())
	o = &tengo.MapIterator{}
	require.True(t, o.IsFalsy())
	o = &tengo.BuiltinFunction{}
	require.False(t, o.IsFalsy())
	o = &tengo.CompiledFunction{}
	require.False(t, o.IsFalsy())
	o = &tengo.Undefined{}
	require.True(t, o.IsFalsy())
	o = &tengo.Error{}
	require.True(t, o.IsFalsy())
	o = &tengo.Bytes{}
	require.True(t, o.IsFalsy())
	o = &tengo.Bytes{Value: []byte{1, 2}}
	require.False(t, o.IsFalsy())
	o = &tengo.Float{Value: math.NaN()}
	require.True(t, o.IsFalsy())
	o = &tengo.Char{Value: 0}
	require.True(t, o.IsFalsy())
	o = &tengo.ImmutableArray{}
	require.True(t, o.IsFalsy())
	o = &tengo.ImmutableArray{Value: []tengo.Object{&tengo.Int{}}}
	require.False(t, o.IsFalsy())
	o = &tengo.ImmutableMap{}
	require.True(t, o.IsFalsy())
	o = &tengo.ImmutableMap{Value: map[string]tengo.Object{"a": &tengo.Int{}}}
	require.False(t, o.IsFalsy())
	o = &tengo.Time{}
	require.True(t, o.IsFalsy())
	o = &tengo.Time{Value: time.Unix(0, 0)}
	require.False(t, o.IsFalsy())
	o = &tengo.ObjectPtr{}
	require.True(t, o.IsFalsy())
	var ptrValue tengo.Object = &tengo.Int{}
	o = &tengo.ObjectPtr{Value: &ptrValue}
	require.False(t, o.IsFalsy())
	o = &tengo.UserFunction{}
	require.False(t, o.IsFalsy())
}

func TestObject_String(t *testing.T) {
	var o tengo.Object = &tengo.Int{Value: 0}
	require.Equal(t, "0", o.String())
	o = &tengo.Int{Value: 1}
	require.Equal(t, "1", o.String())
	o = &tengo.Float{Value: 0}
	require.Equal(t, "0", o.String())
	o = &tengo.Float{Value: 1}
	require.Equal(t, "1", o.String())
	o = &tengo.Char{Value: ' '}
	require.Equal(t, " ", o.String())
	o = &tengo.Char{Value: 'T'}
	require.Equal(t, "T", o.String())
	o = &tengo.String{Value: ""}
	require.Equal(t, `""`, o.String())
	o = &tengo.String{Value: " "}
	require.Equal(t, `" "`, o.String())
	o = &tengo.Array{Value: nil}
	require.Equal(t, "[]", o.String())
	o = &tengo.Map{Value: nil}
	require.Equal(t, "{}", o.String())
	o = &tengo.Error{Value: nil}
	require.Equal(t, "error", o.String())
	o = &tengo.Error{Value: &tengo.String{Value: "error 1"}}
	require.Equal(t, `error: "error 1"`, o.String())
	o = &tengo.StringIterator{}
	require.Equal(t, "<string-iterator>", o.String())
	o = &tengo.ArrayIterator{}
	require.Equal(t, "<array-iterator>", o.String())
	o = &tengo.MapIterator{}
	require.Equal(t, "<map-iterator>", o.String())
	o = &tengo.Undefined{}
	require.Equal(t, "<undefined>", o.String())
	o = &tengo.Bytes{}
	require.Equal(t, "", o.String())
	o = &tengo.Bytes{Value: []byte("foo")}
	require.Equal(t, "foo", o.String())
	o = &tengo.BytesIterator{}
	require.Equal(t, "<bytes-iterator>", o.String())
	o = &tengo.BuiltinFunction{Name: "fn"}
	require.Equal(t, "<builtin-function>", o.String())
	o = &tengo.CompiledFunction{}
	require.Equal(t, "<compiled-function>", o.String())
	o = &tengo.UserFunction{Name: "fn"}
	require.Equal(t, "<user-function>", o.String())
	o = &tengo.ImmutableArray{Value: []tengo.Object{
		&tengo.Int{Value: 1}, &tengo.String{Value: "a"}}}
	require.Equal(t, `[1, "a"]`, o.String())
	o = &tengo.ImmutableMap{Value: map[string]tengo.Object{
		"a": &tengo.Int{Value: 1}}}
	require.Equal(t, "{a: 1}", o.String())
	o = &tengo.Time{Value: time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)}
	require.Equal(t, "2009-11-10 23:00:00 +0000 UTC", o.String())
}

func TestObject_BinaryOp(t *testing.T) {
	var o tengo.Object = &tengo.Char{}
	_, err := o.BinaryOp(token.Add, tengo.UndefinedValue)
	require.Error(t, err)
	o = &tengo.Bool{}
	_, err = o.BinaryOp(token.Add, tengo.UndefinedValue)
	require.Error(t, err)
	o = &tengo.Map{}
	_, err = o.BinaryOp(token.Add, tengo.UndefinedValue)
	require.Error(t, err)
	o = &tengo.ArrayIterator{}
	_, err = o.BinaryOp(token.Add, tengo.UndefinedValue)
	require.Error(t, err)
	o = &tengo.StringIterator{}
	_, err = o.BinaryOp(token.Add, tengo.UndefinedValue)
	require.Error(t, err)
	o = &tengo.MapIterator{}
	_, err = o.BinaryOp(token.Add, tengo.UndefinedValue)
	require.Error(t, err)
	o = &tengo.BuiltinFunction{}
	_, err = o.BinaryOp(token.Add, tengo.UndefinedValue)
	require.Error(t, err)
	o = &tengo.CompiledFunction{}
	_, err = o.BinaryOp(token.Add, tengo.UndefinedValue)
	require.Error(t, err)
	o = &tengo.Undefined{}
	_, err = o.BinaryOp(token.Add, tengo.UndefinedValue)
	require.Error(t, err)
	o = &tengo.Error{}
	_, err = o.BinaryOp(token.Add, tengo.UndefinedValue)
	require.Error(t, err)
}

func TestArray_BinaryOp(t *testing.T) {
	testBinaryOp(t, &tengo.Array{Value: nil}, token.Add,
		&tengo.Array{Value: nil}, &tengo.Array{Value: nil})
	testBinaryOp(t, &tengo.Array{Value: nil}, token.Add,
		&tengo.Array{Value: []tengo.Object{}}, &tengo.Array{Value: nil})
	testBinaryOp(t, &tengo.Array{Value: []tengo.Object{}}, token.Add,
		&tengo.Array{Value: nil}, &tengo.Array{Value: []tengo.Object{}})
	testBinaryOp(t, &tengo.Array{Value: []tengo.Object{}}, token.Add,
		&tengo.Array{Value: []tengo.Object{}},
		&tengo.Array{Value: []tengo.Object{}})
	testBinaryOp(t, &tengo.Array{Value: nil}, token.Add,
		&tengo.Array{Value: []tengo.Object{
			&tengo.Int{Value: 1},
		}}, &tengo.Array{Value: []tengo.Object{
			&tengo.Int{Value: 1},
		}})
	testBinaryOp(t, &tengo.Array{Value: nil}, token.Add,
		&tengo.Array{Value: []tengo.Object{
			&tengo.Int{Value: 1},
			&tengo.Int{Value: 2},
			&tengo.Int{Value: 3},
		}}, &tengo.Array{Value: []tengo.Object{
			&tengo.Int{Value: 1},
			&tengo.Int{Value: 2},
			&tengo.Int{Value: 3},
		}})
	testBinaryOp(t, &tengo.Array{Value: []tengo.Object{
		&tengo.Int{Value: 1},
		&tengo.Int{Value: 2},
		&tengo.Int{Value: 3},
	}}, token.Add, &tengo.Array{Value: nil},
		&tengo.Array{Value: []tengo.Object{
			&tengo.Int{Value: 1},
			&tengo.Int{Value: 2},
			&tengo.Int{Value: 3},
		}})
	testBinaryOp(t, &tengo.Array{Value: []tengo.Object{
		&tengo.Int{Value: 1},
		&tengo.Int{Value: 2},
		&tengo.Int{Value: 3},
	}}, token.Add, &tengo.Array{Value: []tengo.Object{
		&tengo.Int{Value: 4},
		&tengo.Int{Value: 5},
		&tengo.Int{Value: 6},
	}}, &tengo.Array{Value: []tengo.Object{
		&tengo.Int{Value: 1},
		&tengo.Int{Value: 2},
		&tengo.Int{Value: 3},
		&tengo.Int{Value: 4},
		&tengo.Int{Value: 5},
		&tengo.Int{Value: 6},
	}})
}

func TestError_Equals(t *testing.T) {
	err1 := &tengo.Error{Value: &tengo.String{Value: "some error"}}
	err2 := err1
	require.True(t, err1.Equals(err2))
	require.True(t, err2.Equals(err1))

	err2 = &tengo.Error{Value: &tengo.String{Value: "some error"}}
	require.False(t, err1.Equals(err2))
	require.False(t, err2.Equals(err1))
}

func TestFloat_BinaryOp(t *testing.T) {
	// float + float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &tengo.Float{Value: l}, token.Add,
				&tengo.Float{Value: r}, &tengo.Float{Value: l + r})
		}
	}

	// float - float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &tengo.Float{Value: l}, token.Sub,
				&tengo.Float{Value: r}, &tengo.Float{Value: l - r})
		}
	}

	// float * float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &tengo.Float{Value: l}, token.Mul,
				&tengo.Float{Value: r}, &tengo.Float{Value: l * r})
		}
	}

	// float / float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			if r != 0 {
				testBinaryOp(t, &tengo.Float{Value: l}, token.Quo,
					&tengo.Float{Value: r}, &tengo.Float{Value: l / r})
			}
		}
	}

	// float < float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &tengo.Float{Value: l}, token.Less,
				&tengo.Float{Value: r}, boolValue(l < r))
		}
	}

	// float > float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &tengo.Float{Value: l}, token.Greater,
				&tengo.Float{Value: r}, boolValue(l > r))
		}
	}

	// float <= float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &tengo.Float{Value: l}, token.LessEq,
				&tengo.Float{Value: r}, boolValue(l <= r))
		}
	}

	// float >= float
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := float64(-2); r <= 2.1; r += 0.4 {
			testBinaryOp(t, &tengo.Float{Value: l}, token.GreaterEq,
				&tengo.Float{Value: r}, boolValue(l >= r))
		}
	}

	// float + int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Float{Value: l}, token.Add,
				&tengo.Int{Value: r}, &tengo.Float{Value: l + float64(r)})
		}
	}

	// float - int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Float{Value: l}, token.Sub,
				&tengo.Int{Value: r}, &tengo.Float{Value: l - float64(r)})
		}
	}

	// float * int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Float{Value: l}, token.Mul,
				&tengo.Int{Value: r}, &tengo.Float{Value: l * float64(r)})
		}
	}

	// float / int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			if r != 0 {
				testBinaryOp(t, &tengo.Float{Value: l}, token.Quo,
					&tengo.Int{Value: r},
					&tengo.Float{Value: l / float64(r)})
			}
		}
	}

	// float < int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Float{Value: l}, token.Less,
				&tengo.Int{Value: r}, boolValue(l < float64(r)))
		}
	}

	// float > int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Float{Value: l}, token.Greater,
				&tengo.Int{Value: r}, boolValue(l > float64(r)))
		}
	}

	// float <= int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Float{Value: l}, token.LessEq,
				&tengo.Int{Value: r}, boolValue(l <= float64(r)))
		}
	}

	// float >= int
	for l := float64(-2); l <= 2.1; l += 0.4 {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Float{Value: l}, token.GreaterEq,
				&tengo.Int{Value: r}, boolValue(l >= float64(r)))
		}
	}
}

func TestInt_BinaryOp(t *testing.T) {
	// int + int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Add,
				&tengo.Int{Value: r}, &tengo.Int{Value: l + r})
		}
	}

	// int - int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Sub,
				&tengo.Int{Value: r}, &tengo.Int{Value: l - r})
		}
	}

	// int * int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Mul,
				&tengo.Int{Value: r}, &tengo.Int{Value: l * r})
		}
	}

	// int / int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			if r != 0 {
				testBinaryOp(t, &tengo.Int{Value: l}, token.Quo,
					&tengo.Int{Value: r}, &tengo.Int{Value: l / r})
			}
		}
	}

	// int % int
	for l := int64(-4); l <= 4; l++ {
		for r := -int64(-4); r <= 4; r++ {
			if r == 0 {
				testBinaryOp(t, &tengo.Int{Value: l}, token.Rem,
					&tengo.Int{Value: r}, &tengo.Int{Value: l % r})
			}
		}
	}

	// int & int
	testBinaryOp(t,
		&tengo.Int{Value: 0}, token.And, &tengo.Int{Value: 0},
		&tengo.Int{Value: int64(0)})
	testBinaryOp(t,
		&tengo.Int{Value: 1}, token.And, &tengo.Int{Value: 0},
		&tengo.Int{Value: int64(1) & int64(0)})
	testBinaryOp(t,
		&tengo.Int{Value: 0}, token.And, &tengo.Int{Value: 1},
		&tengo.Int{Value: int64(0) & int64(1)})
	testBinaryOp(t,
		&tengo.Int{Value: 1}, token.And, &tengo.Int{Value: 1},
		&tengo.Int{Value: int64(1)})
	testBinaryOp(t,
		&tengo.Int{Value: 0}, token.And, &tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(0) & int64(0xffffffff)})
	testBinaryOp(t,
		&tengo.Int{Value: 1}, token.And, &tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(1) & int64(0xffffffff)})
	testBinaryOp(t,
		&tengo.Int{Value: int64(0xffffffff)}, token.And,
		&tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(0xffffffff)})
	testBinaryOp(t,
		&tengo.Int{Value: 1984}, token.And,
		&tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(1984) & int64(0xffffffff)})
	testBinaryOp(t, &tengo.Int{Value: -1984}, token.And,
		&tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(-1984) & int64(0xffffffff)})

	// int | int
	testBinaryOp(t,
		&tengo.Int{Value: 0}, token.Or, &tengo.Int{Value: 0},
		&tengo.Int{Value: int64(0)})
	testBinaryOp(t,
		&tengo.Int{Value: 1}, token.Or, &tengo.Int{Value: 0},
		&tengo.Int{Value: int64(1) | int64(0)})
	testBinaryOp(t,
		&tengo.Int{Value: 0}, token.Or, &tengo.Int{Value: 1},
		&tengo.Int{Value: int64(0) | int64(1)})
	testBinaryOp(t,
		&tengo.Int{Value: 1}, token.Or, &tengo.Int{Value: 1},
		&tengo.Int{Value: int64(1)})
	testBinaryOp(t,
		&tengo.Int{Value: 0}, token.Or, &tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(0) | int64(0xffffffff)})
	testBinaryOp(t,
		&tengo.Int{Value: 1}, token.Or, &tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(1) | int64(0xffffffff)})
	testBinaryOp(t,
		&tengo.Int{Value: int64(0xffffffff)}, token.Or,
		&tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(0xffffffff)})
	testBinaryOp(t,
		&tengo.Int{Value: 1984}, token.Or,
		&tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(1984) | int64(0xffffffff)})
	testBinaryOp(t,
		&tengo.Int{Value: -1984}, token.Or,
		&tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(-1984) | int64(0xffffffff)})

	// int ^ int
	testBinaryOp(t,
		&tengo.Int{Value: 0}, token.Xor, &tengo.Int{Value: 0},
		&tengo.Int{Value: int64(0)})
	testBinaryOp(t,
		&tengo.Int{Value: 1}, token.Xor, &tengo.Int{Value: 0},
		&tengo.Int{Value: int64(1) ^ int64(0)})
	testBinaryOp(t,
		&tengo.Int{Value: 0}, token.Xor, &tengo.Int{Value: 1},
		&tengo.Int{Value: int64(0) ^ int64(1)})
	testBinaryOp(t,
		&tengo.Int{Value: 1}, token.Xor, &tengo.Int{Value: 1},
		&tengo.Int{Value: int64(0)})
	testBinaryOp(t,
		&tengo.Int{Value: 0}, token.Xor, &tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(0) ^ int64(0xffffffff)})
	testBinaryOp(t,
		&tengo.Int{Value: 1}, token.Xor, &tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(1) ^ int64(0xffffffff)})
	testBinaryOp(t,
		&tengo.Int{Value: int64(0xffffffff)}, token.Xor,
		&tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(0)})
	testBinaryOp(t,
		&tengo.Int{Value: 1984}, token.Xor,
		&tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(1984) ^ int64(0xffffffff)})
	testBinaryOp(t,
		&tengo.Int{Value: -1984}, token.Xor,
		&tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(-1984) ^ int64(0xffffffff)})

	// int &^ int
	testBinaryOp(t,
		&tengo.Int{Value: 0}, token.AndNot, &tengo.Int{Value: 0},
		&tengo.Int{Value: int64(0)})
	testBinaryOp(t,
		&tengo.Int{Value: 1}, token.AndNot, &tengo.Int{Value: 0},
		&tengo.Int{Value: int64(1) &^ int64(0)})
	testBinaryOp(t,
		&tengo.Int{Value: 0}, token.AndNot,
		&tengo.Int{Value: 1}, &tengo.Int{Value: int64(0) &^ int64(1)})
	testBinaryOp(t,
		&tengo.Int{Value: 1}, token.AndNot, &tengo.Int{Value: 1},
		&tengo.Int{Value: int64(0)})
	testBinaryOp(t,
		&tengo.Int{Value: 0}, token.AndNot,
		&tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(0) &^ int64(0xffffffff)})
	testBinaryOp(t,
		&tengo.Int{Value: 1}, token.AndNot,
		&tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(1) &^ int64(0xffffffff)})
	testBinaryOp(t,
		&tengo.Int{Value: int64(0xffffffff)}, token.AndNot,
		&tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(0)})
	testBinaryOp(t,
		&tengo.Int{Value: 1984}, token.AndNot,
		&tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(1984) &^ int64(0xffffffff)})
	testBinaryOp(t,
		&tengo.Int{Value: -1984}, token.AndNot,
		&tengo.Int{Value: int64(0xffffffff)},
		&tengo.Int{Value: int64(-1984) &^ int64(0xffffffff)})

	// int << int
	for s := int64(0); s < 64; s++ {
		testBinaryOp(t,
			&tengo.Int{Value: 0}, token.Shl, &tengo.Int{Value: s},
			&tengo.Int{Value: int64(0) << uint(s)})
		testBinaryOp(t,
			&tengo.Int{Value: 1}, token.Shl, &tengo.Int{Value: s},
			&tengo.Int{Value: int64(1) << uint(s)})
		testBinaryOp(t,
			&tengo.Int{Value: 2}, token.Shl, &tengo.Int{Value: s},
			&tengo.Int{Value: int64(2) << uint(s)})
		testBinaryOp(t,
			&tengo.Int{Value: -1}, token.Shl, &tengo.Int{Value: s},
			&tengo.Int{Value: int64(-1) << uint(s)})
		testBinaryOp(t,
			&tengo.Int{Value: -2}, token.Shl, &tengo.Int{Value: s},
			&tengo.Int{Value: int64(-2) << uint(s)})
		testBinaryOp(t,
			&tengo.Int{Value: int64(0xffffffff)}, token.Shl,
			&tengo.Int{Value: s},
			&tengo.Int{Value: int64(0xffffffff) << uint(s)})
	}

	// int >> int
	for s := int64(0); s < 64; s++ {
		testBinaryOp(t,
			&tengo.Int{Value: 0}, token.Shr, &tengo.Int{Value: s},
			&tengo.Int{Value: int64(0) >> uint(s)})
		testBinaryOp(t,
			&tengo.Int{Value: 1}, token.Shr, &tengo.Int{Value: s},
			&tengo.Int{Value: int64(1) >> uint(s)})
		testBinaryOp(t,
			&tengo.Int{Value: 2}, token.Shr, &tengo.Int{Value: s},
			&tengo.Int{Value: int64(2) >> uint(s)})
		testBinaryOp(t,
			&tengo.Int{Value: -1}, token.Shr, &tengo.Int{Value: s},
			&tengo.Int{Value: int64(-1) >> uint(s)})
		testBinaryOp(t,
			&tengo.Int{Value: -2}, token.Shr, &tengo.Int{Value: s},
			&tengo.Int{Value: int64(-2) >> uint(s)})
		testBinaryOp(t,
			&tengo.Int{Value: int64(0xffffffff)}, token.Shr,
			&tengo.Int{Value: s},
			&tengo.Int{Value: int64(0xffffffff) >> uint(s)})
	}

	// int < int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Less,
				&tengo.Int{Value: r}, boolValue(l < r))
		}
	}

	// int > int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Greater,
				&tengo.Int{Value: r}, boolValue(l > r))
		}
	}

	// int <= int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Int{Value: l}, token.LessEq,
				&tengo.Int{Value: r}, boolValue(l <= r))
		}
	}

	// int >= int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Int{Value: l}, token.GreaterEq,
				&tengo.Int{Value: r}, boolValue(l >= r))
		}
	}

	// int + float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Add,
				&tengo.Float{Value: r},
				&tengo.Float{Value: float64(l) + r})
		}
	}

	// int - float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Sub,
				&tengo.Float{Value: r},
				&tengo.Float{Value: float64(l) - r})
		}
	}

	// int * float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Mul,
				&tengo.Float{Value: r},
				&tengo.Float{Value: float64(l) * r})
		}
	}

	// int / float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			if r != 0 {
				testBinaryOp(t, &tengo.Int{Value: l}, token.Quo,
					&tengo.Float{Value: r},
					&tengo.Float{Value: float64(l) / r})
			}
		}
	}

	// int < float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Less,
				&tengo.Float{Value: r}, boolValue(float64(l) < r))
		}
	}

	// int > float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Greater,
				&tengo.Float{Value: r}, boolValue(float64(l) > r))
		}
	}

	// int <= float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &tengo.Int{Value: l}, token.LessEq,
				&tengo.Float{Value: r}, boolValue(float64(l) <= r))
		}
	}

	// int >= float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &tengo.Int{Value: l}, token.GreaterEq,
				&tengo.Float{Value: r}, boolValue(float64(l) >= r))
		}
	}
}

func TestInt_Equals(t *testing.T) {
	require.True(t, (&tengo.Int{Value: 1}).Equals(&tengo.Int{Value: 1}))
	require.False(t, (&tengo.Int{Value: 1}).Equals(&tengo.Int{Value: 2}))

	// int == float: must be consistent with BinaryOp cross-type ordering
	require.True(t, (&tengo.Int{Value: 1}).Equals(&tengo.Float{Value: 1.0}))
	require.False(t, (&tengo.Int{Value: 1}).Equals(&tengo.Float{Value: 1.5}))
	require.True(t, (&tengo.Int{Value: 0}).Equals(&tengo.Float{Value: 0.0}))
	require.False(t, (&tengo.Int{Value: 1}).Equals(&tengo.Float{Value: -1.0}))

	// trichotomy: neither (a < b) nor (a > b) implies (a == b)
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			lt, _ := (&tengo.Int{Value: l}).BinaryOp(token.Less, &tengo.Float{Value: r})
			gt, _ := (&tengo.Int{Value: l}).BinaryOp(token.Greater, &tengo.Float{Value: r})
			if lt == tengo.FalseValue && gt == tengo.FalseValue {
				require.True(t, (&tengo.Int{Value: l}).Equals(&tengo.Float{Value: r}),
					"expected %d == %g", l, r)
			}
		}
	}
}

func TestFloat_Equals(t *testing.T) {
	require.True(t, (&tengo.Float{Value: 1.5}).Equals(&tengo.Float{Value: 1.5}))
	require.False(t, (&tengo.Float{Value: 1.5}).Equals(&tengo.Float{Value: 2.5}))

	// float == int: must be consistent with BinaryOp cross-type ordering
	require.True(t, (&tengo.Float{Value: 1.0}).Equals(&tengo.Int{Value: 1}))
	require.False(t, (&tengo.Float{Value: 1.5}).Equals(&tengo.Int{Value: 1}))
	require.True(t, (&tengo.Float{Value: 0.0}).Equals(&tengo.Int{Value: 0}))

	// trichotomy: neither (a < b) nor (a > b) implies (a == b)
	for l := float64(-2); l <= 2.1; l += 0.5 {
		for r := int64(-2); r <= 2; r++ {
			lt, _ := (&tengo.Float{Value: l}).BinaryOp(token.Less, &tengo.Int{Value: r})
			gt, _ := (&tengo.Float{Value: l}).BinaryOp(token.Greater, &tengo.Int{Value: r})
			if lt == tengo.FalseValue && gt == tengo.FalseValue {
				require.True(t, (&tengo.Float{Value: l}).Equals(&tengo.Int{Value: r}),
					"expected %g == %d", l, r)
			}
		}
	}
}

func TestMap_Index(t *testing.T) {
	m := &tengo.Map{Value: make(map[string]tengo.Object)}
	k := &tengo.Int{Value: 1}
	v := &tengo.String{Value: "abcdef"}
	err := m.IndexSet(k, v)

	require.NoError(t, err)

	res, err := m.IndexGet(k)
	require.NoError(t, err)
	require.Equal(t, v, res)
}

func TestString_BinaryOp(t *testing.T) {
	lstr := "abcde"
	rstr := "01234"
	for l := 0; l < len(lstr); l++ {
		for r := 0; r < len(rstr); r++ {
			ls := lstr[l:]
			rs := rstr[r:]
			testBinaryOp(t, &tengo.String{Value: ls}, token.Add,
				&tengo.String{Value: rs},
				&tengo.String{Value: ls + rs})

			rc := []rune(rstr)[r]
			testBinaryOp(t, &tengo.String{Value: ls}, token.Add,
				&tengo.Char{Value: rc},
				&tengo.String{Value: ls + string(rc)})
		}
	}
}

func testBinaryOp(
	t *testing.T,
	lhs tengo.Object,
	op token.Token,
	rhs tengo.Object,
	expected tengo.Object,
) {
	t.Helper()
	actual, err := lhs.BinaryOp(op, rhs)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func boolValue(b bool) tengo.Object {
	if b {
		return tengo.TrueValue
	}
	return tengo.FalseValue
}

func TestFromInterface_JSONNumber(t *testing.T) {
	testFromInterface(t, json.Number("123"), &tengo.Int{Value: 123})
	testFromInterface(t, json.Number("-123"), &tengo.Int{Value: -123})
	testFromInterface(t, json.Number("0"), &tengo.Int{Value: 0})
	testFromInterface(t, json.Number("123.45"), &tengo.Float{Value: 123.45})
	testFromInterface(t, json.Number("-123.45"), &tengo.Float{Value: -123.45})
	testFromInterface(t, json.Number("1e5"), &tengo.Float{Value: 100000})
	testFromInterface(t, json.Number("1.5e-3"), &tengo.Float{Value: 0.0015})

	numInt := json.Number("123")
	testFromInterface(t, &numInt, &tengo.Int{Value: 123})
	numFloat := json.Number("123.45")
	testFromInterface(t, &numFloat, &tengo.Float{Value: 123.45})
	var nilNum *json.Number
	testFromInterface(t, nilNum, tengo.UndefinedValue)

	testFromInterface(t, map[string]interface{}{
		"int":   json.Number("123"),
		"float": json.Number("123.45"),
	}, &tengo.Map{
		Value: map[string]tengo.Object{
			"int":   &tengo.Int{Value: 123},
			"float": &tengo.Float{Value: 123.45},
		},
	})

	testFromInterface(t, []interface{}{
		json.Number("123"),
		json.Number("123.45"),
	}, &tengo.Array{
		Value: []tengo.Object{
			&tengo.Int{Value: 123},
			&tengo.Float{Value: 123.45},
		},
	})

	_, err := tengo.FromInterface(json.Number("invalid"))
	require.Error(t, err)

	invalidNum := json.Number("not_a_number")
	_, err = tengo.FromInterface(&invalidNum)
	require.Error(t, err)
}

func TestScript_JSONNumber(t *testing.T) {
	s := tengo.NewScript([]byte(`out := num_int + 1; f_out := num_float * 2.0`))
	err := s.Add("num_int", json.Number("123"))
	require.NoError(t, err)
	err = s.Add("num_float", json.Number("123.45"))
	require.NoError(t, err)

	c, err := s.Run()
	require.NoError(t, err)
	require.Equal(t, int64(124), c.Get("out").Int64())
	require.Equal(t, 246.9, c.Get("f_out").Float())
}

func testFromInterface(t *testing.T, v interface{}, expected tengo.Object) {
	t.Helper()
	actual, err := tengo.FromInterface(v)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestObjectImpl(t *testing.T) {
	o := &tengo.ObjectImpl{}
	requirePanic(t, tengo.ErrNotImplemented, func() { _ = o.TypeName() })
	requirePanic(t, tengo.ErrNotImplemented, func() { _ = o.String() })

	_, err := o.BinaryOp(token.Add, &tengo.Int{Value: 1})
	require.Equal(t, tengo.ErrInvalidOperator, err)
	require.Nil(t, o.Copy())
	require.False(t, o.IsFalsy())
	_, err = o.IndexGet(&tengo.Int{Value: 0})
	require.Equal(t, tengo.ErrNotIndexable, err)
	err = o.IndexSet(&tengo.Int{Value: 0}, &tengo.Int{Value: 1})
	require.Equal(t, tengo.ErrNotIndexAssignable, err)
	require.False(t, o.CanIterate())
	require.Nil(t, o.Iterate())
	require.False(t, o.CanCall())
}

func TestObject_Copy(t *testing.T) {
	require.True(t, tengo.TrueValue.Copy() == tengo.TrueValue)
	require.True(t, tengo.FalseValue.Copy() == tengo.FalseValue)
	require.True(t, tengo.UndefinedValue.Copy() == tengo.UndefinedValue)

	ptr := &tengo.ObjectPtr{}
	require.True(t, ptr.Copy() == ptr)

	tm := time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)
	require.Equal(t, &tengo.Time{Value: tm}, (&tengo.Time{Value: tm}).Copy())

	b := &tengo.Bytes{Value: []byte("abc")}
	bc := b.Copy().(*tengo.Bytes)
	require.Equal(t, []byte("abc"), bc.Value)
	bc.Value[0] = 'x'
	require.Equal(t, []byte("abc"), b.Value)

	e := &tengo.Error{Value: &tengo.Array{Value: []tengo.Object{
		&tengo.Int{Value: 1}}}}
	ec := e.Copy().(*tengo.Error)
	require.Equal(t, e, ec)
	ec.Value.(*tengo.Array).Value[0] = &tengo.Int{Value: 2}
	require.Equal(t, &tengo.Error{Value: &tengo.Array{Value: []tengo.Object{
		&tengo.Int{Value: 1}}}}, e)
}

func TestObject_CopyContainers(t *testing.T) {
	nested := func() []tengo.Object {
		return []tengo.Object{
			&tengo.Int{Value: 1},
			&tengo.Array{Value: []tengo.Object{&tengo.Int{Value: 2}}},
		}
	}
	nestedMap := func() map[string]tengo.Object {
		return map[string]tengo.Object{
			"a": &tengo.Int{Value: 1},
			"b": &tengo.Array{Value: []tengo.Object{&tengo.Int{Value: 2}}},
		}
	}

	arr := &tengo.Array{Value: nested()}
	arrCopy := arr.Copy().(*tengo.Array)
	require.Equal(t, arr, arrCopy)
	arrCopy.Value[0] = &tengo.Int{Value: 9}
	arrCopy.Value[1].(*tengo.Array).Value[0] = &tengo.Int{Value: 9}
	require.Equal(t, &tengo.Array{Value: nested()}, arr)

	iarr := &tengo.ImmutableArray{Value: nested()}
	iarrCopy, ok := iarr.Copy().(*tengo.Array)
	require.True(t, ok, "copy of immutable array must be mutable")
	require.Equal(t, &tengo.Array{Value: nested()}, iarrCopy)
	iarrCopy.Value[1].(*tengo.Array).Value[0] = &tengo.Int{Value: 9}
	require.Equal(t, &tengo.ImmutableArray{Value: nested()}, iarr)

	m := &tengo.Map{Value: nestedMap()}
	mCopy := m.Copy().(*tengo.Map)
	require.Equal(t, m, mCopy)
	mCopy.Value["a"] = &tengo.Int{Value: 9}
	mCopy.Value["b"].(*tengo.Array).Value[0] = &tengo.Int{Value: 9}
	require.Equal(t, &tengo.Map{Value: nestedMap()}, m)

	im := &tengo.ImmutableMap{Value: nestedMap()}
	imCopy, ok := im.Copy().(*tengo.Map)
	require.True(t, ok, "copy of immutable map must be mutable")
	require.Equal(t, &tengo.Map{Value: nestedMap()}, imCopy)
	imCopy.Value["b"].(*tengo.Array).Value[0] = &tengo.Int{Value: 9}
	require.Equal(t, &tengo.ImmutableMap{Value: nestedMap()}, im)
}

func TestObject_CopyFunctions(t *testing.T) {
	strLen := func(args ...tengo.Object) (tengo.Object, error) {
		return &tengo.Int{Value: int64(len(args[0].(*tengo.String).Value))}, nil
	}

	bf := (&tengo.BuiltinFunction{Name: "fn", Value: strLen}).Copy()
	require.True(t, bf.CanCall())
	res, err := bf.Call(&tengo.String{Value: "abc"})
	require.NoError(t, err)
	require.Equal(t, &tengo.Int{Value: 3}, res)

	uf := (&tengo.UserFunction{Name: "fn", Value: strLen}).Copy()
	require.Equal(t, "user-function:fn", uf.TypeName())
	res, err = uf.Call(&tengo.String{Value: "abcd"})
	require.NoError(t, err)
	require.Equal(t, &tengo.Int{Value: 4}, res)

	free := &tengo.ObjectPtr{}
	fn := &tengo.CompiledFunction{
		Instructions:  tengo.MakeInstruction(parser.OpConstant, 1),
		NumLocals:     2,
		NumParameters: 1,
		VarArgs:       true,
		Free:          []*tengo.ObjectPtr{free},
	}
	fnCopy := fn.Copy().(*tengo.CompiledFunction)
	require.Equal(t, fn.Instructions, fnCopy.Instructions)
	require.Equal(t, 2, fnCopy.NumLocals)
	require.Equal(t, 1, fnCopy.NumParameters)
	require.True(t, fnCopy.VarArgs)
	require.Equal(t, 1, len(fnCopy.Free))
	require.True(t, fnCopy.Free[0] == free, "free variables must be shared")

	fnCopy.Instructions[1] = 0xff
	fnCopy.Free[0] = &tengo.ObjectPtr{}
	require.Equal(t, tengo.MakeInstruction(parser.OpConstant, 1),
		fn.Instructions)
	require.True(t, fn.Free[0] == free)
}

func TestObject_Equals(t *testing.T) {
	require.True(t, tengo.TrueValue.Equals(tengo.TrueValue))
	require.False(t, tengo.TrueValue.Equals(tengo.FalseValue))
	require.False(t, tengo.FalseValue.Equals(&tengo.Int{Value: 0}))

	b := &tengo.Bytes{Value: []byte("abc")}
	require.True(t, b.Equals(&tengo.Bytes{Value: []byte("abc")}))
	require.False(t, b.Equals(&tengo.Bytes{Value: []byte("abd")}))
	require.False(t, b.Equals(&tengo.String{Value: "abc"}))

	require.True(t, (&tengo.Char{Value: 'a'}).Equals(&tengo.Char{Value: 'a'}))
	require.False(t, (&tengo.Char{Value: 'a'}).Equals(&tengo.Char{Value: 'b'}))
	require.False(t, (&tengo.Char{Value: 'a'}).Equals(&tengo.Int{Value: 'a'}))

	require.True(t, (&tengo.Float{Value: 1}).Equals(&tengo.Float{Value: 1}))
	require.False(t, (&tengo.Float{Value: 1}).Equals(&tengo.String{Value: "1"}))
	require.False(t, (&tengo.Int{Value: 1}).Equals(&tengo.String{Value: "1"}))

	tm := time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)
	tmTokyo := tm.In(time.FixedZone("UTC+9", 9*60*60))
	require.True(t, (&tengo.Time{Value: tm}).Equals(&tengo.Time{Value: tmTokyo}))
	require.False(t, (&tengo.Time{Value: tm}).Equals(
		&tengo.Time{Value: tm.Add(1)}))
	require.False(t, (&tengo.Time{Value: tm}).Equals(
		&tengo.Int{Value: tm.UnixNano()}))

	require.True(t, tengo.UndefinedValue.Equals(tengo.UndefinedValue))
	require.False(t, tengo.UndefinedValue.Equals(&tengo.Int{Value: 0}))

	var v tengo.Object = &tengo.Int{Value: 1}
	ptr := &tengo.ObjectPtr{Value: &v}
	require.True(t, ptr.Equals(ptr))
	require.False(t, ptr.Equals(&tengo.ObjectPtr{Value: &v}))

	for _, fn := range []tengo.Object{
		&tengo.BuiltinFunction{Name: "fn"},
		&tengo.UserFunction{Name: "fn"},
		&tengo.CompiledFunction{},
	} {
		require.False(t, fn.Equals(fn.Copy()), fn.TypeName())
	}
}

func TestObject_EqualsContainers(t *testing.T) {
	elems := func() []tengo.Object {
		return []tengo.Object{&tengo.Int{Value: 1}, &tengo.String{Value: "a"}}
	}
	for _, lhs := range []tengo.Object{
		&tengo.Array{Value: elems()},
		&tengo.ImmutableArray{Value: elems()},
	} {
		require.True(t, lhs.Equals(&tengo.Array{Value: elems()}), lhs.TypeName())
		require.True(t, lhs.Equals(&tengo.ImmutableArray{Value: elems()}),
			lhs.TypeName())
		require.False(t, lhs.Equals(&tengo.Array{Value: elems()[:1]}),
			lhs.TypeName())
		require.False(t, lhs.Equals(&tengo.Array{Value: []tengo.Object{
			&tengo.Int{Value: 1}, &tengo.String{Value: "b"}}}), lhs.TypeName())
		require.False(t, lhs.Equals(&tengo.Map{}), lhs.TypeName())
	}

	kv := func() map[string]tengo.Object {
		return map[string]tengo.Object{"a": &tengo.Int{Value: 1}}
	}
	for _, lhs := range []tengo.Object{
		&tengo.Map{Value: kv()},
		&tengo.ImmutableMap{Value: kv()},
	} {
		require.True(t, lhs.Equals(&tengo.Map{Value: kv()}), lhs.TypeName())
		require.True(t, lhs.Equals(&tengo.ImmutableMap{Value: kv()}),
			lhs.TypeName())
		require.False(t, lhs.Equals(&tengo.Map{}), lhs.TypeName())
		require.False(t, lhs.Equals(&tengo.Map{Value: map[string]tengo.Object{
			"b": &tengo.Int{Value: 1}}}), lhs.TypeName())
		require.False(t, lhs.Equals(&tengo.Map{Value: map[string]tengo.Object{
			"a": &tengo.Int{Value: 2}}}), lhs.TypeName())
		require.False(t, lhs.Equals(&tengo.Array{}), lhs.TypeName())
	}
}

func TestObject_IndexInvalidIndexType(t *testing.T) {
	for _, o := range []tengo.Object{
		&tengo.Array{Value: []tengo.Object{&tengo.Int{Value: 1}}},
		&tengo.ImmutableArray{Value: []tengo.Object{&tengo.Int{Value: 1}}},
		&tengo.Bytes{Value: []byte("a")},
		&tengo.String{Value: "a"},
	} {
		_, err := o.IndexGet(&tengo.String{Value: "0"})
		require.Equal(t, tengo.ErrInvalidIndexType, err, o.TypeName())
	}

	for _, o := range []tengo.Object{
		&tengo.Map{Value: map[string]tengo.Object{}},
		&tengo.ImmutableMap{Value: map[string]tengo.Object{}},
	} {
		_, err := o.IndexGet(tengo.UndefinedValue)
		require.Equal(t, tengo.ErrInvalidIndexType, err, o.TypeName())
	}

	m := &tengo.Map{Value: map[string]tengo.Object{}}
	err := m.IndexSet(tengo.UndefinedValue, &tengo.Int{Value: 1})
	require.Equal(t, tengo.ErrInvalidIndexType, err)
	require.Equal(t, 0, len(m.Value))
}

func TestObject_BinaryOpInvalidOperator(t *testing.T) {
	for _, tc := range []struct {
		lhs tengo.Object
		op  token.Token
		rhs tengo.Object
	}{
		{&tengo.Array{}, token.Add, &tengo.Int{Value: 1}},
		{&tengo.Array{}, token.Sub, &tengo.Array{}},
		{&tengo.Bytes{}, token.Add, &tengo.String{Value: "a"}},
		{&tengo.Bytes{}, token.Sub, &tengo.Bytes{}},
		{&tengo.Char{Value: 'a'}, token.Mul, &tengo.Char{Value: 'b'}},
		{&tengo.Char{Value: 'a'}, token.Mul, &tengo.Int{Value: 1}},
		{&tengo.Char{Value: 'a'}, token.Add, &tengo.Float{Value: 1}},
		{&tengo.Float{Value: 1}, token.Rem, &tengo.Float{Value: 1}},
		{&tengo.Float{Value: 1}, token.Rem, &tengo.Int{Value: 1}},
		{&tengo.Float{Value: 1}, token.Add, &tengo.Char{Value: 'a'}},
		{&tengo.Int{Value: 1}, token.Mul, &tengo.Char{Value: 'a'}},
		{&tengo.Int{Value: 1}, token.Add, &tengo.String{Value: "a"}},
		{&tengo.String{Value: "a"}, token.Sub, &tengo.String{Value: "a"}},
		{&tengo.String{Value: "a"}, token.Less, &tengo.Int{Value: 1}},
		{&tengo.Time{}, token.Add, &tengo.Time{}},
		{&tengo.Time{}, token.Mul, &tengo.Int{Value: 1}},
		{&tengo.Time{}, token.Less, &tengo.Int{Value: 1}},
		{&tengo.Time{}, token.Add, &tengo.Float{Value: 1}},
	} {
		_, err := tc.lhs.BinaryOp(tc.op, tc.rhs)
		require.Equal(t, tengo.ErrInvalidOperator, err, "%s %s %s",
			tc.lhs.TypeName(), tc.op, tc.rhs.TypeName())
	}
}

func TestChar_BinaryOp(t *testing.T) {
	testBinaryOp(t, &tengo.Char{Value: 'a'}, token.Add,
		&tengo.Char{Value: 2}, &tengo.Char{Value: 'c'})
	testBinaryOp(t, &tengo.Char{Value: 'c'}, token.Sub,
		&tengo.Char{Value: 2}, &tengo.Char{Value: 'a'})
	testBinaryOp(t, &tengo.Char{Value: 'a'}, token.Add,
		&tengo.Char{Value: 0}, &tengo.Char{Value: 'a'})
	testBinaryOp(t, &tengo.Char{Value: 'a'}, token.Sub,
		&tengo.Char{Value: 0}, &tengo.Char{Value: 'a'})
	testBinaryOp(t, &tengo.Char{Value: 'a'}, token.Add,
		&tengo.Int{Value: 2}, &tengo.Char{Value: 'c'})
	testBinaryOp(t, &tengo.Char{Value: 'c'}, token.Sub,
		&tengo.Int{Value: 2}, &tengo.Char{Value: 'a'})
	testBinaryOp(t, &tengo.Char{Value: 'a'}, token.Add,
		&tengo.Int{Value: 0}, &tengo.Char{Value: 'a'})
	testBinaryOp(t, &tengo.Char{Value: 'a'}, token.Sub,
		&tengo.Int{Value: 0}, &tengo.Char{Value: 'a'})
	testBinaryOp(t, &tengo.Int{Value: 2}, token.Add,
		&tengo.Char{Value: 'a'}, &tengo.Char{Value: 'c'})
	testBinaryOp(t, &tengo.Int{Value: 'c'}, token.Sub,
		&tengo.Char{Value: 2}, &tengo.Char{Value: 'a'})

	for l := 'a'; l <= 'c'; l++ {
		for r := 'a'; r <= 'c'; r++ {
			for _, rhs := range []tengo.Object{
				&tengo.Char{Value: r},
				&tengo.Int{Value: int64(r)},
			} {
				lhs := &tengo.Char{Value: l}
				testBinaryOp(t, lhs, token.Less, rhs, boolValue(l < r))
				testBinaryOp(t, lhs, token.Greater, rhs, boolValue(l > r))
				testBinaryOp(t, lhs, token.LessEq, rhs, boolValue(l <= r))
				testBinaryOp(t, lhs, token.GreaterEq, rhs, boolValue(l >= r))
			}

			lhs := &tengo.Int{Value: int64(l)}
			rhs := &tengo.Char{Value: r}
			testBinaryOp(t, lhs, token.Less, rhs, boolValue(l < r))
			testBinaryOp(t, lhs, token.Greater, rhs, boolValue(l > r))
			testBinaryOp(t, lhs, token.LessEq, rhs, boolValue(l <= r))
			testBinaryOp(t, lhs, token.GreaterEq, rhs, boolValue(l >= r))
		}
	}
}

func TestInt_BinaryOpRem(t *testing.T) {
	for l := int64(-4); l <= 4; l++ {
		for r := int64(-4); r <= 4; r++ {
			if r == 0 {
				_, err := (&tengo.Int{Value: l}).BinaryOp(token.Rem,
					&tengo.Int{Value: r})
				require.Equal(t, tengo.ErrDivisionByZero, err)
				continue
			}
			testBinaryOp(t, &tengo.Int{Value: l}, token.Rem,
				&tengo.Int{Value: r}, &tengo.Int{Value: l % r})
		}
	}
}

func TestFloat_BinaryOpIdentity(t *testing.T) {
	testBinaryOp(t, &tengo.Float{Value: 2.5}, token.Mul,
		&tengo.Float{Value: 1}, &tengo.Float{Value: 2.5})
	testBinaryOp(t, &tengo.Float{Value: 2.5}, token.Quo,
		&tengo.Float{Value: 1}, &tengo.Float{Value: 2.5})
	testBinaryOp(t, &tengo.Float{Value: 2.5}, token.Add,
		&tengo.Float{Value: 0}, &tengo.Float{Value: 2.5})
	testBinaryOp(t, &tengo.Float{Value: 2.5}, token.Sub,
		&tengo.Float{Value: 0}, &tengo.Float{Value: 2.5})
	testBinaryOp(t, &tengo.Float{Value: 0}, token.Mul,
		&tengo.Float{Value: 7}, &tengo.Float{Value: 0})
	testBinaryOp(t, &tengo.Float{Value: 2.5}, token.Mul,
		&tengo.Int{Value: 1}, &tengo.Float{Value: 2.5})
	testBinaryOp(t, &tengo.Float{Value: 2.5}, token.Quo,
		&tengo.Int{Value: 1}, &tengo.Float{Value: 2.5})
}

func TestString_BinaryOpCompare(t *testing.T) {
	for _, l := range []string{"", "a", "ab", "b"} {
		for _, r := range []string{"", "a", "ab", "b"} {
			lhs := &tengo.String{Value: l}
			rhs := &tengo.String{Value: r}
			testBinaryOp(t, lhs, token.Less, rhs, boolValue(l < r))
			testBinaryOp(t, lhs, token.Greater, rhs, boolValue(l > r))
			testBinaryOp(t, lhs, token.LessEq, rhs, boolValue(l <= r))
			testBinaryOp(t, lhs, token.GreaterEq, rhs, boolValue(l >= r))
		}
	}
}

func TestTime_BinaryOp(t *testing.T) {
	tm0 := time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)
	tm1 := tm0.Add(time.Second)
	sec := &tengo.Int{Value: int64(time.Second)}

	testBinaryOp(t, &tengo.Time{Value: tm0}, token.Add, sec,
		&tengo.Time{Value: tm1})
	testBinaryOp(t, &tengo.Time{Value: tm1}, token.Sub, sec,
		&tengo.Time{Value: tm0})
	testBinaryOp(t, &tengo.Time{Value: tm0}, token.Add, &tengo.Int{Value: 0},
		&tengo.Time{Value: tm0})
	testBinaryOp(t, &tengo.Time{Value: tm0}, token.Sub, &tengo.Int{Value: 0},
		&tengo.Time{Value: tm0})
	testBinaryOp(t, &tengo.Time{Value: tm1}, token.Sub,
		&tengo.Time{Value: tm0}, sec)
	testBinaryOp(t, &tengo.Time{Value: tm0}, token.Sub,
		&tengo.Time{Value: tm1}, &tengo.Int{Value: -int64(time.Second)})

	tm0Tokyo := tm0.In(time.FixedZone("UTC+9", 9*60*60))
	for _, l := range []time.Time{tm0, tm1, tm0Tokyo} {
		for _, r := range []time.Time{tm0, tm1, tm0Tokyo} {
			lhs := &tengo.Time{Value: l}
			rhs := &tengo.Time{Value: r}
			testBinaryOp(t, lhs, token.Less, rhs, boolValue(l.Before(r)))
			testBinaryOp(t, lhs, token.Greater, rhs, boolValue(l.After(r)))
			testBinaryOp(t, lhs, token.LessEq, rhs, boolValue(!l.After(r)))
			testBinaryOp(t, lhs, token.GreaterEq, rhs,
				boolValue(!l.Before(r)))
		}
	}
	testBinaryOp(t, &tengo.Time{Value: tm0Tokyo}, token.Sub,
		&tengo.Time{Value: tm0}, &tengo.Int{Value: 0})
}

func TestUndefined_Iterate(t *testing.T) {
	require.True(t, tengo.UndefinedValue.CanIterate())
	it := tengo.UndefinedValue.Iterate()
	require.False(t, it.Next())
	require.Equal(t, tengo.UndefinedValue, it.Key())
	require.Equal(t, tengo.UndefinedValue, it.Value())
}

func TestCompiledFunction_Size(t *testing.T) {
	fn := &tengo.CompiledFunction{
		Instructions: concatInsts(
			tengo.MakeInstruction(parser.OpConstant, 0),
			tengo.MakeInstruction(parser.OpReturn, 1)),
		SourceMap: map[int]parser.Pos{0: 1, 3: 5},
		Free:      []*tengo.ObjectPtr{{}},
	}
	require.Equal(t, int64(5+2+1), fn.Size())
	require.Equal(t, int64(0), (&tengo.CompiledFunction{}).Size())
}

func requirePanic(t *testing.T, expected interface{}, fn func()) {
	t.Helper()
	defer func() {
		require.Equal(t, expected, recover())
	}()
	fn()
}

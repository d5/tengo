package tengo_test

import (
	"testing"

	"github.com/d5/tengo/v2"
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

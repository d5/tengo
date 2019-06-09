package tengo_test

import (
	"testing"

	"github.com/d5/tengo"
	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/compiler/token"
)

func TestObject_TypeName(t *testing.T) {
	var o tengo.Object = &tengo.Int{}
	assert.Equal(t, "int", o.TypeName())
	o = &tengo.Float{}
	assert.Equal(t, "float", o.TypeName())
	o = &tengo.Char{}
	assert.Equal(t, "char", o.TypeName())
	o = &tengo.String{}
	assert.Equal(t, "string", o.TypeName())
	o = &tengo.Bool{}
	assert.Equal(t, "bool", o.TypeName())
	o = &tengo.Array{}
	assert.Equal(t, "array", o.TypeName())
	o = &tengo.Map{}
	assert.Equal(t, "map", o.TypeName())
	o = &tengo.ArrayIterator{}
	assert.Equal(t, "array-iterator", o.TypeName())
	o = &tengo.StringIterator{}
	assert.Equal(t, "string-iterator", o.TypeName())
	o = &tengo.MapIterator{}
	assert.Equal(t, "map-iterator", o.TypeName())
	o = &tengo.BuiltinFunction{Name: "fn"}
	assert.Equal(t, "builtin-function:fn", o.TypeName())
	o = &tengo.UserFunction{Name: "fn"}
	assert.Equal(t, "user-function:fn", o.TypeName())
	o = &tengo.CompiledFunction{}
	assert.Equal(t, "compiled-function", o.TypeName())
	o = &tengo.Undefined{}
	assert.Equal(t, "undefined", o.TypeName())
	o = &tengo.Error{}
	assert.Equal(t, "error", o.TypeName())
	o = &tengo.Bytes{}
	assert.Equal(t, "bytes", o.TypeName())
}

func TestObject_IsFalsy(t *testing.T) {
	var o tengo.Object = &tengo.Int{Value: 0}
	assert.True(t, o.IsFalsy())
	o = &tengo.Int{Value: 1}
	assert.False(t, o.IsFalsy())
	o = &tengo.Float{Value: 0}
	assert.False(t, o.IsFalsy())
	o = &tengo.Float{Value: 1}
	assert.False(t, o.IsFalsy())
	o = &tengo.Char{Value: ' '}
	assert.False(t, o.IsFalsy())
	o = &tengo.Char{Value: 'T'}
	assert.False(t, o.IsFalsy())
	o = &tengo.String{Value: ""}
	assert.True(t, o.IsFalsy())
	o = &tengo.String{Value: " "}
	assert.False(t, o.IsFalsy())
	o = &tengo.Array{Value: nil}
	assert.True(t, o.IsFalsy())
	o = &tengo.Array{Value: []tengo.Object{nil}} // nil is not valid but still count as 1 element
	assert.False(t, o.IsFalsy())
	o = &tengo.Map{Value: nil}
	assert.True(t, o.IsFalsy())
	o = &tengo.Map{Value: map[string]tengo.Object{"a": nil}} // nil is not valid but still count as 1 element
	assert.False(t, o.IsFalsy())
	o = &tengo.StringIterator{}
	assert.True(t, o.IsFalsy())
	o = &tengo.ArrayIterator{}
	assert.True(t, o.IsFalsy())
	o = &tengo.MapIterator{}
	assert.True(t, o.IsFalsy())
	o = &tengo.BuiltinFunction{}
	assert.False(t, o.IsFalsy())
	o = &tengo.CompiledFunction{}
	assert.False(t, o.IsFalsy())
	o = &tengo.Undefined{}
	assert.True(t, o.IsFalsy())
	o = &tengo.Error{}
	assert.True(t, o.IsFalsy())
	o = &tengo.Bytes{}
	assert.True(t, o.IsFalsy())
	o = &tengo.Bytes{Value: []byte{1, 2}}
	assert.False(t, o.IsFalsy())
}

func TestObject_String(t *testing.T) {
	var o tengo.Object = &tengo.Int{Value: 0}
	assert.Equal(t, "0", o.String())
	o = &tengo.Int{Value: 1}
	assert.Equal(t, "1", o.String())
	o = &tengo.Float{Value: 0}
	assert.Equal(t, "0", o.String())
	o = &tengo.Float{Value: 1}
	assert.Equal(t, "1", o.String())
	o = &tengo.Char{Value: ' '}
	assert.Equal(t, " ", o.String())
	o = &tengo.Char{Value: 'T'}
	assert.Equal(t, "T", o.String())
	o = &tengo.String{Value: ""}
	assert.Equal(t, `""`, o.String())
	o = &tengo.String{Value: " "}
	assert.Equal(t, `" "`, o.String())
	o = &tengo.Array{Value: nil}
	assert.Equal(t, "[]", o.String())
	o = &tengo.Map{Value: nil}
	assert.Equal(t, "{}", o.String())
	o = &tengo.Error{Value: nil}
	assert.Equal(t, "error", o.String())
	o = &tengo.Error{Value: &tengo.String{Value: "error 1"}}
	assert.Equal(t, `error: "error 1"`, o.String())
	o = &tengo.StringIterator{}
	assert.Equal(t, "<string-iterator>", o.String())
	o = &tengo.ArrayIterator{}
	assert.Equal(t, "<array-iterator>", o.String())
	o = &tengo.MapIterator{}
	assert.Equal(t, "<map-iterator>", o.String())
	o = &tengo.Undefined{}
	assert.Equal(t, "<undefined>", o.String())
	o = &tengo.Bytes{}
	assert.Equal(t, "", o.String())
	o = &tengo.Bytes{Value: []byte("foo")}
	assert.Equal(t, "foo", o.String())
}

func TestObject_BinaryOp(t *testing.T) {
	var o tengo.Object = &tengo.Char{}
	_, err := o.BinaryOp(token.Add, tengo.UndefinedValue)
	assert.Error(t, err)
	o = &tengo.Bool{}
	_, err = o.BinaryOp(token.Add, tengo.UndefinedValue)
	assert.Error(t, err)
	o = &tengo.Map{}
	_, err = o.BinaryOp(token.Add, tengo.UndefinedValue)
	assert.Error(t, err)
	o = &tengo.ArrayIterator{}
	_, err = o.BinaryOp(token.Add, tengo.UndefinedValue)
	assert.Error(t, err)
	o = &tengo.StringIterator{}
	_, err = o.BinaryOp(token.Add, tengo.UndefinedValue)
	assert.Error(t, err)
	o = &tengo.MapIterator{}
	_, err = o.BinaryOp(token.Add, tengo.UndefinedValue)
	assert.Error(t, err)
	o = &tengo.BuiltinFunction{}
	_, err = o.BinaryOp(token.Add, tengo.UndefinedValue)
	assert.Error(t, err)
	o = &tengo.CompiledFunction{}
	_, err = o.BinaryOp(token.Add, tengo.UndefinedValue)
	assert.Error(t, err)
	o = &tengo.Undefined{}
	_, err = o.BinaryOp(token.Add, tengo.UndefinedValue)
	assert.Error(t, err)
	o = &tengo.Error{}
	_, err = o.BinaryOp(token.Add, tengo.UndefinedValue)
	assert.Error(t, err)
}

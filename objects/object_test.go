package objects_test

import (
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/compiler/token"
	"github.com/d5/tengo/objects"
)

func TestObject_TypeName(t *testing.T) {
	var o objects.Object
	o = &objects.Int{}
	assert.Equal(t, "int", o.TypeName())
	o = &objects.Float{}
	assert.Equal(t, "float", o.TypeName())
	o = &objects.Char{}
	assert.Equal(t, "char", o.TypeName())
	o = &objects.String{}
	assert.Equal(t, "string", o.TypeName())
	o = &objects.Bool{}
	assert.Equal(t, "bool", o.TypeName())
	o = &objects.Array{}
	assert.Equal(t, "array", o.TypeName())
	o = &objects.Map{}
	assert.Equal(t, "map", o.TypeName())
	o = &objects.ArrayIterator{}
	assert.Equal(t, "array-iterator", o.TypeName())
	o = &objects.StringIterator{}
	assert.Equal(t, "string-iterator", o.TypeName())
	o = &objects.MapIterator{}
	assert.Equal(t, "map-iterator", o.TypeName())
	o = &objects.Break{}
	assert.Equal(t, "break", o.TypeName())
	o = &objects.Continue{}
	assert.Equal(t, "continue", o.TypeName())
	o = &objects.BuiltinFunction{}
	assert.Equal(t, "builtin-function", o.TypeName())
	o = &objects.Closure{}
	assert.Equal(t, "closure", o.TypeName())
	o = &objects.CompiledFunction{}
	assert.Equal(t, "compiled-function", o.TypeName())
	o = &objects.ReturnValue{}
	assert.Equal(t, "return-value", o.TypeName())
	o = &objects.Undefined{}
	assert.Equal(t, "undefined", o.TypeName())
	o = &objects.Error{}
	assert.Equal(t, "error", o.TypeName())
}

func TestObject_IsFalsy(t *testing.T) {
	var o objects.Object
	o = &objects.Int{Value: 0}
	assert.True(t, o.IsFalsy())
	o = &objects.Int{Value: 1}
	assert.False(t, o.IsFalsy())
	o = &objects.Float{Value: 0}
	assert.False(t, o.IsFalsy())
	o = &objects.Float{Value: 1}
	assert.False(t, o.IsFalsy())
	o = &objects.Char{Value: ' '}
	assert.False(t, o.IsFalsy())
	o = &objects.Char{Value: 'T'}
	assert.False(t, o.IsFalsy())
	o = &objects.String{Value: ""}
	assert.True(t, o.IsFalsy())
	o = &objects.String{Value: " "}
	assert.False(t, o.IsFalsy())
	o = &objects.Array{Value: nil}
	assert.True(t, o.IsFalsy())
	o = &objects.Array{Value: []objects.Object{nil}} // nil is not valid but still count as 1 element
	assert.False(t, o.IsFalsy())
	o = &objects.Map{Value: nil}
	assert.True(t, o.IsFalsy())
	o = &objects.Map{Value: map[string]objects.Object{"a": nil}} // nil is not valid but still count as 1 element
	assert.False(t, o.IsFalsy())
	o = &objects.StringIterator{}
	assert.True(t, o.IsFalsy())
	o = &objects.ArrayIterator{}
	assert.True(t, o.IsFalsy())
	o = &objects.MapIterator{}
	assert.True(t, o.IsFalsy())
	o = &objects.Break{}
	assert.False(t, o.IsFalsy())
	o = &objects.Continue{}
	assert.False(t, o.IsFalsy())
	o = &objects.BuiltinFunction{}
	assert.False(t, o.IsFalsy())
	o = &objects.Closure{}
	assert.False(t, o.IsFalsy())
	o = &objects.CompiledFunction{}
	assert.False(t, o.IsFalsy())
	o = &objects.ReturnValue{}
	assert.False(t, o.IsFalsy())
	o = &objects.Undefined{}
	assert.True(t, o.IsFalsy())
	o = &objects.Error{}
	assert.True(t, o.IsFalsy())
}

func TestObject_String(t *testing.T) {
	var o objects.Object
	o = &objects.Int{Value: 0}
	assert.Equal(t, "0", o.String())
	o = &objects.Int{Value: 1}
	assert.Equal(t, "1", o.String())
	o = &objects.Float{Value: 0}
	assert.Equal(t, "0", o.String())
	o = &objects.Float{Value: 1}
	assert.Equal(t, "1", o.String())
	o = &objects.Char{Value: ' '}
	assert.Equal(t, " ", o.String())
	o = &objects.Char{Value: 'T'}
	assert.Equal(t, "T", o.String())
	o = &objects.String{Value: ""}
	assert.Equal(t, `""`, o.String())
	o = &objects.String{Value: " "}
	assert.Equal(t, `" "`, o.String())
	o = &objects.Array{Value: nil}
	assert.Equal(t, "[]", o.String())
	o = &objects.Map{Value: nil}
	assert.Equal(t, "{}", o.String())
	o = &objects.Error{Value: nil}
	assert.Equal(t, "error", o.String())
	o = &objects.Error{Value: &objects.String{Value: "error 1"}}
	assert.Equal(t, `error: "error 1"`, o.String())
	o = &objects.StringIterator{}
	assert.Equal(t, "<string-iterator>", o.String())
	o = &objects.ArrayIterator{}
	assert.Equal(t, "<array-iterator>", o.String())
	o = &objects.MapIterator{}
	assert.Equal(t, "<map-iterator>", o.String())
	o = &objects.Break{}
	assert.Equal(t, "<break>", o.String())
	o = &objects.Continue{}
	assert.Equal(t, "<continue>", o.String())
	o = &objects.ReturnValue{}
	assert.Equal(t, "<return-value>", o.String())
	o = &objects.Undefined{}
	assert.Equal(t, "<undefined>", o.String())
}

func TestObject_BinaryOp(t *testing.T) {
	var o objects.Object
	o = &objects.Char{}
	_, err := o.BinaryOp(token.Add, objects.UndefinedValue)
	assert.Error(t, err)
	o = &objects.Bool{}
	_, err = o.BinaryOp(token.Add, objects.UndefinedValue)
	assert.Error(t, err)
	o = &objects.Map{}
	_, err = o.BinaryOp(token.Add, objects.UndefinedValue)
	assert.Error(t, err)
	o = &objects.ArrayIterator{}
	_, err = o.BinaryOp(token.Add, objects.UndefinedValue)
	assert.Error(t, err)
	o = &objects.StringIterator{}
	_, err = o.BinaryOp(token.Add, objects.UndefinedValue)
	assert.Error(t, err)
	o = &objects.MapIterator{}
	_, err = o.BinaryOp(token.Add, objects.UndefinedValue)
	assert.Error(t, err)
	o = &objects.Break{}
	_, err = o.BinaryOp(token.Add, objects.UndefinedValue)
	assert.Error(t, err)
	o = &objects.Continue{}
	_, err = o.BinaryOp(token.Add, objects.UndefinedValue)
	assert.Error(t, err)
	o = &objects.BuiltinFunction{}
	_, err = o.BinaryOp(token.Add, objects.UndefinedValue)
	assert.Error(t, err)
	o = &objects.Closure{}
	_, err = o.BinaryOp(token.Add, objects.UndefinedValue)
	assert.Error(t, err)
	o = &objects.CompiledFunction{}
	_, err = o.BinaryOp(token.Add, objects.UndefinedValue)
	assert.Error(t, err)
	o = &objects.ReturnValue{}
	_, err = o.BinaryOp(token.Add, objects.UndefinedValue)
	assert.Error(t, err)
	o = &objects.Undefined{}
	_, err = o.BinaryOp(token.Add, objects.UndefinedValue)
	assert.Error(t, err)
	o = &objects.Error{}
	_, err = o.BinaryOp(token.Add, objects.UndefinedValue)
	assert.Error(t, err)
}

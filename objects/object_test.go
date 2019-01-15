package objects_test

import (
	"testing"

	"github.com/d5/tengo/assert"
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
}

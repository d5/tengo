package tengo_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/require"
)

type VariableTest struct {
	Name        string
	Value       interface{}
	ValueType   string
	IntValue    int
	Int64Value  int64
	FloatValue  float64
	CharValue   rune
	BoolValue   bool
	StringValue string
	Object      tengo.Object
	IsUndefined bool
}

func TestVariable(t *testing.T) {
	vars := []VariableTest{
		{
			Name:        "a",
			Value:       int64(1),
			ValueType:   "int",
			IntValue:    1,
			Int64Value:  1,
			FloatValue:  1.0,
			CharValue:   rune(1),
			BoolValue:   true,
			StringValue: "1",
			Object:      &tengo.Int{Value: 1},
		},
		{
			Name:        "b",
			Value:       "52.11",
			ValueType:   "string",
			FloatValue:  52.11,
			StringValue: "52.11",
			BoolValue:   true,
			Object:      &tengo.String{Value: "52.11"},
		},
		{
			Name:        "c",
			Value:       true,
			ValueType:   "bool",
			IntValue:    1,
			Int64Value:  1,
			FloatValue:  0,
			BoolValue:   true,
			StringValue: "true",
			Object:      tengo.TrueValue,
		},
		{
			Name:        "d",
			Value:       nil,
			ValueType:   "undefined",
			Object:      tengo.UndefinedValue,
			IsUndefined: true,
		},
	}

	for _, tc := range vars {
		v, err := tengo.NewVariable(tc.Name, tc.Value)
		require.NoError(t, err)
		require.Equal(t, tc.Value, v.Value(), "Name: %s", tc.Name)
		require.Equal(t, tc.ValueType, v.ValueType(), "Name: %s", tc.Name)
		require.Equal(t, tc.IntValue, v.Int(), "Name: %s", tc.Name)
		require.Equal(t, tc.Int64Value, v.Int64(), "Name: %s", tc.Name)
		require.Equal(t, tc.FloatValue, v.Float(), "Name: %s", tc.Name)
		require.Equal(t, tc.CharValue, v.Char(), "Name: %s", tc.Name)
		require.Equal(t, tc.BoolValue, v.Bool(), "Name: %s", tc.Name)
		require.Equal(t, tc.StringValue, v.String(), "Name: %s", tc.Name)
		require.Equal(t, tc.Object, v.Object(), "Name: %s", tc.Name)
		require.Equal(t, tc.IsUndefined, v.IsUndefined(), "Name: %s", tc.Name)
	}
}

func TestVariable_Array(t *testing.T) {
	v, err := tengo.NewVariable("a", []interface{}{1, "b", []interface{}{true}})
	require.NoError(t, err)
	expected := []interface{}{int64(1), "b", []interface{}{true}}
	require.True(t, reflect.DeepEqual(expected, v.Array()),
		"expected %#v, got %#v", expected, v.Array())

	for _, value := range []interface{}{1, "abc", nil, map[string]interface{}{}} {
		v, err = tengo.NewVariable("a", value)
		require.NoError(t, err)
		require.True(t, v.Array() == nil, "%#v", value)
	}
}

func TestVariable_Bytes(t *testing.T) {
	v, err := tengo.NewVariable("a", []byte{0, 1, 255})
	require.NoError(t, err)
	require.Equal(t, []byte{0, 1, 255}, v.Bytes())

	v, err = tengo.NewVariable("a", "abc")
	require.NoError(t, err)
	require.Equal(t, []byte("abc"), v.Bytes())

	for _, value := range []interface{}{1, 'a', nil, []interface{}{1}} {
		v, err = tengo.NewVariable("a", value)
		require.NoError(t, err)
		require.Nil(t, v.Bytes(), "%#v", value)
	}
}

func TestVariable_Error(t *testing.T) {
	v, err := tengo.NewVariable("a", errors.New("boom"))
	require.NoError(t, err)
	require.Equal(t, "error", v.ValueType())
	require.Error(t, v.Error())
	require.Equal(t, `error: "boom"`, v.Error().Error())

	for _, value := range []interface{}{1, "boom", nil} {
		v, err = tengo.NewVariable("a", value)
		require.NoError(t, err)
		require.NoError(t, v.Error(), "%#v", value)
	}

	c := compile(t, `e := error("bad input"); n := 1`, nil)
	compiledRun(t, c)
	require.Equal(t, `error: "bad input"`, c.Get("e").Error().Error())
	require.NoError(t, c.Get("n").Error())
}

func TestNewVariable_UnsupportedType(t *testing.T) {
	_, err := tengo.NewVariable("a", make(chan int))
	require.Error(t, err)
}

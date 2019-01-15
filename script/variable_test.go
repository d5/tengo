package script_test

import (
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
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
	Object      objects.Object
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
			StringValue: "1",
			Object:      &objects.Int{Value: 1},
		},
		{
			Name:        "b",
			Value:       "52.11",
			ValueType:   "string",
			FloatValue:  52.11,
			StringValue: "52.11",
			Object:      &objects.String{Value: "52.11"},
		},
		{
			Name:        "c",
			Value:       true,
			ValueType:   "bool",
			IntValue:    1,
			Int64Value:  1,
			FloatValue:  1,
			BoolValue:   true,
			StringValue: "true",
			Object:      &objects.Bool{Value: true},
		},
	}

	for _, tc := range vars {
		v, err := script.NewVariable(tc.Name, tc.Value)
		assert.NoError(t, err)
		assert.Equal(t, tc.Value, v.Value())
		assert.Equal(t, tc.ValueType, v.ValueType())
		assert.Equal(t, tc.IntValue, v.Int())
		assert.Equal(t, tc.Int64Value, v.Int64())
		assert.Equal(t, tc.FloatValue, v.Float())
		assert.Equal(t, tc.CharValue, v.Char())
		assert.Equal(t, tc.BoolValue, v.Bool())
		assert.Equal(t, tc.StringValue, v.String())
		assert.Equal(t, tc.Object, v.Object())
	}

}

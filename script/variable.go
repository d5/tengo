package script

import (
	"strconv"

	"github.com/d5/tengo/objects"
)

// Variable is a user-defined variable for the script.
type Variable struct {
	name  string
	value *objects.Object
}

func NewVariable(name string, value interface{}) (*Variable, error) {
	obj, err := interfaceToObject(value)
	if err != nil {
		return nil, err
	}

	return &Variable{
		name:  name,
		value: &obj,
	}, nil
}

// Name returns the name of the variable.
func (v *Variable) Name() string {
	return v.name
}

// Value returns an empty interface of the variable value.
func (v *Variable) Value() interface{} {
	return objectToInterface(*v.value)
}

// ValueType returns the name of the value type.
func (v *Variable) ValueType() string {
	return (*v.value).TypeName()
}

// Int returns int value of the variable value.
// It returns 0 if the value is not convertible to int.
func (v *Variable) Int() int {
	return int(v.Int64())
}

// Int64 returns int64 value of the variable value.
// It returns 0 if the value is not convertible to int64.
func (v *Variable) Int64() int64 {
	switch val := (*v.value).(type) {
	case *objects.Int:
		return val.Value
	case *objects.Float:
		return int64(val.Value)
	case *objects.Bool:
		if val.Value {
			return 1
		}
		return 0
	case *objects.Char:
		return int64(val.Value)
	case *objects.String:
		n, _ := strconv.ParseInt(val.Value, 10, 64)
		return n
	}

	return 0
}

// Float returns float64 value of the variable value.
// It returns 0.0 if the value is not convertible to float64.
func (v *Variable) Float() float64 {
	switch val := (*v.value).(type) {
	case *objects.Int:
		return float64(val.Value)
	case *objects.Float:
		return val.Value
	case *objects.Bool:
		if val.Value {
			return 1
		}
		return 0
	case *objects.String:
		f, _ := strconv.ParseFloat(val.Value, 64)
		return f
	}

	return 0
}

// Char returns rune value of the variable value.
// It returns 0 if the value is not convertible to rune.
func (v *Variable) Char() rune {
	switch val := (*v.value).(type) {
	case *objects.Char:
		return val.Value
	}

	return 0
}

// Bool returns bool value of the variable value.
// It returns 0 if the value is not convertible to bool.
func (v *Variable) Bool() bool {
	switch val := (*v.value).(type) {
	case *objects.Bool:
		return val.Value
	}

	return false
}

// Array returns []interface value of the variable value.
// It returns 0 if the value is not convertible to []interface.
func (v *Variable) Array() []interface{} {
	switch val := (*v.value).(type) {
	case *objects.Array:
		var arr []interface{}
		for _, e := range val.Value {
			arr = append(arr, objectToInterface(e))
		}
		return arr
	}

	return nil
}

// Map returns map[string]interface{} value of the variable value.
// It returns 0 if the value is not convertible to map[string]interface{}.
func (v *Variable) Map() map[string]interface{} {
	switch val := (*v.value).(type) {
	case *objects.Map:
		kv := make(map[string]interface{})
		for mk, mv := range val.Value {
			kv[mk] = objectToInterface(mv)
		}
		return kv
	}

	return nil
}

// String returns string value of the variable value.
// It returns 0 if the value is not convertible to string.
func (v *Variable) String() string {
	return objectToString(*v.value)
}

// Object returns an underlying Object of the variable value.
// Note that returned Object is a copy of an actual Object used in the script.
func (v *Variable) Object() objects.Object {
	return *v.value
}

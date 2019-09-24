package objects

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/d5/tengo"
)

// ToString will try to convert object o to string value.
func ToString(o Object) (v string, ok bool) {
	if o == UndefinedValue {
		//ok = false
		return
	}

	ok = true

	if str, isStr := o.(*String); isStr {
		v = str.Value
	} else {
		v = o.String()
	}

	return
}

// ToInt will try to convert object o to int value.
func ToInt(o Object) (v int, ok bool) {
	switch o := o.(type) {
	case *Int:
		v = int(o.Value)
		ok = true
	case *Float:
		v = int(o.Value)
		ok = true
	case *Char:
		v = int(o.Value)
		ok = true
	case *Bool:
		if o == TrueValue {
			v = 1
		}
		ok = true
	case *String:
		c, err := strconv.ParseInt(o.Value, 10, 64)
		if err == nil {
			v = int(c)
			ok = true
		}
	}

	//ok = false
	return
}

// ToInt64 will try to convert object o to int64 value.
func ToInt64(o Object) (v int64, ok bool) {
	switch o := o.(type) {
	case *Int:
		v = o.Value
		ok = true
	case *Float:
		v = int64(o.Value)
		ok = true
	case *Char:
		v = int64(o.Value)
		ok = true
	case *Bool:
		if o == TrueValue {
			v = 1
		}
		ok = true
	case *String:
		c, err := strconv.ParseInt(o.Value, 10, 64)
		if err == nil {
			v = c
			ok = true
		}
	}

	//ok = false
	return
}

// ToFloat64 will try to convert object o to float64 value.
func ToFloat64(o Object) (v float64, ok bool) {
	switch o := o.(type) {
	case *Int:
		v = float64(o.Value)
		ok = true
	case *Float:
		v = o.Value
		ok = true
	case *String:
		c, err := strconv.ParseFloat(o.Value, 64)
		if err == nil {
			v = c
			ok = true
		}
	}

	//ok = false
	return
}

// ToBool will try to convert object o to bool value.
func ToBool(o Object) (v bool, ok bool) {
	ok = true
	v = !o.IsFalsy()

	return
}

// ToRune will try to convert object o to rune value.
func ToRune(o Object) (v rune, ok bool) {
	switch o := o.(type) {
	case *Int:
		v = rune(o.Value)
		ok = true
	case *Char:
		v = rune(o.Value)
		ok = true
	}

	//ok = false
	return
}

// ToByteSlice will try to convert object o to []byte value.
func ToByteSlice(o Object) (v []byte, ok bool) {
	switch o := o.(type) {
	case *Bytes:
		v = o.Value
		ok = true
	case *String:
		v = []byte(o.Value)
		ok = true
	}

	//ok = false
	return
}

// ToTime will try to convert object o to time.Time value.
func ToTime(o Object) (v time.Time, ok bool) {
	switch o := o.(type) {
	case *Time:
		v = o.Value
		ok = true
	case *Int:
		v = time.Unix(o.Value, 0)
		ok = true
	}

	//ok = false
	return
}

// ToInterface attempts to convert an object o to an interface{} value
func ToInterface(o Object) (res interface{}) {
	switch o := o.(type) {
	case *Int:
		res = o.Value
	case *String:
		res = o.Value
	case *Float:
		res = o.Value
	case *Bool:
		res = o == TrueValue
	case *Char:
		res = o.Value
	case *Bytes:
		res = o.Value
	case *Array:
		res = make([]interface{}, len(o.Value))
		for i, val := range o.Value {
			res.([]interface{})[i] = ToInterface(val)
		}
	case *ImmutableArray:
		res = make([]interface{}, len(o.Value))
		for i, val := range o.Value {
			res.([]interface{})[i] = ToInterface(val)
		}
	case *Map:
		res = make(map[string]interface{})
		for key, v := range o.Value {
			res.(map[string]interface{})[key] = ToInterface(v)
		}
	case *ImmutableMap:
		res = make(map[string]interface{})
		for key, v := range o.Value {
			res.(map[string]interface{})[key] = ToInterface(v)
		}
	case *Time:
		res = o.Value
	case *Error:
		res = errors.New(o.String())
	case *Undefined:
		res = nil
	case Object:
		return o
	}

	return
}

// FromInterface will attempt to convert an interface{} v to a Tengo Object
func FromInterface(v interface{}) (Object, error) {
	switch v := v.(type) {
	case nil:
		return UndefinedValue, nil
	case string:
		if len(v) > tengo.MaxStringLen {
			return nil, ErrStringLimit
		}
		return &String{Value: v}, nil
	case int64:
		return &Int{Value: v}, nil
	case json.Number:
		if len(v) > tengo.MaxStringLen {
			return nil, ErrStringLimit
		}
		return &String{Value: string(v)}, nil
	case int:
		return &Int{Value: int64(v)}, nil
	case bool:
		if v {
			return TrueValue, nil
		}
		return FalseValue, nil
	case rune:
		return &Char{Value: v}, nil
	case byte:
		return &Char{Value: rune(v)}, nil
	case float64:
		return &Float{Value: v}, nil
	case []byte:
		if len(v) > tengo.MaxBytesLen {
			return nil, ErrBytesLimit
		}
		return &Bytes{Value: v}, nil
	case error:
		return &Error{Value: &String{Value: v.Error()}}, nil
	case map[string]Object:
		return &Map{Value: v}, nil
	case map[string]interface{}:
		kv := make(map[string]Object)
		for vk, vv := range v {
			vo, err := FromInterface(vv)
			if err != nil {
				return nil, err
			}
			kv[vk] = vo
		}
		return &Map{Value: kv}, nil
	case []Object:
		return &Array{Value: v}, nil
	case []interface{}:
		arr := make([]Object, len(v))
		for i, e := range v {
			vo, err := FromInterface(e)
			if err != nil {
				return nil, err
			}

			arr[i] = vo
		}
		return &Array{Value: arr}, nil
	case time.Time:
		return &Time{Value: v}, nil
	case Object:
		return v, nil
	case CallableFunc:
		return &UserFunction{Value: v}, nil
	}

	return nil, fmt.Errorf("cannot convert to object: %T", v)
}

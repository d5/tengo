package tengo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	// MaxStringLen is the maximum byte-length for string value. Note this
	// limit applies to all compiler/VM instances in the process.
	MaxStringLen = 2147483647

	// MaxBytesLen is the maximum length for bytes value. Note this limit
	// applies to all compiler/VM instances in the process.
	MaxBytesLen = 2147483647
)

const (
	// GlobalsSize is the maximum number of global variables for a VM.
	GlobalsSize = 1024

	// StackSize is the maximum stack size for a VM.
	StackSize = 2048

	// MaxFrames is the maximum number of function frames for a VM.
	MaxFrames = 1024

	// SourceFileExtDefault is the default extension for source files.
	SourceFileExtDefault = ".tengo"
)

// CallableFunc is a function signature for the callable functions.
type CallableFunc = func(args ...Object) (Object, error)

// CountObjects returns the number of objects that a given object o contains.
// For scalar value types, it will always be 1. For compound value types,
// this will include its elements and all of their elements recursively.
func CountObjects(o Object) int {
	c := 1
	switch o := o.(type) {
	case *Array:
		for _, v := range o.Value {
			c += CountObjects(v)
		}
	case *ImmutableArray:
		for _, v := range o.Value {
			c += CountObjects(v)
		}
	case *Map:
		for _, v := range o.Value {
			c += CountObjects(v)
		}
	case *ImmutableMap:
		for _, v := range o.Value {
			c += CountObjects(v)
		}
	case *Error:
		c += CountObjects(o.Value)
	}
	return c
}

// ToString will try to convert arg at index i to string value.
func ToString(i int, args ...Object) (string, error) {
	o := args[i]
	if o == UndefinedValue {
		return "", &ErrInvalidArgumentType{
			Index:    i,
			Expected: "not undefined",
			Actual:   UndefinedTN,
		}
	}
	if str, isStr := o.(*String); isStr {
		return str.Value, nil
	}
	return o.String(), nil
}

// ToStringSlice will try to convert object o to string slice value.
func ToStringSlice(i int, args ...Object) ([]string, error) {
	var ss []string
	switch o := args[i].(type) {
	case *Array:
		for idx := range o.Value {
			as, err := ToString(idx, o.Value...)
			if err != nil {
				return nil, err
			}
			ss = append(ss, as)
		}
	case *ImmutableArray:
		for idx := range o.Value {
			as, err := ToString(idx, o.Value...)
			if err != nil {
				return nil, err
			}
			ss = append(ss, as)
		}
	default:
		return nil, ErrInvalidArgumentType{
			Index:    i,
			Expected: ArrayTN,
			Actual:   args[i].TypeName(),
		}
	}
	return ss, nil
}

// ToInt will try to convert a specific arg at index i to an int.
func ToInt(i int, args ...Object) (int, error) {
	switch o := args[i].(type) {
	case *Int:
		return int(o.Value), nil
	case *Float:
		return int(o.Value), nil
	case *Char:
		return int(o.Value), nil
	case *Bool:
		if o == TrueValue {
			return 1, nil
		}
		return 0, nil
	case *String:
		c, err := strconv.ParseInt(o.Value, 10, 64)
		if err == nil {
			return int(c), nil
		}
	}
	return 0, ErrInvalidArgumentType{
		Index:    i,
		Expected: strings.Join(TNs{IntTN, FloatTN, CharTN, BoolTN, StringTN}, "/"),
		Actual:   args[i].TypeName(),
	}
}

// ToIntSlice will try to convert object o to int slice value.
func ToIntSlice(i int, args ...Object) ([]int, error) {
	var is []int
	switch o := args[i].(type) {
	case *Array:
		for idx := range o.Value {
			as, err := ToInt(idx, o.Value...)
			if err != nil {
				return nil, err
			}
			is = append(is, as)
		}
	case *ImmutableArray:
		for idx := range o.Value {
			as, err := ToInt(idx, o.Value...)
			if err != nil {
				return nil, err
			}
			is = append(is, as)
		}
	default:
		return nil, ErrInvalidArgumentType{
			Index:    i,
			Expected: ArrayTN,
			Actual:   args[i].TypeName(),
		}
	}
	return is, nil
}

// ToInt64 will try to convert object o to int64 value.
func ToInt64(i int, args ...Object) (int64, error) {
	switch o := args[i].(type) {
	case *Int:
		return o.Value, nil
	case *Float:
		return int64(o.Value), nil
	case *Char:
		return int64(o.Value), nil
	case *Bool:
		if o == TrueValue {
			return 1, nil
		}
		return 0, nil
	case *String:
		c, err := strconv.ParseInt(o.Value, 10, 64)
		if err == nil {
			return c, nil
		}
	}
	return 0, ErrInvalidArgumentType{
		Index:    i,
		Expected: strings.Join(TNs{IntTN, FloatTN, CharTN, BoolTN, StringTN}, "/"),
		Actual:   args[i].TypeName(),
	}
}

// ToInt64Slice will try to convert object o to int64 slice value.
func ToInt64Slice(i int, args ...Object) ([]int64, error) {
	var is []int64
	switch o := args[i].(type) {
	case *Array:
		for idx := range o.Value {
			as, err := ToInt64(idx, o.Value...)
			if err != nil {
				return nil, err
			}
			is = append(is, as)
		}
	case *ImmutableArray:
		for idx := range o.Value {
			as, err := ToInt64(idx, o.Value...)
			if err != nil {
				return nil, err
			}
			is = append(is, as)
		}
	default:
		return nil, ErrInvalidArgumentType{
			Index:    i,
			Expected: ArrayTN,
			Actual:   args[i].TypeName(),
		}
	}
	return is, nil
}

// ToFloat64 will try to convert object o to float64 value.
func ToFloat64(i int, args ...Object) (float64, error) {
	switch o := args[i].(type) {
	case *Int:
		return float64(o.Value), nil
	case *Float:
		return o.Value, nil
	case *String:
		c, err := strconv.ParseFloat(o.Value, 64)
		if err == nil {
			return c, nil
		}
	}
	return 0, ErrInvalidArgumentType{
		Index:    i,
		Expected: strings.Join(TNs{IntTN, FloatTN, StringTN}, "/"),
		Actual:   args[i].TypeName(),
	}
}

// ToFloat64Slice will try to convert object o to float64 slice value.
func ToFloat64Slice(i int, args ...Object) ([]float64, error) {
	var fs []float64
	switch o := args[i].(type) {
	case *Array:
		for idx := range o.Value {
			as, err := ToFloat64(idx, o.Value...)
			if err != nil {
				return nil, err
			}
			fs = append(fs, as)
		}
	case *ImmutableArray:
		for idx := range o.Value {
			as, err := ToFloat64(idx, o.Value...)
			if err != nil {
				return nil, err
			}
			fs = append(fs, as)
		}
	default:
		return nil, ErrInvalidArgumentType{
			Index:    i,
			Expected: ArrayTN,
			Actual:   args[i].TypeName(),
		}
	}
	return fs, nil
}

// ToBool will try to convert object o to bool value.
func ToBool(i int, args ...Object) bool {
	return !args[i].IsFalsy()
}

// ToRune will try to convert object o to rune value.
func ToRune(i int, args ...Object) (rune, error) {
	switch o := args[i].(type) {
	case *Int:
		return rune(o.Value), nil
	case *Char:
		return o.Value, nil
	}
	return 0, ErrInvalidArgumentType{
		Index:    i,
		Expected: strings.Join(TNs{IntTN, CharTN}, "/"),
		Actual:   args[i].TypeName(),
	}
}

// ToByteSlice will try to convert object o to []byte value.
func ToByteSlice(i int, args ...Object) ([]byte, error) {
	switch o := args[i].(type) {
	case *Bytes:
		return o.Value, nil
	case *String:
		return []byte(o.Value), nil
	}
	return nil, ErrInvalidArgumentType{
		Index:    i,
		Expected: strings.Join(TNs{BytesTN, StringTN}, "/"),
		Actual:   args[i].TypeName(),
	}
}

// ToTime will try to convert object o to time.Time value.
func ToTime(i int, args ...Object) (time.Time, error) {
	switch o := args[i].(type) {
	case *Time:
		return o.Value, nil
	case *Int:
		return time.Unix(o.Value, 0), nil
	}
	return time.Time{}, ErrInvalidArgumentType{
		Index:    i,
		Expected: strings.Join(TNs{TimeTN, IntTN}, "/"),
		Actual:   args[i].TypeName(),
	}
}

// ToInterface attempts to convert an object o to an interface{} value
func ToInterface(o Object) interface{} {
	switch o := o.(type) {
	case *Int:
		return o.Value
	case *String:
		return o.Value
	case *Float:
		return o.Value
	case *Bool:
		return o == TrueValue
	case *Char:
		return o.Value
	case *Bytes:
		return o.Value
	case *Array:
		res := make([]interface{}, len(o.Value))
		for i, val := range o.Value {
			res[i] = ToInterface(val)
		}
		return res
	case *ImmutableArray:
		res := make([]interface{}, len(o.Value))
		for i, val := range o.Value {
			res[i] = ToInterface(val)
		}
		return res
	case *Map:
		res := make(map[string]interface{})
		for key, v := range o.Value {
			res[key] = ToInterface(v)
		}
		return res
	case *ImmutableMap:
		res := make(map[string]interface{})
		for key, v := range o.Value {
			res[key] = ToInterface(v)
		}
		return res
	case *Time:
		return o.Value
	case *Error:
		return errors.New(o.String())
	case *Undefined:
		return nil
	case Object:
		return o
	}
	return nil
}

// FromInterface will attempt to convert an interface{} v to a Tengo Object
func FromInterface(v interface{}) (Object, error) {
	switch v := v.(type) {
	case nil:
		return UndefinedValue, nil
	case string:
		if len(v) > MaxStringLen {
			return nil, ErrStringLimit
		}
		return &String{Value: v}, nil
	case int64:
		return &Int{Value: v}, nil
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
		if len(v) > MaxBytesLen {
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

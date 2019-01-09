package script

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/d5/tengo/objects"
)

func objectToString(o objects.Object) string {
	switch val := o.(type) {
	case *objects.Array:
		var s []string
		for _, e := range val.Value {
			s = append(s, objectToString(e))
		}
		return "[" + strings.Join(s, ", ") + "]"
	case *objects.Map:
		var s []string
		for k, v := range val.Value {
			s = append(s, k+": "+objectToString(v))
		}
		return "{" + strings.Join(s, ", ") + "}"
	case *objects.Int:
		return strconv.FormatInt(val.Value, 10)
	case *objects.Float:
		return strconv.FormatFloat(val.Value, 'f', -1, 64)
	case *objects.Bool:
		if val.Value {
			return "true"
		}
		return "false"
	case *objects.Char:
		return string(val.Value)
	case *objects.String:
		return val.Value
	}

	return ""
}

func objectToInterface(o objects.Object) interface{} {
	switch val := o.(type) {
	case *objects.Array:
		return val.Value
	case *objects.Map:
		return val.Value
	case *objects.Int:
		return val.Value
	case *objects.Float:
		return val.Value
	case *objects.Bool:
		return val.Value
	case *objects.Char:
		return val.Value
	case *objects.String:
		return val.Value
	}

	return nil
}

func interfaceToObject(v interface{}) (objects.Object, error) {
	switch v := v.(type) {
	case string:
		return &objects.String{Value: v}, nil
	case int64:
		return &objects.Int{Value: v}, nil
	case int:
		return &objects.Int{Value: int64(v)}, nil
	case bool:
		return &objects.Bool{Value: v}, nil
	case rune:
		return &objects.Char{Value: v}, nil
	case byte:
		return &objects.Char{Value: rune(v)}, nil
	case float64:
		return &objects.Float{Value: v}, nil
	case map[string]objects.Object:
		return &objects.Map{Value: v}, nil
	case map[string]interface{}:
		kv := make(map[string]objects.Object)
		for vk, vv := range v {
			vo, err := interfaceToObject(vv)
			if err != nil {
				return nil, err
			}

			kv[vk] = vo
		}
		return &objects.Map{Value: kv}, nil
	case []objects.Object:
		return &objects.Array{Value: v}, nil
	case []interface{}:
		arr := make([]objects.Object, len(v), len(v))
		for _, e := range v {
			vo, err := interfaceToObject(e)
			if err != nil {
				return nil, err
			}

			arr = append(arr, vo)
		}
		return &objects.Array{Value: arr}, nil
	case objects.Object:
		return v, nil
	}

	return nil, fmt.Errorf("unsupported value type: %T", v)
}

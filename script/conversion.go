package script

import (
	"fmt"

	"github.com/d5/tengo/objects"
)

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
	case *objects.Bytes:
		return val.Value
	case *objects.Undefined:
		return nil
	}

	return o
}

func interfaceToObject(v interface{}) (objects.Object, error) {
	switch v := v.(type) {
	case nil:
		return undefined, nil
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
	case []byte:
		return &objects.Bytes{Value: v}, nil
	case error:
		return &objects.Error{Value: &objects.String{Value: v.Error()}}, nil
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

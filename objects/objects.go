package objects

import "fmt"

func FromValue(v interface{}) (Object, error) {
	switch v := v.(type) {
	case string:
		return &String{Value: v}, nil
	case int64:
		return &Int{Value: v}, nil
	case int:
		return &Int{Value: int64(v)}, nil
	case bool:
		return &Bool{Value: v}, nil
	case rune:
		return &Char{Value: v}, nil
	case byte:
		return &Char{Value: rune(v)}, nil
	case float64:
		return &Float{Value: v}, nil
	case map[string]Object:
		return &Map{Value: v}, nil
	case map[string]interface{}:
		kv := make(map[string]Object)
		for vk, vv := range v {
			vo, err := FromValue(vv)
			if err != nil {
				return nil, err
			}

			kv[vk] = vo
		}
		return &Map{Value: kv}, nil
	case []Object:
		return &Array{Value: v}, nil
	case []interface{}:
		arr := make([]Object, len(v), len(v))
		for _, e := range v {
			vo, err := FromValue(e)
			if err != nil {
				return nil, err
			}

			arr = append(arr, vo)
		}
		return &Array{Value: arr}, nil
	case Object:
		return v, nil
	}

	return nil, fmt.Errorf("unsupported value type: %T", v)
}

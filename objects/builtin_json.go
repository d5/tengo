package objects

import (
	"encoding/json"
	"fmt"
)

func builtinToJSON(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	res, err := json.Marshal(objectToInterface(args[0]))
	if err != nil {
		return &Error{Value: &String{Value: err.Error()}}, nil
	}

	return &Bytes{Value: res}, nil
}

func builtinFromJSON(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	var target interface{}

	switch o := args[0].(type) {
	case *Bytes:
		err := json.Unmarshal(o.Value, &target)
		if err != nil {
			return &Error{Value: &String{Value: err.Error()}}, nil
		}
	case *String:
		err := json.Unmarshal([]byte(o.Value), &target)
		if err != nil {
			return &Error{Value: &String{Value: err.Error()}}, nil
		}
	default:
		return nil, fmt.Errorf("invalid type for argument 0: %s", o.TypeName())
	}

	res, err := FromInterface(target)
	if err != nil {
		return nil, err
	}

	return res, nil
}

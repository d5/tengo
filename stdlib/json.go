package stdlib

import (
	"encoding/json"

	"github.com/d5/tengo"
	"github.com/d5/tengo/objects"
)

var jsonModule = map[string]objects.Object{
	"parse":     &objects.UserFunction{Name: "parse", Value: jsonParse},
	"stringify": &objects.UserFunction{Name: "stringify", Value: jsonStringify},
}

func jsonParse(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}

	var target interface{}

	switch o := args[0].(type) {
	case *objects.Bytes:
		err := json.Unmarshal(o.Value, &target)
		if err != nil {
			return &objects.Error{Value: &objects.String{Value: err.Error()}}, nil
		}
	case *objects.String:
		err := json.Unmarshal([]byte(o.Value), &target)
		if err != nil {
			return &objects.Error{Value: &objects.String{Value: err.Error()}}, nil
		}
	default:
		return nil, objects.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "bytes/string",
			Found:    args[0].TypeName(),
		}
	}

	res, err := objects.FromInterface(target)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func jsonStringify(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}

	v := objects.ToInterface(args[0])
	if vErr, isErr := v.(error); isErr {
		v = vErr.Error()
	}

	res, err := json.Marshal(v)
	if err != nil {
		return &objects.Error{Value: &objects.String{Value: err.Error()}}, nil
	}

	if len(res) > tengo.MaxBytesLen {
		return nil, objects.ErrBytesLimit
	}

	return &objects.String{Value: string(res)}, nil
}

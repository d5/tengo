package objects

import (
	"encoding/json"
)

func builtinToJSON(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	res, err := json.Marshal(objectToInterface(args[0]))
	if err != nil {
		return nil, err
	}

	return &Bytes{Value: res}, nil
}

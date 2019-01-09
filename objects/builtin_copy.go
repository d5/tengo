package objects

import "errors"

func builtinCopy(args ...Object) (Object, error) {
	// TODO: should multi arguments later?
	if len(args) != 1 {
		return nil, errors.New("wrong number of arguments")
	}

	return args[0].Copy(), nil
}

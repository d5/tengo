package objects

import (
	"errors"
	"strconv"
)

func builtinString(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, errors.New("wrong number of arguments")
	}

	switch arg := args[0].(type) {
	case *String:
		return arg, nil
	case *Undefined:
		return UndefinedValue, nil
	default:
		return &String{Value: arg.String()}, nil
	}
}

func builtinInt(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, errors.New("wrong number of arguments")
	}

	switch arg := args[0].(type) {
	case *Int:
		return arg, nil
	case *Float:
		return &Int{Value: int64(arg.Value)}, nil
	case *Char:
		return &Int{Value: int64(arg.Value)}, nil
	case *Bool:
		if arg.Value {
			return &Int{Value: 1}, nil
		}
		return &Int{Value: 0}, nil
	case *String:
		n, err := strconv.ParseInt(arg.Value, 10, 64)
		if err == nil {
			return &Int{Value: n}, nil
		}
	}

	return UndefinedValue, nil
}

func builtinFloat(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, errors.New("wrong number of arguments")
	}

	switch arg := args[0].(type) {
	case *Float:
		return arg, nil
	case *Int:
		return &Float{Value: float64(arg.Value)}, nil
	case *String:
		f, err := strconv.ParseFloat(arg.Value, 64)
		if err == nil {
			return &Float{Value: f}, nil
		}
	}

	return UndefinedValue, nil
}

func builtinBool(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, errors.New("wrong number of arguments")
	}

	switch arg := args[0].(type) {
	case *Bool:
		return arg, nil
	default:
		return &Bool{Value: !arg.IsFalsy()}, nil
	}
}

func builtinChar(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, errors.New("wrong number of arguments")
	}

	switch arg := args[0].(type) {
	case *Char:
		return arg, nil
	case *Int:
		return &Char{Value: rune(arg.Value)}, nil
	case *String:
		rs := []rune(arg.Value)
		switch len(rs) {
		case 0:
			return &Char{}, nil
		case 1:
			return &Char{Value: rs[0]}, nil
		}
	}

	return UndefinedValue, nil
}

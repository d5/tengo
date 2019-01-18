package objects

func builtinString(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*String); ok {
		return args[0], nil
	}

	v, ok := ToString(args[0])
	if ok {
		return &String{Value: v}, nil
	}

	return UndefinedValue, nil
}

func builtinInt(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Int); ok {
		return args[0], nil
	}

	v, ok := ToInt64(args[0])
	if ok {
		return &Int{Value: v}, nil
	}

	return UndefinedValue, nil
}

func builtinFloat(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Float); ok {
		return args[0], nil
	}

	v, ok := ToFloat64(args[0])
	if ok {
		return &Float{Value: v}, nil
	}

	return UndefinedValue, nil
}

func builtinBool(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Bool); ok {
		return args[0], nil
	}

	v, ok := ToBool(args[0])
	if ok {
		return &Bool{Value: v}, nil
	}

	return UndefinedValue, nil
}

func builtinChar(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Char); ok {
		return args[0], nil
	}

	v, ok := ToRune(args[0])
	if ok {
		return &Char{Value: v}, nil
	}

	return UndefinedValue, nil
}

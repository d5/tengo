package objects

func builtinIsString(rt Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*String); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsInt(rt Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Int); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsFloat(rt Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Float); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsBool(rt Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Bool); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsChar(rt Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Char); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsBytes(rt Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Bytes); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsArray(rt Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Array); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsImmutableArray(rt Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*ImmutableArray); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsMap(rt Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Map); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsImmutableMap(rt Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*ImmutableMap); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsTime(rt Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Time); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsError(rt Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Error); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsUndefined(rt Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if args[0] == UndefinedValue {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsFunction(rt Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	switch args[0].(type) {
	case *CompiledFunction:
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsCallable(rt Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if args[0].CanCall() {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsIterable(rt Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if args[0].CanIterate() {
		return TrueValue, nil
	}

	return FalseValue, nil
}

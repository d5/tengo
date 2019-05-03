package objects

func builtinIsString(_ RuntimeHooks, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*String); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsInt(_ RuntimeHooks, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Int); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsFloat(_ RuntimeHooks, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Float); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsBool(_ RuntimeHooks, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Bool); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsChar(_ RuntimeHooks, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Char); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsBytes(_ RuntimeHooks, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Bytes); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsArray(_ RuntimeHooks, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Array); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsImmutableArray(_ RuntimeHooks, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*ImmutableArray); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsMap(_ RuntimeHooks, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Map); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsImmutableMap(_ RuntimeHooks, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*ImmutableMap); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsTime(_ RuntimeHooks, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Time); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsError(_ RuntimeHooks, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Error); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsUndefined(_ RuntimeHooks, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if args[0] == UndefinedValue {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsFunction(_ RuntimeHooks, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	switch args[0].(type) {
	case *CompiledFunction:
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsCallable(_ RuntimeHooks, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if args[0].CanCall() {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsIterable(_ RuntimeHooks, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if args[0].CanIterate() {
		return TrueValue, nil
	}

	return FalseValue, nil
}

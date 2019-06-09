package tengo

func builtinCopy(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	return args[0].Copy(), nil
}

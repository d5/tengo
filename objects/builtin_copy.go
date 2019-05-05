package objects

func builtinCopy(rt Runtime, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	return args[0].Copy(), nil
}

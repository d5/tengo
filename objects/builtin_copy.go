package objects

func builtinCopy(_ RuntimeHooks, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	return args[0].Copy(), nil
}

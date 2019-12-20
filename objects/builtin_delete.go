package objects

// map delete
func builtinDelete(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, ErrWrongNumArguments
	}

	omap, ok := args[0].(*Map)
	if !ok {
		return nil, ErrInvalidArgumentType{
			Name:     "delete",
			Expected: "map",
			Found:    args[0].TypeName(),
		}
	}

	key, ok := ToString(args[1])
	if !ok {
		return nil, ErrInvalidIndexType
	}

	delete(omap.Value, key)
	return omap, nil
}

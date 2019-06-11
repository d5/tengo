package tengo

// Builtins contains all default builtin functions.
// Use GetBuiltinFunctions instead of accessing Builtins directly.
var Builtins = []*BuiltinFunction{
	{Name: "len", Value: builtinLen},
	{Name: "copy", Value: builtinCopy},
	{Name: "append", Value: builtinAppend},
	{Name: "string", Value: builtinString},
	{Name: "int", Value: builtinInt},
	{Name: "bool", Value: builtinBool},
	{Name: "float", Value: builtinFloat},
	{Name: "char", Value: builtinChar},
	{Name: "bytes", Value: builtinBytes},
	{Name: "time", Value: builtinTime},
	{Name: "is_int", Value: builtinIsInt},
	{Name: "is_float", Value: builtinIsFloat},
	{Name: "is_string", Value: builtinIsString},
	{Name: "is_bool", Value: builtinIsBool},
	{Name: "is_char", Value: builtinIsChar},
	{Name: "is_bytes", Value: builtinIsBytes},
	{Name: "is_array", Value: builtinIsArray},
	{Name: "is_immutable_array", Value: builtinIsImmutableArray},
	{Name: "is_map", Value: builtinIsMap},
	{Name: "is_immutable_map", Value: builtinIsImmutableMap},
	{Name: "is_iterable", Value: builtinIsIterable},
	{Name: "is_time", Value: builtinIsTime},
	{Name: "is_error", Value: builtinIsError},
	{Name: "is_undefined", Value: builtinIsUndefined},
	{Name: "is_function", Value: builtinIsFunction},
	{Name: "is_callable", Value: builtinIsCallable},
	{Name: "type_name", Value: builtinTypeName},
	{Name: "format", Value: builtinFormat},
}

// len(obj object) => int
func builtinLen(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	switch arg := args[0].(type) {
	case *Array:
		return &Int{Value: int64(len(arg.Value))}, nil
	case *ImmutableArray:
		return &Int{Value: int64(len(arg.Value))}, nil
	case *String:
		return &Int{Value: int64(len(arg.Value))}, nil
	case *Bytes:
		return &Int{Value: int64(len(arg.Value))}, nil
	case *Map:
		return &Int{Value: int64(len(arg.Value))}, nil
	case *ImmutableMap:
		return &Int{Value: int64(len(arg.Value))}, nil
	default:
		return nil, ErrInvalidArgumentType{
			Name:     "first",
			Expected: "array/string/bytes/map",
			Found:    arg.TypeName(),
		}
	}
}

// append(arr, items...)
func builtinAppend(_ Interop, args ...Object) (Object, error) {
	if len(args) < 2 {
		return nil, ErrWrongNumArguments
	}

	switch arg := args[0].(type) {
	case *Array:
		return &Array{Value: append(arg.Value, args[1:]...)}, nil
	case *ImmutableArray:
		return &Array{Value: append(arg.Value, args[1:]...)}, nil
	default:
		return nil, ErrInvalidArgumentType{
			Name:     "first",
			Expected: "array",
			Found:    arg.TypeName(),
		}
	}
}

func builtinString(_ Interop, args ...Object) (Object, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*String); ok {
		return args[0], nil
	}

	v, ok := ToString(args[0])
	if ok {
		if len(v) > MaxStringLen {
			return nil, ErrStringLimit
		}

		return &String{Value: v}, nil
	}

	if argsLen == 2 {
		return args[1], nil
	}

	return UndefinedValue, nil
}

func builtinInt(_ Interop, args ...Object) (Object, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Int); ok {
		return args[0], nil
	}

	v, ok := ToInt64(args[0])
	if ok {
		return &Int{Value: v}, nil
	}

	if argsLen == 2 {
		return args[1], nil
	}

	return UndefinedValue, nil
}

func builtinFloat(_ Interop, args ...Object) (Object, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Float); ok {
		return args[0], nil
	}

	v, ok := ToFloat64(args[0])
	if ok {
		return &Float{Value: v}, nil
	}

	if argsLen == 2 {
		return args[1], nil
	}

	return UndefinedValue, nil
}

func builtinBool(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Bool); ok {
		return args[0], nil
	}

	v, ok := ToBool(args[0])
	if ok {
		if v {
			return TrueValue, nil
		}

		return FalseValue, nil
	}

	return UndefinedValue, nil
}

func builtinChar(_ Interop, args ...Object) (Object, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Char); ok {
		return args[0], nil
	}

	v, ok := ToRune(args[0])
	if ok {
		return &Char{Value: v}, nil
	}

	if argsLen == 2 {
		return args[1], nil
	}

	return UndefinedValue, nil
}

func builtinBytes(_ Interop, args ...Object) (Object, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, ErrWrongNumArguments
	}

	// bytes(N) => create a new bytes with given size N
	if n, ok := args[0].(*Int); ok {
		if n.Value > int64(MaxBytesLen) {
			return nil, ErrBytesLimit
		}

		return &Bytes{Value: make([]byte, int(n.Value))}, nil
	}

	v, ok := ToByteSlice(args[0])
	if ok {
		if len(v) > MaxBytesLen {
			return nil, ErrBytesLimit
		}

		return &Bytes{Value: v}, nil
	}

	if argsLen == 2 {
		return args[1], nil
	}

	return UndefinedValue, nil
}

func builtinTime(_ Interop, args ...Object) (Object, error) {
	argsLen := len(args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Time); ok {
		return args[0], nil
	}

	v, ok := ToTime(args[0])
	if ok {
		return &Time{Value: v}, nil
	}

	if argsLen == 2 {
		return args[1], nil
	}

	return UndefinedValue, nil
}

func builtinCopy(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	return args[0].Copy(), nil
}

func builtinFormat(_ Interop, args ...Object) (Object, error) {
	numArgs := len(args)
	if numArgs == 0 {
		return nil, ErrWrongNumArguments
	}

	format, ok := args[0].(*String)
	if !ok {
		return nil, ErrInvalidArgumentType{
			Name:     "format",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}
	if numArgs == 1 {
		// okay to return 'format' directly as String is immutable
		return format, nil
	}

	s, err := Format(format.Value, args[1:]...)
	if err != nil {
		return nil, err
	}

	return &String{Value: s}, nil
}

func builtinTypeName(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	return &String{Value: args[0].TypeName()}, nil
}

func builtinIsString(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*String); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsInt(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Int); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsFloat(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Float); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsBool(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Bool); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsChar(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Char); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsBytes(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Bytes); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsArray(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Array); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsImmutableArray(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*ImmutableArray); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsMap(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Map); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsImmutableMap(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*ImmutableMap); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsTime(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Time); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsError(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if _, ok := args[0].(*Error); ok {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsUndefined(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if args[0] == UndefinedValue {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsFunction(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	switch args[0].(type) {
	case *CompiledFunction:
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsCallable(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if args[0].CanCall() {
		return TrueValue, nil
	}

	return FalseValue, nil
}

func builtinIsIterable(_ Interop, args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}

	if args[0].CanIterate() {
		return TrueValue, nil
	}

	return FalseValue, nil
}

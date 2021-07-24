package tengo

import "fmt"

var builtinFuncs = []*BuiltinFunction{
	{
		Name:  "len",
		Value: builtinLen,
	},
	{
		Name:  "copy",
		Value: builtinCopy,
	},
	{
		Name:  "append",
		Value: builtinAppend,
	},
	{
		Name:  "delete",
		Value: builtinDelete,
	},
	{
		Name:  "splice",
		Value: builtinSplice,
	},
	{
		Name:  "string",
		Value: builtinString,
	},
	{
		Name:  "int",
		Value: builtinInt,
	},
	{
		Name:  "bool",
		Value: builtinBool,
	},
	{
		Name:  "float",
		Value: builtinFloat,
	},
	{
		Name:  "char",
		Value: builtinChar,
	},
	{
		Name:  "bytes",
		Value: builtinBytes,
	},
	{
		Name:  "time",
		Value: builtinTime,
	},
	{
		Name:  "is_int",
		Value: builtinIsInt,
	},
	{
		Name:  "is_float",
		Value: builtinIsFloat,
	},
	{
		Name:  "is_string",
		Value: builtinIsString,
	},
	{
		Name:  "is_bool",
		Value: builtinIsBool,
	},
	{
		Name:  "is_char",
		Value: builtinIsChar,
	},
	{
		Name:  "is_bytes",
		Value: builtinIsBytes,
	},
	{
		Name:  "is_array",
		Value: builtinIsArray,
	},
	{
		Name:  "is_immutable_array",
		Value: builtinIsImmutableArray,
	},
	{
		Name:  "is_map",
		Value: builtinIsMap,
	},
	{
		Name:  "is_immutable_map",
		Value: builtinIsImmutableMap,
	},
	{
		Name:  "is_iterable",
		Value: builtinIsIterable,
	},
	{
		Name:  "is_time",
		Value: builtinIsTime,
	},
	{
		Name:  "is_error",
		Value: builtinIsError,
	},
	{
		Name:  "is_undefined",
		Value: builtinIsUndefined,
	},
	{
		Name:  "is_function",
		Value: builtinIsFunction,
	},
	{
		Name:  "is_callable",
		Value: builtinIsCallable,
	},
	{
		Name:  "type_name",
		Value: builtinTypeName,
	},
	{
		Name:  "format",
		Value: builtinFormat,
	},
	{
		Name:  "range",
		Value: builtinRange,
	},
}

// GetAllBuiltinFunctions returns all builtin function objects.
func GetAllBuiltinFunctions() []*BuiltinFunction {
	return append([]*BuiltinFunction{}, builtinFuncs...)
}

var builtinTypeName = CheckAnyArgs(func(args ...Object) (Object, error) {
	return &String{Value: args[0].TypeName()}, nil
}, 1)

func builtinIsType(typeCheck func(arg Object) bool, args ...Object) CallableFunc {
	return CheckAnyArgs(func(args ...Object) (Object, error) {
		if typeCheck(args[0]) {
			return TrueValue, nil
		}
		return FalseValue, nil
	}, 1)
}

var builtinIsString = builtinIsType(func(arg Object) bool {
	_, ok := arg.(*String)
	return ok
})

var builtinIsInt = builtinIsType(func(arg Object) bool {
	_, ok := arg.(*Int)
	return ok
})

var builtinIsFloat = builtinIsType(func(arg Object) bool {
	_, ok := arg.(*Float)
	return ok
})

var builtinIsBool = builtinIsType(func(arg Object) bool {
	_, ok := arg.(*Bool)
	return ok
})

var builtinIsChar = builtinIsType(func(arg Object) bool {
	_, ok := arg.(*Char)
	return ok
})

var builtinIsBytes = builtinIsType(func(arg Object) bool {
	_, ok := arg.(*Bytes)
	return ok
})

var builtinIsArray = builtinIsType(func(arg Object) bool {
	_, ok := arg.(*Array)
	return ok
})

var builtinIsImmutableArray = builtinIsType(func(arg Object) bool {
	_, ok := arg.(*ImmutableArray)
	return ok
})

var builtinIsMap = builtinIsType(func(arg Object) bool {
	_, ok := arg.(*Map)
	return ok
})

var builtinIsImmutableMap = builtinIsType(func(arg Object) bool {
	_, ok := arg.(*ImmutableMap)
	return ok
})

var builtinIsTime = builtinIsType(func(arg Object) bool {
	_, ok := arg.(*Time)
	return ok
})

var builtinIsError = builtinIsType(func(arg Object) bool {
	_, ok := arg.(*Error)
	return ok
})

var builtinIsUndefined = builtinIsType(func(arg Object) bool {
	_, ok := arg.(*Undefined)
	return ok
})

var builtinIsFunction = builtinIsType(func(arg Object) bool {
	_, ok := arg.(*CompiledFunction)
	return ok
})

var builtinIsCallable = builtinIsType(func(arg Object) bool {
	return arg.CanCall()
})

var builtinIsIterable = builtinIsType(func(arg Object) bool {
	return arg.CanIterate()
})

var builtinLen = CheckAnyArgs(func(args ...Object) (Object, error) {
	if !args[0].HasLen() {
		return nil, fmt.Errorf("arg type %s does not have a length value", args[0].TypeName())
	}
	return &Int{Value: int64(args[0].Len())}, nil
}, 1)

var builtinRange = CheckOptArgs(func(args ...Object) (Object, error) {
	start := args[0].(*Int)
	stop := args[1].(*Int)
	step := &Int{Value: int64(1)}
	if len(args) == 3 {
		step = args[2].(*Int)
		if step.Value <= 0 {
			return nil, ErrInvalidRangeStep
		}
	}
	return buildRange(start.Value, stop.Value, step.Value), nil
}, 2, 3, IntTN, IntTN, IntTN)

func buildRange(start, stop, step int64) *Array {
	array := &Array{}
	if start <= stop {
		for i := start; i < stop; i += step {
			array.Value = append(array.Value, &Int{
				Value: i,
			})
		}
	} else {
		for i := start; i > stop; i -= step {
			array.Value = append(array.Value, &Int{
				Value: i,
			})
		}
	}
	return array
}

var builtinFormat = CheckOptArgs(func(args ...Object) (Object, error) {
	format := args[0].(*String)
	if len(args) == 1 {
		// okay to return 'format' directly as String is immutable
		return format, nil
	}
	s, err := Format(format.Value, args[1:]...)
	if err != nil {
		return nil, err
	}
	return &String{Value: s}, nil
}, 1, -1, StringTN, AnyTN)

var builtinCopy = CheckAnyArgs(func(args ...Object) (Object, error) {
	return args[0].Copy(), nil
}, 1)

var builtinString = CheckAnyArgs(func(args ...Object) (Object, error) {
	if _, ok := args[0].(*String); ok {
		return args[0], nil
	}
	v, err := ToString(0, args...)
	if err == nil {
		if len(v) > MaxStringLen {
			return nil, ErrStringLimit
		}
		return &String{Value: v}, nil
	}
	if len(args) == 2 {
		return args[1], nil
	}
	return UndefinedValue, nil
}, 1, 2)

var builtinInt = CheckAnyArgs(func(args ...Object) (Object, error) {
	if _, ok := args[0].(*Int); ok {
		return args[0], nil
	}
	v, err := ToInt64(0, args...)
	if err == nil {
		return &Int{Value: v}, nil
	}
	if len(args) == 2 {
		return args[1], nil
	}
	return UndefinedValue, nil
}, 1, 2)

var builtinFloat = CheckAnyArgs(func(args ...Object) (Object, error) {
	if _, ok := args[0].(*Float); ok {
		return args[0], nil
	}
	v, err := ToFloat64(0, args...)
	if err == nil {
		return &Float{Value: v}, nil
	}
	if len(args) == 2 {
		return args[1], nil
	}
	return UndefinedValue, nil
}, 1, 2)

var builtinBool = CheckAnyArgs(func(args ...Object) (Object, error) {
	if _, ok := args[0].(*Bool); ok {
		return args[0], nil
	}
	v := ToBool(0, args...)
	if v {
		return TrueValue, nil
	}
	return FalseValue, nil
}, 1)

var builtinChar = CheckAnyArgs(func(args ...Object) (Object, error) {
	if _, ok := args[0].(*Char); ok {
		return args[0], nil
	}
	v, err := ToRune(0, args...)
	if err == nil {
		return &Char{Value: v}, nil
	}
	if len(args) == 2 {
		return args[1], nil
	}
	return UndefinedValue, nil
}, 1, 2)

var builtinBytes = CheckAnyArgs(func(args ...Object) (Object, error) {
	// bytes(N) => create a new bytes with given size N
	if n, ok := args[0].(*Int); ok {
		if n.Value > int64(MaxBytesLen) {
			return nil, ErrBytesLimit
		}
		return &Bytes{Value: make([]byte, int(n.Value))}, nil
	}
	v, err := ToByteSlice(0, args...)
	if err == nil {
		if len(v) > MaxBytesLen {
			return nil, ErrBytesLimit
		}
		return &Bytes{Value: v}, nil
	}
	if len(args) == 2 {
		return args[1], nil
	}
	return UndefinedValue, nil
}, 1, 2)

var builtinTime = CheckAnyArgs(func(args ...Object) (Object, error) {
	if _, ok := args[0].(*Time); ok {
		return args[0], nil
	}
	v, err := ToTime(0, args...)
	if err == nil {
		return &Time{Value: v}, nil
	}
	if len(args) == 2 {
		return args[1], nil
	}
	return UndefinedValue, nil
}, 1, 2)

// append(arr, items...)
var builtinAppend = CheckArgs(func(args ...Object) (Object, error) {
	switch arg := args[0].(type) {
	case *Array:
		return &Array{Value: append(arg.Value, args[1:]...)}, nil
	case *ImmutableArray:
		return &Array{Value: append(arg.Value, args[1:]...)}, nil
	default:
		panic("impossible")
	}
}, 2, -1, TNs{ArrayTN, ImmutableArrayTN}, nil)

// builtinDelete deletes Map keys
// usage: delete(map, "key")
// key must be a string
var builtinDelete = CheckStrictArgs(func(args ...Object) (Object, error) {
	delete(args[0].(*Map).Value, args[1].(*String).Value)
	return UndefinedValue, nil
}, MapTN, StringTN)

// builtinSplice deletes and changes given Array, returns deleted items.
// usage:
// deleted_items := splice(array[,start[,delete_count[,item1[,item2[,...]]]])
var builtinSplice = CheckOptArgs(func(args ...Object) (Object, error) {
	argsLen := len(args)
	array := args[0].(*Array)
	arrayLen := len(array.Value)

	var startIdx int
	if argsLen > 1 {
		arg1 := args[1].(*Int)
		startIdx = int(arg1.Value)
		if startIdx < 0 || startIdx > arrayLen {
			return nil, ErrIndexOutOfBounds
		}
	}

	delCount := len(array.Value)
	if argsLen > 2 {
		arg2 := args[2].(*Int)
		delCount = int(arg2.Value)
		if delCount < 0 {
			return nil, ErrIndexOutOfBounds
		}
	}
	// if count of to be deleted items is bigger than expected, truncate it
	if startIdx+delCount > arrayLen {
		delCount = arrayLen - startIdx
	}
	// delete items
	endIdx := startIdx + delCount
	deleted := append([]Object{}, array.Value[startIdx:endIdx]...)

	head := array.Value[:startIdx]
	var items []Object
	if argsLen > 3 {
		items = make([]Object, 0, argsLen-3)
		for i := 3; i < argsLen; i++ {
			items = append(items, args[i])
		}
	}
	items = append(items, array.Value[endIdx:]...)
	array.Value = append(head, items...)

	// return deleted items
	return &Array{Value: deleted}, nil
}, 1, -1, ArrayTN, IntTN, IntTN, AnyTN)

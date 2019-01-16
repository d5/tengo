package objects

// BuiltinFunc is a function signature for the builtin functions.
type BuiltinFunc func(args ...Object) (ret Object, err error)

// Builtins contains all known builtin functions.
var Builtins = []struct {
	Name string
	Func BuiltinFunc
}{
	{
		Name: "print",
		Func: builtinPrint,
	},
	{
		Name: "len",
		Func: builtinLen,
	},
	{
		Name: "copy",
		Func: builtinCopy,
	},
	{
		Name: "append",
		Func: builtinAppend,
	},
	{
		Name: "string",
		Func: builtinString,
	},
	{
		Name: "int",
		Func: builtinInt,
	},
	{
		Name: "bool",
		Func: builtinBool,
	},
	{
		Name: "float",
		Func: builtinFloat,
	},
	{
		Name: "char",
		Func: builtinChar,
	},
	{
		Name: "is_error",
		Func: builtinIsError,
	},
	{
		Name: "is_undefined",
		Func: builtinIsUndefined,
	},
}

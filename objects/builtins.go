package objects

type BuiltinFunc func(args ...Object) (ret Object, err error)

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
}

package stdlib

import (
	"fmt"

	"github.com/d5/tengo/v2"
)

var fmtModule = map[string]tengo.Object{
	"print":   &tengo.BuiltinFunction{Value: fmtPrint, NeedVMObj: true},
	"printf":  &tengo.BuiltinFunction{Value: fmtPrintf, NeedVMObj: true},
	"println": &tengo.BuiltinFunction{Value: fmtPrintln, NeedVMObj: true},
	"sprintf": &tengo.UserFunction{Name: "sprintf", Value: fmtSprintf},
}

func fmtPrint(args ...tengo.Object) (ret tengo.Object, err error) {
	vm := args[0].(*tengo.VMObj).Value
	args = args[1:] // the first arg is VMObj inserted by VM
	printArgs, err := getPrintArgs(args...)
	if err != nil {
		return nil, err
	}
	fmt.Fprint(vm.Out, printArgs...)
	return nil, nil
}

func fmtPrintf(args ...tengo.Object) (ret tengo.Object, err error) {
	vm := args[0].(*tengo.VMObj).Value
	args = args[1:] // the first arg is VMObj inserted by VM
	numArgs := len(args)
	if numArgs == 0 {
		return nil, tengo.ErrWrongNumArguments
	}

	format, ok := args[0].(*tengo.String)
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "format",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}
	if numArgs == 1 {
		fmt.Fprint(vm.Out, format)
		return nil, nil
	}

	s, err := tengo.Format(format.Value, args[1:]...)
	if err != nil {
		return nil, err
	}
	fmt.Fprint(vm.Out, s)
	return nil, nil
}

func fmtPrintln(args ...tengo.Object) (ret tengo.Object, err error) {
	vm := args[0].(*tengo.VMObj).Value
	args = args[1:] // the first arg is VMObj inserted by VM
	printArgs, err := getPrintArgs(args...)
	if err != nil {
		return nil, err
	}
	printArgs = append(printArgs, "\n")
	fmt.Fprint(vm.Out, printArgs...)
	return nil, nil
}

func fmtSprintf(args ...tengo.Object) (ret tengo.Object, err error) {
	numArgs := len(args)
	if numArgs == 0 {
		return nil, tengo.ErrWrongNumArguments
	}

	format, ok := args[0].(*tengo.String)
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "format",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}
	if numArgs == 1 {
		// okay to return 'format' directly as String is immutable
		return format, nil
	}
	s, err := tengo.Format(format.Value, args[1:]...)
	if err != nil {
		return nil, err
	}
	return &tengo.String{Value: s}, nil
}

func getPrintArgs(args ...tengo.Object) ([]interface{}, error) {
	var printArgs []interface{}
	l := 0
	for _, arg := range args {
		s, _ := tengo.ToString(arg)
		slen := len(s)
		// make sure length does not exceed the limit
		if l+slen > tengo.MaxStringLen {
			return nil, tengo.ErrStringLimit
		}
		l += slen
		printArgs = append(printArgs, s)
	}
	return printArgs, nil
}

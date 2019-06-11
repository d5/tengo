package stdlib

import (
	"fmt"

	"github.com/d5/tengo"
)

var fmtModule = map[string]tengo.Object{
	"print":   &tengo.GoFunction{Name: "print", Value: fmtPrint},
	"printf":  &tengo.GoFunction{Name: "printf", Value: fmtPrintf},
	"println": &tengo.GoFunction{Name: "println", Value: fmtPrintln},
	"sprintf": &tengo.GoFunction{Name: "sprintf", Value: fmtSprintf},
}

func fmtPrint(_ tengo.Interop, args ...tengo.Object) (ret tengo.Object, err error) {
	printArgs, err := getPrintArgs(args...)
	if err != nil {
		return nil, err
	}

	_, _ = fmt.Print(printArgs...)

	return nil, nil
}

func fmtPrintf(_ tengo.Interop, args ...tengo.Object) (ret tengo.Object, err error) {
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
		fmt.Print(format)
		return nil, nil
	}

	s, err := tengo.Format(format.Value, args[1:]...)
	if err != nil {
		return nil, err
	}

	fmt.Print(s)

	return nil, nil
}

func fmtPrintln(_ tengo.Interop, args ...tengo.Object) (ret tengo.Object, err error) {
	printArgs, err := getPrintArgs(args...)
	if err != nil {
		return nil, err
	}

	printArgs = append(printArgs, "\n")
	_, _ = fmt.Print(printArgs...)

	return nil, nil
}

func fmtSprintf(_ tengo.Interop, args ...tengo.Object) (ret tengo.Object, err error) {
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
		return format, nil // okay to return 'format' directly as String is immutable
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
		if l+slen > tengo.MaxStringLen { // make sure length does not exceed the limit
			return nil, tengo.ErrStringLimit
		}
		l += slen

		printArgs = append(printArgs, s)
	}

	return printArgs, nil
}

package stdlib

import (
	"fmt"

	"github.com/d5/tengo/v2"
)

var fmtModule = map[string]tengo.Object{
	"print":   &tengo.UserFunction{Name: "print", Value: fmtPrint},
	"printf":  &tengo.UserFunction{Name: "printf", Value: fmtPrintf},
	"println": &tengo.UserFunction{Name: "println", Value: fmtPrintln},
	"sprintf": &tengo.UserFunction{Name: "sprintf", Value: fmtSprintf},
}

func fmtPrint(args ...tengo.Object) (tengo.Object, error) {
	printArgs, err := getPrintArgs(args...)
	if err != nil {
		return nil, err
	}
	_, _ = fmt.Print(printArgs...)
	return nil, nil
}

var fmtPrintf = tengo.CheckOptArgs(func(args ...tengo.Object) (tengo.Object, error) {
	numArgs := len(args)

	format := args[0].(*tengo.String)
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
}, 1, -1, tengo.StringTN, tengo.AnyTN)

func fmtPrintln(args ...tengo.Object) (tengo.Object, error) {
	printArgs, err := getPrintArgs(args...)
	if err != nil {
		return nil, err
	}
	printArgs = append(printArgs, "\n")
	_, _ = fmt.Print(printArgs...)
	return nil, nil
}

var fmtSprintf = tengo.CheckOptArgs(func(args ...tengo.Object) (tengo.Object, error) {
	numArgs := len(args)

	format := args[0].(*tengo.String)

	if numArgs == 1 {
		// okay to return 'format' directly as String is immutable
		return format, nil
	}
	s, err := tengo.Format(format.Value, args[1:]...)
	if err != nil {
		return nil, err
	}
	return &tengo.String{Value: s}, nil
}, 1, -1, tengo.StringTN, tengo.AnyTN)

func getPrintArgs(args ...tengo.Object) ([]interface{}, error) {
	var printArgs []interface{}
	l := 0
	for i := range args {
		s, _ := tengo.ToString(i, args...)
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

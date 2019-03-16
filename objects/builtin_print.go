package objects

import (
	"fmt"

	"github.com/d5/tengo"
)

func getPrintArgs(args ...Object) ([]interface{}, error) {
	var printArgs []interface{}
	l := 0
	for _, arg := range args {
		s, _ := ToString(arg)
		slen := len(s)
		if l+slen > tengo.MaxStringLen { // make sure length does not exceed the limit
			return nil, ErrStringLimit
		}
		l += slen

		printArgs = append(printArgs, s)
	}

	return printArgs, nil
}

// print(args...)
func builtinPrint(args ...Object) (Object, error) {
	printArgs, err := getPrintArgs(args...)
	
	if err != nil {
		return nil, err
	}

	_, _ = fmt.Print(printArgs...)

	return nil, nil
}

// println(args...)
func builtinPrintln(args ...Object) (Object, error) {
	printArgs, err := getPrintArgs(args...)
	
	if err != nil {
		return nil, err
	}

	_, _ = fmt.Println(printArgs...)

	return nil, nil
}

// printf("format", args...)
func builtinPrintf(args ...Object) (Object, error) {
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
		fmt.Print(format)
		return nil, nil
	}

	formatArgs := make([]interface{}, numArgs-1, numArgs-1)
	for idx, arg := range args[1:] {
		formatArgs[idx] = ToInterface(arg)
	}

	fmt.Printf(format.Value, formatArgs...)

	return nil, nil
}

// sprintf("format", args...)
func builtinSprintf(args ...Object) (Object, error) {
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
		return format, nil // okay to return 'format' directly as String is immutable
	}

	formatArgs := make([]interface{}, numArgs-1, numArgs-1)
	for idx, arg := range args[1:] {
		formatArgs[idx] = ToInterface(arg)
	}

	s := fmt.Sprintf(format.Value, formatArgs...)

	if len(s) > tengo.MaxStringLen {
		return nil, ErrStringLimit
	}

	return &String{Value: s}, nil
}

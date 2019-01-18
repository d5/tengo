package objects

import (
	"fmt"
)

func builtinPrint(args ...Object) (Object, error) {
	for _, arg := range args {
		if str, ok := arg.(*String); ok {
			fmt.Println(str.Value)
		} else {
			fmt.Println(arg.String())
		}
	}

	return nil, nil
}

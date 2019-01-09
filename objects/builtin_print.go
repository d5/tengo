package objects

import (
	"fmt"
)

func builtinPrint(args ...Object) (Object, error) {
	for _, arg := range args {
		fmt.Println(arg.String())
	}

	return nil, nil
}

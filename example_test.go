package tengo_test

import (
	"context"
	"fmt"

	"github.com/d5/tengo/v2"
)

func Example() {
	// Tengo script code
	src := `
each := func(seq, fn) {
    for x in seq { fn(x) }
}

sum := 0
mul := 1
each([a, b, c, d], func(x) {
	sum += x
	mul *= x
})`

	// create a new Script instance
	script := tengo.NewScript([]byte(src))

	// set values
	_ = script.Add("a", 1)
	_ = script.Add("b", 9)
	_ = script.Add("c", 8)
	_ = script.Add("d", 4)

	// run the script
	compiled, err := script.RunContext(context.Background())
	if err != nil {
		panic(err)
	}

	// retrieve values
	sum := compiled.Get("sum")
	mul := compiled.Get("mul")
	fmt.Println(sum, mul)

	// Output:
	// 22 288
}

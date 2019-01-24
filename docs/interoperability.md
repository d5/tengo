# Interoperability

## Table of Contents 

- [Scripts](#scripts)
- [Compiler and VM](#compiler-and-vm)
- [Sandbox Environments](#sandbox-environments)
- [Type Conversion](#type-conversion)

To execute Tengo code in your Go codebase, you should use **Script**. In the simple use cases, all you need is to do is to create a new Script instance and call its `Script.Run()` function.

```golang
import "github.com/d5/tengo/script"

var code = `
reduce := func(seq, fn) {
    s := 0
    for x in seq { fn(x, s) }
    return s
}

print(reduce([1, 2, 3], func(x, s) { s += x }))
`

func main() {
    s := script.New([]byte(code))
    if _, err := s.Run(); err != nil {
        panic(err)
    }
}
```

If you want to compile the source script once and execute it multiple times, you can use `Script.Compile()` function that returns **Compiled** instance.

```golang
import (
	"fmt"

	"github.com/d5/tengo/script"
)

func main() {
	s := script.New([]byte(`a := b + 20`))

	// define variable 'b'
	_ = s.Add("b", 10)

	// compile the source
	c, err := s.Compile()
	if err != nil {
		panic(err)
	}

	// run the compiled bytecode
	// a compiled bytecode 'c' can be executed multiple without re-compiling it
	if err := c.Run(); err != nil {
		panic(err)
	}

	// retrieve value of 'a'
	a := c.Get("a")
	fmt.Println(a.Int())
}
```

In the example above, a variable `b` is defined by the user before compilation using `Script.Add()` function. Then a compiled bytecode `c` is used to execute the bytecode and get the value of global variables. In this example, the value of global variable `a` is read using `Compiled.Get()` function.

One can easily use the custom data types by implementing `objects.Object` interface. See [Interoperability](https://github.com/d5/tengo/wiki/Interoperability) for more details.

As an alternative to using **Script**, you can directly create and interact with the parser, compiler, and, VMs directly. There's no good documentation yet, but, check out Script code if you are interested.

 



## Scripts

...

## Compiler and VM

...

## Sandbox Environments

...

## Type Conversion

...
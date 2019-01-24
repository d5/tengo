# Interoperability  

## Table of Contents 

- [Using Scripts](#using-scripts)
  - [Type Conversion Table](#type-conversion-table)
- [Compiler and VM](#compiler-and-vm)
- [User Types](#user-types)
- [Sandbox Environments](#sandbox-environments)

## Using Scripts

Embedding and executing the Tengo code in Go is very easy. At a high level, this process is like:

- create a [Script](https://godoc.org/github.com/d5/tengo/script#Script) instance with your code,
- _optionally_ add some [Script Variables](https://godoc.org/github.com/d5/tengo/script#Variable) to Script,
- compile or directly run the script,
- retrieve _output_ values from the [Compiled](https://godoc.org/github.com/d5/tengo/script#Compiled) instance.

The following is an example where a Tengo script is compiled and run with no input/output variables.

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

Here's another example where an input variable is added to the script, and, an output variable is accessed through [Variable.Int](https://godoc.org/github.com/d5/tengo/script#Variable.Int) function:

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

A variable `b` is defined by the user before compilation using [Script.Add](https://godoc.org/github.com/d5/tengo/script#Script.Add) function. Then a compiled bytecode `c` is used to execute the bytecode and get the value of global variables. In this example, the value of global variable `a` is read using [Compiled.Get](https://godoc.org/github.com/d5/tengo/script#Compiled.Get) function. See [documentation](https://godoc.org/github.com/d5/tengo/script#Variable) for the full list of variable value functions.

### Type Conversion Table

When adding a Variable _([Script.Add](https://godoc.org/github.com/d5/tengo/script#Script.Add))_, Script converts Go values into Tengo values based on the following conversion table.

| Go Type | Tengo Type | Note |
| :--- | :--- | :--- |
|`nil`|`Undefined`||
|`string`|`String`||
|`int64`|`Int`||
|`int`|`Int`||
|`bool`|`Bool`||
|`rune`|`Char`||
|`byte`|`Char`||
|`float64`|`Float`||
|`[]byte`|`Bytes`||
|`error`|`Error{String}`|use `error.Error()` as String value|
|`map[string]Object`|`Map`||
|`map[string]interface{}`|`Map`|individual elements converted to Tengo objects|
|`[]Object`|`Array`||
|`[]interface{}`|`Array`|individual elements converted to Tengo objects|
|`Object`|`Object`|_(no type conversion performed)_|


## User Types

One can easily add and use customized value types in Tengo code by implementing [Object](https://godoc.org/github.com/d5/tengo/objects#Object) interface. Tengo runtime will treat the user types exactly in the same way it does to the runtime types with no performance overhead. See [Tengo Objects](https://github.com/d5/tengo/blob/master/docs/objects.md) for more details.

## Compiler and VM

Although it's not recommended, you can directly create and run the Tengo [Parser](https://godoc.org/github.com/d5/tengo/compiler/parser#Parser), [Compiler](https://godoc.org/github.com/d5/tengo/compiler#Compiler), and [VM](https://godoc.org/github.com/d5/tengo/runtime#VM) for yourself instead of using Scripts and Script Variables. It's a bit more involved as you have to manage the symbol tables and global variables between them, but, basically that's what Script and Script Variable is doing internally.

_TODO: add more information here_

## Sandbox Environments

In an environment where a _(potentially)_ unsafe script code needs to be executed,   

_TODO: add more information here_


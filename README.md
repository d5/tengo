<p align="center">
  <img src="https://raw.githubusercontent.com/d5/tengolang.com/master/logo_400.png" width="200" height="200">
</p>

# The Tengo Language

[![GoDoc](https://godoc.org/github.com/d5/tengo?status.svg)](https://godoc.org/github.com/d5/tengo/script)
[![Go Report Card](https://goreportcard.com/badge/github.com/d5/tengo)](https://goreportcard.com/report/github.com/d5/tengo)
[![Build Status](https://travis-ci.org/d5/tengo.svg?branch=master)](https://travis-ci.org/d5/tengo)

Tengo is an embeddable script language for Go. 

Tengo is [fast](#benchmark) as it's compiled to bytecode and executed on stack-based VM that's written in native Go.

\>> **Try [Tengo Playground](https://tengolang.com/)** << 

## Features

- Simple and intuitive [syntax](https://github.com/d5/tengo#tengo-syntax-in-5-minutes)
- Dynamically typed with type coercions
- First-class functions and Closures
- Garbage collected _(thanks to Go runtime)_
- Easily extensible using customizable types
- Written in pure Go _(no CGO, no external dependencies)_
- Executable as a standalone language

## Benchmark

| | fib(35) | fibt(35) |  Type  |
| :--- |    ---: |     ---: |  :---: |
| Go | `58ms` | `3ms` | Go (native) |
| [**Tengo**](https://github.com/d5/tengo) | `4,245ms` | `5ms` | VM on Go |
| Lua | `1,739ms` | `3ms` | Lua (native) |
| [go-lua](https://github.com/Shopify/go-lua) | `5,368ms` | `5ms` | Lua VM on Go |
| [GopherLua](https://github.com/yuin/gopher-lua) | `5,408ms` | `5ms` | Lua VM on Go |
| Python | `3,176ms` | `27ms` | Python (native) |
| [starlark-go](https://github.com/google/starlark-go) | `15,400ms` | `5ms` | Python-like Interpreter on Go |
| [gpython](https://github.com/go-python/gpython) | `17,724ms` | `6ms` | Python Interpreter on Go |
| [otto](https://github.com/robertkrimen/otto) | `82,050ms` | `22ms` | JS Interpreter on Go |
| [Anko](https://github.com/mattn/anko) | `98,739ms` | `31ms` | Interpreter on Go |

[fib(35)](https://github.com/d5/tengobench/blob/master/code/fib.tengo) is a function to compute 35th Fibonacci number, and, [fibt(35)](https://github.com/d5/tengobench/blob/master/code/fibtc.tengo) is the [tail-call](https://en.wikipedia.org/wiki/Tail_call) version of the same function.
 
_Please note that **Go** case does not read the source code from a local file, while all other cases do. All shell commands and the source code used in this benchmarking is available [here](https://github.com/d5/tengobench)._

## Tengo Syntax in 5 Minutes

Tengo supports line comments (`//...`) and block comments (`/* ... */`).

```golang
/* 
  multi-line block comments 
*/

a := 5 // line comments
```
> [Run in Playground](https://tengolang.com/?s=02e384399a0397b0a752f08604ccb244d1a6cb37)


Tengo is a dynamically typed language, and, you can initialize the variables using `:=` operator. 

```golang
a := 1984 		// int
b := "aomame"		// string
c := -9.22		// float
d := true		// bool
e := '九'		// char
f := [1, false, "foo"]	// array
g := {			// map
    h: 439,
    i: 12.34,
    j: [0, 9, false]
}
k := func(l, m) {	// function
    return l + m
}
```
> [Run in Playground](https://tengolang.com/?s=f8626a711769502ce20e4560ace65c0e9c1279f4)

After the variable is initialized, it can be re-assigned different value using `=` operator. 

```golang
a := 1928		// int
a = "foo"		// string
f := func() {
    a := false		// 'a' is defined in the function scope
    a = [1, 2, 3]	// and thus does not affect 'a' in global scope.
}
a == "foo" 		// still "foo"
```
> [Run in Playground](https://tengolang.com/?s=1d39bc2af5c51417df82b32db47a0e6a156d48ec)


Type is not directly specified, but, you can use type-coercion functions to convert between types.

```golang
s1 := string(1984)  // "1984"
i2 := int("-999")   // -999
f3 := float(-51)    // -51.0
b4 := bool(1)       // true
c5 := char("X")     // 'X'
```
> [Run in Playground](https://tengolang.com/?s=8d57905b82959eb244e9bbd2111e12ee04a33045)

_See [Variable Types](https://github.com/d5/tengo/wiki/Variable-Types) for more details on the variable types._

You can use the dot selector (`.`) and indexer (`[]`) operator to read or write elements of arrays, strings, or maps.

```golang
["one", "two", "three"][1]	// == "two"

m := {
    a: 1,
    b: [2, 3, 4],
    c: func() { return 10 }
}
m.a				// == 1
m["b"][1]			// == 3
m.c()				// == 10
m.x = 5				// add 'x' to map 'm'
//m.b[5] = 0			// but this is an error: index out of bounds
```
> [Run in Playground](https://tengolang.com/?s=d510c75ed8f06ef1e22c1aaf8a7d4565c793514c)

For sequence types (string, bytes, array), you can use slice operator (`[:]`) too.

```golang
a := [1, 2, 3, 4, 5][1:3]	// == [2, 3]
b := [1, 2, 3, 4, 5][3:]	// == [4, 5]
c := [1, 2, 3, 4, 5][:3]	// == [1, 2, 3]
d := "hello world"[2:10]	// == "llo worl"
```
> [Run in Playground](https://tengolang.com/?s=214ab490bb24549578770984985f6b161aed915d)


In Tengo, functions are first-class citizen, and, it also supports closures, functions that captures variables in outer scopes. In the following example, the function returned from `adder` is capturing `base` variable.

```golang
adder := func(base) {
    return func(x) { return base + x }	// capturing 'base'
}
add5 := adder(5)
nine := add5(4)		// == 9
```
> [Run in Playground](https://tengolang.com/?s=fba79990473d5b38cc944dfa225d38580ddaf422)


For flow control, Tengo currently supports **if-else**, **for**, **for-in** statements.

```golang
// IF-ELSE
if a < 0 {
    // ...
} else if a == 0 {
    // ...
} else {
    // ...
}

// IF with init statement
if a := 0; a < 10 {
    // ...
} else {
    // ...
}

// FOR
for a:=0; a<10; a++ {
    // ...
}

// FOR condition-only (like WHILE in other languages)
for a < 10 {
    // ...
}

// FOR-IN
for x in [1, 2, 3] {		// array: element
    // ...
}
for i, x in [1, 2, 3] {		// array: index and element
    // ...
} 
for k, v in {k1: 1, k2: 2} {	// map: key and value
    // ...
}
```

An error object is created using `error` function-like keyword. An error can have any types of value and the underlying value of the error can be accessed using `.value` selector.
 
```golang
err1 := error("oops")   // error with string value
err2 := error(1+2+3)    // error with int value
if is_error(err1) {     // 'is_error' builtin function
    err_val := err1.value   // get underlying value 
}  
``` 
> [Run in Playground](https://tengolang.com/?s=5eaba4289c9d284d97704dd09cb15f4f03ad05c1)

You can load other scripts as modules using `import` expression.

Main script:
```golang
mod1 := import("./mod1") // assuming mod1.tengo file exists in the current directory 
                         // same as 'import("./mod1.tengo")' or 'import("mod1")'
mod1.func1(a)            // module function 
a += mod1.foo            // module variable
//mod1.foo = 5           // error: module variables are read-only
```

`mod1.tengo` file:

```golang
func1 := func(x) { print(x) }
foo := 2
```

Basically, `import` expression returns all the global variables defined in the module as a Map-like value. One can access the functions or variables defined in the module using `.` selector or `["key"]` indexer, but, module variables are immutable.

Also, you can use `import` to load the [standard libraries](https://github.com/d5/tengo/wiki/Standard-Libraries).

```golang
math := import("math")
a := math.abs(-19.84) // == 19.84
```


## Embedding Tengo in Go

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

## Tengo CLI Tool

Although Tengo is designed as an embedded script language for Go, it can be compiled and executed as native binary using `tengo` tool.

### Installing Tengo Tool

To install `tengo` tool, run:

```bash
go get github.com/d5/tengo/cmd/tengo
```

### Compiling and Executing Tengo Code

You can directly execute the Tengo source code by running `tengo` tool with your Tengo source file (`*.tengo`).

```bash
tengo myapp.tengo
```

Or, you can compile the code into a binary file and execute it later.

```bash
tengo -c -o myapp myapp.tengo   # compile 'myapp.tengo' into binary file 'myapp'
tengo myapp                     # execute the compiled binary `myapp`	
```

### Tengo REPL

You can run Tengo [REPL](https://en.wikipedia.org/wiki/Read–eval–print_loop) if you run `tengo` with no arguments.

```bash
tengo
```

## Roadmap

### v0. _(Current)_

Things are experimental, and, the focus is on the **core language features**, **stability**, **basic interoperability**, and the **performance optimization**.

### v1. Tengo as a Script Language

This will be the first _versioned_ release, and, the main goal for v1 is to make Tengo as a _fast_ embeddable script language for Go, which means Tengo will be comparable to other Go-based script languages such as [Starlark](https://github.com/google/starlark-go), [Lua](https://github.com/Shopify/go-lua) [VM](https://github.com/yuin/gopher-lua)s, and [other](https://github.com/robertkrimen/otto) [interpreter](https://github.com/mattn/anko)s.

- Interoperability with Go code
- Sandbox environment
- More language features such as bound methods and switch-case statements

### v2. Tengo as a Standalone Language

- Language-level concurrency support
- Tengo Standard Libraries
- Native executables compilation
- More language features

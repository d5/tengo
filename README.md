# The Tengo Language

Tengo is an embeddable script language for Go. 

Tengo is [fast](#benchmark) as it's compiled to bytecode and executed on stack-based VM that's written in native Go.

\>> **Try [Tengo Playground](https://tengolang.com/)** << 

## Features

- Simple and intuitive syntax
- Dynamically typed with type coercions
- First-class functions and Closures
- Garbage collected _(thanks to Go runtime)_
- Easily extendible using customizable types
- Written in pure Go _(no CGO, no external dependencies)_
- Excutable as a standalone language _(without writing any Go code)_

## Benchmark

| | fib(35) | fibt(35) |  Type  |
| :--- |    ---: |     ---: |  :---: |
| Go | `59ms` | `4ms` | Go (native) |
| [**Tengo**](https://github.com/d5/tengo) | `4,809ms` | `5ms` | VM on Go |
| Lua | `1,752ms` | `3ms` | Lua (native) |
| [go-lua](https://github.com/Shopify/go-lua) | `5,236ms` | `5ms` | Lua VM on Go |
| [GopherLua](https://github.com/yuin/gopher-lua) | `5,558ms` | `5ms` | Lua VM on Go |
| Python | `3,132ms` | `28ms` | Python (native) |
| [otto](https://github.com/robertkrimen/otto) | `85,765ms` | `22ms` | JS Interpreter on Go |
| [Anko](https://github.com/mattn/anko) | `99,235ms` | `24ms` | Interpreter on Go |

[fib(35)](https://github.com/d5/tengobench/blob/master/code/fib.tengo) is a function to compute 35th Fibonacci number, and, [fibt(35)](https://github.com/d5/tengobench/blob/master/code/fibtc.tengo) is the [tail-call](https://en.wikipedia.org/wiki/Tail_call) version of the same function. You can see all the code used for this test in [tengobench](https://github.com/d5/tengobench).

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


Type is not explicitly specified, but, you can use type coercion functions to convert between types.

```golang
s1 := string(1984)  // "1984"
i2 := int("-999")   // -999
f3 := float(-51)    // -51.0
b4 := bool(1)       // true
c5 := char("X")     // 'X'
```
> [Run in Playground](https://tengolang.com/?s=8d57905b82959eb244e9bbd2111e12ee04a33045)

You can use dot selector (`.`) and indexer (`[]`) operator to read or write elemens of arrays or maps.

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

For sequence types (string or array), you can use slice operator (`[:]`) too.

```golang
a := [1, 2, 3, 4, 5][1:3]	// == [2, 3]
b := [1, 2, 3, 4, 5][3:]	// == [4, 5]
c := [1, 2, 3, 4, 5][:3]	// == [1, 2, 3]
d := "hello world"[2:10]	// == "llo worl"
```
> [Run in Playground](https://tengolang.com/?s=214ab490bb24549578770984985f6b161aed915d)


In Tengo, functions are first-class citizen and be treated like any other variables. Tengo also supports closures, functions that captures variables in outer scopes. In the following example, the function that's being returned from `adder` function is capturing `base` variable.

```golang
adder := func(base) {
    return func(x) { return base + x }	// capturing 'base'
}
add5 := adder(5)
nine := add5(4)		// nine
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

In the example above, a variable `b` is defined by the user before compiliation using `Script.Add()` function. Then a compiled bytecode `c` is used to execute the bytecode and get the value of global variables. In thie example, the value of global variable `a` is read using `Compiled.Get()` function.

If you need the custom data types (outside Tengo's primitive types), you can define your own `struct` that implements `objects.Object` interface _(and optinoally `objects.Callable` if you want to make function-like invokable objects)_.

```golang
import (
	"errors"
	"fmt"

	"github.com/d5/tengo/compiler/token"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
)

type Counter struct {
	value int64
}

func (o *Counter) TypeName() string {
	return "counter"
}

func (o *Counter) String() string {
	return fmt.Sprintf("Counter(%d)", o.value)
}

func (o *Counter) BinaryOp(op token.Token, rhs objects.Object) (objects.Object, error) {
	switch rhs := rhs.(type) {
	case *Counter:
		switch op {
		case token.Add:
			return &Counter{value: o.value + rhs.value}, nil
		case token.Sub:
			return &Counter{value: o.value - rhs.value}, nil
		}
	case *objects.Int:
		switch op {
		case token.Add:
			return &Counter{value: o.value + rhs.Value}, nil
		case token.Sub:
			return &Counter{value: o.value - rhs.Value}, nil
		}
	}

	return nil, errors.New("invalid operator")
}

func (o *Counter) IsFalsy() bool {
	return o.value == 0
}

func (o *Counter) Equals(t objects.Object) bool {
	if tc, ok := t.(*Counter); ok {
		return o.value == tc.value
	}

	return false
}

func (o *Counter) Copy() objects.Object {
	return &Counter{value: o.value}
}

func (o *Counter) Call(args ...objects.Object) (objects.Object, error) {
	return &objects.Int{Value: o.value}, nil
}

var code = []byte(`
arr := [1, 2, 3, 4]
for x in arr {
	c1 += x
}
out := c1()`)

func main() {
	s := script.New(code)

	// define variable 'c1'
	_ = s.Add("c1", &Counter{value: 5})

	// compile the source
	c, err := s.Run()
	if err != nil {
		panic(err)
	}

	// retrieve value of 'out'
	out := c.Get("out")
	fmt.Println(out.Int()) // prints "15" ( = 5 + (1 + 2 + 3 + 4) )
}

```

As an alternative to using **Script**, you can directly create and interact with the parser, compiler and VMs directly. There's no good documentations yet, but, check out Script code if you are interested.


## Tengo as a Standalone Language

Although Tengo is designed as an embedded script language for Go, it can be compiled and executed as native binary without any Go code using `tengo` tool.

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

Development roadmap for Tengo:

- Module system _(or packages)_
- Standard libraries
- Better documentations
- More language constructs such as error handling, object methods, switch-case statements
- Native executables compilation
- Performance improvements
- Syntax highlighter for IDEs


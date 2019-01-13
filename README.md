# The Tengo Language

Tengo is an embedded script language for Go.

\>> **Try Tengo in online [Playground](https://tengolang.com/)** << 

## Language Features

- Simple and intuitive syntax
- Dynamically typed with type coercions
- Bytecode compiled _(see the [benchmark](#benchmark) results)_
- First-class functions and Closures
- Garbage collected _(thanks to Go runtime)_
- Easily extendible using customizable types
- Written in pure Go _(no CGO, no external dependencies)_
- Excutable as a standalone language _(without writing any Go code)_

## Benchmark

| | fib(35) | fibt(35) |  Type  |
| :--- |    ---: |     ---: |  :---: |
| Go | `68,713,331` | `3,264,992` | Go (native) |
| [**Tengo**](https://github.com/d5/tengo) | `6,811,234,411` | `4,699,512` | Go-VM |
| Lua | `1,946,451,017` | `3,220,991` | Lua (native) |
| [go-lua](https://github.com/Shopify/go-lua) | `5,658,423,479` | `4,247,160` | Go-Lua-VM |
| [GopherLua](https://github.com/yuin/gopher-lua) | `6,301,424,553` | `5,194,735` | Go-Lua-VM |
| Python | `3,159,870,102` | `28,512,040` | Python (native) |
| [otto](https://github.com/robertkrimen/otto) | `91,616,109,035` | `13,780,650` | Go-JS-Interpreter |
| [Anko](https://github.com/mattn/anko) | `119,395,411,432` | `22,266,008` | Go-Interpreter |

_*Nanoseconds_

`fib(35)` is a function to calculate 35th Fibonacci number.

```golang
fib := func(x) {
	if x == 0 {
		return 0
	} else if x == 1 {
		return 1
	} else {
		return fib(x-1) + fib(x-2)
	}
}
fib(35)
```

`fibt(35)` is a [tail-call](https://en.wikipedia.org/wiki/Tail_call) version of `fib(35)`.

```golang
fibt := func(x, a, b) {
	if x == 0 {
		return a
	} else if x == 1 {
		return b
	} else {
		return fibt(x-1, b, a+b)
	}
}
fibt(35, 0, 1)
```

Please see [tengobench](https://github.com/d5/tengobench) for more details.

## Tengo Syntax in 5 Minutes

Tengo supports line comments (`//...`) and block comments (`/* ... */`).

```golang
/* 
  multi-line block comments 
*/

a := 5 // line comments
```

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
After the variable is initialized, it can be re-assigned different value using `=` operator. 

```golang
a := 1928		// int
a = "foo"		// string
f := func() {
    a := false		// 'a' is defined in the function scope
    a = [1, 2, 3]	// and thus does not affect 'a' in global scope.
}
print(a) 		// still "foo"
```

Type is not explicitly specified, but, you can use type coercion functions to convert between types.

```golang
s1 := string(1984)  // "1984"
i2 := int("-999")   // -999
f3 := float(-51)    // -51.0
b4 := bool(1)       // true
c5 := char("X")     // 'X'
```

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
m.b[5] = 0			// but this is an error: index out of bounds
```

For sequence types (string or array), you can use slice operator (`[:]`) too.

```golang
[1, 2, 3, 4, 5][1:3]	// == [2, 3]
[1, 2, 3, 4, 5][3:]	// == [4, 5]
[1, 2, 3, 4, 5][:3]	// == [1, 2, 3]
"hello world"[2:10]	// == "llo worl"
```

In Tengo, functions are first-class citizen and be treated like any other variables. Tengo also supports closures, functions that captures variables in outer scopes. In the following example, the function that's being returned from `adder` function is capturing `base` variable.

```golang
adder := func(base) {
    return func(x) { return base + x }	// capturing 'base'
}
add5 := adder(5)
nine := add5(4)		// nine
```

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

## Embed Tengo in Go

...

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

...


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

The next big features planned include:

- Module system _(or packages)_
- Standard libraries
- Better documentations
- More language constructs such as error handling, object methods, switch-case statements
- Native executables compilation
- Performance improvements
- Syntax highlighter for IDEs


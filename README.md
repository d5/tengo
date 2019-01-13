# The Tengo Language

Tengo is an embedded script language for Go.

Try Tengo language **[here](https://tengolang.com/)**!

```golang
each := func(seq, fn) {
    for x in seq { // array iteration
        fn(x) 
    }
}

sum := func(seq) {
   s := 0
   each(seq, func(x) { s += x }) // closure: capturing variable 's'
   return s
}

print(sum([1, 2, 3])) // "6"
```

## Language Features

Tengo, as a programming language, has the following features:

- Simple and intuitive syntax
- Dynamically typed with type coercions
- Bytecode compiled _(see the [benchmark](#benchmark) results)_
- First-class functions and Closures
- Garbage collected _(thanks to Go runtime)_
- Easily extendible using customizable types
- Written in pure Go _(no CGO, no external dependencies)_
- _(Can be)_ a standalone language _(without writing any Go code)_

## Benchmark

| | fib(35) | fibt(35) |  Type  |
| :--- |    ---: |     ---: |  :---: |
| Go | `75,245,201` | `527` | Go (native) |
| **Tengo** | `6,716,413,970` | `4,338,042` | Go-VM |
| Lua | `1,839,627,814` | `3,768,932` | Lua (native) |
| go-lua | `5,466,274,012` | `4,508,039` | Go-Lua-VM |
| GopherLua | `5,740,626,066` | `5,027,486` | Go-Lua-VM |
| Python | `3,021,823,532` | `22,829,440` | Python (native) |
| otto | `92,194,561,922` | `13,250,326` | Go-JS-Interpreter |
| Anko | `125,748,242,982` | `15,296,442` | Go-Interpreter |

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

## Syntax

### Variables and Types

...

### Flow Control

...

### Functions

...

## Tengo as an Embedded Script

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

- Module system _(or packages)
- Standard libraries _(most likely with modules)_
- More language constructs such as error handling, object methods, switch-case statements
- Tengo tool to compile into native executables
- Improvements on compilation and execution performance 

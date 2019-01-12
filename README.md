# Tengo Language

Tengo is an embeddable script language for Go.

## Features

The Tengo language, as a programming language, has the following features:

- Dynamically typed
- Bytecode compiled: very fast (see [benchmark](#benchmark) results below)
- First-class functions and closures
- Type coercions

Also Tengo is ...

- Garbage collected (free bonus as a Go embedded script)
- Easily extendible using customizable types
- Written in native Go (no CGO or any external dependencies)
- _(Can be)_ compiled and executed as a standalone language (without any Go code)

## Compiling and Execution

## Benchmark

| | Fib(35) | FibTC(35) | Type |
|--|--|--|--|
| [Go](https://github.com/yuin/gopher-lua) | 61.215ms | 381ns | Go |
| [Lua](https://github.com/yuin/gopher-lua) | 100s | 100s | Lua |
| [Python3](https://github.com/yuin/gopher-lua) | 100s | 100s | Python |
| [GopherLua](https://github.com/yuin/gopher-lua) | 100s | 100s | Go-VM |
| [go-lua](https://github.com/Shopify/go-lua) | 100s | 100s | Go-VM |
| [otto](https://github.com/robertkrimen/otto) | 100s | 100s | Go-Interpreter |
| [Anko](https://github.com/mattn/anko) | 100s | 100s | Go-Interpreter |
| **Tengo** | **6.852s** | **78.728µs** | Go-VM |


`Fib(35)` is a function to calculate 35th Fibonacci number.

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
```

`FibTC(35)` is a [tail-call](https://en.wikipedia.org/wiki/Tail_call) version of `Fib(35)`.

```golang
fibtc := func(x, a, b) {
	if x == 0 {
		return a
	} else if x == 1 {
		return b
	} else {
		return fibtc(x-1, b, a+b)
	}
}
```

Please see this [Wiki](https://github.com/d5/tengo/wiki/Benchmarks) for more details.

## Binary Compilation and Execution

Although Tengo language is designed as an embedded script for Go, it can also be compiled and executed as native binary without any Go code.

### Tengo tool

To install `tengo` tool, run:

```bash
go get github.com/d5/tengo/cmd/tengo
```

_(In the future release, prebuilt binaries for `tengo` tool will be provided so the users don't need `go` tool.)_

To compile a Tengo source code, use `-c` or `-compile` flag:

```bash
tengo -c myapp.tengo
```

This will compile the source code (`myapp.tengo`) and generate a compiled binary `myapp.out`. You can use `-o` flag to override the output file name:

```bash
tengo -c -o myapp myapp.tengo
```

Now the compiled binary can be executed using the same `tengo` tool:

```bash
tengo myapp
```

_(`tengo` tool is still needed for execution like you need `java` tool to execute Java applications. In the future release, `tengo` compiler might be able to generate native executables directly.)_

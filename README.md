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

| Lang | fib(35) | fibt(35) |  Type  |
| :--- |    ---: |     ---: |  :---: |
| Go | `79,335,297` | `825` | Go |
| Python | `3,041,348,388` | `23,416,612` | Python |
| Lua | 1,898,257,490 | 2,979,080 | Lua |
| GopherLua | 6,332,378,472 | 4,547,143 | Go-Lua-VM |
| go-lua | 5,561,849,390 | 4,018,411 | Go-Lua-VM |
| Anko | 114,157,535,022 | 13,238,407 | Go-Interpreter |
| otto | 91,862,077,318 | 13,043,798 | Go-Interpreter |
| **Tengo** | `6,646,884,016`| `4,396,387` | Go-VM |

_*All units in nanoseconds._

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
```

Please see [tengobench](https://github.com/d5/tengobench) for more details.

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

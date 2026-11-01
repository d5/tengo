# The Tengo Language

[![GoDoc](https://godoc.org/github.com/d5/tengo/v2?status.svg)](https://godoc.org/github.com/d5/tengo/v2)
![test](https://github.com/d5/tengo/workflows/test/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/d5/tengo)](https://goreportcard.com/report/github.com/d5/tengo)

**Tengo is a small, dynamic, fast, secure script language for Go.** 

Tengo is **[fast](#benchmark)** and secure because it's compiled/executed as
bytecode on stack-based VM that's written in native Go.

```golang
/* The Tengo Language */
fmt := import("fmt")

each := func(seq, fn) {
    for x in seq { fn(x) }
}

sum := func(init, seq) {
    each(seq, func(x) { init += x })
    return init
}

fmt.println(sum(0, [1, 2, 3]))   // "6"
fmt.println(sum("", [1, 2, 3]))  // "123"
```

> Test this Tengo code in the
> [Tengo Playground](https://tengolang.com/?s=0c8d5d0d88f2795a7093d7f35ae12c3afa17bea3)

## Features

- Simple and highly readable
  [Syntax](https://github.com/d5/tengo/blob/master/docs/tutorial.md)
  - Dynamic typing with type coercion
  - Higher-order functions and closures
  - Immutable values
- [Securely Embeddable](https://github.com/d5/tengo/blob/master/docs/interoperability.md)
  and [Extensible](https://github.com/d5/tengo/blob/master/docs/objects.md)
- Compiler/runtime written in native Go _(no external deps or cgo)_
- Executable as a
  [standalone](https://github.com/d5/tengo/blob/master/docs/tengo-cli.md)
  language / REPL
- Use cases: rules engine, [state machine](https://github.com/d5/go-fsm),
  data pipeline, [transpiler](https://github.com/d5/tengo2lua)

## Benchmark

| | fib(35) | fibt(35) |  Language (Type)  |
| :--- |    ---: |     ---: |  :---: |
| [**Tengo**](https://github.com/d5/tengo) | `2,315ms` | `3ms` | Tengo (VM) |
| [go-lua](https://github.com/Shopify/go-lua) | `4,028ms` | `3ms` | Lua (VM) |
| [GopherLua](https://github.com/yuin/gopher-lua) | `4,409ms` | `3ms` | Lua (VM) |
| [goja](https://github.com/dop251/goja) | `5,194ms` | `4ms` | JavaScript (VM) |
| [starlark-go](https://github.com/google/starlark-go) | `6,954ms` | `3ms` | Starlark (Interpreter) |
| [gpython](https://github.com/go-python/gpython) | `11,324ms` | `4ms` | Python (Interpreter) |
| [Yaegi](https://github.com/containous/yaegi) | `11,715ms` | `10ms` | Yaegi (Interpreter) |
| [otto](https://github.com/robertkrimen/otto) | `48,539ms` | `6ms` | JavaScript (Interpreter) |
| [Anko](https://github.com/mattn/anko) | `52,821ms` | `6ms` | Anko (Interpreter) |
| - | - | - | - |
| Go | `47ms` | `2ms` | Go (Native) |
| Lua | `756ms` | `2ms` | Lua (Native) |
| Python | `1,907ms` | `14ms` | Python2 (Native) |

_* [fib(35)](https://github.com/d5/tengobench/blob/master/code/fib.tengo):
Fibonacci(35)_  
_* [fibt(35)](https://github.com/d5/tengobench/blob/master/code/fibtc.tengo):
[tail-call](https://en.wikipedia.org/wiki/Tail_call) version of Fibonacci(35)_  
_* **Go** does not read the source code from file, while all other cases do_  
_* See [here](https://github.com/d5/tengobench) for commands/codes used_

`bench_test.go` holds the Go benchmarks every change is measured against:
recursive and tail calls, loops, closures, arrays, maps and strings, stdlib
calls, and compilation. A change must not slow any of them beyond
[benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)'s noise.

```
git checkout master && make bench > old.txt
git checkout my-branch && make bench > new.txt
go run golang.org/x/perf/cmd/benchstat@latest old.txt new.txt
```

`make bench BENCH=Fib COUNT=20` runs a subset.

Where `(*VM).run` lands in the binary can move Fib, FibTailCall, Loop and
Closure by a few percent on some machines. A change touches a benchmark's
path when it changes the VM or the bytecode the compiler emits for that
benchmark's script. If a benchmark moves after a change that touches
neither, compare the two binaries, and treat the move as noise when only
the function's address changed:

```
git checkout master && go test -c -o old.test .
git checkout my-branch && go test -c -o new.test .
go tool nm -size old.test new.test | grep 'VM).run$'
```

## Quick Start

```
go get github.com/d5/tengo/v2
```

A simple Go example code that compiles/runs Tengo script code with some input/output values:

```golang
package main

import (
	"context"
	"fmt"

	"github.com/d5/tengo/v2"
)

func main() {
	// create a new Script instance
	script := tengo.NewScript([]byte(
`each := func(seq, fn) {
    for x in seq { fn(x) }
}

sum := 0
mul := 1
each([a, b, c, d], func(x) {
    sum += x
    mul *= x
})`))

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
	fmt.Println(sum, mul) // "22 288"
}
```

Or, if you need to evaluate a simple expression, you can use [Eval](https://pkg.go.dev/github.com/d5/tengo/v2#Eval) function instead:


```golang
res, err := tengo.Eval(ctx,
	`input ? "success" : "fail"`,
	map[string]interface{}{"input": 1})
if err != nil {
	panic(err)
}
fmt.Println(res) // "success"
```

## References

- [Language Syntax](https://github.com/d5/tengo/blob/master/docs/tutorial.md)
- [Object Types](https://github.com/d5/tengo/blob/master/docs/objects.md)
- [Runtime Types](https://github.com/d5/tengo/blob/master/docs/runtime-types.md)
  and [Operators](https://github.com/d5/tengo/blob/master/docs/operators.md)
- [Builtin Functions](https://github.com/d5/tengo/blob/master/docs/builtins.md)
- [Interoperability](https://github.com/d5/tengo/blob/master/docs/interoperability.md)
- [Tengo CLI](https://github.com/d5/tengo/blob/master/docs/tengo-cli.md)
- [Standard Library](https://github.com/d5/tengo/blob/master/docs/stdlib.md)
- Syntax Highlighters: [VSCode](https://github.com/lissein/vscode-tengo), [Atom](https://github.com/d5/tengo-atom), [Vim](https://github.com/geseq/tengo-vim), [Emacs](https://github.com/CsBigDataHub/tengo-mode)
- **Why the name Tengo?** It's from [1Q84](https://en.wikipedia.org/wiki/1Q84).



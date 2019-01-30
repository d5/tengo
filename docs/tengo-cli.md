# Tengo CLI Tool

Tengo is designed as an embedding script language for Go, but, it can also be compiled and executed as native binary using `tengo` CLI tool.

## Installing Tengo CLI

To install `tengo` tool, run:

```bash
go get github.com/d5/tengo/cmd/tengo
```

Or, you can download the precompiled binaries from [here](https://github.com/d5/tengo/releases/latest). 

## Compiling and Executing Tengo Code

You can directly execute the Tengo source code by running `tengo` tool with your Tengo source file (`*.tengo`).

```bash
tengo myapp.tengo
```

Or, you can compile the code into a binary file and execute it later.

```bash
tengo -o myapp myapp.tengo   # compile 'myapp.tengo' into binary file 'myapp'
tengo myapp                  # execute the compiled binary `myapp`	
```

## Tengo REPL

You can run Tengo [REPL](https://en.wikipedia.org/wiki/Read–eval–print_loop) if you run `tengo` with no arguments.

```bash
tengo
```
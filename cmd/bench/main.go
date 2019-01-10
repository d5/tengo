package main

import (
	"fmt"
	"time"

	"github.com/d5/tengo/ast"
	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/parser"
	"github.com/d5/tengo/scanner"
	"github.com/d5/tengo/vm"
)

func main() {
	runFib(35)
}

func runFib(n int) {
	start := time.Now()
	nativeResult := fib(n)
	nativeTime := time.Since(start)

	input := `
fib := func(x) {
	if x == 0 {
		return 0
	} else if x == 1 {
		return 1
	} else {
		return fib(x-1) + fib(x-2)
	}
}
` + fmt.Sprintf("out = fib(%d)", n)

	parseTime, compileTime, runTime, result, err := runBench([]byte(input))
	if err != nil {
		panic(err)
	}

	if nativeResult != int(result.(*objects.Int).Value) {
		panic(fmt.Errorf("wrong result: %d != %d", nativeResult, int(result.(*objects.Int).Value)))
	}

	fmt.Printf("fib(%d) = %d\n", n, nativeResult)
	fmt.Println("-------------------------------------")
	fmt.Printf("Go:      %s\n", nativeTime)
	fmt.Printf("Parser:  %s\n", parseTime)
	fmt.Printf("Compile: %s\n", compileTime)
	fmt.Printf("VM:      %s\n", runTime)
}

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}

func runBench(input []byte) (parseTime time.Duration, compileTime time.Duration, runTime time.Duration, result objects.Object, err error) {
	var astFile *ast.File
	parseTime, astFile, err = parse(input)
	if err != nil {
		return
	}

	var bytecode *compiler.Bytecode
	compileTime, bytecode, err = compileFile(astFile)
	if err != nil {
		return
	}

	runTime, result, err = runVM(bytecode)

	return
}

func parse(input []byte) (time.Duration, *ast.File, error) {
	fileSet := scanner.NewFileSet()
	inputFile := fileSet.AddFile("test", -1, len(input))

	start := time.Now()

	file, err := parser.ParseFile(inputFile, input, nil)
	if err != nil {
		return time.Since(start), nil, err
	}

	return time.Since(start), file, nil
}

func compileFile(file *ast.File) (time.Duration, *compiler.Bytecode, error) {
	symTable := compiler.NewSymbolTable()
	symTable.Define("out")

	start := time.Now()

	c := compiler.NewCompiler(symTable, nil)
	if err := c.Compile(file); err != nil {
		return time.Since(start), nil, err
	}

	return time.Since(start), c.Bytecode(), nil
}

func runVM(bytecode *compiler.Bytecode) (time.Duration, objects.Object, error) {
	globals := make([]*objects.Object, vm.GlobalsSize)

	start := time.Now()

	v := vm.NewVM(bytecode, globals)
	if err := v.Run(); err != nil {
		return time.Since(start), nil, err
	}

	return time.Since(start), *globals[0], nil
}

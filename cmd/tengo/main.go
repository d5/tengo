package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/d5/tengo"
	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/runtime"
)

const (
	sourceFileExt = ".tengo"
)

var (
	compile    bool
	outputFile = flag.String("o", "", "Output file")
)

func init() {
	flag.BoolVar(&compile, "compile", false, "Compile input file")
	flag.BoolVar(&compile, "c", false, "Compile input file")
	flag.Parse()
}

func main() {
	inputFile := flag.Arg(0)
	if inputFile == "" {
		doHelp()
		os.Exit(2)
	}

	inputData, err := ioutil.ReadFile(inputFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error reading input file: %s", err.Error())
		os.Exit(1)
	}

	if compile {
		if err := doCompile(inputData, inputFile, *outputFile); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	} else if filepath.Ext(inputFile) == sourceFileExt {
		if err := doCompileRun(inputData, inputFile, *outputFile); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	} else {
		if err := doRun(inputData, inputFile, *outputFile); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}

func doHelp() {
	fmt.Println("Tengo is a tool to compile Tengo source code or execute compiled Tengo binary.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("	tengo [flags] input-file")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println()
	fmt.Println("	-compile/-c compile source file")
	fmt.Println("	-o          output file")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println()
	fmt.Println("	tengo -c -o program.out program.tengo")
	fmt.Println("	            : Compile program.tengo and write compiled binary to program.out file")
	fmt.Println()
	fmt.Println("	tengo program.out")
	fmt.Println("	            : Execute compiled binary program.out")
	fmt.Println()
	fmt.Println()
}

func doCompile(data []byte, inputFile, outputFile string) (err error) {
	bytecode, err := tengo.Compile(data, filepath.Base(inputFile))
	if err != nil {
		return
	}

	if outputFile == "" {
		outputFile = basename(inputFile) + ".out"
	}

	out, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = out.Close()
		} else {
			err = out.Close()
		}
	}()

	err = bytecode.Encode(out)
	if err != nil {
		return
	}

	fmt.Println(outputFile)

	return
}

func doCompileRun(data []byte, inputFile, _ string) (err error) {
	bytecode, err := tengo.Compile(data, filepath.Base(inputFile))
	if err != nil {
		return
	}

	machine := runtime.NewVM(bytecode, nil)

	err = machine.Run()
	if err != nil {
		return
	}

	return
}

func doRun(data []byte, _, _ string) (err error) {
	bytecode := &compiler.Bytecode{}
	err = bytecode.Decode(bytes.NewReader(data))
	if err != nil {
		return
	}

	machine := runtime.NewVM(bytecode, nil)

	err = machine.Run()
	if err != nil {
		return
	}

	return
}

func basename(s string) string {
	s = filepath.Base(s)

	n := strings.LastIndexByte(s, '.')
	if n > 0 {
		return s[:n]
	}

	return s
}

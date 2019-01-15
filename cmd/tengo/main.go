package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/compiler/ast"
	"github.com/d5/tengo/compiler/parser"
	"github.com/d5/tengo/compiler/source"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/runtime"
)

const (
	sourceFileExt = ".tengo"
	replPrompt    = ">> "
)

var (
	compile    bool
	showHelp   bool
	outputFile = flag.String("o", "", "Output file")
)

func init() {
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&compile, "compile", false, "Compile input file")
	flag.BoolVar(&compile, "c", false, "Compile input file")
	flag.Parse()
}

func main() {
	if showHelp {
		doHelp()
		os.Exit(2)
	}

	inputFile := flag.Arg(0)
	if inputFile == "" {
		// REPL
		doRepl(os.Stdin, os.Stdout)
		return
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
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("	tengo [flags] {input-file}")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println()
	fmt.Println("	-c/-compile compile the input and produce bytecode file")
	fmt.Println("	-o          output")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println()
	fmt.Println("	tengo")
	fmt.Println("	            : Start Tengo REPL")
	fmt.Println()
	fmt.Println("	tengo myapp.tengo")
	fmt.Println("	            : Compile and execute source file (myapp.tengo)")
	fmt.Println()
	fmt.Println("	tengo -c myapp myapp.tengo")
	fmt.Println("	            : Compile source file (myapp.tengo) and produce bytecode file (myapp)")
	fmt.Println()
	fmt.Println("	tengo myapp")
	fmt.Println("	            : Execute bytecode file (myapp)")
	fmt.Println()
	fmt.Println()
}

func doCompile(data []byte, inputFile, outputFile string) (err error) {
	bytecode, err := compileSrc(data, filepath.Base(inputFile))
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
	bytecode, err := compileSrc(data, filepath.Base(inputFile))
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

func doRepl(in io.Reader, out io.Writer) {
	stdin := bufio.NewScanner(in)

	fileSet := source.NewFileSet()
	globals := make([]*objects.Object, runtime.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()

	for {
		_, _ = fmt.Fprintf(out, replPrompt)

		scanned := stdin.Scan()
		if !scanned {
			return
		}

		line := stdin.Text()

		file, err := parser.ParseFile(fileSet.AddFile("test", -1, len(line)), []byte(line), nil)
		if err != nil {
			_, _ = fmt.Fprintf(out, "error: %s\n", err.Error())
			continue
		}

		file = addPrints(file)

		c := compiler.NewCompiler(symbolTable, nil)
		if err := c.Compile(file); err != nil {
			_, _ = fmt.Fprintf(out, "Compilation error:\n %s\n", err.Error())
			continue
		}

		machine := runtime.NewVM(c.Bytecode(), globals)
		if err != nil {
			_, _ = fmt.Fprintf(out, "VM error:\n %s\n", err.Error())
			continue
		}
		if err := machine.Run(); err != nil {
			_, _ = fmt.Fprintf(out, "Execution error:\n %s\n", err.Error())
			continue
		}
	}
}

func compileSrc(src []byte, filename string) (*compiler.Bytecode, error) {
	fileSet := source.NewFileSet()

	p := parser.NewParser(fileSet.AddFile(filename, -1, len(src)), src, nil)
	file, err := p.ParseFile()
	if err != nil {
		return nil, err
	}

	c := compiler.NewCompiler(nil, nil)
	if err := c.Compile(file); err != nil {
		return nil, err
	}

	return c.Bytecode(), nil
}

func addPrints(file *ast.File) *ast.File {
	var stmts []ast.Stmt
	for _, s := range file.Stmts {
		switch s := s.(type) {
		case *ast.ExprStmt:
			stmts = append(stmts, &ast.ExprStmt{
				Expr: &ast.CallExpr{
					Func: &ast.Ident{
						Name: "print",
					},
					Args: []ast.Expr{s.Expr},
				},
			})

		case *ast.AssignStmt:
			stmts = append(stmts, s)

			stmts = append(stmts, &ast.ExprStmt{
				Expr: &ast.CallExpr{
					Func: &ast.Ident{
						Name: "print",
					},
					Args: s.LHS,
				},
			})

		default:
			stmts = append(stmts, s)
		}
	}

	return &ast.File{
		InputFile: file.InputFile,
		Stmts:     stmts,
	}
}

func basename(s string) string {
	s = filepath.Base(s)

	n := strings.LastIndexByte(s, '.')
	if n > 0 {
		return s[:n]
	}

	return s
}

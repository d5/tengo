package repl

import (
	"bufio"
	"bytes"
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


//REPL Options
type Options struct {
	//Compile output file
	CompileOutput  string

	//Show help flag
	ShowHelp       bool

	//Show version flag
	ShowVersion    bool

	//Input file
	InputFile      string

	//Builtin modules
	BuiltinModules map[string]objects.Object
}

var (
	bm             map[string]bool
	builtinModules map[string]objects.Object
	version        = "dev"
)

//Run REPL
func Run(options *Options) {
	if options.ShowHelp {
		doHelp()
		os.Exit(2)
	} else if options.ShowVersion {
		fmt.Println(version)
		return
	}

	builtinModules = options.BuiltinModules
	bm = make(map[string]bool, len(builtinModules))
	for k := range builtinModules {
		bm[k] = true
	}

	if options.InputFile == "" {
		// REPL
		runREPL(os.Stdin, os.Stdout)
		return
	}

	inputData, err := ioutil.ReadFile(options.InputFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error reading input file: %s", err.Error())
		os.Exit(1)
	}

	if options.CompileOutput != "" {
		if err := compileOnly(inputData, options.InputFile, options.CompileOutput); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	} else if filepath.Ext(options.InputFile) == sourceFileExt {
		if err := compileAndRun(inputData, options.InputFile); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	} else {
		if err := runCompiled(inputData); err != nil {
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
	fmt.Println("	-o        compile output file")
	fmt.Println("	-version  show version")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println()
	fmt.Println("	tengo")
	fmt.Println()
	fmt.Println("	          Start Tengo REPL")
	fmt.Println()
	fmt.Println("	tengo myapp.tengo")
	fmt.Println()
	fmt.Println("	          Compile and run source file (myapp.tengo)")
	fmt.Println("	          Source file must have .tengo extension")
	fmt.Println()
	fmt.Println("	tengo -o myapp myapp.tengo")
	fmt.Println()
	fmt.Println("	          Compile source file (myapp.tengo) into bytecode file (myapp)")
	fmt.Println()
	fmt.Println("	tengo myapp")
	fmt.Println()
	fmt.Println("	          Run bytecode file (myapp)")
	fmt.Println()
	fmt.Println()
}

func compileOnly(data []byte, inputFile, outputFile string) (err error) {
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

func compileAndRun(data []byte, inputFile string) (err error) {
	bytecode, err := compileSrc(data, filepath.Base(inputFile))
	if err != nil {
		return
	}

	machine := runtime.NewVM(bytecode, nil, nil, builtinModules, -1)

	err = machine.Run()
	if err != nil {
		return
	}

	return
}

func runCompiled(data []byte) (err error) {
	bytecode := &compiler.Bytecode{}
	err = bytecode.Decode(bytes.NewReader(data))
	if err != nil {
		return
	}

	machine := runtime.NewVM(bytecode, nil, nil, builtinModules, -1)

	err = machine.Run()
	if err != nil {
		return
	}

	return
}

func runREPL(in io.Reader, out io.Writer) {
	stdin := bufio.NewScanner(in)

	fileSet := source.NewFileSet()
	globals := make([]objects.Object, runtime.GlobalsSize)

	symbolTable := compiler.NewSymbolTable()
	for idx, fn := range objects.Builtins {
		symbolTable.DefineBuiltin(idx, fn.Name)
	}

	var constants []objects.Object

	for {
		_, _ = fmt.Fprint(out, replPrompt)

		scanned := stdin.Scan()
		if !scanned {
			return
		}

		line := stdin.Text()

		srcFile := fileSet.AddFile("repl", -1, len(line))
		p := parser.NewParser(srcFile, []byte(line), nil)
		file, err := p.ParseFile()
		if err != nil {
			_, _ = fmt.Fprintln(out, err.Error())
			continue
		}

		file = addPrints(file)

		c := compiler.NewCompiler(srcFile, symbolTable, constants, bm, nil)
		if err := c.Compile(file); err != nil {
			_, _ = fmt.Fprintln(out, err.Error())
			continue
		}

		bytecode := c.Bytecode()

		machine := runtime.NewVM(bytecode, globals, nil, builtinModules, -1)
		if err := machine.Run(); err != nil {
			_, _ = fmt.Fprintln(out, err.Error())
			continue
		}

		constants = bytecode.Constants
	}
}

func compileSrc(src []byte, filename string) (*compiler.Bytecode, error) {
	fileSet := source.NewFileSet()
	srcFile := fileSet.AddFile(filename, -1, len(src))

	p := parser.NewParser(srcFile, src, nil)
	file, err := p.ParseFile()
	if err != nil {
		return nil, err
	}

	c := compiler.NewCompiler(srcFile, nil, nil, bm, nil)
	if err := c.Compile(file); err != nil {
		return nil, err
	}

	bytecode := c.Bytecode()
	bytecode.RemoveDuplicates()

	return bytecode, nil
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

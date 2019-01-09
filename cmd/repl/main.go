package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/d5/tengo/ast"
	"github.com/d5/tengo/compiler"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/parser"
	"github.com/d5/tengo/scanner"
	"github.com/d5/tengo/vm"
)

const (
	Prompt = ">> "
)

func main() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello, %s! This is the Ghost programming language!\n", currentUser.Name)

	startRepl(os.Stdin, os.Stdout)
}

func startRepl(in io.Reader, out io.Writer) {
	stdin := bufio.NewScanner(in)

	fileSet := scanner.NewFileSet()
	globals := make([]*objects.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()

	for {
		_, _ = fmt.Fprintf(out, Prompt)

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

		machine := vm.NewVM(c.Bytecode(), globals)
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
					Args: s.Lhs,
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

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/d5/tengo/v2/parser"
	"github.com/d5/tengo/v2/require"
	"github.com/d5/tengo/v2/stdlib"
)

func TestBasename(t *testing.T) {
	for _, tc := range []struct {
		in, out string
	}{
		{"app.tengo", "app"},
		{"dir/app.tengo", "app"},
		{"/abs/dir/app.v1.tengo", "app.v1"},
		{"dir.d/app", "app"},
		{".hidden", ".hidden"},
		{"noext", "noext"},
	} {
		require.Equal(t, tc.out, basename(tc.in), tc.in)
	}
}

func TestAddPrints(t *testing.T) {
	fileSet := parser.NewFileSet()
	src := []byte("a := 1\nb\nif true {}")
	srcFile := fileSet.AddFile("test", -1, len(src))
	file, err := parser.NewParser(srcFile, src, nil).ParseFile()
	require.NoError(t, err)

	got := addPrints(file)
	require.True(t, got.InputFile == file.InputFile)
	require.Equal(t, 4, len(got.Stmts))
	require.True(t, got.Stmts[0] == file.Stmts[0])
	expectPrintCall(t, got.Stmts[1], file.Stmts[0].(*parser.AssignStmt).LHS)
	expectPrintCall(t, got.Stmts[2],
		[]parser.Expr{file.Stmts[1].(*parser.ExprStmt).Expr})
	require.True(t, got.Stmts[3] == file.Stmts[2])
}

func expectPrintCall(t *testing.T, stmt parser.Stmt, args []parser.Expr) {
	call := stmt.(*parser.ExprStmt).Expr.(*parser.CallExpr)
	require.Equal(t, "__repl_println__", call.Func.(*parser.Ident).Name)
	require.Equal(t, len(args), len(call.Args))
	for i := range args {
		require.True(t, call.Args[i] == args[i])
	}
}

func TestRunREPL(t *testing.T) {
	in := strings.Join([]string{
		`a := 1`,
		`a + 2`,
		`s := "x"`,
		`s + "y"`,
		`undefined`,
		`1 +`,
		`missing`,
		`a / 0`,
	}, "\n")
	out := &strings.Builder{}
	printed := captureStdout(t, func() {
		RunREPL(stdlib.GetModuleMap(), strings.NewReader(in), out)
	})

	require.Equal(t, "1\n3\nx\nxy\n<undefined>\n", printed)
	require.Equal(t, 9, strings.Count(out.String(), replPrompt))
	require.True(t, strings.Contains(out.String(), "Parse Error"),
		out.String())
	require.True(t, strings.Contains(out.String(),
		"unresolved reference 'missing'"), out.String())
	require.True(t, strings.Contains(out.String(), "division by zero"),
		out.String())
}

func TestCompileAndRun(t *testing.T) {
	modules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)

	printed := captureStdout(t, func() {
		err := CompileAndRun(modules,
			[]byte(`fmt := import("fmt"); fmt.print("ok")`), "test.tengo")
		require.NoError(t, err)
	})
	require.Equal(t, "ok", printed)

	err := CompileAndRun(modules, []byte(`a :=`), "test.tengo")
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "Parse Error"), err)

	err = CompileAndRun(modules, []byte(`a = 1`), "test.tengo")
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(),
		"unresolved reference 'a'"), err)

	err = CompileAndRun(modules, []byte(`a := 1 / 0`), "test.tengo")
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "division by zero"), err)
}

func TestCompileOnlyAndRunCompiled(t *testing.T) {
	modules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	dir := t.TempDir()
	inputFile := filepath.Join(dir, "app.tengo")
	outputFile := filepath.Join(dir, "app.bin")
	src := []byte(`fmt := import("fmt"); fmt.print("compiled")`)

	printed := captureStdout(t, func() {
		require.NoError(t, CompileOnly(modules, src, inputFile, outputFile))
	})
	require.Equal(t, outputFile+"\n", printed)

	data, err := ioutil.ReadFile(outputFile)
	require.NoError(t, err)
	printed = captureStdout(t, func() {
		require.NoError(t, RunCompiled(modules, data))
	})
	require.Equal(t, "compiled", printed)

	require.Error(t, RunCompiled(modules, []byte("not bytecode")))

	err = CompileOnly(modules, []byte(`a = 1`), inputFile,
		filepath.Join(dir, "bad.bin"))
	require.Error(t, err)
	_, err = os.Stat(filepath.Join(dir, "bad.bin"))
	require.True(t, os.IsNotExist(err), err)
}

func TestCompileOnlyDefaultOutput(t *testing.T) {
	dir := t.TempDir()
	wd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	defer func() { require.NoError(t, os.Chdir(wd)) }()

	printed := captureStdout(t, func() {
		require.NoError(t, CompileOnly(stdlib.GetModuleMap(),
			[]byte(`a := 1`), filepath.Join("src", "app.tengo"), ""))
	})
	require.Equal(t, "app.out\n", printed)
	_, err = os.Stat(filepath.Join(dir, "app.out"))
	require.NoError(t, err)
}

func TestCompileSrcResolvePath(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, ioutil.WriteFile(filepath.Join(dir, "mod.tengo"),
		[]byte(`export 42`), 0644))
	inputFile := filepath.Join(dir, "main.tengo")
	src := []byte(`out := import("./mod")`)

	defer func(v bool) { resolvePath = v }(resolvePath)

	resolvePath = false
	_, err := compileSrc(stdlib.GetModuleMap(), src, inputFile)
	require.Error(t, err)

	resolvePath = true
	bytecode, err := compileSrc(stdlib.GetModuleMap(), src, inputFile)
	require.NoError(t, err)
	require.NotNil(t, bytecode)
}

func captureStdout(t *testing.T, fn func()) string {
	r, w, err := os.Pipe()
	require.NoError(t, err)
	stdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = stdout }()

	fn()
	require.NoError(t, w.Close())
	out, err := ioutil.ReadAll(r)
	require.NoError(t, err)
	return string(out)
}

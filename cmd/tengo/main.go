package main

import (
	"flag"

	"github.com/d5/tengo/repl"
	"github.com/d5/tengo/stdlib"
	"github.com/d5/tengo/objects"
)

const (
	sourceFileExt = ".tengo"
	replPrompt    = ">> "
)

var (
	compileOutput  string
	showHelp       bool
	showVersion    bool
	version        = "dev"
	bm             map[string]bool
	builtinModules map[string]objects.Object
)

func init() {
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.StringVar(&compileOutput, "o", "", "Compile output file")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.Parse()
}

func main() {
	builtinModules = make(map[string]objects.Object, len(stdlib.Modules))
	for k, mod := range stdlib.Modules {
		builtinModules[k] = mod
	}

	repl.ShowHelp = showHelp
	repl.ShowVersion = showVersion
	repl.CompileOutput = compileOutput
	repl.BuiltinModules = builtinModules
	repl.Run()
}
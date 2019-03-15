package main

import (
	"flag"
	
	"github.com/d5/tengo/repl"
)

var (
	compileOutput string
	showHelp      bool
	showVersion   bool
	version       = "dev"
)

func init() {
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.StringVar(&compileOutput, "o", "", "Compile output file")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.Parse()
}

func main() {
	repl.Run(&repl.Options {
		ShowHelp: showHelp,
		ShowVersion: showVersion,
		CompileOutput: compileOutput,
		InputFile: flag.Arg(0),
	})
}
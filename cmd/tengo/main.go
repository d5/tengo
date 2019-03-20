package main

import (
	"flag"

	"github.com/d5/tengo/cli"
	"github.com/d5/tengo/stdlib"
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
	cli.Run(&cli.Options{
		ShowHelp:      showHelp,
		ShowVersion:   showVersion,
		Version:       version,
		CompileOutput: compileOutput,
		Modules:       stdlib.GetModuleMap(stdlib.AllModuleNames()...),
		InputFile:     flag.Arg(0),
	})
}

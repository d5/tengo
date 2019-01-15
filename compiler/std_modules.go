package compiler

import (
	"github.com/d5/tengo/compiler/stdmods"
	"github.com/d5/tengo/objects"
)

var compiledStdMods map[string]*objects.CompiledModule

func init() {
	compiledStdMods = make(map[string]*objects.CompiledModule)

	for _, stdmod := range stdmods.All {
		mod, err := compileModule("module:"+stdmod.Name, []byte(stdmods.Math.Source), stdmods.Math.Globals)
		if err != nil {
			panic(err)
		}

		compiledStdMods[stdmod.Name] = mod
	}
}

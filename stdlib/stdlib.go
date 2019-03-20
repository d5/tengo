package stdlib

//go:generate go run gensrcmods.go

import "github.com/d5/tengo/objects"

// AllModuleNames returns a list of all default module names.
func AllModuleNames() []string {
	var names []string
	for name := range BuiltinModules {
		names = append(names, name)
	}
	for name := range SourceModules {
		names = append(names, name)
	}
	return names
}

// GetModules returns the modules for the given names.
// Duplicate names and invalid names are ignore.
func GetModules(names ...string) *objects.ModuleMap {
	modules := objects.NewModuleMap()

	for _, name := range names {
		if mod := BuiltinModules[name]; mod != nil {
			modules.AddBuiltinModule(name, mod)
		}
		if mod := SourceModules[name]; mod != "" {
			modules.AddSourceModule(name, []byte(mod))
		}
	}

	return modules
}

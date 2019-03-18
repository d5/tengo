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
func GetModules(names ...string) map[string]objects.Importable {
	modules := make(map[string]objects.Importable)
	for _, name := range names {
		if mod := BuiltinModules[name]; mod != nil {
			modules[name] = mod
		}
		if mod := SourceModules[name]; mod != nil {
			modules[name] = mod
		}
	}

	return modules
}

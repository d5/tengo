package stdlib

//go:generate go run genmods.go

import "github.com/d5/tengo/objects"

// BuiltinModules are builtin module standard libraries.
var BuiltinModules = map[string]*objects.ImmutableMap{
	"math":  {Value: mathModule},
	"os":    {Value: osModule},
	"text":  {Value: textModule},
	"times": {Value: timesModule},
	"rand":  {Value: randModule},
}

// AllModuleNames returns a list of all default module names.
func AllModuleNames() []string {
	var names []string
	for name := range BuiltinModules {
		names = append(names, name)
	}
	for name := range CompiledModules {
		names = append(names, name)
	}
	return names
}

// GetModules returns the modules for the given names.
// Duplicate names and invalid names are ignore.
func GetModules(names ...string) map[string]*objects.ImmutableMap {
	modules := make(map[string]*objects.ImmutableMap)
	for _, name := range names {
		if mod := Modules[name]; mod != nil {
			modules[name] = mod
		}
	}
	return modules
}

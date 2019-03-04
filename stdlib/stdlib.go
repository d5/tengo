package stdlib

import "github.com/d5/tengo/objects"

// Modules contain the standard modules.
var Modules = map[string]*objects.ImmutableMap{
	"math":  &objects.ImmutableMap{Value: mathModule},
	"os":    &objects.ImmutableMap{Value: osModule},
	"text":  &objects.ImmutableMap{Value: textModule},
	"times": &objects.ImmutableMap{Value: timesModule},
	"rand":  &objects.ImmutableMap{Value: randModule},
}

// AllModuleNames returns a list of all default module names.
func AllModuleNames() []string {
	var names []string
	for name := range Modules {
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

func objectPtr(o objects.Object) *objects.Object {
	return &o
}

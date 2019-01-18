package stdmods

import "github.com/d5/tengo/objects"

// Modules contain the standard modules.
var Modules = map[string]*objects.ModuleMap{
	"math": {Value: mathModule},
	"os":   {Value: osModule},
}

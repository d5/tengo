package stdlib

import "github.com/d5/tengo/objects"

// Modules contain the standard modules.
var Modules = map[string]*objects.ImmutableMap{
	"math": {Value: mathModule},
	"os":   {Value: osModule},
}

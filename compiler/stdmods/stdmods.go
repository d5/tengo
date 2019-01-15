package stdmods

import "github.com/d5/tengo/objects"

// Module represents a standard module.
type Module struct {
	Name    string
	Source  string
	Globals map[string]objects.Object
}

// All contains all the standard modules.
var All = []Module{
	Math,
}

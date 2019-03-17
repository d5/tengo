package stdlib

import "github.com/d5/tengo/objects"

// BuiltinModules are builtin type standard library modules.
var BuiltinModules = map[string]*objects.BuiltinModule{
	"math":  {Attrs: mathModule},
	"os":    {Attrs: osModule},
	"text":  {Attrs: textModule},
	"times": {Attrs: timesModule},
	"rand":  {Attrs: randModule},
	"fmt":   {Attrs: fmtModule},
	"json":  {Attrs: jsonModule},
}

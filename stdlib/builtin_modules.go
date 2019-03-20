package stdlib

import "github.com/d5/tengo/objects"

// BuiltinModules are builtin type standard library modules.
var BuiltinModules = map[string]*objects.BuiltinModule{
	"math":  {Name: "math", Attrs: mathModule},
	"os":    {Name: "os", Attrs: osModule},
	"text":  {Name: "text", Attrs: textModule},
	"times": {Name: "times", Attrs: timesModule},
	"rand":  {Name: "rand", Attrs: randModule},
	"fmt":   {Name: "fmt", Attrs: fmtModule},
	"json":  {Name: "json", Attrs: jsonModule},
}

package stdlib

import "github.com/d5/tengo"

// BuiltinModules are builtin type standard library modules.
var BuiltinModules = map[string]map[string]tengo.Object{
	"math":  mathModule,
	"os":    osModule,
	"text":  textModule,
	"times": timesModule,
	"rand":  randModule,
	"fmt":   fmtModule,
	"json":  jsonModule,
	"enum":  enumModule,
}

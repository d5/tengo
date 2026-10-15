package tengo_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
)

func benchModuleGraph(numModules, funcsPerModule int) (*tengo.ModuleMap, string) {
	modules := tengo.NewModuleMap()
	for m := 0; m < numModules; m++ {
		var b strings.Builder
		b.WriteString("fmt := import(\"fmt\")\ntext := import(\"text\")\n")
		for d := 1; d <= 3 && m-d >= 0; d++ {
			fmt.Fprintf(&b, "dep%d := import(\"mod%d\")\n", d, m-d)
		}
		for f := 0; f < funcsPerModule; f++ {
			fmt.Fprintf(&b, `
fn%[1]d := func(a, b, ...rest) {
	if is_undefined(a) || !is_int(b) {
		return undefined
	}
	acc := {count: len(rest), items: []}
	for i := 0; i < b; i++ {
		if i %% 3 == 0 {
			acc.items = append(acc.items, fmt.sprintf("%%d-%%v", i, a))
		} else if text.has_prefix(string(a), "x") {
			acc.count += i * 2
		} else {
			acc.count -= 1
		}
	}
	return acc
}
`, f)
		}
		exports := make([]string, funcsPerModule)
		for f := range exports {
			exports[f] = fmt.Sprintf("fn%[1]d: fn%[1]d", f)
		}
		b.WriteString("export {" + strings.Join(exports, ", ") + "}\n")
		modules.AddSourceModule(fmt.Sprintf("mod%d", m), []byte(b.String()))
	}
	for name, mod := range stdlib.BuiltinModules {
		modules.AddBuiltinModule(name, mod)
	}
	main := fmt.Sprintf("m := import(\"mod%d\")\nout := m.fn0(\"x\", 3)\n", numModules-1)
	return modules, main
}

func BenchmarkCompile(b *testing.B) {
	modules, main := benchModuleGraph(40, 25)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s := tengo.NewScript([]byte(main))
		s.SetImports(modules)
		if _, err := s.Compile(); err != nil {
			b.Fatal(err)
		}
	}
}

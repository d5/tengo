package tengo_test

import (
	"testing"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
)

var vmBenchScripts = map[string]string{
	"ArithLoop": `
s := 0
for i := 0; i < 300000; i++ {
	s += (i * 3 + 1) % 7
}
out := s`,
	"FloatLoop": `
zr := 0.0
zi := 0.0
n := 0
for i := 0; i < 200000; i++ {
	t := zr*zr - zi*zi + 0.1
	zi = 2.0*zr*zi + 0.2
	zr = t
	if zr*zr + zi*zi > 4.0 {
		zr = 0.0
		zi = 0.0
		n++
	}
}
out := n`,
	"CompareLoop": `
c := 0
for i := 0; i < 300000; i++ {
	if i % 2 == 0 && i != 7 || i > 299000 {
		c++
	}
}
out := c`,
	"Fib": `
fib := func(n) {
	if n < 2 { return n }
	return fib(n-1) + fib(n-2)
}
out := fib(22)`,
	"ClosureCalls": `
mk := func() {
	c := 0
	return func(d) { c += d; return c }
}
f := mk()
for i := 0; i < 200000; i++ {
	f(1)
}
out := f(0)`,
	"MapAccess": `
m := {a: 1, b: 2, c: 3}
s := 0
for i := 0; i < 200000; i++ {
	s += m.a + m["b"] + m.c
	m.a = i
}
out := s`,
	"MethodCalls": `
obj := {v: 0}
obj.inc = func(d) { obj.v += d; return obj.v }
for i := 0; i < 200000; i++ {
	obj.inc(1)
}
out := obj.v`,
	"ArrayIter": `
arr := []
for i := 0; i < 1000; i++ { arr = append(arr, i) }
s := 0
for k := 0; k < 200; k++ {
	for x in arr { s += x }
}
out := s`,
	"ArrayIndex": `
arr := []
for i := 0; i < 1000; i++ { arr = append(arr, i) }
s := 0
for i := 0; i < 200000; i++ {
	s += arr[i % 1000]
	arr[i % 1000] = s % 1000
}
out := s`,
	"StringOps": `
c := 0
for i := 0; i < 50000; i++ {
	s := "key" + string(i % 10)
	if s == "key3" || len(s) > 10 {
		c++
	}
}
out := c`,
	"Builtins": `
arr := [1, 2, 3]
c := 0
for i := 0; i < 200000; i++ {
	if is_int(i) && len(arr) == 3 && !is_undefined(arr[1]) {
		c++
	}
}
out := c`,
	"MapIter": `
m := {}
for i := 0; i < 50; i++ { m["key" + string(i)] = i }
c := 0
for k := 0; k < 4000; k++ {
	for key, v in m { c += v }
}
out := c`,
	"StdlibCalls": `
text := import("text")
math := import("math")
c := 0
for i := 0; i < 50000; i++ {
	if text.has_prefix("prefix-value", "prefix") { c += math.abs(-1) }
}
out := c`,
}

func BenchmarkVM(b *testing.B) {
	for _, name := range []string{"ArithLoop", "FloatLoop", "CompareLoop", "Fib",
		"ClosureCalls", "MapAccess", "MethodCalls", "ArrayIter", "ArrayIndex",
		"StringOps", "Builtins", "StdlibCalls", "MapIter"} {
		src := vmBenchScripts[name]
		b.Run(name, func(b *testing.B) {
			s := tengo.NewScript([]byte(src))
			s.SetImports(stdlib.GetModuleMap(stdlib.AllModuleNames()...))
			c, err := s.Compile()
			if err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := c.Run(); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

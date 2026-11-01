package tengo_test

import (
	"testing"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
)

const benchFib = `
fib := func(x) {
	if x == 0 {
		return 0
	} else if x == 1 {
		return 1
	}
	return fib(x-1) + fib(x-2)
}
out := fib(20)
`

const benchFibTailCall = `
fib := func(x, a, b) {
	if x == 0 {
		return a
	} else if x == 1 {
		return b
	}
	return fib(x-1, b, a+b)
}
out := 0
for i := 0; i < 1000; i++ {
	out = fib(50, 0, 1)
}
`

const benchLoop = `
sum := 0
for i := 0; i < 100000; i++ {
	if i % 3 == 0 {
		sum += i
	} else {
		sum -= 1
	}
}
out := sum
`

const benchClosure = `
counter := func() {
	n := 0
	return func(d) {
		n += d
		return n
	}
}
c := counter()
for i := 0; i < 20000; i++ {
	c(2)
}
out := c(0)
`

const benchCollections = `
a := []
for i := 0; i < 1000; i++ {
	a = append(a, i)
}
m := {}
for i, v in a {
	m["k" + string(i)] = v * 2
}
sum := 0
for k, v in m {
	sum += v + len(k)
}
for i := 0; i < len(a); i++ {
	sum += a[i]
}
s := ""
for i := 0; i < 500; i++ {
	s += string(i % 10)
}
for c in s {
	if c == '5' {
		sum++
	}
}
out := sum
`

const benchStdlib = `
json := import("json")
text := import("text")
fmt := import("fmt")
doc := {name: "tengo", tags: ["vm", "script"], n: 0}
out := 0
for i := 0; i < 200; i++ {
	doc.n = i
	d := json.decode(json.encode(doc))
	parts := text.split(fmt.sprintf("%s,%d,%v", d.name, d.n, d.tags), ",")
	out += len(parts) + len(text.to_upper(d.name))
}
`

func BenchmarkFib(b *testing.B) {
	benchRun(b, benchFib, 6765)
}

func BenchmarkFibTailCall(b *testing.B) {
	benchRun(b, benchFibTailCall, 12586269025)
}

func BenchmarkLoop(b *testing.B) {
	benchRun(b, benchLoop, 1666616667)
}

func BenchmarkClosure(b *testing.B) {
	benchRun(b, benchClosure, 40000)
}

func BenchmarkCollections(b *testing.B) {
	benchRun(b, benchCollections, 1502440)
}

func BenchmarkStdlib(b *testing.B) {
	benchRun(b, benchStdlib, 1800)
}

// BenchmarkCompile parses and compiles every script above.
func BenchmarkCompile(b *testing.B) {
	srcs := [][]byte{
		[]byte(benchFib), []byte(benchFibTailCall), []byte(benchLoop),
		[]byte(benchClosure), []byte(benchCollections), []byte(benchStdlib),
	}
	modules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, src := range srcs {
			s := tengo.NewScript(src)
			s.SetImports(modules)
			if _, err := s.Compile(); err != nil {
				b.Fatal(err)
			}
		}
	}
}

// benchRun times Run alone, after checking that one run of src sets out to
// want.
func benchRun(b *testing.B, src string, want int64) {
	s := tengo.NewScript([]byte(src))
	s.SetImports(stdlib.GetModuleMap(stdlib.AllModuleNames()...))
	c, err := s.Compile()
	if err != nil {
		b.Fatal(err)
	}
	if err := c.Run(); err != nil {
		b.Fatal(err)
	}
	if got := c.Get("out").Int64(); got != want {
		b.Fatalf("out = %d, want %d", got, want)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := c.Run(); err != nil {
			b.Fatal(err)
		}
	}
}

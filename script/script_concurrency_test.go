package script_test

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
)

func TestScriptConcurrency(t *testing.T) {
	solve := func(a, b, c int) (d, e int) {
		a += 2
		b += c
		a += b * 2
		d = a + b + c
		e = 0
		for i := 1; i <= d; i++ {
			e += i
		}
		e *= 2
		return
	}

	code := []byte(`
mod1 := import("mod1")

a += 2
b += c
a += b * 2

arr := [a, b, c]
arrstr := stringify(arr)
map := {a: a, b: b, c: c}

d := a + b + c
s := 0

for i:=1; i<=d; i++ {
	s += i
}

e := mod1.double(s)
`)
	mod1 := &objects.ImmutableMap{
		Value: map[string]objects.Object{
			"double": &objects.UserFunction{
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					arg0, _ := objects.ToInt64(args[0])
					ret = &objects.Int{Value: arg0 * 2}
					return
				},
			},
		},
	}

	scr := script.New(code)
	_ = scr.Add("a", 0)
	_ = scr.Add("b", 0)
	_ = scr.Add("c", 0)
	scr.SetBuiltinModules(map[string]*objects.ImmutableMap{
		"mod1": mod1,
	})
	scr.SetBuiltinFunctions([]*objects.BuiltinFunction{
		{
			Name: "stringify",
			Value: func(args ...objects.Object) (ret objects.Object, err error) {
				ret = &objects.String{Value: args[0].String()}
				return
			},
		},
	})
	compiled, err := scr.Compile()
	assert.NoError(t, err)

	executeFn := func(compiled *script.Compiled, a, b, c int) (d, e int) {
		_ = compiled.Set("a", a)
		_ = compiled.Set("b", b)
		_ = compiled.Set("c", c)
		err := compiled.Run()
		assert.NoError(t, err)
		d = compiled.Get("d").Int()
		e = compiled.Get("e").Int()
		return
	}

	concurrency := 500
	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(compiled *script.Compiled) {
			time.Sleep(time.Duration(rand.Int63n(50)) * time.Millisecond)
			defer wg.Done()

			a := rand.Intn(10)
			b := rand.Intn(10)
			c := rand.Intn(10)

			d, e := executeFn(compiled, a, b, c)
			expectedD, expectedE := solve(a, b, c)

			assert.Equal(t, expectedD, d, "input: %d, %d, %d", a, b, c)
			assert.Equal(t, expectedE, e, "input: %d, %d, %d", a, b, c)
		}(compiled.Clone())
	}
	wg.Wait()
}

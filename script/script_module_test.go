package script_test

import (
	"strings"
	"testing"

	"github.com/d5/tengo"
	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/script"
)

func TestScriptSourceModule(t *testing.T) {
	// script1 imports "mod1"
	scr := script.New([]byte(`out := import("mod")`))
	mods := tengo.NewModuleMap()
	mods.AddSourceModule("mod", []byte(`export 5`))
	scr.SetImports(mods)
	c, err := scr.Run()
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, int64(5), c.Get("out").Value())

	// executing module function
	scr = script.New([]byte(`fn := import("mod"); out := fn()`))
	mods = tengo.NewModuleMap()
	mods.AddSourceModule("mod", []byte(`a := 3; export func() { return a + 5 }`))
	scr.SetImports(mods)
	c, err = scr.Run()
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, int64(8), c.Get("out").Value())

	scr = script.New([]byte(`out := import("mod")`))
	mods = tengo.NewModuleMap()
	mods.AddSourceModule("mod", []byte(`text := import("text"); export text.title("foo")`))
	mods.AddBuiltinModule("text", map[string]tengo.Object{
		"title": &tengo.UserFunction{Name: "title", Value: func(_ tengo.Interop, args ...tengo.Object) (ret tengo.Object, err error) {
			s, _ := tengo.ToString(args[0])
			return &tengo.String{Value: strings.Title(s)}, nil
		}},
	})
	scr.SetImports(mods)
	c, err = scr.Run()
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "Foo", c.Get("out").Value())
	scr.SetImports(nil)
	_, err = scr.Run()
	if !assert.Error(t, err) {
		return
	}
}

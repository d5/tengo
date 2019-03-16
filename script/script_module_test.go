package script_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
)

func TestScript_SetUserModuleLoader(t *testing.T) {
	// script1 imports "mod1"
	scr := script.New([]byte(`out := import("mod")`))
	scr.SetUserModuleLoader(func(name string) ([]byte, error) {
		return []byte(`export 5`), nil
	})
	c, err := scr.Run()
	assert.Equal(t, int64(5), c.Get("out").Value())

	// executing module function
	scr = script.New([]byte(`fn := import("mod"); out := fn()`))
	scr.SetUserModuleLoader(func(name string) ([]byte, error) {
		return []byte(`a := 3; export func() { return a + 5 }`), nil
	})
	c, err = scr.Run()
	assert.NoError(t, err)
	assert.Equal(t, int64(8), c.Get("out").Value())

	scr = script.New([]byte(`out := import("mod")`))
	scr.SetBuiltinModules(map[string]*objects.ImmutableMap{
		"text": objectPtr(&objects.ImmutableMap{
			Value: map[string]objects.Object{
				"title": &objects.UserFunction{Name: "title", Value: func(args ...objects.Object) (ret objects.Object, err error) {
					s, _ := objects.ToString(args[0])
					return &objects.String{Value: strings.Title(s)}, nil
				}},
			},
		}),
	})
	scr.SetUserModuleLoader(func(name string) ([]byte, error) {
		if name == "mod" {
			return []byte(`text := import("text"); export text.title("foo")`), nil
		}
		return nil, errors.New("module not found")
	})
	c, err = scr.Run()
	assert.NoError(t, err)
	assert.Equal(t, "Foo", c.Get("out").Value())
	scr.SetBuiltinModules(nil)
	_, err = scr.Run()
	assert.Error(t, err)

}

package script_test

import (
	"strings"
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
)

func TestScriptSourceModule(t *testing.T) {
	// script1 imports "mod1"
	scr := script.New([]byte(`out := import("mod")`))
	scr.SetImports(map[string]objects.Importable{
		"mod": &objects.SourceModule{Src: []byte(`export 5`)},
	})
	c, err := scr.Run()
	assert.Equal(t, int64(5), c.Get("out").Value())

	// executing module function
	scr = script.New([]byte(`fn := import("mod"); out := fn()`))
	scr.SetImports(map[string]objects.Importable{
		"mod": &objects.SourceModule{Src: []byte(`a := 3; export func() { return a + 5 }`)},
	})
	c, err = scr.Run()
	assert.NoError(t, err)
	assert.Equal(t, int64(8), c.Get("out").Value())

	scr = script.New([]byte(`out := import("mod")`))
	scr.SetImports(map[string]objects.Importable{
		"text": &objects.BuiltinModule{
			Attrs: map[string]objects.Object{
				"title": &objects.UserFunction{Name: "title", Value: func(args ...objects.Object) (ret objects.Object, err error) {
					s, _ := objects.ToString(args[0])
					return &objects.String{Value: strings.Title(s)}, nil
				}},
			},
		},
		"mod": &objects.SourceModule{Src: []byte(`text := import("text"); export text.title("foo")`)},
	})
	c, err = scr.Run()
	assert.NoError(t, err)
	assert.Equal(t, "Foo", c.Get("out").Value())
	scr.SetImports(nil)
	_, err = scr.Run()
	assert.Error(t, err)

}

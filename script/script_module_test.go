package script_test

import (
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/script"
)

func TestScript_AddModule(t *testing.T) {
	// mod1 module
	mod1 := script.New([]byte(`a := 5`))

	// script1 imports "mod1"
	scr1 := script.New([]byte(`mod1 := import("mod1"); out := mod1.a`))
	scr1.AddModule("mod1", mod1)
	c, err := scr1.Run()
	assert.Equal(t, int64(5), c.Get("out").Value())

	// mod2 module imports "mod1"
	mod2 := script.New([]byte(`mod1 := import("mod1"); b := mod1.a * 2`))
	mod2.AddModule("mod1", mod1)

	// script2 imports "mod2" (which imports "mod1")
	scr2 := script.New([]byte(`mod2 := import("mod2"); out := mod2.b`))
	scr2.AddModule("mod2", mod2)
	c, err = scr2.Run()
	assert.NoError(t, err)
	assert.Equal(t, int64(10), c.Get("out").Value())
}

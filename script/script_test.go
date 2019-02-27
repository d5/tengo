package script_test

import (
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
	"github.com/d5/tengo/stdlib"
)

func TestScript_Add(t *testing.T) {
	s := script.New([]byte(`a := b`))
	assert.NoError(t, s.Add("b", 5))     // b = 5
	assert.NoError(t, s.Add("b", "foo")) // b = "foo"  (re-define before compilation)
	c, err := s.Compile()
	assert.NoError(t, err)
	assert.NoError(t, c.Run())
	assert.Equal(t, "foo", c.Get("a").Value())
	assert.Equal(t, "foo", c.Get("b").Value())
}

func TestScript_Remove(t *testing.T) {
	s := script.New([]byte(`a := b`))
	err := s.Add("b", 5)
	assert.NoError(t, err)
	assert.True(t, s.Remove("b")) // b is removed
	_, err = s.Compile()          // should not compile because b is undefined
	assert.Error(t, err)
}

func TestScript_Run(t *testing.T) {
	s := script.New([]byte(`a := b`))
	err := s.Add("b", 5)
	assert.NoError(t, err)
	c, err := s.Run()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	compiledGet(t, c, "a", int64(5))
}

func TestScript_SetBuiltinFunctions(t *testing.T) {
	s := script.New([]byte(`a := len([1, 2, 3])`))
	c, err := s.Run()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	compiledGet(t, c, "a", int64(3))

	s = script.New([]byte(`a := len([1, 2, 3])`))
	s.SetBuiltinFunctions(map[string]*objects.BuiltinFunction{"len": &objects.Builtins[3]})
	c, err = s.Run()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	compiledGet(t, c, "a", int64(3))

	s.SetBuiltinFunctions(map[string]*objects.BuiltinFunction{"print": &objects.Builtins[0]})
	_, err = s.Run()
	assert.Error(t, err)

	s.SetBuiltinFunctions(nil)
	_, err = s.Run()
	assert.Error(t, err)
}

func TestScript_SetBuiltinModules(t *testing.T) {
	s := script.New([]byte(`math := import("math"); a := math.abs(-19.84)`))
	c, err := s.Run()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	compiledGet(t, c, "a", 19.84)

	s.SetBuiltinModules(map[string]*objects.ImmutableMap{"math": objectPtr(*stdlib.Modules["math"])})
	c, err = s.Run()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	compiledGet(t, c, "a", 19.84)

	s.SetBuiltinModules(map[string]*objects.ImmutableMap{"os": objectPtr(*stdlib.Modules["os"])})
	_, err = s.Run()
	assert.Error(t, err)

	s.SetBuiltinModules(nil)
	_, err = s.Run()
	assert.Error(t, err)
}

func objectPtr(o objects.Object) *objects.ImmutableMap {
	return o.(*objects.ImmutableMap)
}

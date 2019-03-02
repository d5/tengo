package script_test

import (
	"errors"
	"math"
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
)

func TestScript_Add(t *testing.T) {
	s := script.New([]byte(`a := b; c := test(b); d := test(5)`))
	assert.NoError(t, s.Add("b", 5))     // b = 5
	assert.NoError(t, s.Add("b", "foo")) // b = "foo"  (re-define before compilation)
	assert.NoError(t, s.Add("test", func(args ...objects.Object) (ret objects.Object, err error) {
		if len(args) > 0 {
			switch arg := args[0].(type) {
			case *objects.Int:
				return &objects.Int{Value: arg.Value + 1}, nil
			}
		}

		return &objects.Int{Value: 0}, nil
	}))
	c, err := s.Compile()
	assert.NoError(t, err)
	assert.NoError(t, c.Run())
	assert.Equal(t, "foo", c.Get("a").Value())
	assert.Equal(t, "foo", c.Get("b").Value())
	assert.Equal(t, int64(0), c.Get("c").Value())
	assert.Equal(t, int64(6), c.Get("d").Value())
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
	s.SetBuiltinFunctions([]*objects.BuiltinFunction{&objects.Builtins[3]})
	c, err = s.Run()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	compiledGet(t, c, "a", int64(3))

	s.SetBuiltinFunctions([]*objects.BuiltinFunction{&objects.Builtins[0]})
	_, err = s.Run()
	assert.Error(t, err)

	s.SetBuiltinFunctions(nil)
	_, err = s.Run()
	assert.Error(t, err)

	s = script.New([]byte(`a := import("b")`))
	s.SetUserModuleLoader(func(name string) ([]byte, error) {
		if name == "b" {
			return []byte(`export import("c")`), nil
		} else if name == "c" {
			return []byte("export len([1, 2, 3])"), nil
		}
		return nil, errors.New("module not found")
	})
	s.SetBuiltinFunctions([]*objects.BuiltinFunction{&objects.Builtins[3]})
	c, err = s.Run()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	compiledGet(t, c, "a", int64(3))
}

func TestScript_SetBuiltinModules(t *testing.T) {
	s := script.New([]byte(`math := import("math"); a := math.abs(-19.84)`))
	s.SetBuiltinModules(map[string]*objects.ImmutableMap{
		"math": objectPtr(&objects.ImmutableMap{
			Value: map[string]objects.Object{
				"abs": &objects.UserFunction{Name: "abs", Value: func(args ...objects.Object) (ret objects.Object, err error) {
					v, _ := objects.ToFloat64(args[0])
					return &objects.Float{Value: math.Abs(v)}, nil
				}},
			},
		}),
	})
	c, err := s.Run()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	compiledGet(t, c, "a", 19.84)

	c, err = s.Run()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	compiledGet(t, c, "a", 19.84)

	s.SetBuiltinModules(map[string]*objects.ImmutableMap{"os": objectPtr(&objects.ImmutableMap{Value: map[string]objects.Object{}})})
	_, err = s.Run()
	assert.Error(t, err)

	s.SetBuiltinModules(nil)
	_, err = s.Run()
	assert.Error(t, err)
}

func objectPtr(o objects.Object) *objects.ImmutableMap {
	return o.(*objects.ImmutableMap)
}

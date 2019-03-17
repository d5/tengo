package script_test

import (
	"math"
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
	"github.com/d5/tengo/stdlib"
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

func TestScript_BuiltinModules(t *testing.T) {
	s := script.New([]byte(`math := import("math"); a := math.abs(-19.84)`))
	s.SetImports(map[string]objects.Importable{
		"math": &objects.BuiltinModule{
			Attrs: map[string]objects.Object{
				"abs": &objects.UserFunction{Name: "abs", Value: func(args ...objects.Object) (ret objects.Object, err error) {
					v, _ := objects.ToFloat64(args[0])
					return &objects.Float{Value: math.Abs(v)}, nil
				}},
			},
		},
	})
	c, err := s.Run()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	compiledGet(t, c, "a", 19.84)

	c, err = s.Run()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	compiledGet(t, c, "a", 19.84)

	s.SetImports(map[string]objects.Importable{"os": &objects.BuiltinModule{Attrs: map[string]objects.Object{}}})
	_, err = s.Run()
	assert.Error(t, err)

	s.SetImports(nil)
	_, err = s.Run()
	assert.Error(t, err)
}

func TestScript_SourceModules(t *testing.T) {
	s := script.New([]byte(`
enum := import("enum")
a := 0
enum.each([1,2,3], func(_, v) { 
	a += v
})
`))
	s.SetImports(stdlib.GetModules("enum"))
	c, err := s.Run()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	compiledGet(t, c, "a", int64(6))

	s.SetImports(nil)
	_, err = s.Run()
	assert.Error(t, err)
}

func TestScript_SetMaxConstObjects(t *testing.T) {
	// one constant '5'
	s := script.New([]byte(`a := 5`))
	s.SetMaxConstObjects(1) // limit = 1
	_, err := s.Compile()
	assert.NoError(t, err)
	s.SetMaxConstObjects(0) // limit = 0
	_, err = s.Compile()
	assert.Equal(t, "exceeding constant objects limit: 1", err.Error())

	// two constants '5' and '1'
	s = script.New([]byte(`a := 5 + 1`))
	s.SetMaxConstObjects(2) // limit = 2
	_, err = s.Compile()
	assert.NoError(t, err)
	s.SetMaxConstObjects(1) // limit = 1
	_, err = s.Compile()
	assert.Equal(t, "exceeding constant objects limit: 2", err.Error())

	// duplicates will be removed
	s = script.New([]byte(`a := 5 + 5`))
	s.SetMaxConstObjects(1) // limit = 1
	_, err = s.Compile()
	assert.NoError(t, err)
	s.SetMaxConstObjects(0) // limit = 0
	_, err = s.Compile()
	assert.Equal(t, "exceeding constant objects limit: 1", err.Error())

	// no limit set
	s = script.New([]byte(`a := 1 + 2 + 3 + 4 + 5`))
	_, err = s.Compile()
	assert.NoError(t, err)
}

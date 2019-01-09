package script_test

import (
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/script"
)

type M map[string]interface{}

func TestCompiled_Get(t *testing.T) {
	// simple script
	c := compile(t, `a = 5`, nil)
	compiledRun(t, c)
	compiledGet(t, c, "a", int64(5))

	// user-defined variables
	compileError(t, `a = b`, nil)          // compile error because "b" is not defined
	c = compile(t, `a = b`, M{"b": "foo"}) // now compile with b = "foo" defined
	compiledGet(t, c, "a", nil)            // a = undefined; because it's before Compiled.Run()
	compiledRun(t, c)                      // Compiled.Run()
	compiledGet(t, c, "a", "foo")          // a = "foo"
}

func TestCompiled_GetAll(t *testing.T) {
	c := compile(t, `a = 5`, nil)
	compiledRun(t, c)
	compiledGetAll(t, c, M{"a": int64(5)})

	c = compile(t, `a = b`, M{"b": "foo"})
	compiledRun(t, c)
	compiledGetAll(t, c, M{"a": "foo", "b": "foo"})

	c = compile(t, `a = b; b = 5`, M{"b": "foo"})
	compiledRun(t, c)
	compiledGetAll(t, c, M{"a": "foo", "b": int64(5)})
}

func TestCompiled_IsDefined(t *testing.T) {
	c := compile(t, `a = 5`, nil)
	compiledIsDefined(t, c, "a", false) // a is not defined before Run()
	compiledRun(t, c)
	compiledIsDefined(t, c, "a", true)
	compiledIsDefined(t, c, "b", false)
}

func compile(t *testing.T, input string, vars M) *script.Compiled {
	s := script.New([]byte(input))
	for vn, vv := range vars {
		err := s.Add(vn, vv)
		if !assert.NoError(t, err) {
			return nil
		}
	}

	c, err := s.Compile()
	if !assert.NoError(t, err) || !assert.NotNil(t, c) {
		return nil
	}

	return c
}

func compileError(t *testing.T, input string, vars M) bool {
	s := script.New([]byte(input))
	for vn, vv := range vars {
		err := s.Add(vn, vv)
		if !assert.NoError(t, err) {
			return false
		}
	}

	_, err := s.Compile()

	return assert.Error(t, err)
}

func compiledRun(t *testing.T, c *script.Compiled) bool {
	err := c.Run()

	return assert.NoError(t, err)
}

func compiledGet(t *testing.T, c *script.Compiled, name string, expected interface{}) bool {
	v := c.Get(name)
	if !assert.NotNil(t, v) {
		return false
	}

	return assert.Equal(t, expected, v.Value())
}

func compiledGetAll(t *testing.T, c *script.Compiled, expected M) bool {
	vars := c.GetAll()

	if !assert.Equal(t, len(expected), len(vars)) {
		return false
	}

	for k, v := range expected {
		var found bool
		for _, e := range vars {
			if e.Name() == k {
				if !assert.Equal(t, v, e.Value()) {
					return false
				}
				found = true
			}
		}
		if !found {
			assert.Fail(t, "variable '%s' not found", k)
		}
	}

	return true
}

func compiledIsDefined(t *testing.T, c *script.Compiled, name string, expected bool) bool {
	return assert.Equal(t, expected, c.IsDefined(name))
}

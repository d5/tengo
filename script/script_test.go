package script_test

import (
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/script"
)

func TestScript_Add(t *testing.T) {
	s := script.New([]byte(`a = b`))
	assert.NoError(t, s.Add("b", 5))     // b = 5
	assert.NoError(t, s.Add("b", "foo")) // b = "foo"  (re-define before compilation)
	c, err := s.Compile()
	assert.NoError(t, err)
	assert.NoError(t, c.Run())
	assert.Equal(t, "foo", c.Get("a").Value())
	assert.Equal(t, "foo", c.Get("b").Value())
}

func TestScript_Remove(t *testing.T) {
	s := script.New([]byte(`a = b`))
	err := s.Add("b", 5)
	assert.NoError(t, err)
	assert.True(t, s.Remove("b")) // b is removed
	_, err = s.Compile()          // should not compile because b is undefined
	assert.Error(t, err)
}

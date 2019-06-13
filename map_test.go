package tengo_test

import (
	"testing"

	"github.com/d5/tengo"
	"github.com/d5/tengo/assert"
)

func TestMap_Index(t *testing.T) {
	m := &tengo.Map{Value: make(map[string]tengo.Object)}
	k := &tengo.Int{Value: 1}
	v := &tengo.String{Value: "abcdef"}
	err := m.IndexSet(k, v)

	assert.NoError(t, err)

	res, err := m.IndexGet(k)
	assert.NoError(t, err)
	assert.Equal(t, v, res)
}

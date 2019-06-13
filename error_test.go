package tengo_test

import (
	"testing"

	"github.com/d5/tengo"
	"github.com/d5/tengo/assert"
)

func TestError_Equals(t *testing.T) {
	err1 := &tengo.Error{Value: &tengo.String{Value: "some error"}}
	err2 := err1
	assert.True(t, err1.Equals(err2))
	assert.True(t, err2.Equals(err1))

	err2 = &tengo.Error{Value: &tengo.String{Value: "some error"}}
	assert.False(t, err1.Equals(err2))
	assert.False(t, err2.Equals(err1))
}

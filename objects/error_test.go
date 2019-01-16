package objects_test

import (
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/objects"
)

func TestError_Equals(t *testing.T) {
	err1 := &objects.Error{Value: &objects.String{Value: "some error"}}
	err2 := err1
	assert.True(t, err1.Equals(err2))
	assert.True(t, err2.Equals(err1))

	err2 = &objects.Error{Value: &objects.String{Value: "some error"}}
	assert.False(t, err1.Equals(err2))
	assert.False(t, err2.Equals(err1))
}

package objects_test

import (
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/compiler/token"
	"github.com/d5/tengo/objects"
)

func testBinaryOp(t *testing.T, lhs objects.Object, op token.Token, rhs objects.Object, expected objects.Object) bool {
	t.Helper()

	actual, err := lhs.BinaryOp(op, rhs)

	return assert.NoError(t, err) && assert.Equal(t, expected, actual)
}

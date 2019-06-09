package tengo_test

import (
	"testing"

	"github.com/d5/tengo"
	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/compiler/token"
)

func testBinaryOp(t *testing.T, lhs tengo.Object, op token.Token, rhs tengo.Object, expected tengo.Object) bool {
	t.Helper()

	actual, err := lhs.BinaryOp(op, rhs)

	return assert.NoError(t, err) && assert.Equal(t, expected, actual)
}

func boolValue(b bool) tengo.Object {
	if b {
		return tengo.TrueValue
	}

	return tengo.FalseValue
}

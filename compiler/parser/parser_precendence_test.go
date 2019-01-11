package parser_test

import (
	"testing"
)

func TestPrecedence(t *testing.T) {
	expectString(t, `a + b + c`, `((a + b) + c)`)
	expectString(t, `a + b * c`, `(a + (b * c))`)
	expectString(t, `x = 2 * 1 + 3 / 4`, `x = ((2 * 1) + (3 / 4))`)
}

package parser_test

import (
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/compiler/parser"
	"github.com/d5/tengo/compiler/source"
)

func TestError_Error(t *testing.T) {
	err := &parser.Error{Pos: source.FilePos{Offset: 10, Line: 1, Column: 10}, Msg: "test"}
	assert.Equal(t, "1:10: test", err.Error())
}

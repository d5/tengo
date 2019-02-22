package parser_test

import (
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/compiler/parser"
	"github.com/d5/tengo/compiler/source"
)

func TestErrorList_Sort(t *testing.T) {
	var list parser.ErrorList
	list.Add(source.FilePos{Offset: 20, Line: 2, Column: 10}, "error 2")
	list.Add(source.FilePos{Offset: 30, Line: 3, Column: 10}, "error 3")
	list.Add(source.FilePos{Offset: 10, Line: 1, Column: 10}, "error 1")
	list.Sort()
	assert.Equal(t, "Parse Error: error 1\n\tat 1:10 (and 2 more errors)", list.Error())
}

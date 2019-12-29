package parser_test

import (
	"testing"

	"github.com/d5/tengo/v2/parser"
)

func TestIdentListString(t *testing.T) {
	identListVar := &parser.IdentList{
		List: []*parser.Ident{
			{Name: "a"},
			{Name: "b"},
			{Name: "c"},
		},
		VarArgs: true,
	}

	expectedVar := "(a, b, ...c)"
	if str := identListVar.String(); str != expectedVar {
		t.Fatalf("expected string of %#v to be %s, got %s",
			identListVar, expectedVar, str)
	}

	identList := &parser.IdentList{
		List: []*parser.Ident{
			{Name: "a"},
			{Name: "b"},
			{Name: "c"},
		},
		VarArgs: false,
	}

	expected := "(a, b, c)"
	if str := identList.String(); str != expected {
		t.Fatalf("expected string of %#v to be %s, got %s",
			identList, expected, str)
	}
}

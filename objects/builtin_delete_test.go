package objects

import (
	"testing"
)

func TestBuiltinDelete(t *testing.T) {
	m := &Map{Value: make(map[string]Object)}

	k1 := &String{Value: "foo"}
	v1 := &String{Value: "bar"}
	m.IndexSet(k1, v1)

	k2 := &String{Value: "hello"}
	v2 := &String{Value: "world"}
	m.IndexSet(k2, v2)

	builtinDelete(m, &String{Value: "foo"})

	world, _ := m.IndexGet(&String{Value: "hello"})
	val := world.String()
	if `"world"` != val {
		t.Fatal(val)
	}

	bar, _ := m.IndexGet(&String{Value: "foo"})
	if bar.String() != "<undefined>" {
		t.Fatal(bar)
	}
}

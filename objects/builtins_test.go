package objects_test

import (
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/objects"
)

func TestGetBuiltinFunctions(t *testing.T) {
	testGetBuiltinFunctions(t)
	testGetBuiltinFunctions(t, "print")
	testGetBuiltinFunctions(t, "int", "float")
	testGetBuiltinFunctions(t, "int", "float", "printf")
	testGetBuiltinFunctions(t, "int", "int") // duplicate names ignored
}

func TestGetAllBuiltinFunctions(t *testing.T) {
	funcs := objects.GetAllBuiltinFunctions()
	if !assert.Equal(t, len(objects.Builtins), len(funcs)) {
		return
	}

	namesM := make(map[string]bool)
	for _, bf := range objects.Builtins {
		namesM[bf.Name] = true
	}

	for _, bf := range funcs {
		assert.True(t, namesM[bf.Name], "name: %s", bf.Name)
	}
}

func TestAllBuiltinFunctionNames(t *testing.T) {
	names := objects.AllBuiltinFunctionNames()
	if !assert.Equal(t, len(objects.Builtins), len(names)) {
		return
	}

	namesM := make(map[string]bool)
	for _, name := range names {
		namesM[name] = true
	}

	for _, bf := range objects.Builtins {
		assert.True(t, namesM[bf.Name], "name: %s", bf.Name)
	}
}

func testGetBuiltinFunctions(t *testing.T, names ...string) {
	// remove duplicates
	namesM := make(map[string]bool)
	for _, name := range names {
		namesM[name] = true
	}

	funcs := objects.GetBuiltinFunctions(names...)
	if !assert.Equal(t, len(namesM), len(funcs)) {
		return
	}

	for _, bf := range funcs {
		assert.True(t, namesM[bf.Name], "name: %s", bf.Name)
	}
}

package tengo_test

import (
	"testing"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/require"
)

func TestModuleMap(t *testing.T) {
	m := tengo.NewModuleMap()
	require.Equal(t, 0, m.Len())

	m.AddSourceModule("src", []byte(`export 1`))
	m.AddBuiltinModule("builtin", map[string]tengo.Object{
		"a": &tengo.Int{Value: 1},
	})
	require.Equal(t, 2, m.Len())

	src := m.GetSourceModule("src")
	require.NotNil(t, src)
	require.Equal(t, []byte(`export 1`), src.Src)
	require.Nil(t, m.GetSourceModule("builtin"))
	require.Nil(t, m.GetSourceModule("missing"))
	require.NotNil(t, m.GetBuiltinModule("builtin"))
	require.Nil(t, m.GetBuiltinModule("src"))

	m.AddSourceModule("src", []byte(`export 2`))
	require.Equal(t, 2, m.Len())
	require.Equal(t, []byte(`export 2`), m.GetSourceModule("src").Src)

	c := m.Copy()
	c.Remove("src")
	require.Equal(t, 1, c.Len())
	require.Nil(t, c.Get("src"))
	require.Equal(t, 2, m.Len())
	require.NotNil(t, m.Get("src"))

	m.Remove("src")
	require.Equal(t, 1, m.Len())
	require.Nil(t, m.Get("src"))
	require.Nil(t, m.GetSourceModule("src"))

	m.Remove("missing")
	require.Equal(t, 1, m.Len())

	m.Remove("builtin")
	require.Equal(t, 0, m.Len())
}

func TestModuleMap_RemoveFromScript(t *testing.T) {
	m := tengo.NewModuleMap()
	m.AddSourceModule("mod", []byte(`export 42`))

	s := tengo.NewScript([]byte(`out := import("mod")`))
	s.SetImports(m)
	c, err := s.Run()
	require.NoError(t, err)
	compiledGet(t, c, "out", int64(42))

	m.Remove("mod")
	_, err = s.Compile()
	require.Error(t, err)
}

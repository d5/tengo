package tengo_test

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/require"
)

func TestMultiFS_Open(t *testing.T) {
	fsA := fstest.MapFS{
		"a.tengo": &fstest.MapFile{Data: []byte(`export "from-a"`)},
	}
	fsB := fstest.MapFS{
		"b.tengo": &fstest.MapFile{Data: []byte(`export "from-b"`)},
	}

	mfs := tengo.MultiFS{fsA, fsB}

	t.Run("finds file in first FS", func(t *testing.T) {
		f, err := mfs.Open("a.tengo")
		require.NoError(t, err)
		require.NoError(t, f.Close())
	})

	t.Run("finds file in second FS", func(t *testing.T) {
		f, err := mfs.Open("b.tengo")
		require.NoError(t, err)
		require.NoError(t, f.Close())
	})

	t.Run("first FS shadows second for same path", func(t *testing.T) {
		overlap := tengo.MultiFS{
			fstest.MapFS{"x.tengo": &fstest.MapFile{Data: []byte(`export "first"`)}},
			fstest.MapFS{"x.tengo": &fstest.MapFile{Data: []byte(`export "second"`)}},
		}
		data, err := fs.ReadFile(overlap, "x.tengo")
		require.NoError(t, err)
		require.Equal(t, `export "first"`, string(data))
	})

	t.Run("missing file returns ErrNotExist", func(t *testing.T) {
		_, err := mfs.Open("missing.tengo")
		require.Error(t, err)
		require.True(t, errors.Is(err, fs.ErrNotExist))
	})

	t.Run("non-ErrNotExist error stops search immediately", func(t *testing.T) {
		// errFS always returns a sentinel error that is not ErrNotExist.
		// We place it first so that MultiFS must stop rather than try fsB.
		sentinel := errors.New("storage unavailable")
		errFS := errorFS{err: sentinel}
		layered := tengo.MultiFS{errFS, fsB}
		_, err := layered.Open("b.tengo")
		require.Error(t, err)
		require.True(t, errors.Is(err, sentinel))
		// b.tengo exists in fsB but must never be reached.
		require.False(t, errors.Is(err, fs.ErrNotExist))
	})

	t.Run("empty MultiFS returns ErrNotExist", func(t *testing.T) {
		_, err := tengo.MultiFS{}.Open("any.tengo")
		require.Error(t, err)
		require.True(t, errors.Is(err, fs.ErrNotExist))
	})
}

// errorFS is a test helper that always returns a fixed error from Open.
type errorFS struct{ err error }

func (e errorFS) Open(_ string) (fs.File, error) { return nil, e.err }

func TestMultiFS_ScriptIntegration(t *testing.T) {
	// Two independent FSes, each holding different modules.
	fsA := fstest.MapFS{
		"greet.tengo": &fstest.MapFile{
			Data: []byte(`export func(name) { return "hello, " + name }`),
		},
	}
	fsB := fstest.MapFS{
		"math/add.tengo": &fstest.MapFile{
			Data: []byte(`export func(a, b) { return a + b }`),
		},
		// Relative import: math/mul.tengo imports its sibling ./add.
		"math/mul.tengo": &fstest.MapFile{
			Data: []byte(`
add := import("./add")
export func(a, b) {
    result := 0
    for i := 0; i < b; i++ { result = add(result, a) }
    return result
}`),
		},
	}

	s := tengo.NewScript([]byte(`
greet   := import("greet")
add     := import("math/add")
mul     := import("math/mul")
out_greet := greet("world")
out_add   := add(3, 4)
out_mul   := mul(6, 7)
`))
	s.SetImportFS(tengo.MultiFS{fsA, fsB})

	compiled, err := s.Run()
	require.NoError(t, err)
	require.Equal(t, "hello, world", compiled.Get("out_greet").String())
	require.Equal(t, 7, compiled.Get("out_add").Int())
	require.Equal(t, 42, compiled.Get("out_mul").Int())
}

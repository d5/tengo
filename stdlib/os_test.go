package stdlib_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/require"
)

func TestReadFile(t *testing.T) {
	content := []byte("the quick brown fox jumps over the lazy dog")
	tf, err := ioutil.TempFile("", "test")
	require.NoError(t, err)
	defer func() { _ = os.Remove(tf.Name()) }()

	_, err = tf.Write(content)
	require.NoError(t, err)
	_ = tf.Close()

	module(t, "os").call("read_file", tf.Name()).
		expect(&tengo.Bytes{Value: content})
}

func TestReadFileArgs(t *testing.T) {
	module(t, "os").call("read_file").expectError()
}
func TestFileStatArgs(t *testing.T) {
	module(t, "os").call("stat").expectError()
}

func TestFileStatFile(t *testing.T) {
	content := []byte("the quick brown fox jumps over the lazy dog")
	tf, err := ioutil.TempFile("", "test")
	require.NoError(t, err)
	defer func() { _ = os.Remove(tf.Name()) }()

	_, err = tf.Write(content)
	require.NoError(t, err)
	_ = tf.Close()

	stat, err := os.Stat(tf.Name())
	if err != nil {
		t.Logf("could not get tmp file stat: %s", err)
		return
	}

	module(t, "os").call("stat", tf.Name()).expect(&tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"name":      &tengo.String{Value: stat.Name()},
			"mtime":     &tengo.Time{Value: stat.ModTime()},
			"size":      &tengo.Int{Value: stat.Size()},
			"mode":      &tengo.Int{Value: int64(stat.Mode())},
			"directory": tengo.FalseValue,
		},
	})
}

func TestFileStatDir(t *testing.T) {
	td, err := ioutil.TempDir("", "test")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(td) }()

	stat, err := os.Stat(td)
	require.NoError(t, err)

	module(t, "os").call("stat", td).expect(&tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"name":      &tengo.String{Value: stat.Name()},
			"mtime":     &tengo.Time{Value: stat.ModTime()},
			"size":      &tengo.Int{Value: stat.Size()},
			"mode":      &tengo.Int{Value: int64(stat.Mode())},
			"directory": tengo.TrueValue,
		},
	})
}

func TestOSExpandEnv(t *testing.T) {
	curMaxStringLen := tengo.MaxStringLen
	defer func() { tengo.MaxStringLen = curMaxStringLen }()
	tengo.MaxStringLen = 12

	_ = os.Setenv("TENGO", "FOO BAR")
	module(t, "os").call("expand_env", "$TENGO").expect("FOO BAR")

	_ = os.Setenv("TENGO", "FOO")
	module(t, "os").call("expand_env", "$TENGO $TENGO").expect("FOO FOO")

	_ = os.Setenv("TENGO", "123456789012")
	module(t, "os").call("expand_env", "$TENGO").expect("123456789012")

	_ = os.Setenv("TENGO", "1234567890123")
	module(t, "os").call("expand_env", "$TENGO").expectError()

	_ = os.Setenv("TENGO", "123456")
	module(t, "os").call("expand_env", "$TENGO$TENGO").expect("123456123456")

	_ = os.Setenv("TENGO", "123456")
	module(t, "os").call("expand_env", "${TENGO}${TENGO}").
		expect("123456123456")

	_ = os.Setenv("TENGO", "123456")
	module(t, "os").call("expand_env", "$TENGO $TENGO").expectError()

	_ = os.Setenv("TENGO", "123456")
	module(t, "os").call("expand_env", "${TENGO} ${TENGO}").expectError()
}

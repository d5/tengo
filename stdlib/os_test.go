package stdlib_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/d5/tengo"
	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/objects"
)

func TestReadFile(t *testing.T) {
	content := []byte("the quick brown fox jumps over the lazy dog")
	tf, err := ioutil.TempFile("", "test")
	if !assert.NoError(t, err) {
		return
	}
	defer func() { _ = os.Remove(tf.Name()) }()

	_, err = tf.Write(content)
	if !assert.NoError(t, err) {
		return
	}
	_ = tf.Close()

	module(t, "os").call("read_file", mockRuntime{}, tf.Name()).expect(&objects.Bytes{Value: content})
}

func TestReadFileArgs(t *testing.T) {
	module(t, "os").call("read_file", mockRuntime{}).expectError()
}
func TestFileStatArgs(t *testing.T) {
	module(t, "os").call("stat", mockRuntime{}).expectError()
}

func TestFileStatFile(t *testing.T) {
	content := []byte("the quick brown fox jumps over the lazy dog")
	tf, err := ioutil.TempFile("", "test")
	if !assert.NoError(t, err) {
		return
	}
	defer func() { _ = os.Remove(tf.Name()) }()

	_, err = tf.Write(content)
	if !assert.NoError(t, err) {
		return
	}
	_ = tf.Close()

	stat, err := os.Stat(tf.Name())
	if err != nil {
		t.Logf("could not get tmp file stat: %s", err)
		return
	}

	module(t, "os").call("stat", mockRuntime{}, tf.Name()).expect(&objects.ImmutableMap{
		Value: map[string]objects.Object{
			"name":      &objects.String{Value: stat.Name()},
			"mtime":     &objects.Time{Value: stat.ModTime()},
			"size":      &objects.Int{Value: stat.Size()},
			"mode":      &objects.Int{Value: int64(stat.Mode())},
			"directory": objects.FalseValue,
		},
	})
}

func TestFileStatDir(t *testing.T) {
	td, err := ioutil.TempDir("", "test")
	if !assert.NoError(t, err) {
		return
	}
	defer func() { _ = os.RemoveAll(td) }()

	stat, err := os.Stat(td)
	if !assert.NoError(t, err) {
		return
	}

	module(t, "os").call("stat", mockRuntime{}, td).expect(&objects.ImmutableMap{
		Value: map[string]objects.Object{
			"name":      &objects.String{Value: stat.Name()},
			"mtime":     &objects.Time{Value: stat.ModTime()},
			"size":      &objects.Int{Value: stat.Size()},
			"mode":      &objects.Int{Value: int64(stat.Mode())},
			"directory": objects.TrueValue,
		},
	})
}

func TestOSExpandEnv(t *testing.T) {
	curMaxStringLen := tengo.MaxStringLen
	defer func() { tengo.MaxStringLen = curMaxStringLen }()
	tengo.MaxStringLen = 12

	_ = os.Setenv("TENGO", "FOO BAR")
	module(t, "os").call("expand_env", mockRuntime{}, "$TENGO").expect("FOO BAR")

	_ = os.Setenv("TENGO", "FOO")
	module(t, "os").call("expand_env", mockRuntime{}, "$TENGO $TENGO").expect("FOO FOO")

	_ = os.Setenv("TENGO", "123456789012")
	module(t, "os").call("expand_env", mockRuntime{}, "$TENGO").expect("123456789012")

	_ = os.Setenv("TENGO", "1234567890123")
	module(t, "os").call("expand_env", mockRuntime{}, "$TENGO").expectError()

	_ = os.Setenv("TENGO", "123456")
	module(t, "os").call("expand_env", mockRuntime{}, "$TENGO$TENGO").expect("123456123456")

	_ = os.Setenv("TENGO", "123456")
	module(t, "os").call("expand_env", mockRuntime{}, "${TENGO}${TENGO}").expect("123456123456")

	_ = os.Setenv("TENGO", "123456")
	module(t, "os").call("expand_env", mockRuntime{}, "$TENGO $TENGO").expectError()

	_ = os.Setenv("TENGO", "123456")
	module(t, "os").call("expand_env", mockRuntime{}, "${TENGO} ${TENGO}").expectError()
}

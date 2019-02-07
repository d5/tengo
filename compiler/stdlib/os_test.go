package stdlib_test

import (
	"io/ioutil"
	"os"
	"testing"

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

	module(t, "os").call("read_file", tf.Name()).expect(&objects.Bytes{Value: content})
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

	module(t, "os").call("stat", tf.Name()).expect(&objects.ImmutableMap{
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

	module(t, "os").call("stat", td).expect(&objects.ImmutableMap{
		Value: map[string]objects.Object{
			"name":      &objects.String{Value: stat.Name()},
			"mtime":     &objects.Time{Value: stat.ModTime()},
			"size":      &objects.Int{Value: stat.Size()},
			"mode":      &objects.Int{Value: int64(stat.Mode())},
			"directory": objects.TrueValue,
		},
	})
}

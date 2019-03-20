package cli_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/cli"
	"github.com/d5/tengo/stdlib"
)

func TestCLICompileAndRun(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "tengo_tests")
	_ = os.MkdirAll(tempDir, os.ModePerm)
	binFile := filepath.Join(tempDir, "cli_bin")
	outFile := filepath.Join(tempDir, "cli_out")
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	src := []byte(`
os := import("os")
rand := import("rand")
times := import("times")

rand.seed(times.time_nanosecond(times.now()))

rand_num := func() {
	return rand.intn(100)
}

file := os.create("` + outFile + `")
file.write_string("random number is " + rand_num())
file.close()
`)

	mods := stdlib.GetModuleMap(stdlib.AllModuleNames()...)

	err := cli.CompileOnly(mods, src, "src", binFile)
	if !assert.NoError(t, err) {
		return
	}

	compiledBin, err := ioutil.ReadFile(binFile)
	if !assert.NoError(t, err) {
		return
	}

	err = cli.RunCompiled(mods, compiledBin)
	if !assert.NoError(t, err) {
		return
	}

	read, err := ioutil.ReadFile(outFile)
	if !assert.NoError(t, err) {
		return
	}
	ok, err := regexp.Match(`^random number is \d+$`, read)
	assert.NoError(t, err)
	assert.True(t, ok, string(read))
}

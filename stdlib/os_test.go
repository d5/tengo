package stdlib_test

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
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

func TestOSOpenFile(t *testing.T) {
	name := filepath.Join(t.TempDir(), "file.txt")
	module(t, "os").call("open_file", name, os.O_RDONLY, 0).
		expectErrorObject()

	f := module(t, "os").call("open_file", name, os.O_CREATE|os.O_RDWR, 0600)
	require.NoError(t, f.e)
	f.call("name").expect(name)
	f.call("write_string", "hello ").expect(6)
	f.call("write", []byte("world")).expect(5)
	f.call("sync").expect(true)

	f.call("seek", 0, io.SeekStart).expect(0)
	buf := &tengo.Bytes{Value: make([]byte, 11)}
	f.call("read", buf).expect(11)
	require.Equal(t, "hello world", string(buf.Value))
	f.call("read", buf).expectErrorObject()
	f.call("seek", -5, io.SeekEnd).expect(6)
	buf = &tengo.Bytes{Value: make([]byte, 5)}
	f.call("read", buf).expect(5)
	require.Equal(t, "world", string(buf.Value))
	f.call("seek", 0, 99).expectErrorObject()
	f.call("readdirnames", -1).expectErrorObject()

	f.call("chmod", 0640).expect(true)
	f.call("chown", os.Getuid(), os.Getgid()).expect(true)
	stat, err := os.Stat(name)
	require.NoError(t, err)
	require.Equal(t, int64(0640), int64(stat.Mode().Perm()))
	f.call("stat").expect(&tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"name":      &tengo.String{Value: "file.txt"},
			"mtime":     &tengo.Time{Value: stat.ModTime()},
			"size":      &tengo.Int{Value: 11},
			"mode":      &tengo.Int{Value: int64(stat.Mode())},
			"directory": tengo.FalseValue,
		},
	})

	f.call("close").expect(true)
	f.call("close").expectErrorObject()
	f.call("write_string", "x").expectErrorObject()

	module(t, "os").
		call("open_file", name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600).
		expectErrorObject()
	f = module(t, "os").call("open_file", name, os.O_APPEND|os.O_WRONLY, 0)
	f.call("write_string", "!").expect(1)
	f.call("close").expect(true)

	data, err := ioutil.ReadFile(name)
	require.NoError(t, err)
	require.Equal(t, "hello world!", string(data))
}

func TestOSFileDir(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, ioutil.WriteFile(filepath.Join(dir, "a.txt"), nil, 0600))

	d := module(t, "os").call("open", dir)
	require.NoError(t, d.e)
	defer d.call("close")
	d.call("readdirnames", -1).expect(ARR{"a.txt"})
	require.Equal(t, tengo.TrueValue,
		d.call("stat").o.(*tengo.ImmutableMap).Value["directory"])

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { require.NoError(t, os.Chdir(wd)) }()
	d.call("chdir").expect(true)
	cwd, err := os.Getwd()
	require.NoError(t, err)
	want, err := filepath.EvalSymlinks(dir)
	require.NoError(t, err)
	got, err := filepath.EvalSymlinks(cwd)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestOSMkdirChmod(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	nested := filepath.Join(dir, "a", "b")

	module(t, "os").call("mkdir", sub, 0700).expect(true)
	fi, err := os.Stat(sub)
	require.NoError(t, err)
	require.True(t, fi.IsDir())
	module(t, "os").call("mkdir", sub, 0700).expectErrorObject()

	module(t, "os").call("mkdir_all", nested, 0700).expect(true)
	module(t, "os").call("mkdir_all", nested, 0700).expect(true)
	fi, err = os.Stat(nested)
	require.NoError(t, err)
	require.True(t, fi.IsDir())

	module(t, "os").call("chmod", sub, 0500).expect(true)
	fi, err = os.Stat(sub)
	require.NoError(t, err)
	require.Equal(t, 0500, int(fi.Mode().Perm()))
	module(t, "os").call("chmod", filepath.Join(dir, "missing"), 0700).
		expectErrorObject()
}

func TestOSFileArgErrors(t *testing.T) {
	f := module(t, "os").call("create", filepath.Join(t.TempDir(), "f"))
	require.NoError(t, f.e)
	defer f.call("close")

	for _, c := range []argErrorCase{
		{"chmod", ARR{}, tengo.ErrWrongNumArguments},
		{"chmod", ARR{"x"}, invalidArg("first", "int(compatible)", "string")},
		{"seek", ARR{0}, tengo.ErrWrongNumArguments},
		{"seek", ARR{"x", 0}, invalidArg("first", "int(compatible)", "string")},
		{"seek", ARR{0, "x"}, invalidArg("second", "int(compatible)", "string")},
		{"stat", ARR{1}, tengo.ErrWrongNumArguments},
	} {
		f.call(c.fn, c.args...).expectErr(c.err, "file.%s%v", c.fn, c.args)
	}
}

func TestOSArgs(t *testing.T) {
	var want ARR
	for _, arg := range os.Args {
		want = append(want, arg)
	}
	module(t, "os").call("args").expect(want)

	curMaxStringLen := tengo.MaxStringLen
	defer func() { tengo.MaxStringLen = curMaxStringLen }()
	tengo.MaxStringLen = len(os.Args[0]) - 1
	module(t, "os").call("args").expectErr(tengo.ErrStringLimit)
}

func TestOSLookupEnv(t *testing.T) {
	const key = "TENGO_TEST_LOOKUP_ENV"
	t.Setenv(key, "value")
	module(t, "os").call("lookup_env", key).expect("value")
	t.Setenv(key, "")
	module(t, "os").call("lookup_env", key).expect("")
	require.NoError(t, os.Unsetenv(key))
	module(t, "os").call("lookup_env", key).expect(false)

	curMaxStringLen := tengo.MaxStringLen
	defer func() { tengo.MaxStringLen = curMaxStringLen }()
	tengo.MaxStringLen = 5
	t.Setenv(key, "12345")
	module(t, "os").call("lookup_env", key).expect("12345")
	t.Setenv(key, "123456")
	module(t, "os").call("lookup_env", key).expectErr(tengo.ErrStringLimit)
}

func TestOSFindProcess(t *testing.T) {
	p := module(t, "os").call("find_process", os.Getpid())
	require.NoError(t, p.e)
	for _, c := range []argErrorCase{
		{"signal", ARR{}, tengo.ErrWrongNumArguments},
		{"signal", ARR{"x"}, invalidArg("first", "int(compatible)", "string")},
		{"wait", ARR{1}, tengo.ErrWrongNumArguments},
	} {
		p.call(c.fn, c.args...).expectErr(c.err, "process.%s%v", c.fn, c.args)
	}
	p.call("signal", 0).expect(true)
	p.call("release").expect(true)
	p.call("signal", 0).expectErrorObject()
}

func TestOSStartProcess(t *testing.T) {
	dir := t.TempDir()
	p := module(t, "os").call("start_process", "/bin/sh",
		ARR{"sh", "-c", `printf "$FOO" > out.txt; exit 3`}, dir, ARR{"FOO=bar"})
	require.NoError(t, p.e)
	state := p.call("wait")
	require.NoError(t, state.e)
	state.call("exited").expect(true)
	state.call("success").expect(false)
	state.call("string").expect("exit status 3")
	pid := state.call("pid").o.(*tengo.Int).Value
	require.True(t, pid > 0 && pid != int64(os.Getpid()), pid)
	data, err := ioutil.ReadFile(filepath.Join(dir, "out.txt"))
	require.NoError(t, err)
	require.Equal(t, "bar", string(data))
	p.call("wait").expectErrorObject()

	p = module(t, "os").call("start_process", "/bin/sh",
		IARR{"sh", "-c", "exit 0"}, dir, IARR{})
	state = p.call("wait")
	state.call("success").expect(true)
	state.call("string").expect("exit status 0")

	sleep := ARR{"sh", "-c", "exec sleep 10"}
	env := ARR{"PATH=" + os.Getenv("PATH")}
	p = module(t, "os").call("start_process", "/bin/sh", sleep, dir, env)
	p.call("signal", int(syscall.SIGTERM)).expect(true)
	state = p.call("wait")
	state.call("exited").expect(false)
	state.call("success").expect(false)
	state.call("string").expect("signal: terminated")

	p = module(t, "os").call("start_process", "/bin/sh", sleep, dir, env)
	p.call("kill").expect(true)
	p.call("wait").call("string").expect("signal: killed")

	module(t, "os").call("start_process", filepath.Join(dir, "missing"),
		ARR{"missing"}, dir, ARR{}).expectErrorObject()
}

func TestOSArgErrors(t *testing.T) {
	wrongNum := tengo.ErrWrongNumArguments
	undefined := tengo.UndefinedValue
	invalidString := func(name string) error {
		return invalidArg(name, "string(compatible)", "undefined")
	}
	invalidInt := func(name string) error {
		return invalidArg(name, "int(compatible)", "string")
	}
	invalidElem := func(name string) error {
		return invalidArg(name, "string", "int")
	}
	expectArgErrors(t, "os", []argErrorCase{
		{"open", ARR{}, wrongNum},
		{"open", ARR{undefined}, invalidString("first")},
		{"create", ARR{}, wrongNum},
		{"create", ARR{undefined}, invalidString("first")},
		{"open_file", ARR{"f", 0}, wrongNum},
		{"open_file", ARR{undefined, 0, 0}, invalidString("first")},
		{"open_file", ARR{"f", "x", 0}, invalidInt("second")},
		{"open_file", ARR{"f", 0, "x"}, invalidInt("third")},
		{"mkdir", ARR{"d"}, wrongNum},
		{"mkdir", ARR{undefined, 0700}, invalidString("first")},
		{"mkdir", ARR{"d", "x"}, invalidInt("second")},
		{"mkdir_all", ARR{"d"}, wrongNum},
		{"chmod", ARR{"f", 0700, 0}, wrongNum},
		{"chmod", ARR{undefined, 0700}, invalidString("first")},
		{"chmod", ARR{"f", "x"}, invalidInt("second")},
		{"args", ARR{1}, wrongNum},
		{"lookup_env", ARR{}, wrongNum},
		{"lookup_env", ARR{undefined}, invalidString("first")},
		{"find_process", ARR{}, wrongNum},
		{"find_process", ARR{"x"}, invalidInt("first")},
		{"start_process", ARR{"sh", ARR{}, "."}, wrongNum},
		{"start_process", ARR{undefined, ARR{}, ".", ARR{}},
			invalidString("first")},
		{"start_process", ARR{"sh", "x", ".", ARR{}},
			invalidArg("second", "array", "string")},
		{"start_process", ARR{"sh", ARR{"sh", 1}, ".", ARR{}},
			invalidElem("second[1]")},
		{"start_process", ARR{"sh", IARR{1}, ".", ARR{}},
			invalidElem("second[0]")},
		{"start_process", ARR{"sh", ARR{}, undefined, ARR{}},
			invalidString("third")},
		{"start_process", ARR{"sh", ARR{}, ".", "x"},
			invalidArg("fourth", "array", "string")},
		{"start_process", ARR{"sh", ARR{}, ".", ARR{"A=b", 1}},
			invalidElem("fourth[1]")},
		{"start_process", ARR{"sh", ARR{}, ".", IARR{1}},
			invalidElem("fourth[0]")},
		{"exec", ARR{}, wrongNum},
		{"exec", ARR{undefined}, invalidString("first")},
		{"exec", ARR{"echo", "a", undefined}, invalidString("args[1]")},
	})
}

func TestOSExecCommand(t *testing.T) {
	dir := t.TempDir()
	sh := func(script string) callres {
		cmd := module(t, "os").call("exec", "/bin/sh", "-c", script)
		require.NoError(t, cmd.e)
		return cmd
	}

	sh("echo out; echo err 1>&2").call("combined_output").
		expect([]byte("out\nerr\n"))
	sh("exit 2").call("output").expectErrorObject()
	sh("exit 0").call("run").expect(true)
	sh("exit 1").call("run").
		expect(&tengo.Error{Value: &tengo.String{Value: "exit status 1"}})

	cmd := sh("echo hi > out.txt")
	cmd.call("set_dir", dir).expect(tengo.UndefinedValue)
	cmd.call("run").expect(true)
	data, err := ioutil.ReadFile(filepath.Join(dir, "out.txt"))
	require.NoError(t, err)
	require.Equal(t, "hi\n", string(data))

	cmd = sh("echo $FOO$BAR")
	cmd.call("set_env", ARR{"FOO=a", "BAR=b"}).expect(tengo.UndefinedValue)
	cmd.call("output").expect([]byte("ab\n"))
	cmd = sh("echo $FOO")
	cmd.call("set_env", IARR{"FOO=c"}).expect(tengo.UndefinedValue)
	cmd.call("output").expect([]byte("c\n"))

	cmd = module(t, "os").
		call("exec", filepath.Join(dir, "missing"), "-c", "echo replaced")
	cmd.call("set_path", "/bin/sh").expect(tengo.UndefinedValue)
	cmd.call("output").expect([]byte("replaced\n"))

	cmd = sh("exec sleep 10")
	cmd.call("set_env", ARR{"PATH=" + os.Getenv("PATH")})
	cmd.call("start").expect(true)
	cmd.call("process").call("kill").expect(true)
	cmd.call("wait").
		expect(&tengo.Error{Value: &tengo.String{Value: "signal: killed"}})

	cmd = sh("true")
	for _, c := range []argErrorCase{
		{"set_path", ARR{}, tengo.ErrWrongNumArguments},
		{"set_path", ARR{tengo.UndefinedValue},
			invalidArg("first", "string(compatible)", "undefined")},
		{"set_dir", ARR{}, tengo.ErrWrongNumArguments},
		{"set_dir", ARR{tengo.UndefinedValue},
			invalidArg("first", "string(compatible)", "undefined")},
		{"set_env", ARR{}, tengo.ErrWrongNumArguments},
		{"set_env", ARR{"x"}, invalidArg("first", "array", "string")},
		{"set_env", ARR{ARR{1}}, invalidArg("first[0]", "string", "int")},
		{"set_env", ARR{IARR{1}}, invalidArg("first[0]", "string", "int")},
		{"process", ARR{1}, tengo.ErrWrongNumArguments},
	} {
		cmd.call(c.fn, c.args...).expectErr(c.err, "cmd.%s%v", c.fn, c.args)
	}
}

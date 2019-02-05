package stdlib

import (
	"io"
	"os"
	"os/exec"

	"github.com/d5/tengo/objects"
)

var osModule = map[string]objects.Object{
	"o_rdonly":            &objects.Int{Value: int64(os.O_RDONLY)},
	"o_wronly":            &objects.Int{Value: int64(os.O_WRONLY)},
	"o_rdwr":              &objects.Int{Value: int64(os.O_RDWR)},
	"o_append":            &objects.Int{Value: int64(os.O_APPEND)},
	"o_create":            &objects.Int{Value: int64(os.O_CREATE)},
	"o_excl":              &objects.Int{Value: int64(os.O_EXCL)},
	"o_sync":              &objects.Int{Value: int64(os.O_SYNC)},
	"o_trunc":             &objects.Int{Value: int64(os.O_TRUNC)},
	"mode_dir":            &objects.Int{Value: int64(os.ModeDir)},
	"mode_append":         &objects.Int{Value: int64(os.ModeAppend)},
	"mode_exclusive":      &objects.Int{Value: int64(os.ModeExclusive)},
	"mode_temporary":      &objects.Int{Value: int64(os.ModeTemporary)},
	"mode_symlink":        &objects.Int{Value: int64(os.ModeSymlink)},
	"mode_device":         &objects.Int{Value: int64(os.ModeDevice)},
	"mode_named_pipe":     &objects.Int{Value: int64(os.ModeNamedPipe)},
	"mode_socket":         &objects.Int{Value: int64(os.ModeSocket)},
	"mode_setuid":         &objects.Int{Value: int64(os.ModeSetuid)},
	"mode_setgui":         &objects.Int{Value: int64(os.ModeSetgid)},
	"mode_char_device":    &objects.Int{Value: int64(os.ModeCharDevice)},
	"mode_sticky":         &objects.Int{Value: int64(os.ModeSticky)},
	"mode_type":           &objects.Int{Value: int64(os.ModeType)},
	"mode_perm":           &objects.Int{Value: int64(os.ModePerm)},
	"path_separator":      &objects.Char{Value: os.PathSeparator},
	"path_list_separator": &objects.Char{Value: os.PathListSeparator},
	"dev_null":            &objects.String{Value: os.DevNull},
	"seek_set":            &objects.Int{Value: int64(io.SeekStart)},
	"seek_cur":            &objects.Int{Value: int64(io.SeekCurrent)},
	"seek_end":            &objects.Int{Value: int64(io.SeekEnd)},
	"args":                &objects.UserFunction{Value: osArgs},         // args() => array(string)
	"chdir":               FuncASRE(os.Chdir),                           // chdir(dir string) => error
	"chmod":               osFuncASFmRE(os.Chmod),                       // chmod(name string, mode int) => error
	"chown":               FuncASIIRE(os.Chown),                         // chown(name string, uid int, gid int) => error
	"clearenv":            FuncAR(os.Clearenv),                          // clearenv()
	"environ":             FuncARSs(os.Environ),                         // environ() => array(string)
	"exit":                FuncAIR(os.Exit),                             // exit(code int)
	"expand_env":          FuncASRS(os.ExpandEnv),                       // expand_env(s string) => string
	"getegid":             FuncARI(os.Getegid),                          // getegid() => int
	"getenv":              FuncASRS(os.Getenv),                          // getenv(s string) => string
	"geteuid":             FuncARI(os.Geteuid),                          // geteuid() => int
	"getgid":              FuncARI(os.Getgid),                           // getgid() => int
	"getgroups":           FuncARIsE(os.Getgroups),                      // getgroups() => array(string)/error
	"getpagesize":         FuncARI(os.Getpagesize),                      // getpagesize() => int
	"getpid":              FuncARI(os.Getpid),                           // getpid() => int
	"getppid":             FuncARI(os.Getppid),                          // getppid() => int
	"getuid":              FuncARI(os.Getuid),                           // getuid() => int
	"getwd":               FuncARSE(os.Getwd),                           // getwd() => string/error
	"hostname":            FuncARSE(os.Hostname),                        // hostname() => string/error
	"lchown":              FuncASIIRE(os.Lchown),                        // lchown(name string, uid int, gid int) => error
	"link":                FuncASSRE(os.Link),                           // link(oldname string, newname string) => error
	"lookup_env":          &objects.UserFunction{Value: osLookupEnv},    // lookup_env(key string) => string/false
	"mkdir":               osFuncASFmRE(os.Mkdir),                       // mkdir(name string, perm int) => error
	"mkdir_all":           osFuncASFmRE(os.MkdirAll),                    // mkdir_all(name string, perm int) => error
	"readlink":            FuncASRSE(os.Readlink),                       // readlink(name string) => string/error
	"remove":              FuncASRE(os.Remove),                          // remove(name string) => error
	"remove_all":          FuncASRE(os.RemoveAll),                       // remove_all(name string) => error
	"rename":              FuncASSRE(os.Rename),                         // rename(oldpath string, newpath string) => error
	"setenv":              FuncASSRE(os.Setenv),                         // setenv(key string, value string) => error
	"symlink":             FuncASSRE(os.Symlink),                        // symlink(oldname string newname string) => error
	"temp_dir":            FuncARS(os.TempDir),                          // temp_dir() => string
	"truncate":            FuncASI64RE(os.Truncate),                     // truncate(name string, size int) => error
	"unsetenv":            FuncASRE(os.Unsetenv),                        // unsetenv(key string) => error
	"create":              &objects.UserFunction{Value: osCreate},       // create(name string) => imap(file)/error
	"open":                &objects.UserFunction{Value: osOpen},         // open(name string) => imap(file)/error
	"open_file":           &objects.UserFunction{Value: osOpenFile},     // open_file(name string, flag int, perm int) => imap(file)/error
	"find_process":        &objects.UserFunction{Value: osFindProcess},  // find_process(pid int) => imap(process)/error
	"start_process":       &objects.UserFunction{Value: osStartProcess}, // start_process(name string, argv array(string), dir string, env array(string)) => imap(process)/error
	"exec_look_path":      FuncASRSE(exec.LookPath),                     // exec_look_path(file) => string/error
	"exec":                &objects.UserFunction{Value: osExec},         // exec(name, args...) => command
	"stat":                &objects.UserFunction{Value: osStat},         // stat(name) => imap(fileinfo)/error
}

func osStat(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}

	fname, ok := objects.ToString(args[0])
	if !ok {
		return nil, objects.ErrInvalidTypeConversion
	}

	stat, err := os.Stat(fname)
	if err != nil {
		return wrapError(err), nil
	}

	fstat := &objects.ImmutableMap{
		Value: map[string]objects.Object{
			"name":  &objects.String{Value: stat.Name()},
			"mtime": &objects.Time{Value: stat.ModTime()},
			"size":  &objects.Int{Value: stat.Size()},
			"mode":  &objects.Int{Value: int64(stat.Mode())},
		},
	}

	if stat.IsDir() {
		fstat.Value["directory"] = objects.TrueValue
	} else {
		fstat.Value["directory"] = objects.FalseValue
	}

	return fstat, nil
}

func osCreate(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}

	s1, ok := objects.ToString(args[0])
	if !ok {
		return nil, objects.ErrInvalidTypeConversion
	}

	res, err := os.Create(s1)
	if err != nil {
		return wrapError(err), nil
	}

	return makeOSFile(res), nil
}

func osOpen(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}

	s1, ok := objects.ToString(args[0])
	if !ok {
		return nil, objects.ErrInvalidTypeConversion
	}

	res, err := os.Open(s1)
	if err != nil {
		return wrapError(err), nil
	}

	return makeOSFile(res), nil
}

func osOpenFile(args ...objects.Object) (objects.Object, error) {
	if len(args) != 3 {
		return nil, objects.ErrWrongNumArguments
	}

	s1, ok := objects.ToString(args[0])
	if !ok {
		return nil, objects.ErrInvalidTypeConversion
	}

	i2, ok := objects.ToInt(args[1])
	if !ok {
		return nil, objects.ErrInvalidTypeConversion
	}

	i3, ok := objects.ToInt(args[2])
	if !ok {
		return nil, objects.ErrInvalidTypeConversion
	}

	res, err := os.OpenFile(s1, i2, os.FileMode(i3))
	if err != nil {
		return wrapError(err), nil
	}

	return makeOSFile(res), nil
}

func osArgs(args ...objects.Object) (objects.Object, error) {
	if len(args) != 0 {
		return nil, objects.ErrWrongNumArguments
	}

	arr := &objects.Array{}
	for _, osArg := range os.Args {
		arr.Value = append(arr.Value, &objects.String{Value: osArg})
	}

	return arr, nil
}

func osFuncASFmRE(fn func(string, os.FileMode) error) *objects.UserFunction {
	return &objects.UserFunction{
		Value: func(args ...objects.Object) (objects.Object, error) {
			if len(args) != 2 {
				return nil, objects.ErrWrongNumArguments
			}

			s1, ok := objects.ToString(args[0])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}
			i2, ok := objects.ToInt64(args[1])
			if !ok {
				return nil, objects.ErrInvalidTypeConversion
			}

			return wrapError(fn(s1, os.FileMode(i2))), nil
		},
	}
}

func osLookupEnv(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}

	s1, ok := objects.ToString(args[0])
	if !ok {
		return nil, objects.ErrInvalidTypeConversion
	}

	res, ok := os.LookupEnv(s1)
	if !ok {
		return objects.FalseValue, nil
	}

	return &objects.String{Value: res}, nil
}

func osExec(args ...objects.Object) (objects.Object, error) {
	if len(args) == 0 {
		return nil, objects.ErrWrongNumArguments
	}

	name, ok := objects.ToString(args[0])
	if !ok {
		return nil, objects.ErrInvalidTypeConversion
	}

	var execArgs []string
	for _, arg := range args[1:] {
		execArg, ok := objects.ToString(arg)
		if !ok {
			return nil, objects.ErrInvalidTypeConversion
		}

		execArgs = append(execArgs, execArg)
	}

	return makeOSExecCommand(exec.Command(name, execArgs...)), nil
}

func osFindProcess(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}

	i1, ok := objects.ToInt(args[0])
	if !ok {
		return nil, objects.ErrInvalidTypeConversion
	}

	proc, err := os.FindProcess(i1)
	if err != nil {
		return wrapError(err), nil
	}

	return makeOSProcess(proc), nil
}

func osStartProcess(args ...objects.Object) (objects.Object, error) {
	if len(args) != 4 {
		return nil, objects.ErrWrongNumArguments
	}

	name, ok := objects.ToString(args[0])
	if !ok {
		return nil, objects.ErrInvalidTypeConversion
	}

	argv, err := stringArray(args[1])
	if err != nil {
		return nil, err
	}

	dir, ok := objects.ToString(args[2])
	if !ok {
		return nil, objects.ErrInvalidTypeConversion
	}

	env, err := stringArray(args[3])
	if err != nil {
		return nil, err
	}

	proc, err := os.StartProcess(name, argv, &os.ProcAttr{
		Dir: dir,
		Env: env,
	})
	if err != nil {
		return wrapError(err), nil
	}

	return makeOSProcess(proc), nil
}

func stringArray(o objects.Object) ([]string, error) {
	arr, ok := o.(*objects.Array)
	if !ok {
		return nil, objects.ErrInvalidTypeConversion
	}

	var sarr []string
	for _, elem := range arr.Value {
		str, ok := elem.(*objects.String)
		if !ok {
			return nil, objects.ErrInvalidTypeConversion
		}

		sarr = append(sarr, str.Value)
	}

	return sarr, nil
}

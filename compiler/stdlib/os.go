package stdlib

import (
	"io"
	"os"

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
	// args() => array(string)
	"args": &objects.UserFunction{Value: osArgs},
	// chdir(dir string) => error
	"chdir": FuncASRE(os.Chdir),
	// chmod(name string, mode int) => error
	"chmod": osFuncASFmRE(os.Chmod),
	// chown(name string, uid int, gid int) => error
	"chown": FuncASIIRE(os.Chown),
	// clearenv()
	"clearenv": FuncAR(os.Clearenv),
	// environ() => array(string)
	"environ": FuncARSs(os.Environ),
	// exit(code int)
	"exit": FuncAIR(os.Exit),
	// expand_env(s string) => string
	"expand_env": FuncASRS(os.ExpandEnv),
	// getegid() => int
	"getegid": FuncARI(os.Getegid),
	// getenv(s string) => string
	"getenv": FuncASRS(os.Getenv),
	// geteuid() => int
	"geteuid": FuncARI(os.Geteuid),
	// getgid() => int
	"getgid": FuncARI(os.Getgid),
	// getgroups() => array(string)/error
	"getgroups": FuncARIsE(os.Getgroups),
	// getpagesize() => int
	"getpagesize": FuncARI(os.Getpagesize),
	// getpid() => int
	"getpid": FuncARI(os.Getpid),
	// getppid() => int
	"getppid": FuncARI(os.Getppid),
	// getuid() => int
	"getuid": FuncARI(os.Getuid),
	// getwd() => string/error
	"getwd": FuncARSE(os.Getwd),
	// hostname() => string/error
	"hostname": FuncARSE(os.Hostname),
	// lchown(name string, uid int, gid int) => error
	"lchown": FuncASIIRE(os.Lchown),
	// link(oldname string, newname string) => error
	"link": FuncASSRE(os.Link),
	// lookup_env(key string) => string/false
	"lookup_env": &objects.UserFunction{Value: osLookupEnv},
	// mkdir(name string, perm int) => error
	"mkdir": osFuncASFmRE(os.Mkdir),
	// mkdir_all(name string, perm int) => error
	"mkdir_all": osFuncASFmRE(os.MkdirAll),
	// readlink(name string) => string/error
	"readlink": FuncASRSE(os.Readlink),
	// remove(name string) => error
	"remove": FuncASRE(os.Remove),
	// remove_all(name string) => error
	"remove_all": FuncASRE(os.RemoveAll),
	// rename(oldpath string, newpath string) => error
	"rename": FuncASSRE(os.Rename),
	// setenv(key string, value string) => error
	"setenv": FuncASSRE(os.Setenv),
	// symlink(oldname string newname string) => error
	"symlink": FuncASSRE(os.Symlink),
	// temp_dir() => string
	"temp_dir": FuncARS(os.TempDir),
	// truncate(name string, size int) => error
	"truncate": FuncASI64RE(os.Truncate),
	// unsetenv(key string) => error
	"unsetenv": FuncASRE(os.Unsetenv),
	// create(name string) => imap(file)/error
	"create": &objects.UserFunction{Value: osCreate},
	// open(name string) => imap(file)/error
	"open": &objects.UserFunction{Value: osOpen},
	// open_file(name string, flag int, perm int) => imap(file)/error
	"open_file": &objects.UserFunction{Value: osOpenFile},
	// find_process(pid int) => imap(process)/error
	"find_process": &objects.UserFunction{Value: osFindProcess},
	// start_process(name string, argv array(string), dir string, env array(string)) => imap(process)/error
	"start_process": &objects.UserFunction{Value: osStartProcess},
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

func osExecutable(args ...objects.Object) (objects.Object, error) {
	if len(args) != 0 {
		return nil, objects.ErrWrongNumArguments
	}

	res, err := os.Executable()
	if err != nil {
		return wrapError(err), nil
	}

	return &objects.String{Value: res}, nil
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

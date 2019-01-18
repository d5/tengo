package stdlib

import (
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
	"mode_irregular":      &objects.Int{Value: int64(os.ModeIrregular)},
	"mode_type":           &objects.Int{Value: int64(os.ModeType)},
	"mode_perm":           &objects.Int{Value: int64(os.ModePerm)},
	"path_separator":      &objects.Char{Value: os.PathSeparator},
	"path_list_separator": &objects.Char{Value: os.PathListSeparator},
	"dev_null":            &objects.String{Value: os.DevNull},
	"args":                &objects.UserFunction{Value: osArgs},
	"chdir":               FuncASRE(os.Chdir),
	"chmod":               osFuncASFmRE(os.Chmod),
	"chown":               FuncASIIRE(os.Chown),
	"clearenv":            FuncAR(os.Clearenv),
	"environ":             FuncARSs(os.Environ),
	"executable":          &objects.UserFunction{Value: osExecutable},
	"exit":                FuncAIR(os.Exit),
	"expand_env":          FuncASRS(os.ExpandEnv),
	"getegid":             FuncARI(os.Getegid),
	"getenv":              FuncASRS(os.Getenv),
	"geteuid":             FuncARI(os.Geteuid),
	"getgid":              FuncARI(os.Getgid),
	"getgroups":           FuncARIsE(os.Getgroups),
	"getpagesize":         FuncARI(os.Getpagesize),
	"getpid":              FuncARI(os.Getpid),
	"getppid":             FuncARI(os.Getppid),
	"getuid":              FuncARI(os.Getuid),
	"getwd":               FuncARSE(os.Getwd),
	"hostname":            FuncARSE(os.Hostname),
	"lchown":              FuncASIIRE(os.Lchown),
	"link":                FuncASSRE(os.Link),
	"lookup_env":          &objects.UserFunction{Value: osLookupEnv},
	"mkdir":               osFuncASFmRE(os.Mkdir),
	"mkdir_all":           osFuncASFmRE(os.MkdirAll),
	"readlink":            FuncASRSE(os.Readlink),
	"remove":              FuncASRE(os.Remove),
	"remove_all":          FuncASRE(os.RemoveAll),
	"rename":              FuncASSRE(os.Rename),
	"setenv":              FuncASSRE(os.Setenv),
	"symlink":             FuncASSRE(os.Symlink),
	"temp_dir":            FuncARS(os.TempDir),
	"truncate":            FuncASI64RE(os.Truncate),
	"unsetenv":            FuncASRE(os.Unsetenv),
	"user_cache_dir":      FuncARSE(os.UserCacheDir),
	"create":              &objects.UserFunction{Value: osCreate},
	"open":                &objects.UserFunction{Value: osOpen},
	"open_file":           &objects.UserFunction{Value: osOpenFile},

	// TODO: not implemented yet
	//"stdin":         nil,
	//"stdout":        nil,
	//"stderr":        nil,
	//"chtimes":       nil,
	//"expand":        nil,
	//"is_exists":     nil,
	//"is_not_exist":  nil,
	//"is_path_separator": nil,
	//"is_permission": nil,
	//"is_timeout": nil,
	//"new_syscall_error": nil,
	//"pipe": nil,
	//"same_file": nil,
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

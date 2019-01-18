package stdmods

import (
	"os"

	"github.com/d5/tengo/objects"
)

var osModule = map[string]objects.Object{
	"o_rdonly": &objects.Int{Value: int64(os.O_RDONLY)},
	"o_wronly": &objects.Int{Value: int64(os.O_WRONLY)},
	"o_rdwr":   &objects.Int{Value: int64(os.O_RDWR)},
	"o_append": &objects.Int{Value: int64(os.O_APPEND)},
	"o_create": &objects.Int{Value: int64(os.O_CREATE)},
	"o_excl":   &objects.Int{Value: int64(os.O_EXCL)},
	"o_sync":   &objects.Int{Value: int64(os.O_SYNC)},
	"o_trunc":  &objects.Int{Value: int64(os.O_TRUNC)},

	"mode_dir":         &objects.Int{Value: int64(os.ModeDir)},
	"mode_append":      &objects.Int{Value: int64(os.ModeAppend)},
	"mode_exclusive":   &objects.Int{Value: int64(os.ModeExclusive)},
	"mode_temporary":   &objects.Int{Value: int64(os.ModeTemporary)},
	"mode_symlink":     &objects.Int{Value: int64(os.ModeSymlink)},
	"mode_device":      &objects.Int{Value: int64(os.ModeDevice)},
	"mode_named_pipe":  &objects.Int{Value: int64(os.ModeNamedPipe)},
	"mode_socket":      &objects.Int{Value: int64(os.ModeSocket)},
	"mode_setuid":      &objects.Int{Value: int64(os.ModeSetuid)},
	"mode_setgui":      &objects.Int{Value: int64(os.ModeSetgid)},
	"mode_char_device": &objects.Int{Value: int64(os.ModeCharDevice)},
	"mode_sticky":      &objects.Int{Value: int64(os.ModeSticky)},
	"mode_irregular":   &objects.Int{Value: int64(os.ModeIrregular)},
	"mode_type":        &objects.Int{Value: int64(os.ModeType)},
	"mode_perm":        &objects.Int{Value: int64(os.ModePerm)},

	"path_separator":      &objects.Char{Value: os.PathSeparator},
	"path_list_separator": &objects.Char{Value: os.PathListSeparator},
	"dev_null":            &objects.String{Value: os.DevNull},

	"args":     &objects.UserFunction{Value: osArgs},
	"chdir":    FuncASRE(os.Chdir),
	"chmod":    &objects.UserFunction{Value: osChmod},
	"chown":    &objects.UserFunction{Value: osChown},
	"clearenv": FuncAR(os.Clearenv),
	"environ":  FuncARSs(os.Environ),

	// TODO: system errors
	//"err_invalid":         &objects.Error{Value: os.ErrInvalid.Error()},
	// TODO: STDIN, STDOUT, STDERR
	// "stdin": nil,
	// "stdout": nil,
	// "stderr": nil,
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

func osChmod(args ...objects.Object) (objects.Object, error) {
	if len(args) != 2 {
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

	return wrapError(os.Chmod(s1, os.FileMode(i2))), nil
}

func osChown(args ...objects.Object) (objects.Object, error) {
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

	return wrapError(os.Chown(s1, i2, i3)), nil
}

package stdlib

import (
	"os"

	"github.com/d5/tengo/objects"
)

func osFileImmutableMap(file *os.File) *objects.ImmutableMap {
	return &objects.ImmutableMap{
		Value: map[string]objects.Object{
			// chdir() => true/error
			"chdir": FuncARE(file.Chdir),
			// chown(uid int, gid int) => true/error
			"chown": FuncAIIRE(file.Chown),
			// close() => error
			"close": FuncARE(file.Close),
			// name() => string
			"name": FuncARS(file.Name),
			// readdirnames(n int) => array(string)/error
			"readdirnames": FuncAIRSsE(file.Readdirnames),
			// sync() => error
			"sync": FuncARE(file.Sync),
			// write(bytes) => int/error
			"write": FuncAYRIE(file.Write),
			// write(string) => int/error
			"write_string": FuncASRIE(file.WriteString),
			// read(bytes) => int/error
			"read": FuncAYRIE(file.Read),
			// chmod(mode int) => error
			"chmod": &objects.UserFunction{
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					if len(args) != 1 {
						return nil, objects.ErrWrongNumArguments
					}

					i1, ok := objects.ToInt64(args[0])
					if !ok {
						return nil, objects.ErrInvalidTypeConversion
					}

					return wrapError(file.Chmod(os.FileMode(i1))), nil
				},
			},
			// seek(offset int, whence int) => int/error
			"seek": &objects.UserFunction{
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					if len(args) != 2 {
						return nil, objects.ErrWrongNumArguments
					}

					i1, ok := objects.ToInt64(args[0])
					if !ok {
						return nil, objects.ErrInvalidTypeConversion
					}
					i2, ok := objects.ToInt(args[1])
					if !ok {
						return nil, objects.ErrInvalidTypeConversion
					}

					res, err := file.Seek(i1, i2)
					if err != nil {
						return wrapError(err), nil
					}

					return &objects.Int{Value: res}, nil
				},
			},
			// TODO: implement more functions
			//"fd": nil,
			//"read_at": nil,
			//"readdir": nil,
			//"set_deadline": nil,
			//"set_read_deadline": nil,
			//"set_write_deadline": nil,
			//"stat": nil,
			//"write_at": nil,
		},
	}
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

	return osFileImmutableMap(res), nil
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

	return osFileImmutableMap(res), nil
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

	return osFileImmutableMap(res), nil
}

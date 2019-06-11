package stdlib

import (
	"os"

	"github.com/d5/tengo"
)

func makeOSFile(file *os.File) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			// chdir() => true/error
			"chdir": &tengo.GoFunction{Name: "chdir", Value: FuncARE(file.Chdir)}, //
			// chown(uid int, gid int) => true/error
			"chown": &tengo.GoFunction{Name: "chown", Value: FuncAIIRE(file.Chown)}, //
			// close() => error
			"close": &tengo.GoFunction{Name: "close", Value: FuncARE(file.Close)}, //
			// name() => string
			"name": &tengo.GoFunction{Name: "name", Value: FuncARS(file.Name)}, //
			// readdirnames(n int) => array(string)/error
			"readdirnames": &tengo.GoFunction{Name: "readdirnames", Value: FuncAIRSsE(file.Readdirnames)}, //
			// sync() => error
			"sync": &tengo.GoFunction{Name: "sync", Value: FuncARE(file.Sync)}, //
			// write(bytes) => int/error
			"write": &tengo.GoFunction{Name: "write", Value: FuncAYRIE(file.Write)}, //
			// write(string) => int/error
			"write_string": &tengo.GoFunction{Name: "write_string", Value: FuncASRIE(file.WriteString)}, //
			// read(bytes) => int/error
			"read": &tengo.GoFunction{Name: "read", Value: FuncAYRIE(file.Read)}, //
			// chmod(mode int) => error
			"chmod": &tengo.GoFunction{
				Name: "chmod",
				Value: func(_ tengo.Interop, args ...tengo.Object) (ret tengo.Object, err error) {
					if len(args) != 1 {
						return nil, tengo.ErrWrongNumArguments
					}

					i1, ok := tengo.ToInt64(args[0])
					if !ok {
						return nil, tengo.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "int(compatible)",
							Found:    args[0].TypeName(),
						}
					}

					return wrapError(file.Chmod(os.FileMode(i1))), nil
				},
			},
			// seek(offset int, whence int) => int/error
			"seek": &tengo.GoFunction{
				Name: "seek",
				Value: func(_ tengo.Interop, args ...tengo.Object) (ret tengo.Object, err error) {
					if len(args) != 2 {
						return nil, tengo.ErrWrongNumArguments
					}

					i1, ok := tengo.ToInt64(args[0])
					if !ok {
						return nil, tengo.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "int(compatible)",
							Found:    args[0].TypeName(),
						}
					}
					i2, ok := tengo.ToInt(args[1])
					if !ok {
						return nil, tengo.ErrInvalidArgumentType{
							Name:     "second",
							Expected: "int(compatible)",
							Found:    args[1].TypeName(),
						}
					}

					res, err := file.Seek(i1, i2)
					if err != nil {
						return wrapError(err), nil
					}

					return &tengo.Int{Value: res}, nil
				},
			},
			// stat() => imap(fileinfo)/error
			"stat": &tengo.GoFunction{
				Name: "start",
				Value: func(rt tengo.Interop, args ...tengo.Object) (ret tengo.Object, err error) {
					if len(args) != 0 {
						return nil, tengo.ErrWrongNumArguments
					}

					return osStat(rt, &tengo.String{Value: file.Name()})
				},
			},
		},
	}
}

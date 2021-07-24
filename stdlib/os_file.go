package stdlib

import (
	"os"

	"github.com/d5/tengo/v2"
)

func makeOSFile(file *os.File) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			// chdir() => true/error
			"chdir": &tengo.UserFunction{
				Name:  "chdir",
				Value: FuncARE(file.Chdir),
			}, //
			// chown(uid int, gid int) => true/error
			"chown": &tengo.UserFunction{
				Name:  "chown",
				Value: FuncAIIRE(file.Chown),
			}, //
			// close() => error
			"close": &tengo.UserFunction{
				Name:  "close",
				Value: FuncARE(file.Close),
			}, //
			// name() => string
			"name": &tengo.UserFunction{
				Name:  "name",
				Value: FuncARS(file.Name),
			}, //
			// readdirnames(n int) => array(string)/error
			"readdirnames": &tengo.UserFunction{
				Name:  "readdirnames",
				Value: FuncAIRSsE(file.Readdirnames),
			}, //
			// sync() => error
			"sync": &tengo.UserFunction{
				Name:  "sync",
				Value: FuncARE(file.Sync),
			}, //
			// write(bytes) => int/error
			"write": &tengo.UserFunction{
				Name:  "write",
				Value: FuncAYRIE(file.Write),
			}, //
			// write(string) => int/error
			"write_string": &tengo.UserFunction{
				Name:  "write_string",
				Value: FuncASRIE(file.WriteString),
			}, //
			// read(bytes) => int/error
			"read": &tengo.UserFunction{
				Name:  "read",
				Value: FuncAYRIE(file.Read),
			}, //
			// chmod(mode int) => error
			"chmod": &tengo.UserFunction{
				Name: "chmod",
				Value: tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
					i1, err := tengo.ToInt64(0, args...)
					if err != nil {
						return nil, err
					}
					return wrapError(file.Chmod(os.FileMode(i1))), nil
				}, 1),
			},
			// seek(offset int, whence int) => int/error
			"seek": &tengo.UserFunction{
				Name: "seek",
				Value: tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
					i1, err := tengo.ToInt64(0, args...)
					if err != nil {
						return nil, err
					}
					i2, err := tengo.ToInt(1, args...)
					if err != nil {
						return nil, err
					}
					res, err := file.Seek(i1, i2)
					if err != nil {
						return wrapError(err), nil
					}
					return &tengo.Int{Value: res}, nil
				}, 2),
			},
			// stat() => imap(fileinfo)/error
			"stat": &tengo.UserFunction{
				Name: "stat",
				Value: tengo.CheckStrictArgs(func(args ...tengo.Object) (tengo.Object, error) {
					return osStat(&tengo.String{Value: file.Name()})
				}),
			},
		},
	}
}

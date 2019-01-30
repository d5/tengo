package stdlib

import (
	"os/exec"

	"github.com/d5/tengo/objects"
)

func makeOSExecCommand(cmd *exec.Cmd) *objects.ImmutableMap {
	return &objects.ImmutableMap{
		Value: map[string]objects.Object{
			// combined_output() => bytes/error
			"combined_output": FuncARYE(cmd.CombinedOutput),
			// output() => bytes/error
			"output": FuncARYE(cmd.Output),
			// run() => error
			"run": FuncARE(cmd.Run),
			// start() => error
			"start": FuncARE(cmd.Start),
			// wait() => error
			"wait": FuncARE(cmd.Wait),
			// set_path(path string)
			"set_path": &objects.UserFunction{
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					if len(args) != 1 {
						return nil, objects.ErrWrongNumArguments
					}

					s1, ok := objects.ToString(args[0])
					if !ok {
						return nil, objects.ErrInvalidTypeConversion
					}

					cmd.Path = s1

					return objects.UndefinedValue, nil
				},
			},
			// set_dir(dir string)
			"set_dir": &objects.UserFunction{
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					if len(args) != 1 {
						return nil, objects.ErrWrongNumArguments
					}

					s1, ok := objects.ToString(args[0])
					if !ok {
						return nil, objects.ErrInvalidTypeConversion
					}

					cmd.Dir = s1

					return objects.UndefinedValue, nil
				},
			},
			// set_env(env array(string))
			"set_env": &objects.UserFunction{
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					if len(args) != 1 {
						return nil, objects.ErrWrongNumArguments
					}

					envs, err := stringArray(args[0])
					if err != nil {
						return nil, err
					}

					cmd.Env = envs

					return objects.UndefinedValue, nil
				},
			},
			// process() => imap(process)
			"process": &objects.UserFunction{
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					if len(args) != 0 {
						return nil, objects.ErrWrongNumArguments
					}

					return makeOSProcess(cmd.Process), nil
				},
			},
		},
	}
}

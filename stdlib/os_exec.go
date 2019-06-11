package stdlib

import (
	"os/exec"

	"github.com/d5/tengo"
)

func makeOSExecCommand(cmd *exec.Cmd) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			// combined_output() => bytes/error
			"combined_output": &tengo.GoFunction{Name: "combined_output", Value: FuncARYE(cmd.CombinedOutput)}, //
			// output() => bytes/error
			"output": &tengo.GoFunction{Name: "output", Value: FuncARYE(cmd.Output)}, //
			// run() => error
			"run": &tengo.GoFunction{Name: "run", Value: FuncARE(cmd.Run)}, //
			// start() => error
			"start": &tengo.GoFunction{Name: "start", Value: FuncARE(cmd.Start)}, //
			// wait() => error
			"wait": &tengo.GoFunction{Name: "wait", Value: FuncARE(cmd.Wait)}, //
			// set_path(path string)
			"set_path": &tengo.GoFunction{
				Name: "set_path",
				Value: func(_ tengo.Interop, args ...tengo.Object) (ret tengo.Object, err error) {
					if len(args) != 1 {
						return nil, tengo.ErrWrongNumArguments
					}

					s1, ok := tengo.ToString(args[0])
					if !ok {
						return nil, tengo.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "string(compatible)",
							Found:    args[0].TypeName(),
						}
					}

					cmd.Path = s1

					return tengo.UndefinedValue, nil
				},
			},
			// set_dir(dir string)
			"set_dir": &tengo.GoFunction{
				Name: "set_dir",
				Value: func(_ tengo.Interop, args ...tengo.Object) (ret tengo.Object, err error) {
					if len(args) != 1 {
						return nil, tengo.ErrWrongNumArguments
					}

					s1, ok := tengo.ToString(args[0])
					if !ok {
						return nil, tengo.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "string(compatible)",
							Found:    args[0].TypeName(),
						}
					}

					cmd.Dir = s1

					return tengo.UndefinedValue, nil
				},
			},
			// set_env(env array(string))
			"set_env": &tengo.GoFunction{
				Name: "set_env",
				Value: func(_ tengo.Interop, args ...tengo.Object) (tengo.Object, error) {
					if len(args) != 1 {
						return nil, tengo.ErrWrongNumArguments
					}

					var env []string
					var err error
					switch arg0 := args[0].(type) {
					case *tengo.Array:
						env, err = stringArray(arg0.Value, "first")
						if err != nil {
							return nil, err
						}
					case *tengo.ImmutableArray:
						env, err = stringArray(arg0.Value, "first")
						if err != nil {
							return nil, err
						}
					default:
						return nil, tengo.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "array",
							Found:    arg0.TypeName(),
						}
					}

					cmd.Env = env

					return tengo.UndefinedValue, nil
				},
			},
			// process() => imap(process)
			"process": &tengo.GoFunction{
				Name: "process",
				Value: func(_ tengo.Interop, args ...tengo.Object) (ret tengo.Object, err error) {
					if len(args) != 0 {
						return nil, tengo.ErrWrongNumArguments
					}

					return makeOSProcess(cmd.Process), nil
				},
			},
		},
	}
}

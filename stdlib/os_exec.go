package stdlib

import (
	"os/exec"

	"github.com/d5/tengo/v2"
)

func makeOSExecCommand(cmd *exec.Cmd) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			// combined_output() => bytes/error
			"combined_output": &tengo.UserFunction{
				Name:  "combined_output",
				Value: FuncARYE(cmd.CombinedOutput),
			},
			// output() => bytes/error
			"output": &tengo.UserFunction{
				Name:  "output",
				Value: FuncARYE(cmd.Output),
			}, //
			// run() => error
			"run": &tengo.UserFunction{
				Name:  "run",
				Value: FuncARE(cmd.Run),
			}, //
			// start() => error
			"start": &tengo.UserFunction{
				Name:  "start",
				Value: FuncARE(cmd.Start),
			}, //
			// wait() => error
			"wait": &tengo.UserFunction{
				Name:  "wait",
				Value: FuncARE(cmd.Wait),
			}, //
			// set_path(path string)
			"set_path": &tengo.UserFunction{
				Name: "set_path",
				Value: tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
					s1, err := tengo.ToString(0, args...)
					if err != nil {
						return nil, err
					}
					cmd.Path = s1
					return tengo.UndefinedValue, nil
				}, 1),
			},
			// set_dir(dir string)
			"set_dir": &tengo.UserFunction{
				Name: "set_dir",
				Value: tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
					s1, err := tengo.ToString(0, args...)
					if err != nil {
						return nil, err
					}
					cmd.Dir = s1
					return tengo.UndefinedValue, nil
				}, 1),
			},
			// set_env(env array(string))
			"set_env": &tengo.UserFunction{
				Name: "set_env",
				Value: tengo.CheckArgs(func(args ...tengo.Object) (tengo.Object, error) {
					var env []string
					var err error
					switch arg0 := args[0].(type) {
					case *tengo.Array:
						env, err = stringArray(arg0.Value, 0)
						if err != nil {
							return nil, err
						}
					case *tengo.ImmutableArray:
						env, err = stringArray(arg0.Value, 0)
						if err != nil {
							return nil, err
						}
					default:
						panic("impossible")
					}
					cmd.Env = env
					return tengo.UndefinedValue, nil
				}, 1, 1, tengo.TNs{tengo.ArrayTN, tengo.ImmutableArrayTN}),
			},
			// process() => imap(process)
			"process": &tengo.UserFunction{
				Name: "process",
				Value: tengo.CheckStrictArgs(func(args ...tengo.Object) (tengo.Object, error) {
					return makeOSProcess(cmd.Process), nil
				}),
			},
		},
	}
}

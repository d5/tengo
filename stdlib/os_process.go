package stdlib

import (
	"os"
	"syscall"

	"github.com/d5/tengo/v2"
)

func makeOSProcessState(state *os.ProcessState) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"exited": &tengo.UserFunction{
				Name:  "exited",
				Value: FuncARB(state.Exited),
			},
			"pid": &tengo.UserFunction{
				Name:  "pid",
				Value: FuncARI(state.Pid),
			},
			"string": &tengo.UserFunction{
				Name:  "string",
				Value: FuncARS(state.String),
			},
			"success": &tengo.UserFunction{
				Name:  "success",
				Value: FuncARB(state.Success),
			},
		},
	}
}

func makeOSProcess(proc *os.Process) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"kill": &tengo.UserFunction{
				Name:  "kill",
				Value: FuncARE(proc.Kill),
			},
			"release": &tengo.UserFunction{
				Name:  "release",
				Value: FuncARE(proc.Release),
			},
			"signal": &tengo.UserFunction{
				Name: "signal",
				Value: tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
					i1, err := tengo.ToInt64(0, args...)
					if err != nil {
						return nil, err
					}
					return wrapError(proc.Signal(syscall.Signal(i1))), nil
				}, 1),
			},
			"wait": &tengo.UserFunction{
				Name: "wait",
				Value: tengo.CheckStrictArgs(func(args ...tengo.Object) (tengo.Object, error) {
					state, err := proc.Wait()
					if err != nil {
						return wrapError(err), nil
					}
					return makeOSProcessState(state), nil
				}),
			},
		},
	}
}

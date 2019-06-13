package stdlib

import (
	"os"
	"syscall"

	"github.com/d5/tengo"
)

func makeOSProcessState(state *os.ProcessState) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"exited":  &tengo.GoFunction{Name: "exited", Value: FuncARB(state.Exited)},   //
			"pid":     &tengo.GoFunction{Name: "pid", Value: FuncARI(state.Pid)},         //
			"string":  &tengo.GoFunction{Name: "string", Value: FuncARS(state.String)},   //
			"success": &tengo.GoFunction{Name: "success", Value: FuncARB(state.Success)}, //
		},
	}
}

func makeOSProcess(proc *os.Process) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"kill":    &tengo.GoFunction{Name: "kill", Value: FuncARE(proc.Kill)},       //
			"release": &tengo.GoFunction{Name: "release", Value: FuncARE(proc.Release)}, //
			"signal": &tengo.GoFunction{
				Name: "signal",
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

					return wrapError(proc.Signal(syscall.Signal(i1))), nil
				},
			},
			"wait": &tengo.GoFunction{
				Name: "wait",
				Value: func(_ tengo.Interop, args ...tengo.Object) (ret tengo.Object, err error) {
					if len(args) != 0 {
						return nil, tengo.ErrWrongNumArguments
					}

					state, err := proc.Wait()
					if err != nil {
						return wrapError(err), nil
					}

					return makeOSProcessState(state), nil
				},
			},
		},
	}
}

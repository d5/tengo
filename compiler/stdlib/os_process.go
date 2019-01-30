package stdlib

import (
	"os"
	"syscall"

	"github.com/d5/tengo/objects"
)

func makeOSProcessState(state *os.ProcessState) *objects.ImmutableMap {
	return &objects.ImmutableMap{
		Value: map[string]objects.Object{
			"exited":  FuncARB(state.Exited),
			"pid":     FuncARI(state.Pid),
			"string":  FuncARS(state.String),
			"success": FuncARB(state.Success),
		},
	}
}

func makeOSProcess(proc *os.Process) *objects.ImmutableMap {
	return &objects.ImmutableMap{
		Value: map[string]objects.Object{
			"kill":    FuncARE(proc.Kill),
			"release": FuncARE(proc.Release),
			"signal": &objects.UserFunction{
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					if len(args) != 1 {
						return nil, objects.ErrWrongNumArguments
					}

					i1, ok := objects.ToInt64(args[0])
					if !ok {
						return nil, objects.ErrInvalidTypeConversion
					}

					return wrapError(proc.Signal(syscall.Signal(i1))), nil
				},
			},
			"wait": &objects.UserFunction{
				Value: func(args ...objects.Object) (ret objects.Object, err error) {
					if len(args) != 0 {
						return nil, objects.ErrWrongNumArguments
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

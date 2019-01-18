package stdlib

import (
	"os"
	"syscall"

	"github.com/d5/tengo/objects"
)

func osProcessStateImmutableMap(state *os.ProcessState) *objects.ImmutableMap {
	return &objects.ImmutableMap{
		Value: map[string]objects.Object{
			"exited":  FuncARB(state.Exited),
			"pid":     FuncARI(state.Pid),
			"string":  FuncARS(state.String),
			"success": FuncARB(state.Success),
		},
	}
}

func osProcessImmutableMap(proc *os.Process) *objects.ImmutableMap {
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

					return osProcessStateImmutableMap(state), nil
				},
			},
		},
	}
}

func osFindProcess(args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}

	i1, ok := objects.ToInt(args[0])
	if !ok {
		return nil, objects.ErrInvalidTypeConversion
	}

	proc, err := os.FindProcess(i1)
	if err != nil {
		return wrapError(err), nil
	}

	return osProcessImmutableMap(proc), nil
}

func osStartProcess(args ...objects.Object) (objects.Object, error) {
	if len(args) != 4 {
		return nil, objects.ErrWrongNumArguments
	}

	name, ok := objects.ToString(args[0])
	if !ok {
		return nil, objects.ErrInvalidTypeConversion
	}

	argv, err := stringArray(args[1])
	if err != nil {
		return nil, err
	}

	dir, ok := objects.ToString(args[2])
	if !ok {
		return nil, objects.ErrInvalidTypeConversion
	}

	env, err := stringArray(args[3])
	if err != nil {
		return nil, err
	}

	proc, err := os.StartProcess(name, argv, &os.ProcAttr{
		Dir: dir,
		Env: env,
	})
	if err != nil {
		return wrapError(err), nil
	}

	return osProcessImmutableMap(proc), nil
}

func stringArray(o objects.Object) ([]string, error) {
	arr, ok := o.(*objects.Array)
	if !ok {
		return nil, objects.ErrInvalidTypeConversion
	}

	var sarr []string
	for _, elem := range arr.Value {
		str, ok := elem.(*objects.String)
		if !ok {
			return nil, objects.ErrInvalidTypeConversion
		}

		sarr = append(sarr, str.Value)
	}

	return sarr, nil
}

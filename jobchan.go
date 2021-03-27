package tengo

import (
	"time"
)

func init() {
	addBuiltinFunction("go", builtinGo)
	addBuiltinFunction("makechan", builtinMakechan)
}

type result struct {
	retVal Object
	err    error
}

type job struct {
	result
	vm       *VM
	waitChan chan result
	done     bool
}

// Start a goroutine which run fn(arg1, arg2, ...) in a new VM cloned from the current running VM.
// Return a job object that has wait, result, abort methods.
func builtinGo(args ...Object) (Object, error) {
	vm := args[0].(*VMObj).Value
	args = args[1:] // the first arg is VMObj inserted by VM
	if len(args) == 0 {
		return nil, ErrWrongNumArguments
	}
	fn, ok := args[0].(*CompiledFunction)
	if !ok {
		return nil, ErrInvalidArgumentType{
			Name:     "first",
			Expected: "func",
			Found:    args[0].TypeName(),
		}
	}

	newVM := vm.ShallowClone()
	jb := &job{
		vm:       newVM,
		waitChan: make(chan result),
	}

	go func() {
		retVal, err := jb.vm.RunCompiled(fn, args[1:]...)
		jb.waitChan <- result{retVal, err}
	}()

	obj := map[string]Object{
		"result": &BuiltinFunction{Value: jb.getResult},
		"wait":   &BuiltinFunction{Value: jb.waitTimeout},
		"abort":  &BuiltinFunction{Value: jb.abort},
	}
	return &Map{Value: obj}, nil
}

// Return true if job is done
func (jb *job) wait(seconds int64) bool {
	if jb.done {
		return true
	}

	if seconds <= 0 {
		seconds = 3153600000 // 100 years
	}

	select {
	case jb.result = <-jb.waitChan:
		jb.done = true
	case <-time.After(time.Duration(seconds) * time.Second):
		return false
	}

	return true
}

// Wait the job to complete.
// Wait can have optional timeout in seconds if the first arg is int.
// Wait forever if the optional timeout not specified, or timeout <= 0
// Return true if the job exited(successfully or not) within the timeout peroid.
func (jb *job) waitTimeout(args ...Object) (Object, error) {
	args = args[1:] // the first arg is VMObj inserted by VM
	if len(args) > 1 {
		return nil, ErrWrongNumArguments
	}
	timeOut := -1
	if len(args) == 1 {
		t, ok := ToInt(args[0])
		if !ok {
			return nil, ErrInvalidArgumentType{
				Name:     "first",
				Expected: "int(compatible)",
				Found:    args[0].TypeName(),
			}
		}
		timeOut = t
	}

	if jb.wait(int64(timeOut)) {
		return TrueValue, nil
	}
	return FalseValue, nil
}

func (jb *job) abort(args ...Object) (Object, error) {
	args = args[1:] // the first arg is VMObj inserted by VM
	if len(args) != 0 {
		return nil, ErrWrongNumArguments
	}

	jb.vm.Abort()
	return nil, nil
}

// Wait job to complete, return Error value when jb.err is present,
// otherwise return the result value of fn(arg1, arg2, ...)
func (jb *job) getResult(args ...Object) (Object, error) {
	args = args[1:] // the first arg is VMObj inserted by VM
	if len(args) != 0 {
		return nil, ErrWrongNumArguments
	}

	jb.wait(-1)
	if jb.err != nil {
		return &Error{Value: &String{Value: jb.err.Error()}}, nil
	}

	return jb.retVal, nil
}

type objchan chan Object

// Make a channel to send/receive object
// Return a chan object that has send, recv, close methods.
func builtinMakechan(args ...Object) (Object, error) {
	args = args[1:] // the first arg is VMObj inserted by VM
	var size int
	switch len(args) {
	case 0:
	case 1:
		n, ok := ToInt(args[0])
		if !ok {
			return nil, ErrInvalidArgumentType{
				Name:     "first",
				Expected: "int(compatible)",
				Found:    args[0].TypeName(),
			}
		}
		size = n
	default:
		return nil, ErrWrongNumArguments
	}

	oc := make(objchan, size)
	obj := map[string]Object{
		"send":  &BuiltinFunction{Value: oc.send},
		"recv":  &BuiltinFunction{Value: oc.recv},
		"close": &BuiltinFunction{Value: oc.close},
	}
	return &Map{Value: obj}, nil
}

// Send an obj to the channel, will block until channel is not full or (*VM).Abort() has been called.
// Send to a closed channel causes panic.
func (oc objchan) send(args ...Object) (Object, error) {
	vm := args[0].(*VMObj).Value
	args = args[1:] // the first arg is VMObj inserted by VM
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}
	select {
	case <-vm.ctx.Done():
		return nil, vm.ctx.Err()
		//return &String{Value: vm.ctx.Err().Error()}, nil
	case oc <- args[0]:
	}
	return nil, nil
}

// Receive an obj from the channel, will block until channel is not empty or (*VM).Abort() has been called.
// Receive from a closed channel returns undefined value.
func (oc objchan) recv(args ...Object) (Object, error) {
	vm := args[0].(*VMObj).Value
	args = args[1:] // the first arg is VMObj inserted by VM
	if len(args) != 0 {
		return nil, ErrWrongNumArguments
	}
	select {
	case <-vm.ctx.Done():
		return nil, vm.ctx.Err()
		//return &String{Value: vm.ctx.Err().Error()}, nil
	case obj, ok := <-oc:
		if ok {
			return obj, nil
		}
	}
	return nil, nil
}

// Close the channel.
func (oc objchan) close(args ...Object) (Object, error) {
	args = args[1:] // the first arg is VMObj inserted by VM
	if len(args) != 0 {
		return nil, ErrWrongNumArguments
	}
	close(oc)
	return nil, nil
}

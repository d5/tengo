package tengo

import (
	"sync/atomic"
	"time"
)

func init() {
	addBuiltinFunction("go", builtinGo)
	addBuiltinFunction("makechan", builtinMakechan)
}

type ret struct {
	val Object
	err error
}

type goroutineVM struct {
	*VM
	ret      // return value of (*VM).RunCompiled()
	waitChan chan ret
	done     int64
}

// Start a goroutine which run fn(arg1, arg2, ...) in a new VM cloned from the current running VM.
// Return a goroutineVM object that has wait, result, abort methods.
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
	gvm := &goroutineVM{
		VM:       newVM,
		waitChan: make(chan ret),
	}

	go func() {
		val, err := gvm.RunCompiled(fn, args[1:]...)
		gvm.waitChan <- ret{val, err}
	}()

	obj := map[string]Object{
		"result": &BuiltinFunction{Value: gvm.getRet},
		"wait":   &BuiltinFunction{Value: gvm.waitTimeout},
		"abort":  &BuiltinFunction{Value: gvm.abort},
	}
	return &Map{Value: obj}, nil
}

// Return true if the goroutineVM is done
func (gvm *goroutineVM) wait(seconds int64) bool {
	if atomic.LoadInt64(&gvm.done) == 1 {
		return true
	}

	if seconds <= 0 {
		seconds = 3153600000 // 100 years
	}

	select {
	case gvm.ret = <-gvm.waitChan:
		atomic.StoreInt64(&gvm.done, 1)
	case <-time.After(time.Duration(seconds) * time.Second):
		return false
	}

	return true
}

// Wait for the goroutineVM to complete.
// Wait can have optional timeout in seconds if the first arg is int.
// Wait forever if the optional timeout not specified, or timeout <= 0
// Return true if the goroutineVM exited(successfully or not) within the timeout peroid.
func (gvm *goroutineVM) waitTimeout(args ...Object) (Object, error) {
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

	if gvm.wait(int64(timeOut)) {
		return TrueValue, nil
	}
	return FalseValue, nil
}

// Terminate the execution of the goroutineVM.
func (gvm *goroutineVM) abort(args ...Object) (Object, error) {
	args = args[1:] // the first arg is VMObj inserted by VM
	if len(args) != 0 {
		return nil, ErrWrongNumArguments
	}

	gvm.Abort()
	return nil, nil
}

// Wait the goroutineVM to complete, return Error object if any runtime error occurred
// during the execution, otherwise return the result value of fn(arg1, arg2, ...)
func (gvm *goroutineVM) getRet(args ...Object) (Object, error) {
	args = args[1:] // the first arg is VMObj inserted by VM
	if len(args) != 0 {
		return nil, ErrWrongNumArguments
	}

	gvm.wait(-1)
	if gvm.ret.err != nil {
		return &Error{Value: &String{Value: gvm.ret.err.Error()}}, nil
	}

	return gvm.ret.val, nil
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

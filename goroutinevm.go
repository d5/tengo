package tengo

import (
	"sync/atomic"
	"time"
)

func init() {
	addBuiltinFunction("govm", builtinGovm, true)
	addBuiltinFunction("abort", builtinAbort, true)
	addBuiltinFunction("makechan", builtinMakechan, false)
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

// Starts a goroutine which runs fn(arg1, arg2, ...) in a new VM cloned from the current running VM.
// Returns a goroutineVM object that has wait, result, abort methods.
//
// The goroutineVM will not exit unless:
//  1. All its descendant VMs exit
//  2. It calls abort()
//  3. Its goroutineVM object abort() is called on behalf of its parent VM
// The latter 2 cases will trigger aborting procedure of all the descendant VMs, which will
// further result in #1 above.
func builtinGovm(args ...Object) (Object, error) {
	vm := args[0].(*vmObj).Value
	args = args[1:] // the first arg is vmObj inserted by VM
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
		waitChan: make(chan ret, 1),
	}

	vm.addChildVM(gvm.VM)
	go func() {
		val, err := gvm.RunCompiled(fn, args[1:]...)
		gvm.waitChan <- ret{val, err}
		vm.delChildVM(gvm.VM)
	}()

	obj := map[string]Object{
		"result": &BuiltinFunction{Value: gvm.getRet},
		"wait":   &BuiltinFunction{Value: gvm.waitTimeout},
		"abort":  &BuiltinFunction{Value: gvm.abort},
	}
	return &Map{Value: obj}, nil
}

// Terminates the current VM and all its descendant VMs.
// Calling abort() will always result the current VM returns ErrVMAborted.
func builtinAbort(args ...Object) (Object, error) {
	vm := args[0].(*vmObj).Value
	args = args[1:] // the first arg is vmObj inserted by VM
	if len(args) != 0 {
		return nil, ErrWrongNumArguments
	}
	vm.Abort() // aborts self and all descendant VMs
	return nil, nil
}

// Returns true if the goroutineVM is done
func (gvm *goroutineVM) wait(seconds int64) bool {
	if atomic.LoadInt64(&gvm.done) == 1 {
		return true
	}

	if seconds < 0 {
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

// Waits for the goroutineVM to complete in timeout seconds.
// Returns true if the goroutineVM exited(successfully or not) within the timeout peroid.
// Waits forever if the optional timeout not specified, or timeout < 0.
func (gvm *goroutineVM) waitTimeout(args ...Object) (Object, error) {
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

// Terminates the execution of the current VM and all its descendant VMs.
func (gvm *goroutineVM) abort(args ...Object) (Object, error) {
	if len(args) != 0 {
		return nil, ErrWrongNumArguments
	}

	gvm.Abort()
	return nil, nil
}

// Waits the goroutineVM to complete, return Error object if any runtime error occurred
// during the execution, otherwise return the result value of fn(arg1, arg2, ...)
func (gvm *goroutineVM) getRet(args ...Object) (Object, error) {
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

// Makes a channel to send/receive object
// Returns a chan object that has send, recv, close methods.
func builtinMakechan(args ...Object) (Object, error) {
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
		"send":  &BuiltinFunction{Value: oc.send, needvmObj: true},
		"recv":  &BuiltinFunction{Value: oc.recv, needvmObj: true},
		"close": &BuiltinFunction{Value: oc.close},
	}
	return &Map{Value: obj}, nil
}

// Sends an obj to the channel, will block if channel is full and the VM has not been aborted.
// Sends to a closed channel causes panic.
func (oc objchan) send(args ...Object) (Object, error) {
	vm := args[0].(*vmObj).Value
	args = args[1:] // the first arg is vmObj inserted by VM
	if len(args) != 1 {
		return nil, ErrWrongNumArguments
	}
	select {
	case <-vm.abortChan:
		return nil, ErrVMAborted
	case oc <- args[0]:
	}
	return nil, nil
}

// Receives an obj from the channel, will block if channel is empty and the VM has not been aborted.
// Receives from a closed channel returns undefined value.
func (oc objchan) recv(args ...Object) (Object, error) {
	vm := args[0].(*vmObj).Value
	args = args[1:] // the first arg is vmObj inserted by VM
	if len(args) != 0 {
		return nil, ErrWrongNumArguments
	}
	select {
	case <-vm.abortChan:
		return nil, ErrVMAborted
	case obj, ok := <-oc:
		if ok {
			return obj, nil
		}
	}
	return nil, nil
}

// Closes the channel.
func (oc objchan) close(args ...Object) (Object, error) {
	if len(args) != 0 {
		return nil, ErrWrongNumArguments
	}
	close(oc)
	return nil, nil
}

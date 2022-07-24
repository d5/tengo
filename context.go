package tengo

import (
	"context"
	"fmt"
)

// VmContext represents the VM Context
type VmContext struct {
	context.Context
	VM *VM
}

// CallContext represents the Go function Call Context
type CallContext struct {
	// VM current VM
	VM *VM
	// Args called args
	Args []Object
	// Kwargs called keyword arguments
	Kwargs map[string]Object
}

// GetArgs destructure args into dest
func (this *CallContext) GetArgs(dest ...*Object) error {
	if len(dest) != len(this.Args) {
		return fmt.Errorf(
			"wrong number of arguments: want=%d, got=%d",
			len(dest), len(this.Args))
	}
	for i, v := range dest {
		*v = this.Args[i]
	}
	return nil
}

// GetArgsVar destructure args into dest and other args into argVar
func (this *CallContext) GetArgsVar(argVar *[]Object, dest ...*Object) error {
	if len(dest) != len(this.Args) {
		return fmt.Errorf(
			"wrong number of arguments: want=>%d, got=%d",
			len(dest), len(this.Args))
	}
	for i, v := range this.Args {
		*dest[i] = v
	}
	*argVar = this.Args[len(dest):]
	return nil
}

// GetKwargs destructure kwargs into dest
func (this *CallContext) GetKwargs(dest map[string]*Object) error {
	for k, v := range this.Kwargs {
		if _, ok := dest[k]; ok {
			*dest[k] = v
		} else {
			return fmt.Errorf("unexpected kwarg: %q", k)
		}
	}
	return nil
}

// GetKwargsVar destructure kwargs into dest and other kwargs into argVar
func (this *CallContext) GetKwargsVar(argVar map[string]Object, dest map[string]*Object) {
	for k, v := range this.Kwargs {
		if _, ok := dest[k]; ok {
			if v != UndefinedValue {
				*dest[k] = v
			}
		} else {
			argVar[k] = v
		}
	}
}

package objects

import (
	"github.com/d5/tengo/compiler/token"
)

// CompiledModuleGlobal represents a global variable in the compiled module.
type CompiledModuleGlobal struct {
	Index int
	Value *Object
}

// CompiledModule represents a compiled function.
type CompiledModule struct {
	Instructions []byte
	Constants    []Object
	Globals      map[string]CompiledModuleGlobal
}

// TypeName returns the name of the type.
func (o *CompiledModule) TypeName() string {
	return "compiled-function"
}

func (o *CompiledModule) String() string {
	return "<compiled-function>"
}

// BinaryOp returns another object that is the result of
// a given binary operator and a right-hand side object.
func (o *CompiledModule) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *CompiledModule) Copy() Object {
	var constants []Object
	for _, c := range o.Constants {
		constants = append(constants, c.Copy())
	}
	globals := make(map[string]CompiledModuleGlobal, len(o.Globals))
	for name, obj := range o.Globals {
		globals[name] = obj
	}

	return &CompiledModule{
		Instructions: append([]byte{}, o.Instructions...),
		Constants:    constants,
		Globals:      globals,
	}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *CompiledModule) IsFalsy() bool {
	return false
}

// Equals returns true if the value of the type
// is equal to the value of another object.
func (o *CompiledModule) Equals(x Object) bool {
	return false
}

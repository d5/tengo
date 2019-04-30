package objects

import (
	"github.com/d5/tengo/compiler/source"
	"github.com/d5/tengo/compiler/token"
)

// CompiledFunction represents a compiled function.
type CompiledFunction struct {
	ObjectImpl
	Instructions  []byte
	NumLocals     int // number of local variables (including function parameters)
	NumParameters int
	SourceMap     map[int]source.Pos
	Fn            CallableFunc // set by the runtime
}

// TypeName returns the name of the type.
func (o *CompiledFunction) TypeName() string {
	return "compiled-function"
}

func (o *CompiledFunction) String() string {
	return "<compiled-function>"
}

// BinaryOp returns another object that is the result of
// a given binary operator and a right-hand side object.
func (o *CompiledFunction) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *CompiledFunction) Copy() Object {
	return &CompiledFunction{
		Instructions:  append([]byte{}, o.Instructions...),
		NumLocals:     o.NumLocals,
		NumParameters: o.NumParameters,
	}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *CompiledFunction) IsFalsy() bool {
	return false
}

// Equals returns true if the value of the type
// is equal to the value of another object.
func (o *CompiledFunction) Equals(x Object) bool {
	return false
}

// SourcePos returns the source position of the instruction at ip.
func (o *CompiledFunction) SourcePos(ip int) source.Pos {
	for ip >= 0 {
		if p, ok := o.SourceMap[ip]; ok {
			return p
		}
		ip--
	}
	return source.NoPos
}

// Call takes an arbitrary number of arguments
// and returns a return value and/or an error.
func (o *CompiledFunction) Call(args ...Object) (Object, error) {
	return o.Fn(args...)
}

// CanCall returns whether the Object can be Called.
func (o *CompiledFunction) CanCall() bool {
	return true
}

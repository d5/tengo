package compiler

import (
	"encoding/gob"
	"fmt"
	"io"
	"reflect"

	"github.com/d5/tengo"
	"github.com/d5/tengo/compiler/source"
)

// Bytecode is a compiled instructions and constants.
type Bytecode struct {
	FileSet      *source.FileSet
	MainFunction *tengo.CompiledFunction
	Constants    []tengo.Object
}

// Encode writes Bytecode data to the writer.
func (b *Bytecode) Encode(w io.Writer) error {
	enc := gob.NewEncoder(w)

	if err := enc.Encode(b.FileSet); err != nil {
		return err
	}

	if err := enc.Encode(b.MainFunction); err != nil {
		return err
	}

	// constants
	return enc.Encode(b.Constants)
}

// CountObjects returns the number of objects found in Constants.
func (b *Bytecode) CountObjects() int {
	n := 0

	for _, c := range b.Constants {
		n += tengo.CountObjects(c)
	}

	return n
}

// FormatInstructions returns human readable string representations of
// compiled instructions.
func (b *Bytecode) FormatInstructions() []string {
	return FormatInstructions(b.MainFunction.Instructions, 0)
}

// FormatConstants returns human readable string representations of
// compiled constants.
func (b *Bytecode) FormatConstants() (output []string) {
	for cidx, cn := range b.Constants {
		switch cn := cn.(type) {
		case *tengo.CompiledFunction:
			output = append(output, fmt.Sprintf("[% 3d] (Compiled Function|%p)", cidx, &cn))
			for _, l := range FormatInstructions(cn.Instructions, 0) {
				output = append(output, fmt.Sprintf("     %s", l))
			}
		default:
			output = append(output, fmt.Sprintf("[% 3d] %s (%s|%p)", cidx, cn, reflect.TypeOf(cn).Elem().Name(), &cn))
		}
	}

	return
}

func init() {
	gob.Register(&source.FileSet{})
	gob.Register(&source.File{})
	gob.Register(&tengo.Array{})
	gob.Register(&tengo.Bool{})
	gob.Register(&tengo.Bytes{})
	gob.Register(&tengo.Char{})
	gob.Register(&tengo.CompiledFunction{})
	gob.Register(&tengo.Error{})
	gob.Register(&tengo.Float{})
	gob.Register(&tengo.ImmutableArray{})
	gob.Register(&tengo.ImmutableMap{})
	gob.Register(&tengo.Int{})
	gob.Register(&tengo.Map{})
	gob.Register(&tengo.String{})
	gob.Register(&tengo.Time{})
	gob.Register(&tengo.Undefined{})
	gob.Register(&tengo.UserFunction{})
}

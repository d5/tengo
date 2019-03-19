package compiler

import (
	"encoding/gob"
	"fmt"
	"io"
	"reflect"

	"github.com/d5/tengo/compiler/source"
	"github.com/d5/tengo/objects"
)

// Bytecode is a compiled instructions and constants.
type Bytecode struct {
	FileSet      *source.FileSet
	MainFunction *objects.CompiledFunction
	Constants    []objects.Object
}

// Decode reads Bytecode data from the reader.
func (b *Bytecode) Decode(r io.Reader) error {
	dec := gob.NewDecoder(r)

	if err := dec.Decode(&b.FileSet); err != nil {
		return err
	}
	// TODO: files in b.FileSet.File does not have their 'set' field properly set to b.FileSet
	// as it's private field and not serialized by gob encoder/decoder.

	if err := dec.Decode(&b.MainFunction); err != nil {
		return err
	}

	if err := dec.Decode(&b.Constants); err != nil {
		return err
	}

	return nil
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
		n += objects.CountObjects(c)
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
		case *objects.CompiledFunction:
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
	gob.Register(&objects.Char{})
	gob.Register(&objects.CompiledFunction{})
	gob.Register(&objects.ImmutableMap{})
	gob.Register(&objects.Float{})
	gob.Register(&objects.Int{})
	gob.Register(&objects.String{})
}

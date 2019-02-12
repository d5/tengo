package compiler

import (
	"encoding/gob"
	"fmt"
	"io"
	"reflect"

	"github.com/d5/tengo/objects"
)

// Bytecode is a compiled instructions and constants.
type Bytecode struct {
	Instructions []byte
	Constants    []objects.Object
}

// Decode reads Bytecode data from the reader.
func (b *Bytecode) Decode(r io.Reader) error {
	dec := gob.NewDecoder(r)

	if err := dec.Decode(&b.Instructions); err != nil {
		return err
	}

	if err := dec.Decode(&b.Constants); err != nil {
		return err
	}

	// replace Bool and Undefined with known value
	for i, v := range b.Constants {
		b.Constants[i] = cleanupObjects(v)
	}

	return nil
}

// Encode writes Bytecode data to the writer.
func (b *Bytecode) Encode(w io.Writer) error {
	enc := gob.NewEncoder(w)

	if err := enc.Encode(b.Instructions); err != nil {
		return err
	}

	// constants
	return enc.Encode(b.Constants)
}

// FormatInstructions returns human readable string representations of
// compiled instructions.
func (b *Bytecode) FormatInstructions() []string {
	return FormatInstructions(b.Instructions, 0)
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

func cleanupObjects(o objects.Object) objects.Object {
	switch o := o.(type) {
	case *objects.Bool:
		if o.IsFalsy() {
			return objects.FalseValue
		}
		return objects.TrueValue
	case *objects.Undefined:
		return objects.UndefinedValue
	case *objects.Array:
		for i, v := range o.Value {
			o.Value[i] = cleanupObjects(v)
		}
	case *objects.Map:
		for k, v := range o.Value {
			o.Value[k] = cleanupObjects(v)
		}
	}

	return o
}

func init() {
	gob.Register(&objects.Int{})
	gob.Register(&objects.Float{})
	gob.Register(&objects.String{})
	gob.Register(&objects.Bool{})
	gob.Register(&objects.Char{})
	gob.Register(&objects.Array{})
	gob.Register(&objects.ImmutableArray{})
	gob.Register(&objects.Map{})
	gob.Register(&objects.ImmutableMap{})
	gob.Register(&objects.CompiledFunction{})
	gob.Register(&objects.Undefined{})
	gob.Register(&objects.Error{})
	gob.Register(&objects.Bytes{})
	gob.Register(&objects.StringIterator{})
	gob.Register(&objects.MapIterator{})
	gob.Register(&objects.ArrayIterator{})
	gob.Register(&objects.Time{})
}

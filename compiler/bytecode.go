package compiler

import (
	"encoding/gob"
	"io"

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

	return nil
}

// Encode writes Bytecode data to the writer.
func (b *Bytecode) Encode(w io.Writer) error {
	enc := gob.NewEncoder(w)

	if err := enc.Encode(b.Instructions); err != nil {
		return err
	}

	if err := enc.Encode(b.Constants); err != nil {
		return err
	}

	return nil
}

func init() {
	gob.Register(&objects.Int{})
	gob.Register(&objects.Float{})
	gob.Register(&objects.String{})
	gob.Register(&objects.Bool{})
	gob.Register(&objects.Char{})
	gob.Register(&objects.Array{})
	gob.Register(&objects.Map{})
	gob.Register(&objects.CompiledFunction{})
	gob.Register(&objects.Undefined{})
}

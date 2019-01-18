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

	return dec.Decode(&b.Constants)
}

// Encode writes Bytecode data to the writer.
func (b *Bytecode) Encode(w io.Writer) error {
	enc := gob.NewEncoder(w)

	if err := enc.Encode(b.Instructions); err != nil {
		return err
	}

	return enc.Encode(b.Constants)
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
	gob.Register(&objects.Error{})
	gob.Register(&objects.ImmutableMap{})
	gob.Register(&objects.Bytes{})
	gob.Register(&objects.StringIterator{})
	gob.Register(&objects.MapIterator{})
	gob.Register(&objects.ImmutableMapIterator{})
	gob.Register(&objects.ArrayIterator{})
}

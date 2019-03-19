package compiler

import (
	"encoding/gob"
	"io"

	"github.com/d5/tengo/objects"
)

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
	for i, v := range b.Constants {
		b.Constants[i] = fixDecoded(v)
	}

	return nil
}

func fixDecoded(o objects.Object) objects.Object {
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
			o.Value[i] = fixDecoded(v)
		}
	case *objects.ImmutableArray:
		for i, v := range o.Value {
			o.Value[i] = fixDecoded(v)
		}
	case *objects.Map:
		for k, v := range o.Value {
			o.Value[k] = fixDecoded(v)
		}
	case *objects.ImmutableMap:
		for k, v := range o.Value {
			o.Value[k] = fixDecoded(v)
		}
	}

	return o
}

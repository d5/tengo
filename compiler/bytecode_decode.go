package compiler

import (
	"encoding/gob"
	"fmt"
	"io"

	"github.com/d5/tengo/objects"
)

// Decode reads Bytecode data from the reader.
func (b *Bytecode) Decode(r io.Reader, builtinModules map[string]*objects.BuiltinModule) error {
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
		fv, err := fixDecoded(v, builtinModules)
		if err != nil {
			return err
		}
		b.Constants[i] = fv
	}

	return nil
}

func fixDecoded(o objects.Object, builtinModules map[string]*objects.BuiltinModule) (objects.Object, error) {
	switch o := o.(type) {
	case *objects.Bool:
		if o.IsFalsy() {
			return objects.FalseValue, nil
		}
		return objects.TrueValue, nil
	case *objects.Undefined:
		return objects.UndefinedValue, nil
	case *objects.Array:
		for i, v := range o.Value {
			fv, err := fixDecoded(v, builtinModules)
			if err != nil {
				return nil, err
			}
			o.Value[i] = fv
		}
	case *objects.ImmutableArray:
		for i, v := range o.Value {
			fv, err := fixDecoded(v, builtinModules)
			if err != nil {
				return nil, err
			}
			o.Value[i] = fv
		}
	case *objects.Map:
		for k, v := range o.Value {
			fv, err := fixDecoded(v, builtinModules)
			if err != nil {
				return nil, err
			}
			o.Value[k] = fv
		}
	case *objects.ImmutableMap:
		modName := moduleName(o)
		if mod, ok := builtinModules[modName]; ok {
			return mod.AsImmutableMap(), nil
		}

		for k, v := range o.Value {
			// encoding of user function not supported
			if _, isUserFunction := v.(*objects.UserFunction); isUserFunction {
				return nil, fmt.Errorf("user function not decodable")
			}

			fv, err := fixDecoded(v, builtinModules)
			if err != nil {
				return nil, err
			}
			o.Value[k] = fv
		}
	}

	return o, nil
}

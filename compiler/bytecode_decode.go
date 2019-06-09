package compiler

import (
	"encoding/gob"
	"fmt"
	"io"

	"github.com/d5/tengo"
)

// Decode reads Bytecode data from the reader.
func (b *Bytecode) Decode(r io.Reader, modules *tengo.ModuleMap) error {
	if modules == nil {
		modules = tengo.NewModuleMap()
	}

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
		fv, err := fixDecoded(v, modules)
		if err != nil {
			return err
		}
		b.Constants[i] = fv
	}

	return nil
}

func fixDecoded(o tengo.Object, modules *tengo.ModuleMap) (tengo.Object, error) {
	switch o := o.(type) {
	case *tengo.Bool:
		if o.IsFalsy() {
			return tengo.FalseValue, nil
		}
		return tengo.TrueValue, nil
	case *tengo.Undefined:
		return tengo.UndefinedValue, nil
	case *tengo.Array:
		for i, v := range o.Value {
			fv, err := fixDecoded(v, modules)
			if err != nil {
				return nil, err
			}
			o.Value[i] = fv
		}
	case *tengo.ImmutableArray:
		for i, v := range o.Value {
			fv, err := fixDecoded(v, modules)
			if err != nil {
				return nil, err
			}
			o.Value[i] = fv
		}
	case *tengo.Map:
		for k, v := range o.Value {
			fv, err := fixDecoded(v, modules)
			if err != nil {
				return nil, err
			}
			o.Value[k] = fv
		}
	case *tengo.ImmutableMap:
		modName := moduleName(o)
		if mod := modules.GetBuiltinModule(modName); mod != nil {
			return mod.AsImmutableMap(modName), nil
		}

		for k, v := range o.Value {
			// encoding of user function not supported
			if _, isUserFunction := v.(*tengo.UserFunction); isUserFunction {
				return nil, fmt.Errorf("user function not decodable")
			}

			fv, err := fixDecoded(v, modules)
			if err != nil {
				return nil, err
			}
			o.Value[k] = fv
		}
	}

	return o, nil
}

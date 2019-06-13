package stdlib

import (
	"bytes"
	gojson "encoding/json"

	"github.com/d5/tengo"
	"github.com/d5/tengo/stdlib/json"
)

var jsonModule = map[string]tengo.Object{
	"decode":      &tengo.GoFunction{Name: "decode", Value: jsonDecode},
	"encode":      &tengo.GoFunction{Name: "encode", Value: jsonEncode},
	"indent":      &tengo.GoFunction{Name: "encode", Value: jsonIndent},
	"html_escape": &tengo.GoFunction{Name: "html_escape", Value: jsonHTMLEscape},
}

func jsonDecode(_ tengo.Interop, args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}

	switch o := args[0].(type) {
	case *tengo.Bytes:
		v, err := json.Decode(o.Value)
		if err != nil {
			return &tengo.Error{Value: &tengo.String{Value: err.Error()}}, nil
		}
		return v, nil
	case *tengo.String:
		v, err := json.Decode([]byte(o.Value))
		if err != nil {
			return &tengo.Error{Value: &tengo.String{Value: err.Error()}}, nil
		}
		return v, nil
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "bytes/string",
			Found:    args[0].TypeName(),
		}
	}
}

func jsonEncode(_ tengo.Interop, args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}

	b, err := json.Encode(args[0])
	if err != nil {
		return &tengo.Error{Value: &tengo.String{Value: err.Error()}}, nil
	}

	return &tengo.Bytes{Value: b}, nil
}

func jsonIndent(_ tengo.Interop, args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 3 {
		return nil, tengo.ErrWrongNumArguments
	}

	prefix, ok := tengo.ToString(args[1])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "prefix",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	indent, ok := tengo.ToString(args[2])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "indent",
			Expected: "string(compatible)",
			Found:    args[2].TypeName(),
		}
	}

	switch o := args[0].(type) {
	case *tengo.Bytes:
		var dst bytes.Buffer
		err := gojson.Indent(&dst, o.Value, prefix, indent)
		if err != nil {
			return &tengo.Error{Value: &tengo.String{Value: err.Error()}}, nil
		}
		return &tengo.Bytes{Value: dst.Bytes()}, nil
	case *tengo.String:
		var dst bytes.Buffer
		err := gojson.Indent(&dst, []byte(o.Value), prefix, indent)
		if err != nil {
			return &tengo.Error{Value: &tengo.String{Value: err.Error()}}, nil
		}
		return &tengo.Bytes{Value: dst.Bytes()}, nil
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "bytes/string",
			Found:    args[0].TypeName(),
		}
	}
}

func jsonHTMLEscape(_ tengo.Interop, args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}

	switch o := args[0].(type) {
	case *tengo.Bytes:
		var dst bytes.Buffer
		gojson.HTMLEscape(&dst, o.Value)
		return &tengo.Bytes{Value: dst.Bytes()}, nil
	case *tengo.String:
		var dst bytes.Buffer
		gojson.HTMLEscape(&dst, []byte(o.Value))
		return &tengo.Bytes{Value: dst.Bytes()}, nil
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "bytes/string",
			Found:    args[0].TypeName(),
		}
	}
}

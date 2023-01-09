package stdlib

import (
	"bytes"
	gojson "encoding/json"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib/json"
)

var jsonModule = map[string]tengo.Object{
	"decode": &tengo.UserFunction{
		Name:  "decode",
		Value: jsonDecode,
	},
	"encode": &tengo.UserFunction{
		Name:  "encode",
		Value: jsonEncode,
	},
	"indent": &tengo.UserFunction{
		Name:  "encode",
		Value: jsonIndent,
	},
	"html_escape": &tengo.UserFunction{
		Name:  "html_escape",
		Value: jsonHTMLEscape,
	},
}

var jsonDecode = tengo.CheckArgs(func(args ...tengo.Object) (tengo.Object, error) {
	switch o := args[0].(type) {
	case *tengo.Bytes:
		v, err := json.Decode(o.Value)
		if err != nil {
			return &tengo.Error{
				Value: &tengo.String{Value: err.Error()},
			}, nil
		}
		return v, nil
	case *tengo.String:
		v, err := json.Decode([]byte(o.Value))
		if err != nil {
			return &tengo.Error{
				Value: &tengo.String{Value: err.Error()},
			}, nil
		}
		return v, nil
	default:
		panic("impossible")
	}
}, 1, 1, tengo.TNs{tengo.BytesTN, tengo.StringTN})

var jsonEncode = tengo.CheckAnyArgs(func(args ...tengo.Object) (tengo.Object, error) {
	b, err := json.Encode(args[0])
	if err != nil {
		return &tengo.Error{Value: &tengo.String{Value: err.Error()}}, nil
	}

	return &tengo.Bytes{Value: b}, nil
}, 1)

var jsonIndent = tengo.CheckArgs(func(args ...tengo.Object) (tengo.Object, error) {
	prefix, err := tengo.ToString(1, args...)
	if err != nil {
		return nil, err
	}

	indent, err := tengo.ToString(2, args...)
	if err != nil {
		return nil, err
	}

	switch o := args[0].(type) {
	case *tengo.Bytes:
		var dst bytes.Buffer
		err := gojson.Indent(&dst, o.Value, prefix, indent)
		if err != nil {
			return &tengo.Error{
				Value: &tengo.String{Value: err.Error()},
			}, nil
		}
		return &tengo.Bytes{Value: dst.Bytes()}, nil
	case *tengo.String:
		var dst bytes.Buffer
		err := gojson.Indent(&dst, []byte(o.Value), prefix, indent)
		if err != nil {
			return &tengo.Error{
				Value: &tengo.String{Value: err.Error()},
			}, nil
		}
		return &tengo.Bytes{Value: dst.Bytes()}, nil
	default:
		panic("impossible")
	}
}, 3, 3, nil, nil, tengo.TNs{tengo.BytesTN, tengo.StringTN})

var jsonHTMLEscape = tengo.CheckArgs(func(args ...tengo.Object) (tengo.Object, error) {
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
		panic("impossible")
	}
}, 1, 1, tengo.TNs{tengo.BytesTN, tengo.StringTN})

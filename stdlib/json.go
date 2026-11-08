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

func jsonDecode(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}

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
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "bytes/string",
			Found:    args[0].TypeName(),
		}
	}
}

func jsonEncode(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}

	b, err := json.Encode(args[0])
	if err == tengo.ErrBytesLimit {
		return nil, err
	}
	if err != nil {
		return &tengo.Error{Value: &tengo.String{Value: err.Error()}}, nil
	}

	return &tengo.Bytes{Value: b}, nil
}

func jsonIndent(args ...tengo.Object) (ret tengo.Object, err error) {
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

	var src []byte
	switch o := args[0].(type) {
	case *tengo.Bytes:
		src = o.Value
	case *tengo.String:
		src = []byte(o.Value)
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "bytes/string",
			Found:    args[0].TypeName(),
		}
	}

	if indentedLen(src, prefix, indent) > int64(tengo.MaxBytesLen) {
		return nil, tengo.ErrBytesLimit
	}
	var dst bytes.Buffer
	if err := gojson.Indent(&dst, src, prefix, indent); err != nil {
		return &tengo.Error{
			Value: &tengo.String{Value: err.Error()},
		}, nil
	}
	return &tengo.Bytes{Value: dst.Bytes()}, nil
}

// indentedLen returns how long gojson.Indent makes src, counting each space
// outside a string as kept. The output grows with depth times len(indent),
// so it must be counted before Indent allocates it. It stops counting past
// MaxBytesLen.
func indentedLen(src []byte, prefix, indent string) int64 {
	var n int64
	depth := 0
	needIndent, inString, escaped := false, false, false
	newline := func() int64 {
		return 1 + int64(len(prefix)) + int64(depth)*int64(len(indent))
	}
	for _, c := range src {
		if n > int64(tengo.MaxBytesLen) {
			break
		}
		if inString {
			n++
			switch {
			case escaped:
				escaped = false
			case c == '\\':
				escaped = true
			case c == '"':
				inString = false
			}
			continue
		}
		if needIndent && c != '}' && c != ']' && !isJSONSpace(c) {
			needIndent = false
			depth++
			n += newline()
		}
		switch c {
		case '"':
			inString = true
			n++
		case '{', '[':
			needIndent = true
			n++
		case ',':
			n += 1 + newline()
		case ':':
			n += 2
		case '}', ']':
			if needIndent {
				needIndent = false
			} else {
				depth--
				if depth < 0 {
					// Indent stops at this error
					return n
				}
				n += newline()
			}
			n++
		default:
			n++
		}
	}
	return n
}

func isJSONSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func jsonHTMLEscape(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}

	var src []byte
	switch o := args[0].(type) {
	case *tengo.Bytes:
		src = o.Value
	case *tengo.String:
		src = []byte(o.Value)
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "bytes/string",
			Found:    args[0].TypeName(),
		}
	}

	// HTMLEscape writes <, > and & as \u003c and the like, 6 bytes for 1,
	// and U+2028 and U+2029 as \u2028, 6 bytes for 3.
	n := int64(len(src))
	for i, c := range src {
		switch {
		case c == '<' || c == '>' || c == '&':
			n += 5
		case c == 0xE2 && i+2 < len(src) && src[i+1] == 0x80 &&
			src[i+2]&^1 == 0xA8:
			n += 3
		}
	}
	if n > int64(tengo.MaxBytesLen) {
		return nil, tengo.ErrBytesLimit
	}
	var dst bytes.Buffer
	gojson.HTMLEscape(&dst, src)
	return &tengo.Bytes{Value: dst.Bytes()}, nil
}

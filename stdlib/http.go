package stdlib

import (
	"bytes"
	"io"
	"net/http"

	"github.com/d5/tengo/v2"
)

var httpModule = map[string]tengo.Object{
	"do": &tengo.UserFunction{
		Name: "do",
		Value: (func(args ...tengo.Object) (tengo.Object, error) {
			numArgs := len(args)
			if numArgs < 2 || numArgs > 4 {
				return nil, tengo.ErrWrongNumArguments
			}

			// build req from method, url [, headers[, body]]
			method, ok := args[0].(*tengo.String)
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     "1/method",
					Expected: "string",
					Found:    args[0].TypeName(),
				}
			}

			url, ok := args[1].(*tengo.String)
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     "2/url",
					Expected: "string",
					Found:    args[1].TypeName(),
				}
			}
			var body io.Reader
			if len(args) > 3 {
				bs, ok := args[3].(*tengo.Bytes)
				if !ok {
					return nil, tengo.ErrInvalidArgumentType{
						Name:     "4/body",
						Expected: "bytes",
						Found:    args[3].TypeName(),
					}
				}
				body = bytes.NewBuffer(bs.Value)
			}
			req, err := http.NewRequest(method.Value, url.Value, body)
			if err != nil {
				return wrapError(err), nil
			}
			// add headers
			if len(args) > 2 {
				m, ok := args[2].(*tengo.Map)
				if !ok {
					return nil, tengo.ErrInvalidArgumentType{
						Name:     "3/headers",
						Expected: "map",
						Found:    args[2].TypeName(),
					}
				}
				for k, v := range m.Value {
					s, ok := tengo.ToString(v)
					if !ok {
						return nil, tengo.ErrInvalidArgumentType{
							Name:     "headers",
							Expected: "string",
							Found:    v.TypeName(),
						}
					}
					if err != nil {
						return nil, err
					}
					req.Header.Add(k, s)
				}
			}

			// do req
			res, err := http.DefaultClient.Do(req)
			if res != nil && res.Body != nil {
				// ensure to always close body no matter what
				defer res.Body.Close()
			}
			if err != nil {
				return wrapError(err), nil
			}
			if res.ContentLength > int64(tengo.MaxBytesLen) {
				// don't allow going over byte limit
				return nil, tengo.ErrBytesLimit
			}

			// read full body, with byte limit on it
			bs, err := io.ReadAll(io.LimitReader(res.Body, int64(tengo.MaxBytesLen)))
			if err != nil {
				return wrapError(err), nil
			}
			resHeaders := &tengo.Map{Value: map[string]tengo.Object{}}
			for k := range res.Header {
				resHeaders.Value[k] = &tengo.String{Value: res.Header.Get(k)}
			}
			return &tengo.Map{
				Value: map[string]tengo.Object{
					"code":    &tengo.Int{Value: int64(res.StatusCode)},
					"status":  &tengo.String{Value: res.Status},
					"headers": resHeaders,
					"body": &tengo.Bytes{
						Value: bs,
					},
				},
			}, nil
		}),
	},
}

# Module - "base64"

```golang
base64 := import("base64")
```

## Functions

- `encode(src)`: returns the base64 encoding of src.
- `decode(s)`: returns the bytes represented by the base64 string s.
- `raw_encode(src)`: returns the base64 encoding of src but omits the padding.
- `raw_decode(s)`: returns the bytes represented by the base64 string s which
  omits the padding.
- `url_encode(src)`: returns the url-base64 encoding of src.
- `url_decode(s)`: returns the bytes represented by the url-base64 string s.
- `raw_url_encode(src)`: returns the url-base64 encoding of src but omits the
  padding.
- `raw_url_decode(s)`: returns the bytes represented by the url-base64 string
  s which omits the padding.

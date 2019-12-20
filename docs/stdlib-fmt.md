# Module - "fmt"

```golang
fmt := import("fmt")
```

## Functions

- `print(args...)`: Prints a string representation of the given variable to the
  standard output. Unlike Go's `fmt.Print` function, no spaces are added between
  the operands.
- `println(args...)`: Prints a string representation of the given variable to
  the standard output with a newline appended. Unlike Go's `fmt.Println`
  function, no spaces are added between the operands.
- `printf(format, args...)`: Prints a formatted string to the standard output.
  It does not append the newline character at the end. The first argument must
  a String object. See
  [this](https://github.com/d5/tengo/blob/master/docs/formatting.md) for more
  details on formatting.
- `sprintf(format, args...)`: Returns a formatted string. Alias of the builtin
  function `format`. The first argument must be a String object. See
  [this](https://github.com/d5/tengo/blob/master/docs/formatting.md) for more
  details on formatting.

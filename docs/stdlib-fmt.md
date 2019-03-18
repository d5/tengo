# Module - "fmt"

```golang
fmt := import("fmt")
```

## Functions

- `print(args...)`: Prints a string representation of the given variable to the standard output. Unlike Go's `fmt.Print` function, no spaces are added between the operands.
- `println(args...)`: Prints a string representation of the given variable to the standard output with a newline appended. Unlike Go's `fmt.Println` function, no spaces are added between the operands.
- `printf(format, args...)`: Prints a formatted string to the standard output. It does not append the newline character at the end. The first argument must a String object. It's same as Go's `fmt.Printf`.
- `sprintf(format, args...)`: Returns a formatted string. The first argument must be a String object. It's the same as Go's `fmt.Sprintf`.

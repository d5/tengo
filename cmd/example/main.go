package main

import "github.com/d5/tengo/script"

var code = `
reduce := func(seq, fn) {
	s := 0
    for x in seq { fn(x, s) }
	return s
}

print(reduce([1, 2, 3], func(x, s) { s += x }))
`

func main() {
	s := script.New([]byte(code))
	if _, err := s.Run(); err != nil {
		panic(err)
	}
}

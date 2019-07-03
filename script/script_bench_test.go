package script_test

import (
	"testing"

	"github.com/d5/tengo/script"
)

func BenchmarkArrayIndex(b *testing.B) {
	bench(b.N, `a := [1, 2, 3, 4, 5, 6, 7, 8, 9];
        for i := 0; i < 1000; i++ {
            a[0]; a[1]; a[2]; a[3]; a[4]; a[5]; a[6]; a[7]; a[7];
        }
    `)
}

func BenchmarkArrayIndexCompare(b *testing.B) {
	bench(b.N, `a := [1, 2, 3, 4, 5, 6, 7, 8, 9];
        for i := 0; i < 1000; i++ {
            1; 2; 3; 4; 5; 6; 7; 8; 9;
        }
    `)
}

func bench(n int, input string) {
	s := script.New([]byte(input))
	c, err := s.Compile()
	if err != nil {
		panic(err)
	}

	for i := 0; i < n; i++ {
		if err := c.Run(); err != nil {
			panic(err)
		}
	}
}

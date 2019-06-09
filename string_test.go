package tengo_test

import (
	"testing"

	"github.com/d5/tengo"
	"github.com/d5/tengo/compiler/token"
)

func TestString_BinaryOp(t *testing.T) {
	lstr := "abcde"
	rstr := "01234"
	for l := 0; l < len(lstr); l++ {
		for r := 0; r < len(rstr); r++ {
			ls := lstr[l:]
			rs := rstr[r:]
			testBinaryOp(t, &tengo.String{Value: ls}, token.Add, &tengo.String{Value: rs}, &tengo.String{Value: ls + rs})

			rc := []rune(rstr)[r]
			testBinaryOp(t, &tengo.String{Value: ls}, token.Add, &tengo.Char{Value: rc}, &tengo.String{Value: ls + string(rc)})
		}
	}
}

package objects_test

import (
	"testing"

	"github.com/d5/tengo/compiler/token"
	"github.com/d5/tengo/objects"
)

func TestString_BinaryOp(t *testing.T) {
	lstr := "abcde"
	rstr := "01234"
	for l := 0; l < len(lstr); l++ {
		for r := 0; r < len(rstr); r++ {
			ls := lstr[l:]
			rs := rstr[r:]
			testBinaryOp(t, &objects.String{Value: ls}, token.Add, &objects.String{Value: rs}, &objects.String{Value: ls + rs})

			rc := []rune(rstr)[r]
			testBinaryOp(t, &objects.String{Value: ls}, token.Add, &objects.Char{Value: rc}, &objects.String{Value: ls + string(rc)})
		}
	}
}

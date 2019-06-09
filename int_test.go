package tengo_test

import (
	"testing"

	"github.com/d5/tengo"
	"github.com/d5/tengo/compiler/token"
)

func TestInt_BinaryOp(t *testing.T) {
	// int + int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Add, &tengo.Int{Value: r}, &tengo.Int{Value: l + r})
		}
	}

	// int - int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Sub, &tengo.Int{Value: r}, &tengo.Int{Value: l - r})
		}
	}

	// int * int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Mul, &tengo.Int{Value: r}, &tengo.Int{Value: l * r})
		}
	}

	// int / int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			if r != 0 {
				testBinaryOp(t, &tengo.Int{Value: l}, token.Quo, &tengo.Int{Value: r}, &tengo.Int{Value: l / r})
			}
		}
	}

	// int % int
	for l := int64(-4); l <= 4; l++ {
		for r := -int64(-4); r <= 4; r++ {
			if r == 0 {
				testBinaryOp(t, &tengo.Int{Value: l}, token.Rem, &tengo.Int{Value: r}, &tengo.Int{Value: l % r})
			}
		}
	}

	// int & int
	testBinaryOp(t, &tengo.Int{Value: 0}, token.And, &tengo.Int{Value: 0}, &tengo.Int{Value: int64(0) & int64(0)})
	testBinaryOp(t, &tengo.Int{Value: 1}, token.And, &tengo.Int{Value: 0}, &tengo.Int{Value: int64(1) & int64(0)})
	testBinaryOp(t, &tengo.Int{Value: 0}, token.And, &tengo.Int{Value: 1}, &tengo.Int{Value: int64(0) & int64(1)})
	testBinaryOp(t, &tengo.Int{Value: 1}, token.And, &tengo.Int{Value: 1}, &tengo.Int{Value: int64(1) & int64(1)})
	testBinaryOp(t, &tengo.Int{Value: 0}, token.And, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(0) & int64(0xffffffff)})
	testBinaryOp(t, &tengo.Int{Value: 1}, token.And, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(1) & int64(0xffffffff)})
	testBinaryOp(t, &tengo.Int{Value: int64(0xffffffff)}, token.And, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(0xffffffff) & int64(0xffffffff)})
	testBinaryOp(t, &tengo.Int{Value: 1984}, token.And, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(1984) & int64(0xffffffff)})
	testBinaryOp(t, &tengo.Int{Value: -1984}, token.And, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(-1984) & int64(0xffffffff)})

	// int | int
	testBinaryOp(t, &tengo.Int{Value: 0}, token.Or, &tengo.Int{Value: 0}, &tengo.Int{Value: int64(0) | int64(0)})
	testBinaryOp(t, &tengo.Int{Value: 1}, token.Or, &tengo.Int{Value: 0}, &tengo.Int{Value: int64(1) | int64(0)})
	testBinaryOp(t, &tengo.Int{Value: 0}, token.Or, &tengo.Int{Value: 1}, &tengo.Int{Value: int64(0) | int64(1)})
	testBinaryOp(t, &tengo.Int{Value: 1}, token.Or, &tengo.Int{Value: 1}, &tengo.Int{Value: int64(1) | int64(1)})
	testBinaryOp(t, &tengo.Int{Value: 0}, token.Or, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(0) | int64(0xffffffff)})
	testBinaryOp(t, &tengo.Int{Value: 1}, token.Or, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(1) | int64(0xffffffff)})
	testBinaryOp(t, &tengo.Int{Value: int64(0xffffffff)}, token.Or, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(0xffffffff) | int64(0xffffffff)})
	testBinaryOp(t, &tengo.Int{Value: 1984}, token.Or, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(1984) | int64(0xffffffff)})
	testBinaryOp(t, &tengo.Int{Value: -1984}, token.Or, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(-1984) | int64(0xffffffff)})

	// int ^ int
	testBinaryOp(t, &tengo.Int{Value: 0}, token.Xor, &tengo.Int{Value: 0}, &tengo.Int{Value: int64(0) ^ int64(0)})
	testBinaryOp(t, &tengo.Int{Value: 1}, token.Xor, &tengo.Int{Value: 0}, &tengo.Int{Value: int64(1) ^ int64(0)})
	testBinaryOp(t, &tengo.Int{Value: 0}, token.Xor, &tengo.Int{Value: 1}, &tengo.Int{Value: int64(0) ^ int64(1)})
	testBinaryOp(t, &tengo.Int{Value: 1}, token.Xor, &tengo.Int{Value: 1}, &tengo.Int{Value: int64(1) ^ int64(1)})
	testBinaryOp(t, &tengo.Int{Value: 0}, token.Xor, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(0) ^ int64(0xffffffff)})
	testBinaryOp(t, &tengo.Int{Value: 1}, token.Xor, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(1) ^ int64(0xffffffff)})
	testBinaryOp(t, &tengo.Int{Value: int64(0xffffffff)}, token.Xor, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(0xffffffff) ^ int64(0xffffffff)})
	testBinaryOp(t, &tengo.Int{Value: 1984}, token.Xor, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(1984) ^ int64(0xffffffff)})
	testBinaryOp(t, &tengo.Int{Value: -1984}, token.Xor, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(-1984) ^ int64(0xffffffff)})

	// int &^ int
	testBinaryOp(t, &tengo.Int{Value: 0}, token.AndNot, &tengo.Int{Value: 0}, &tengo.Int{Value: int64(0) &^ int64(0)})
	testBinaryOp(t, &tengo.Int{Value: 1}, token.AndNot, &tengo.Int{Value: 0}, &tengo.Int{Value: int64(1) &^ int64(0)})
	testBinaryOp(t, &tengo.Int{Value: 0}, token.AndNot, &tengo.Int{Value: 1}, &tengo.Int{Value: int64(0) &^ int64(1)})
	testBinaryOp(t, &tengo.Int{Value: 1}, token.AndNot, &tengo.Int{Value: 1}, &tengo.Int{Value: int64(1) &^ int64(1)})
	testBinaryOp(t, &tengo.Int{Value: 0}, token.AndNot, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(0) &^ int64(0xffffffff)})
	testBinaryOp(t, &tengo.Int{Value: 1}, token.AndNot, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(1) &^ int64(0xffffffff)})
	testBinaryOp(t, &tengo.Int{Value: int64(0xffffffff)}, token.AndNot, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(0xffffffff) &^ int64(0xffffffff)})
	testBinaryOp(t, &tengo.Int{Value: 1984}, token.AndNot, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(1984) &^ int64(0xffffffff)})
	testBinaryOp(t, &tengo.Int{Value: -1984}, token.AndNot, &tengo.Int{Value: int64(0xffffffff)}, &tengo.Int{Value: int64(-1984) &^ int64(0xffffffff)})

	// int << int
	for s := int64(0); s < 64; s++ {
		testBinaryOp(t, &tengo.Int{Value: 0}, token.Shl, &tengo.Int{Value: s}, &tengo.Int{Value: int64(0) << uint(s)})
		testBinaryOp(t, &tengo.Int{Value: 1}, token.Shl, &tengo.Int{Value: s}, &tengo.Int{Value: int64(1) << uint(s)})
		testBinaryOp(t, &tengo.Int{Value: 2}, token.Shl, &tengo.Int{Value: s}, &tengo.Int{Value: int64(2) << uint(s)})
		testBinaryOp(t, &tengo.Int{Value: -1}, token.Shl, &tengo.Int{Value: s}, &tengo.Int{Value: int64(-1) << uint(s)})
		testBinaryOp(t, &tengo.Int{Value: -2}, token.Shl, &tengo.Int{Value: s}, &tengo.Int{Value: int64(-2) << uint(s)})
		testBinaryOp(t, &tengo.Int{Value: int64(0xffffffff)}, token.Shl, &tengo.Int{Value: s}, &tengo.Int{Value: int64(0xffffffff) << uint(s)})
	}

	// int >> int
	for s := int64(0); s < 64; s++ {
		testBinaryOp(t, &tengo.Int{Value: 0}, token.Shr, &tengo.Int{Value: s}, &tengo.Int{Value: int64(0) >> uint(s)})
		testBinaryOp(t, &tengo.Int{Value: 1}, token.Shr, &tengo.Int{Value: s}, &tengo.Int{Value: int64(1) >> uint(s)})
		testBinaryOp(t, &tengo.Int{Value: 2}, token.Shr, &tengo.Int{Value: s}, &tengo.Int{Value: int64(2) >> uint(s)})
		testBinaryOp(t, &tengo.Int{Value: -1}, token.Shr, &tengo.Int{Value: s}, &tengo.Int{Value: int64(-1) >> uint(s)})
		testBinaryOp(t, &tengo.Int{Value: -2}, token.Shr, &tengo.Int{Value: s}, &tengo.Int{Value: int64(-2) >> uint(s)})
		testBinaryOp(t, &tengo.Int{Value: int64(0xffffffff)}, token.Shr, &tengo.Int{Value: s}, &tengo.Int{Value: int64(0xffffffff) >> uint(s)})
	}

	// int < int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Less, &tengo.Int{Value: r}, boolValue(l < r))
		}
	}

	// int > int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Greater, &tengo.Int{Value: r}, boolValue(l > r))
		}
	}

	// int <= int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Int{Value: l}, token.LessEq, &tengo.Int{Value: r}, boolValue(l <= r))
		}
	}

	// int >= int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &tengo.Int{Value: l}, token.GreaterEq, &tengo.Int{Value: r}, boolValue(l >= r))
		}
	}

	// int + float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Add, &tengo.Float{Value: r}, &tengo.Float{Value: float64(l) + r})
		}
	}

	// int - float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Sub, &tengo.Float{Value: r}, &tengo.Float{Value: float64(l) - r})
		}
	}

	// int * float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Mul, &tengo.Float{Value: r}, &tengo.Float{Value: float64(l) * r})
		}
	}

	// int / float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			if r != 0 {
				testBinaryOp(t, &tengo.Int{Value: l}, token.Quo, &tengo.Float{Value: r}, &tengo.Float{Value: float64(l) / r})
			}
		}
	}

	// int < float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Less, &tengo.Float{Value: r}, boolValue(float64(l) < r))
		}
	}

	// int > float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &tengo.Int{Value: l}, token.Greater, &tengo.Float{Value: r}, boolValue(float64(l) > r))
		}
	}

	// int <= float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &tengo.Int{Value: l}, token.LessEq, &tengo.Float{Value: r}, boolValue(float64(l) <= r))
		}
	}

	// int >= float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &tengo.Int{Value: l}, token.GreaterEq, &tengo.Float{Value: r}, boolValue(float64(l) >= r))
		}
	}
}

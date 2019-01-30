package objects_test

import (
	"testing"

	"github.com/d5/tengo/compiler/token"
	"github.com/d5/tengo/objects"
)

func TestInt_BinaryOp(t *testing.T) {
	// int + int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &objects.Int{Value: l}, token.Add, &objects.Int{Value: r}, &objects.Int{Value: l + r})
		}
	}

	// int - int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &objects.Int{Value: l}, token.Sub, &objects.Int{Value: r}, &objects.Int{Value: l - r})
		}
	}

	// int * int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &objects.Int{Value: l}, token.Mul, &objects.Int{Value: r}, &objects.Int{Value: l * r})
		}
	}

	// int / int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			if r != 0 {
				testBinaryOp(t, &objects.Int{Value: l}, token.Quo, &objects.Int{Value: r}, &objects.Int{Value: l / r})
			}
		}
	}

	// int % int
	for l := int64(-4); l <= 4; l++ {
		for r := -int64(-4); r <= 4; r++ {
			if r == 0 {
				testBinaryOp(t, &objects.Int{Value: l}, token.Rem, &objects.Int{Value: r}, &objects.Int{Value: l % r})
			}
		}
	}

	// int & int
	testBinaryOp(t, &objects.Int{Value: 0}, token.And, &objects.Int{Value: 0}, &objects.Int{Value: int64(0) & int64(0)})
	testBinaryOp(t, &objects.Int{Value: 1}, token.And, &objects.Int{Value: 0}, &objects.Int{Value: int64(1) & int64(0)})
	testBinaryOp(t, &objects.Int{Value: 0}, token.And, &objects.Int{Value: 1}, &objects.Int{Value: int64(0) & int64(1)})
	testBinaryOp(t, &objects.Int{Value: 1}, token.And, &objects.Int{Value: 1}, &objects.Int{Value: int64(1) & int64(1)})
	testBinaryOp(t, &objects.Int{Value: 0}, token.And, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(0) & int64(0xffffffff)})
	testBinaryOp(t, &objects.Int{Value: 1}, token.And, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(1) & int64(0xffffffff)})
	testBinaryOp(t, &objects.Int{Value: int64(0xffffffff)}, token.And, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(0xffffffff) & int64(0xffffffff)})
	testBinaryOp(t, &objects.Int{Value: 1984}, token.And, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(1984) & int64(0xffffffff)})
	testBinaryOp(t, &objects.Int{Value: -1984}, token.And, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(-1984) & int64(0xffffffff)})

	// int | int
	testBinaryOp(t, &objects.Int{Value: 0}, token.Or, &objects.Int{Value: 0}, &objects.Int{Value: int64(0) | int64(0)})
	testBinaryOp(t, &objects.Int{Value: 1}, token.Or, &objects.Int{Value: 0}, &objects.Int{Value: int64(1) | int64(0)})
	testBinaryOp(t, &objects.Int{Value: 0}, token.Or, &objects.Int{Value: 1}, &objects.Int{Value: int64(0) | int64(1)})
	testBinaryOp(t, &objects.Int{Value: 1}, token.Or, &objects.Int{Value: 1}, &objects.Int{Value: int64(1) | int64(1)})
	testBinaryOp(t, &objects.Int{Value: 0}, token.Or, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(0) | int64(0xffffffff)})
	testBinaryOp(t, &objects.Int{Value: 1}, token.Or, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(1) | int64(0xffffffff)})
	testBinaryOp(t, &objects.Int{Value: int64(0xffffffff)}, token.Or, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(0xffffffff) | int64(0xffffffff)})
	testBinaryOp(t, &objects.Int{Value: 1984}, token.Or, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(1984) | int64(0xffffffff)})
	testBinaryOp(t, &objects.Int{Value: -1984}, token.Or, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(-1984) | int64(0xffffffff)})

	// int ^ int
	testBinaryOp(t, &objects.Int{Value: 0}, token.Xor, &objects.Int{Value: 0}, &objects.Int{Value: int64(0) ^ int64(0)})
	testBinaryOp(t, &objects.Int{Value: 1}, token.Xor, &objects.Int{Value: 0}, &objects.Int{Value: int64(1) ^ int64(0)})
	testBinaryOp(t, &objects.Int{Value: 0}, token.Xor, &objects.Int{Value: 1}, &objects.Int{Value: int64(0) ^ int64(1)})
	testBinaryOp(t, &objects.Int{Value: 1}, token.Xor, &objects.Int{Value: 1}, &objects.Int{Value: int64(1) ^ int64(1)})
	testBinaryOp(t, &objects.Int{Value: 0}, token.Xor, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(0) ^ int64(0xffffffff)})
	testBinaryOp(t, &objects.Int{Value: 1}, token.Xor, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(1) ^ int64(0xffffffff)})
	testBinaryOp(t, &objects.Int{Value: int64(0xffffffff)}, token.Xor, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(0xffffffff) ^ int64(0xffffffff)})
	testBinaryOp(t, &objects.Int{Value: 1984}, token.Xor, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(1984) ^ int64(0xffffffff)})
	testBinaryOp(t, &objects.Int{Value: -1984}, token.Xor, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(-1984) ^ int64(0xffffffff)})

	// int &^ int
	testBinaryOp(t, &objects.Int{Value: 0}, token.AndNot, &objects.Int{Value: 0}, &objects.Int{Value: int64(0) &^ int64(0)})
	testBinaryOp(t, &objects.Int{Value: 1}, token.AndNot, &objects.Int{Value: 0}, &objects.Int{Value: int64(1) &^ int64(0)})
	testBinaryOp(t, &objects.Int{Value: 0}, token.AndNot, &objects.Int{Value: 1}, &objects.Int{Value: int64(0) &^ int64(1)})
	testBinaryOp(t, &objects.Int{Value: 1}, token.AndNot, &objects.Int{Value: 1}, &objects.Int{Value: int64(1) &^ int64(1)})
	testBinaryOp(t, &objects.Int{Value: 0}, token.AndNot, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(0) &^ int64(0xffffffff)})
	testBinaryOp(t, &objects.Int{Value: 1}, token.AndNot, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(1) &^ int64(0xffffffff)})
	testBinaryOp(t, &objects.Int{Value: int64(0xffffffff)}, token.AndNot, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(0xffffffff) &^ int64(0xffffffff)})
	testBinaryOp(t, &objects.Int{Value: 1984}, token.AndNot, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(1984) &^ int64(0xffffffff)})
	testBinaryOp(t, &objects.Int{Value: -1984}, token.AndNot, &objects.Int{Value: int64(0xffffffff)}, &objects.Int{Value: int64(-1984) &^ int64(0xffffffff)})

	// int << int
	for s := int64(0); s < 64; s++ {
		testBinaryOp(t, &objects.Int{Value: 0}, token.Shl, &objects.Int{Value: s}, &objects.Int{Value: int64(0) << uint(s)})
		testBinaryOp(t, &objects.Int{Value: 1}, token.Shl, &objects.Int{Value: s}, &objects.Int{Value: int64(1) << uint(s)})
		testBinaryOp(t, &objects.Int{Value: 2}, token.Shl, &objects.Int{Value: s}, &objects.Int{Value: int64(2) << uint(s)})
		testBinaryOp(t, &objects.Int{Value: -1}, token.Shl, &objects.Int{Value: s}, &objects.Int{Value: int64(-1) << uint(s)})
		testBinaryOp(t, &objects.Int{Value: -2}, token.Shl, &objects.Int{Value: s}, &objects.Int{Value: int64(-2) << uint(s)})
		testBinaryOp(t, &objects.Int{Value: int64(0xffffffff)}, token.Shl, &objects.Int{Value: s}, &objects.Int{Value: int64(0xffffffff) << uint(s)})
	}

	// int >> int
	for s := int64(0); s < 64; s++ {
		testBinaryOp(t, &objects.Int{Value: 0}, token.Shr, &objects.Int{Value: s}, &objects.Int{Value: int64(0) >> uint(s)})
		testBinaryOp(t, &objects.Int{Value: 1}, token.Shr, &objects.Int{Value: s}, &objects.Int{Value: int64(1) >> uint(s)})
		testBinaryOp(t, &objects.Int{Value: 2}, token.Shr, &objects.Int{Value: s}, &objects.Int{Value: int64(2) >> uint(s)})
		testBinaryOp(t, &objects.Int{Value: -1}, token.Shr, &objects.Int{Value: s}, &objects.Int{Value: int64(-1) >> uint(s)})
		testBinaryOp(t, &objects.Int{Value: -2}, token.Shr, &objects.Int{Value: s}, &objects.Int{Value: int64(-2) >> uint(s)})
		testBinaryOp(t, &objects.Int{Value: int64(0xffffffff)}, token.Shr, &objects.Int{Value: s}, &objects.Int{Value: int64(0xffffffff) >> uint(s)})
	}

	// int < int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &objects.Int{Value: l}, token.Less, &objects.Int{Value: r}, boolValue(l < r))
		}
	}

	// int > int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &objects.Int{Value: l}, token.Greater, &objects.Int{Value: r}, boolValue(l > r))
		}
	}

	// int <= int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &objects.Int{Value: l}, token.LessEq, &objects.Int{Value: r}, boolValue(l <= r))
		}
	}

	// int >= int
	for l := int64(-2); l <= 2; l++ {
		for r := int64(-2); r <= 2; r++ {
			testBinaryOp(t, &objects.Int{Value: l}, token.GreaterEq, &objects.Int{Value: r}, boolValue(l >= r))
		}
	}

	// int + float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &objects.Int{Value: l}, token.Add, &objects.Float{Value: r}, &objects.Float{Value: float64(l) + r})
		}
	}

	// int - float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &objects.Int{Value: l}, token.Sub, &objects.Float{Value: r}, &objects.Float{Value: float64(l) - r})
		}
	}

	// int * float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &objects.Int{Value: l}, token.Mul, &objects.Float{Value: r}, &objects.Float{Value: float64(l) * r})
		}
	}

	// int / float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			if r != 0 {
				testBinaryOp(t, &objects.Int{Value: l}, token.Quo, &objects.Float{Value: r}, &objects.Float{Value: float64(l) / r})
			}
		}
	}

	// int < float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &objects.Int{Value: l}, token.Less, &objects.Float{Value: r}, boolValue(float64(l) < r))
		}
	}

	// int > float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &objects.Int{Value: l}, token.Greater, &objects.Float{Value: r}, boolValue(float64(l) > r))
		}
	}

	// int <= float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &objects.Int{Value: l}, token.LessEq, &objects.Float{Value: r}, boolValue(float64(l) <= r))
		}
	}

	// int >= float
	for l := int64(-2); l <= 2; l++ {
		for r := float64(-2); r <= 2.1; r += 0.5 {
			testBinaryOp(t, &objects.Int{Value: l}, token.GreaterEq, &objects.Float{Value: r}, boolValue(float64(l) >= r))
		}
	}
}

package tengo_test

import (
	"testing"

	"github.com/d5/tengo"
	"github.com/d5/tengo/compiler/token"
)

func TestArray_BinaryOp(t *testing.T) {
	testBinaryOp(t, &tengo.Array{Value: nil}, token.Add, &tengo.Array{Value: nil}, &tengo.Array{Value: nil})
	testBinaryOp(t, &tengo.Array{Value: nil}, token.Add, &tengo.Array{Value: []tengo.Object{}}, &tengo.Array{Value: nil})
	testBinaryOp(t, &tengo.Array{Value: []tengo.Object{}}, token.Add, &tengo.Array{Value: nil}, &tengo.Array{Value: []tengo.Object{}})
	testBinaryOp(t, &tengo.Array{Value: []tengo.Object{}}, token.Add, &tengo.Array{Value: []tengo.Object{}}, &tengo.Array{Value: []tengo.Object{}})
	testBinaryOp(t, &tengo.Array{Value: nil}, token.Add, &tengo.Array{Value: []tengo.Object{
		&tengo.Int{Value: 1},
	}}, &tengo.Array{Value: []tengo.Object{
		&tengo.Int{Value: 1},
	}})
	testBinaryOp(t, &tengo.Array{Value: nil}, token.Add, &tengo.Array{Value: []tengo.Object{
		&tengo.Int{Value: 1},
		&tengo.Int{Value: 2},
		&tengo.Int{Value: 3},
	}}, &tengo.Array{Value: []tengo.Object{
		&tengo.Int{Value: 1},
		&tengo.Int{Value: 2},
		&tengo.Int{Value: 3},
	}})
	testBinaryOp(t, &tengo.Array{Value: []tengo.Object{
		&tengo.Int{Value: 1},
		&tengo.Int{Value: 2},
		&tengo.Int{Value: 3},
	}}, token.Add, &tengo.Array{Value: nil}, &tengo.Array{Value: []tengo.Object{
		&tengo.Int{Value: 1},
		&tengo.Int{Value: 2},
		&tengo.Int{Value: 3},
	}})
	testBinaryOp(t, &tengo.Array{Value: []tengo.Object{
		&tengo.Int{Value: 1},
		&tengo.Int{Value: 2},
		&tengo.Int{Value: 3},
	}}, token.Add, &tengo.Array{Value: []tengo.Object{
		&tengo.Int{Value: 4},
		&tengo.Int{Value: 5},
		&tengo.Int{Value: 6},
	}}, &tengo.Array{Value: []tengo.Object{
		&tengo.Int{Value: 1},
		&tengo.Int{Value: 2},
		&tengo.Int{Value: 3},
		&tengo.Int{Value: 4},
		&tengo.Int{Value: 5},
		&tengo.Int{Value: 6},
	}})
}

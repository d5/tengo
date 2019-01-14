package objects_test

import (
	"testing"

	"github.com/d5/tengo/compiler/token"
	"github.com/d5/tengo/objects"
)

func TestArray_BinaryOp(t *testing.T) {
	testBinaryOp(t, &objects.Array{Value: nil}, token.Add, &objects.Array{Value: nil}, &objects.Array{Value: nil})
	testBinaryOp(t, &objects.Array{Value: nil}, token.Add, &objects.Array{Value: []objects.Object{}}, &objects.Array{Value: nil})
	testBinaryOp(t, &objects.Array{Value: []objects.Object{}}, token.Add, &objects.Array{Value: nil}, &objects.Array{Value: []objects.Object{}})
	testBinaryOp(t, &objects.Array{Value: []objects.Object{}}, token.Add, &objects.Array{Value: []objects.Object{}}, &objects.Array{Value: []objects.Object{}})
	testBinaryOp(t, &objects.Array{Value: nil}, token.Add, &objects.Array{Value: []objects.Object{
		&objects.Int{Value: 1},
	}}, &objects.Array{Value: []objects.Object{
		&objects.Int{Value: 1},
	}})
	testBinaryOp(t, &objects.Array{Value: nil}, token.Add, &objects.Array{Value: []objects.Object{
		&objects.Int{Value: 1},
		&objects.Int{Value: 2},
		&objects.Int{Value: 3},
	}}, &objects.Array{Value: []objects.Object{
		&objects.Int{Value: 1},
		&objects.Int{Value: 2},
		&objects.Int{Value: 3},
	}})
	testBinaryOp(t, &objects.Array{Value: []objects.Object{
		&objects.Int{Value: 1},
		&objects.Int{Value: 2},
		&objects.Int{Value: 3},
	}}, token.Add, &objects.Array{Value: nil}, &objects.Array{Value: []objects.Object{
		&objects.Int{Value: 1},
		&objects.Int{Value: 2},
		&objects.Int{Value: 3},
	}})
	testBinaryOp(t, &objects.Array{Value: []objects.Object{
		&objects.Int{Value: 1},
		&objects.Int{Value: 2},
		&objects.Int{Value: 3},
	}}, token.Add, &objects.Array{Value: []objects.Object{
		&objects.Int{Value: 4},
		&objects.Int{Value: 5},
		&objects.Int{Value: 6},
	}}, &objects.Array{Value: []objects.Object{
		&objects.Int{Value: 1},
		&objects.Int{Value: 2},
		&objects.Int{Value: 3},
		&objects.Int{Value: 4},
		&objects.Int{Value: 5},
		&objects.Int{Value: 6},
	}})
}

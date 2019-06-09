package tengo_test

import (
	"testing"
	"time"

	"github.com/d5/tengo"
	"github.com/d5/tengo/assert"
)

func TestNumObjects(t *testing.T) {
	testCountObjects(t, &tengo.Array{}, 1)
	testCountObjects(t, &tengo.Array{Value: []tengo.Object{
		&tengo.Int{Value: 1},
		&tengo.Int{Value: 2},
		&tengo.Array{Value: []tengo.Object{
			&tengo.Int{Value: 3},
			&tengo.Int{Value: 4},
			&tengo.Int{Value: 5},
		}},
	}}, 7)
	testCountObjects(t, tengo.TrueValue, 1)
	testCountObjects(t, tengo.FalseValue, 1)
	testCountObjects(t, &tengo.BuiltinFunction{}, 1)
	testCountObjects(t, &tengo.Bytes{Value: []byte("foobar")}, 1)
	testCountObjects(t, &tengo.Char{Value: '가'}, 1)
	testCountObjects(t, &tengo.CompiledFunction{}, 1)
	testCountObjects(t, &tengo.Error{Value: &tengo.Int{Value: 5}}, 2)
	testCountObjects(t, &tengo.Float{Value: 19.84}, 1)
	testCountObjects(t, &tengo.ImmutableArray{Value: []tengo.Object{
		&tengo.Int{Value: 1},
		&tengo.Int{Value: 2},
		&tengo.ImmutableArray{Value: []tengo.Object{
			&tengo.Int{Value: 3},
			&tengo.Int{Value: 4},
			&tengo.Int{Value: 5},
		}},
	}}, 7)
	testCountObjects(t, &tengo.ImmutableMap{Value: map[string]tengo.Object{
		"k1": &tengo.Int{Value: 1},
		"k2": &tengo.Int{Value: 2},
		"k3": &tengo.Array{Value: []tengo.Object{
			&tengo.Int{Value: 3},
			&tengo.Int{Value: 4},
			&tengo.Int{Value: 5},
		}},
	}}, 7)
	testCountObjects(t, &tengo.Int{Value: 1984}, 1)
	testCountObjects(t, &tengo.Map{Value: map[string]tengo.Object{
		"k1": &tengo.Int{Value: 1},
		"k2": &tengo.Int{Value: 2},
		"k3": &tengo.Array{Value: []tengo.Object{
			&tengo.Int{Value: 3},
			&tengo.Int{Value: 4},
			&tengo.Int{Value: 5},
		}},
	}}, 7)
	testCountObjects(t, &tengo.String{Value: "foo bar"}, 1)
	testCountObjects(t, &tengo.Time{Value: time.Now()}, 1)
	testCountObjects(t, tengo.UndefinedValue, 1)
}

func testCountObjects(t *testing.T, o tengo.Object, expected int) {
	assert.Equal(t, expected, tengo.CountObjects(o))
}

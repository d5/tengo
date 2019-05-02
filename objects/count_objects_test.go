package objects_test

import (
	"testing"
	"time"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/objects"
)

func TestNumObjects(t *testing.T) {
	testCountObjects(t, &objects.Array{}, 1)
	testCountObjects(t, &objects.Array{Value: []objects.Object{
		&objects.Int{Value: 1},
		&objects.Int{Value: 2},
		&objects.Array{Value: []objects.Object{
			&objects.Int{Value: 3},
			&objects.Int{Value: 4},
			&objects.Int{Value: 5},
		}},
	}}, 7)
	testCountObjects(t, objects.TrueValue, 1)
	testCountObjects(t, objects.FalseValue, 1)
	testCountObjects(t, &objects.BuiltinFunction{}, 1)
	testCountObjects(t, &objects.Bytes{Value: []byte("foobar")}, 1)
	testCountObjects(t, &objects.Char{Value: '가'}, 1)
	testCountObjects(t, &objects.CompiledFunction{}, 1)
	testCountObjects(t, &objects.Error{Value: &objects.Int{Value: 5}}, 2)
	testCountObjects(t, &objects.Float{Value: 19.84}, 1)
	testCountObjects(t, &objects.ImmutableArray{Value: []objects.Object{
		&objects.Int{Value: 1},
		&objects.Int{Value: 2},
		&objects.ImmutableArray{Value: []objects.Object{
			&objects.Int{Value: 3},
			&objects.Int{Value: 4},
			&objects.Int{Value: 5},
		}},
	}}, 7)
	testCountObjects(t, &objects.ImmutableMap{Value: map[string]objects.Object{
		"k1": &objects.Int{Value: 1},
		"k2": &objects.Int{Value: 2},
		"k3": &objects.Array{Value: []objects.Object{
			&objects.Int{Value: 3},
			&objects.Int{Value: 4},
			&objects.Int{Value: 5},
		}},
	}}, 7)
	testCountObjects(t, &objects.Int{Value: 1984}, 1)
	testCountObjects(t, &objects.Map{Value: map[string]objects.Object{
		"k1": &objects.Int{Value: 1},
		"k2": &objects.Int{Value: 2},
		"k3": &objects.Array{Value: []objects.Object{
			&objects.Int{Value: 3},
			&objects.Int{Value: 4},
			&objects.Int{Value: 5},
		}},
	}}, 7)
	testCountObjects(t, &objects.String{Value: "foo bar"}, 1)
	testCountObjects(t, &objects.Time{Value: time.Now()}, 1)
	testCountObjects(t, objects.UndefinedValue, 1)
}

func testCountObjects(t *testing.T, o objects.Object, expected int) {
	assert.Equal(t, expected, objects.CountObjects(o))
}

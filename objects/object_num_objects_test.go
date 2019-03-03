package objects_test

import (
	"testing"
	"time"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/objects"
)

func TestNumObjects(t *testing.T) {
	testNumObjects(t, &objects.Array{}, 1)
	testNumObjects(t, &objects.Array{Value: []objects.Object{
		&objects.Int{Value: 1},
		&objects.Int{Value: 2},
		&objects.Array{Value: []objects.Object{
			&objects.Int{Value: 3},
			&objects.Int{Value: 4},
			&objects.Int{Value: 5},
		}},
	}}, 7)
	testNumObjects(t, objects.TrueValue, 0)
	testNumObjects(t, objects.FalseValue, 0)
	testNumObjects(t, &objects.BuiltinFunction{}, 0)
	testNumObjects(t, &objects.Bytes{Value: []byte("foobar")}, 1)
	testNumObjects(t, &objects.Char{Value: '가'}, 1)
	testNumObjects(t, &objects.Closure{}, 0)
	testNumObjects(t, &objects.CompiledFunction{}, 0)
	testNumObjects(t, &objects.Error{Value: &objects.Int{Value: 5}}, 2)
	testNumObjects(t, &objects.Float{Value: 19.84}, 1)
	testNumObjects(t, &objects.ImmutableArray{Value: []objects.Object{
		&objects.Int{Value: 1},
		&objects.Int{Value: 2},
		&objects.ImmutableArray{Value: []objects.Object{
			&objects.Int{Value: 3},
			&objects.Int{Value: 4},
			&objects.Int{Value: 5},
		}},
	}}, 7)
	testNumObjects(t, &objects.ImmutableMap{Value: map[string]objects.Object{
		"k1": &objects.Int{Value: 1},
		"k2": &objects.Int{Value: 2},
		"k3": &objects.Array{Value: []objects.Object{
			&objects.Int{Value: 3},
			&objects.Int{Value: 4},
			&objects.Int{Value: 5},
		}},
	}}, 7)
	testNumObjects(t, &objects.Int{Value: 1984}, 1)
	testNumObjects(t, &objects.Map{Value: map[string]objects.Object{
		"k1": &objects.Int{Value: 1},
		"k2": &objects.Int{Value: 2},
		"k3": &objects.Array{Value: []objects.Object{
			&objects.Int{Value: 3},
			&objects.Int{Value: 4},
			&objects.Int{Value: 5},
		}},
	}}, 7)
	testNumObjects(t, &objects.String{Value: "foo bar"}, 1)
	testNumObjects(t, &objects.Time{Value: time.Now()}, 1)
	testNumObjects(t, objects.UndefinedValue, 0)
}

func testNumObjects(t *testing.T, o objects.Object, expected int64) {
	assert.Equal(t, expected, o.NumObjects())
}

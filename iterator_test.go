package tengo_test

import (
	"testing"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/require"
)

func TestIterator(t *testing.T) {
	ints := func(vs ...int64) []tengo.Object {
		var objs []tengo.Object
		for _, v := range vs {
			objs = append(objs, &tengo.Int{Value: v})
		}
		return objs
	}

	for _, tc := range []struct {
		obj    tengo.Object
		keys   []tengo.Object
		values []tengo.Object
	}{
		{&tengo.Array{}, nil, nil},
		{
			&tengo.Array{Value: []tengo.Object{
				&tengo.Int{Value: 10}, &tengo.String{Value: "x"}}},
			ints(0, 1),
			[]tengo.Object{&tengo.Int{Value: 10}, &tengo.String{Value: "x"}},
		},
		{
			&tengo.ImmutableArray{Value: []tengo.Object{
				&tengo.Int{Value: 10}, &tengo.String{Value: "x"}}},
			ints(0, 1),
			[]tengo.Object{&tengo.Int{Value: 10}, &tengo.String{Value: "x"}},
		},
		{&tengo.Bytes{}, nil, nil},
		{&tengo.Bytes{Value: []byte{7, 255}}, ints(0, 1), ints(7, 255)},
		{&tengo.String{}, nil, nil},
		{
			&tengo.String{Value: "가b"},
			ints(0, 1),
			[]tengo.Object{&tengo.Char{Value: '가'}, &tengo.Char{Value: 'b'}},
		},
	} {
		require.True(t, tc.obj.CanIterate(), "%s", tc.obj)
		it := tc.obj.Iterate()
		keys, values := drainIterator(it)
		require.Equal(t, tc.keys, keys, "%s", tc.obj)
		require.Equal(t, tc.values, values, "%s", tc.obj)
		require.False(t, it.Next(), "%s", tc.obj)
	}
}

func TestMapIterator(t *testing.T) {
	kv := map[string]tengo.Object{
		"a": &tengo.Int{Value: 1},
		"b": &tengo.String{Value: "x"},
		"c": tengo.TrueValue,
	}
	for _, obj := range []tengo.Object{
		&tengo.Map{Value: kv},
		&tengo.ImmutableMap{Value: kv},
	} {
		it := obj.Iterate()
		seen := make(map[string]tengo.Object)
		for it.Next() {
			key, ok := it.Key().(*tengo.String)
			require.True(t, ok, obj.TypeName())
			seen[key.Value] = it.Value()
		}
		require.Equal(t, &tengo.Map{Value: kv}, &tengo.Map{Value: seen},
			obj.TypeName())
		require.False(t, it.Next(), obj.TypeName())
	}

	_, values := drainIterator((&tengo.Map{}).Iterate())
	require.Equal(t, 0, len(values))
}

func TestIterator_Copy(t *testing.T) {
	for _, obj := range []tengo.Object{
		&tengo.Array{Value: []tengo.Object{
			&tengo.Int{Value: 1}, &tengo.Int{Value: 2}, &tengo.Int{Value: 3}}},
		&tengo.Bytes{Value: []byte("abc")},
		&tengo.String{Value: "abc"},
		&tengo.Map{Value: map[string]tengo.Object{
			"a": &tengo.Int{Value: 1},
			"b": &tengo.Int{Value: 2},
			"c": &tengo.Int{Value: 3},
		}},
	} {
		it := obj.Iterate()
		require.True(t, it.Next(), obj.TypeName())

		c, ok := it.Copy().(tengo.Iterator)
		require.True(t, ok, obj.TypeName())
		require.Equal(t, it.TypeName(), c.TypeName())
		require.Equal(t, it.Key(), c.Key(), obj.TypeName())
		require.Equal(t, it.Value(), c.Value(), obj.TypeName())

		keys, values := drainIterator(it)
		require.Equal(t, 2, len(keys), obj.TypeName())
		copyKeys, copyValues := drainIterator(c)
		require.Equal(t, keys, copyKeys, obj.TypeName())
		require.Equal(t, values, copyValues, obj.TypeName())
	}
}

func TestIterator_Equals(t *testing.T) {
	for _, obj := range []tengo.Object{
		&tengo.Array{Value: []tengo.Object{&tengo.Int{Value: 1}}},
		&tengo.Bytes{Value: []byte("a")},
		&tengo.String{Value: "a"},
		&tengo.Map{Value: map[string]tengo.Object{"a": &tengo.Int{Value: 1}}},
	} {
		it := obj.Iterate()
		require.False(t, it.Equals(it.Copy()), it.TypeName())
		require.False(t, it.Equals(obj.Iterate()), it.TypeName())
	}
}

func drainIterator(it tengo.Iterator) (keys, values []tengo.Object) {
	for it.Next() {
		keys = append(keys, it.Key())
		values = append(values, it.Value())
	}
	return
}

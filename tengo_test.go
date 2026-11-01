package tengo_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/parser"
	"github.com/d5/tengo/v2/require"
)

func TestInstructions_String(t *testing.T) {
	assertInstructionString(t,
		[][]byte{
			tengo.MakeInstruction(parser.OpConstant, 1),
			tengo.MakeInstruction(parser.OpConstant, 2),
			tengo.MakeInstruction(parser.OpConstant, 65535),
		},
		`0000 CONST   1    
0003 CONST   2    
0006 CONST   65535`)

	assertInstructionString(t,
		[][]byte{
			tengo.MakeInstruction(parser.OpBinaryOp, 11),
			tengo.MakeInstruction(parser.OpConstant, 2),
			tengo.MakeInstruction(parser.OpConstant, 65535),
		},
		`0000 BINARYOP 11   
0002 CONST   2    
0005 CONST   65535`)

	assertInstructionString(t,
		[][]byte{
			tengo.MakeInstruction(parser.OpBinaryOp, 11),
			tengo.MakeInstruction(parser.OpGetLocal, 1),
			tengo.MakeInstruction(parser.OpConstant, 2),
			tengo.MakeInstruction(parser.OpConstant, 65535),
		},
		`0000 BINARYOP 11   
0002 GETL    1    
0004 CONST   2    
0007 CONST   65535`)
}

func TestMakeInstruction(t *testing.T) {
	makeInstruction(t, []byte{parser.OpConstant, 0, 0},
		parser.OpConstant, 0)
	makeInstruction(t, []byte{parser.OpConstant, 0, 1},
		parser.OpConstant, 1)
	makeInstruction(t, []byte{parser.OpConstant, 255, 254},
		parser.OpConstant, 65534)
	makeInstruction(t, []byte{parser.OpPop}, parser.OpPop)
	makeInstruction(t, []byte{parser.OpTrue}, parser.OpTrue)
	makeInstruction(t, []byte{parser.OpFalse}, parser.OpFalse)
}

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
	testCountObjects(t, &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
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
	require.Equal(t, expected, tengo.CountObjects(o))
}

func assertInstructionString(
	t *testing.T,
	instructions [][]byte,
	expected string,
) {
	concatted := make([]byte, 0)
	for _, e := range instructions {
		concatted = append(concatted, e...)
	}
	require.Equal(t, expected, strings.Join(
		tengo.FormatInstructions(concatted, 0), "\n"))
}

func makeInstruction(
	t *testing.T,
	expected []byte,
	opcode parser.Opcode,
	operands ...int,
) {
	inst := tengo.MakeInstruction(opcode, operands...)
	require.Equal(t, expected, inst)
}

func TestToInt(t *testing.T) {
	for _, tc := range []struct {
		o        tengo.Object
		expected int
		ok       bool
	}{
		{&tengo.Int{Value: -7}, -7, true},
		{&tengo.Float{Value: 1.9}, 1, true},
		{&tengo.Float{Value: -1.9}, -1, true},
		{&tengo.Char{Value: 'a'}, 97, true},
		{tengo.TrueValue, 1, true},
		{tengo.FalseValue, 0, true},
		{&tengo.String{Value: "-42"}, -42, true},
		{&tengo.String{Value: "1.5"}, 0, false},
		{&tengo.String{Value: "abc"}, 0, false},
		{&tengo.Bytes{Value: []byte("1")}, 0, false},
		{&tengo.Array{}, 0, false},
		{tengo.UndefinedValue, 0, false},
	} {
		v, ok := tengo.ToInt(tc.o)
		require.Equal(t, tc.ok, ok, "%s", tc.o)
		require.Equal(t, tc.expected, v, "%s", tc.o)
	}
}

func TestToRune(t *testing.T) {
	for _, tc := range []struct {
		o        tengo.Object
		expected rune
		ok       bool
	}{
		{&tengo.Int{Value: 97}, 'a', true},
		{&tengo.Char{Value: '가'}, '가', true},
		{&tengo.String{Value: "a"}, 0, false},
		{&tengo.Float{Value: 97}, 0, false},
		{tengo.TrueValue, 0, false},
		{tengo.UndefinedValue, 0, false},
	} {
		v, ok := tengo.ToRune(tc.o)
		require.Equal(t, tc.ok, ok, "%s", tc.o)
		require.Equal(t, tc.expected, v, "%s", tc.o)
	}
}

func TestToByteSlice(t *testing.T) {
	v, ok := tengo.ToByteSlice(&tengo.Bytes{Value: []byte{0, 1, 255}})
	require.True(t, ok)
	require.Equal(t, []byte{0, 1, 255}, v)

	v, ok = tengo.ToByteSlice(&tengo.String{Value: "가a"})
	require.True(t, ok)
	require.Equal(t, []byte("가a"), v)

	for _, o := range []tengo.Object{
		&tengo.Int{Value: 1},
		&tengo.Char{Value: 'a'},
		&tengo.Array{Value: []tengo.Object{&tengo.Int{Value: 1}}},
		tengo.UndefinedValue,
	} {
		v, ok = tengo.ToByteSlice(o)
		require.False(t, ok, "%s", o)
		require.Nil(t, v, "%s", o)
	}
}

func TestToTime(t *testing.T) {
	tm := time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)
	v, ok := tengo.ToTime(&tengo.Time{Value: tm})
	require.True(t, ok)
	require.True(t, v.Equal(tm))

	v, ok = tengo.ToTime(&tengo.Int{Value: tm.Unix()})
	require.True(t, ok)
	require.True(t, v.Equal(tm))

	for _, o := range []tengo.Object{
		&tengo.String{Value: "2009-11-10T23:00:00Z"},
		&tengo.Float{Value: float64(tm.Unix())},
		tengo.UndefinedValue,
	} {
		v, ok = tengo.ToTime(o)
		require.False(t, ok, "%s", o)
		require.True(t, v.IsZero(), "%s", o)
	}
}

func TestToInterface(t *testing.T) {
	tm := time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)
	builtin := &tengo.BuiltinFunction{Name: "fn"}
	for _, tc := range []struct {
		o        tengo.Object
		expected interface{}
	}{
		{&tengo.Int{Value: 1}, int64(1)},
		{&tengo.String{Value: "a"}, "a"},
		{&tengo.Float{Value: 1.5}, 1.5},
		{tengo.TrueValue, true},
		{tengo.FalseValue, false},
		{&tengo.Char{Value: 'a'}, 'a'},
		{&tengo.Bytes{Value: []byte("a")}, []byte("a")},
		{&tengo.Time{Value: tm}, tm},
		{tengo.UndefinedValue, nil},
		{builtin, builtin},
		{&tengo.Array{}, []interface{}{}},
		{
			&tengo.Array{Value: []tengo.Object{
				&tengo.Int{Value: 1},
				&tengo.Array{Value: []tengo.Object{&tengo.String{Value: "a"}}},
				tengo.UndefinedValue,
			}},
			[]interface{}{int64(1), []interface{}{"a"}, nil},
		},
		{
			&tengo.ImmutableArray{Value: []tengo.Object{
				&tengo.Int{Value: 1},
				&tengo.Map{Value: map[string]tengo.Object{
					"a": tengo.TrueValue}},
			}},
			[]interface{}{int64(1), map[string]interface{}{"a": true}},
		},
		{&tengo.Map{}, map[string]interface{}{}},
		{
			&tengo.Map{Value: map[string]tengo.Object{
				"a": &tengo.Int{Value: 1},
				"b": &tengo.ImmutableArray{Value: []tengo.Object{
					&tengo.Char{Value: 'x'}}},
			}},
			map[string]interface{}{"a": int64(1), "b": []interface{}{'x'}},
		},
		{
			&tengo.ImmutableMap{Value: map[string]tengo.Object{
				"a": &tengo.Float{Value: 2.5},
			}},
			map[string]interface{}{"a": 2.5},
		},
	} {
		actual := tengo.ToInterface(tc.o)
		require.True(t, reflect.DeepEqual(tc.expected, actual),
			"%s: expected %#v, got %#v", tc.o.TypeName(), tc.expected, actual)
	}

	err, ok := tengo.ToInterface(&tengo.Error{
		Value: &tengo.String{Value: "boom"}}).(error)
	require.True(t, ok)
	require.Equal(t, `error: "boom"`, err.Error())
}

func TestFromInterface(t *testing.T) {
	tm := time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)
	obj := &tengo.Int{Value: 5}
	objs := map[string]tengo.Object{"a": &tengo.Int{Value: 1}}
	arr := []tengo.Object{&tengo.Int{Value: 1}, &tengo.String{Value: "a"}}

	testFromInterface(t, nil, tengo.UndefinedValue)
	testFromInterface(t, "abc", &tengo.String{Value: "abc"})
	testFromInterface(t, int64(-1), &tengo.Int{Value: -1})
	testFromInterface(t, 7, &tengo.Int{Value: 7})
	testFromInterface(t, true, tengo.TrueValue)
	testFromInterface(t, false, tengo.FalseValue)
	testFromInterface(t, '가', &tengo.Char{Value: '가'})
	testFromInterface(t, byte('a'), &tengo.Char{Value: 'a'})
	testFromInterface(t, 1.5, &tengo.Float{Value: 1.5})
	testFromInterface(t, []byte("abc"), &tengo.Bytes{Value: []byte("abc")})
	testFromInterface(t, errors.New("boom"),
		&tengo.Error{Value: &tengo.String{Value: "boom"}})
	testFromInterface(t, objs, &tengo.Map{Value: objs})
	testFromInterface(t, arr, &tengo.Array{Value: arr})
	testFromInterface(t, tm, &tengo.Time{Value: tm})
	testFromInterface(t, map[string]interface{}{
		"a": 1,
		"b": []interface{}{"x", true, nil},
	}, &tengo.Map{Value: map[string]tengo.Object{
		"a": &tengo.Int{Value: 1},
		"b": &tengo.Array{Value: []tengo.Object{
			&tengo.String{Value: "x"}, tengo.TrueValue, tengo.UndefinedValue,
		}},
	}})
	testFromInterface(t, []interface{}{}, &tengo.Array{Value: []tengo.Object{}})

	res, err := tengo.FromInterface(obj)
	require.NoError(t, err)
	require.True(t, res == obj, "objects must not be converted")

	res, err = tengo.FromInterface(func(args ...tengo.Object) (tengo.Object, error) {
		return &tengo.Int{Value: int64(len(args))}, nil
	})
	require.NoError(t, err)
	fn, ok := res.(*tengo.UserFunction)
	require.True(t, ok)
	ret, err := fn.Call(tengo.TrueValue, tengo.FalseValue)
	require.NoError(t, err)
	require.Equal(t, &tengo.Int{Value: 2}, ret)

	for _, v := range []interface{}{
		make(chan int),
		struct{}{},
		map[string]interface{}{"a": struct{}{}},
		[]interface{}{1, struct{}{}},
		map[string]interface{}{"a": []interface{}{struct{}{}}},
	} {
		_, err = tengo.FromInterface(v)
		require.Error(t, err, "%#v", v)
		require.True(t,
			strings.HasPrefix(err.Error(), "cannot convert to object: "),
			"%#v: %s", v, err)
	}
	_, err = tengo.FromInterface(float32(1))
	require.Equal(t, "cannot convert to object: float32", err.Error())
}

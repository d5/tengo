package json_test

import (
	gojson "encoding/json"
	"math"
	"testing"
	"time"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/require"
	"github.com/d5/tengo/v2/stdlib/json"
)

type ARR = []interface{}
type MAP = map[string]interface{}

func TestJSON(t *testing.T) {
	testJSONEncodeDecode(t, nil)

	testJSONEncodeDecode(t, 0)
	testJSONEncodeDecode(t, 1)
	testJSONEncodeDecode(t, -1)
	testJSONEncodeDecode(t, 1984)
	testJSONEncodeDecode(t, -1984)

	testJSONEncodeDecode(t, 0.0)
	testJSONEncodeDecode(t, 1.0)
	testJSONEncodeDecode(t, -1.0)
	testJSONEncodeDecode(t, 19.84)
	testJSONEncodeDecode(t, -19.84)

	testJSONEncodeDecode(t, "")
	testJSONEncodeDecode(t, "foo")
	testJSONEncodeDecode(t, "foo bar")
	testJSONEncodeDecode(t, "foo \"bar\"")
	// See: https://github.com/d5/tengo/issues/268
	testJSONEncodeDecode(t, "1\u001C04")
	testJSONEncodeDecode(t, "çığöşü")
	testJSONEncodeDecode(t, "ç1\u001C04IĞÖŞÜ")
	testJSONEncodeDecode(t, "错误测试")

	testJSONEncodeDecode(t, true)
	testJSONEncodeDecode(t, false)

	testJSONEncodeDecode(t, ARR{})
	testJSONEncodeDecode(t, ARR{0})
	testJSONEncodeDecode(t, ARR{false})
	testJSONEncodeDecode(t, ARR{1, 2, 3,
		"four", false})
	testJSONEncodeDecode(t, ARR{1, 2, 3,
		"four", false, MAP{"a": 0, "b": "bee", "bool": true}})

	testJSONEncodeDecode(t, MAP{})
	testJSONEncodeDecode(t, MAP{"a": 0})
	testJSONEncodeDecode(t, MAP{"a": 0, "b": "bee"})
	testJSONEncodeDecode(t, MAP{"a": 0, "b": "bee", "bool": true})

	testJSONEncodeDecode(t, MAP{"a": 0, "b": "bee",
		"arr": ARR{1, 2, 3, "four"}})
	testJSONEncodeDecode(t, MAP{"a": 0, "b": "bee",
		"arr": ARR{1, 2, 3, MAP{"a": false, "b": 109.4}}})

	testJSONEncodeDecode(t, MAP{"id1": 7075984636689534001, "id2": 7075984636689534002})
	testJSONEncodeDecode(t, ARR{1e3, 1e7})
}

func TestDecode(t *testing.T) {
	testDecodeError(t, `{`)
	testDecodeError(t, `}`)
	testDecodeError(t, `{}a`)
	testDecodeError(t, `{{}`)
	testDecodeError(t, `{}}`)
	testDecodeError(t, `[`)
	testDecodeError(t, `]`)
	testDecodeError(t, `[]a`)
	testDecodeError(t, `[[]`)
	testDecodeError(t, `[]]`)
	testDecodeError(t, `"`)
	testDecodeError(t, `"abc`)
	testDecodeError(t, `abc"`)
	testDecodeError(t, `.123`)
	testDecodeError(t, `123.`)
	testDecodeError(t, `1.2.3`)
	testDecodeError(t, `'a'`)
	testDecodeError(t, `true, false`)
	testDecodeError(t, `{"a:"b"}`)
	testDecodeError(t, `{a":"b"}`)
	testDecodeError(t, `{"a":"b":"c"}`)
}

func TestEncodeUnsupported(t *testing.T) {
	fn := &tengo.CompiledFunction{}
	for _, o := range []tengo.Object{
		fn,
		&tengo.Map{Value: map[string]tengo.Object{"a": fn, "b": &tengo.Int{Value: 1}}},
		&tengo.Array{Value: []tengo.Object{fn, &tengo.Int{Value: 1}}},
		&tengo.Error{Value: &tengo.String{Value: "x"}},
	} {
		b, err := json.Encode(o)
		require.Error(t, err, string(b))
	}
}

func TestDecodeFloatOverflow(t *testing.T) {
	_, err := json.Decode([]byte("1e400"))
	require.Error(t, err)
	_, err = json.Decode([]byte("[-1e400]"))
	require.Error(t, err)

	// underflow rounds to zero, as encoding/json does
	o, err := json.Decode([]byte("1e-400"))
	require.NoError(t, err)
	require.Equal(t, 0.0, o.(*tengo.Float).Value)
}

func TestDecodeIntOverflow(t *testing.T) {
	// ParseInt reports overflow and also returns the extreme int64.
	// Dropping that error made 2^63 decode as MaxInt64.
	_, err := json.Decode([]byte("9223372036854775808"))
	require.Error(t, err)
	_, err = json.Decode([]byte("-9223372036854775809"))
	require.Error(t, err)

	o, err := json.Decode([]byte("9223372036854775807"))
	require.NoError(t, err)
	require.Equal(t, int64(9223372036854775807), o.(*tengo.Int).Value)

	o, err = json.Decode([]byte("-9223372036854775808"))
	require.NoError(t, err)
	require.Equal(t, int64(-9223372036854775808), o.(*tengo.Int).Value)

	o, err = json.Decode([]byte("1.5"))
	require.NoError(t, err)
	require.Equal(t, 1.5, o.(*tengo.Float).Value)
}

func testDecodeError(t *testing.T, input string) {
	_, err := json.Decode([]byte(input))
	require.Error(t, err)
}

func testJSONEncodeDecode(t *testing.T, v interface{}) {
	o, err := tengo.FromInterface(v)
	require.NoError(t, err)

	b, err := json.Encode(o)
	require.NoError(t, err)

	a, err := json.Decode(b)
	require.NoError(t, err, string(b))

	vj, err := gojson.Marshal(v)
	require.NoError(t, err)

	aj, err := gojson.Marshal(tengo.ToInterface(a))
	require.NoError(t, err)

	require.Equal(t, vj, aj)
}

func TestDecodeNumbers(t *testing.T) {
	for input, expected := range map[string]tengo.Object{
		"-0":     &tengo.Int{Value: 0},
		"-7":     &tengo.Int{Value: -7},
		"-0.5":   &tengo.Float{Value: -0.5},
		"1e3":    &tengo.Float{Value: 1000},
		"1E3":    &tengo.Float{Value: 1000},
		"1e+3":   &tengo.Float{Value: 1000},
		"1e-3":   &tengo.Float{Value: 0.001},
		"0e0":    &tengo.Float{Value: 0},
		"-1.5e2": &tengo.Float{Value: -150},
		"0.5E-1": &tengo.Float{Value: 0.05},
		"12e10":  &tengo.Float{Value: 12e10},
		"[1e2,-3]": &tengo.Array{Value: []tengo.Object{
			&tengo.Float{Value: 100}, &tengo.Int{Value: -3}}},
		`{"a":2.5e1}`: &tengo.Map{Value: map[string]tengo.Object{
			"a": &tengo.Float{Value: 25}}},
	} {
		o, err := json.Decode([]byte(input))
		require.NoError(t, err, input)
		require.Equal(t, expected, o, input)
	}
}

func TestDecodeSyntaxError(t *testing.T) {
	for _, c := range []struct {
		input  string
		msg    string
		offset int64
	}{
		{"1ex", "invalid character 'x' in exponent of numeric literal", 3},
		{"1e+x", "invalid character 'x' in exponent of numeric literal", 4},
		{"1.5e-]", "invalid character ']' in exponent of numeric literal", 6},
		{"-x", "invalid character 'x' in numeric literal", 2},
		{"-.5", "invalid character '.' in numeric literal", 2},
		{"1.x", "invalid character 'x' after decimal point in numeric literal", 3},
		{"1.e5", "invalid character 'e' after decimal point in numeric literal", 3},
		{"1.5x", "invalid character 'x' after top-level value", 4},
		{"1e5x", "invalid character 'x' after top-level value", 4},
		{"[1e5}", "invalid character '}' after array element", 5},
		{"{", "unexpected end of JSON input", 1},
	} {
		_, err := json.Decode([]byte(c.input))
		require.Error(t, err, c.input)
		se, ok := err.(*json.SyntaxError)
		require.True(t, ok, "%s: %T", c.input, err)
		require.Equal(t, c.msg, se.Error(), c.input)
		require.Equal(t, c.offset, se.Offset, c.input)
	}
}

func TestDecodeStringEscapes(t *testing.T) {
	for input, expected := range map[string]string{
		`"\b\f\n\r\t"`:     "\b\f\n\r\t",
		`"\/\\\""`:         `/\"`,
		`"a\tb"`:           "a\tb",
		`"\n\u65e5\u672c"`: "\n日本",
		`"\n日本"`:           "\n日本",
		`"\u0041\u00e9"`:   "Aé",
		`"\u00E9\u00c9"`:   "éÉ",
		`"\ud83d\ude00"`:   "\U0001F600",
		`"\ud83d"`:         "\uFFFD",
		`"\ude00"`:         "\uFFFD",
		`"\ud83d\u0041"`:   "\uFFFDA",
		`"\ud83dx"`:        "\uFFFDx",
		"\"\xff\"":         "\uFFFD",
		"\"\\n\xff\"":      "\n\uFFFD",
	} {
		o, err := json.Decode([]byte(input))
		require.NoError(t, err, input)
		require.Equal(t, &tengo.String{Value: expected}, o, input)
	}

	o, err := json.Decode([]byte(`{"a\nb":"\u0063"}`))
	require.NoError(t, err)
	require.Equal(t, &tengo.Map{Value: map[string]tengo.Object{
		"a\nb": &tengo.String{Value: "c"}}}, o)
}

func TestEncodeTypes(t *testing.T) {
	inf := &tengo.Float{Value: math.Inf(1)}
	for _, c := range []struct {
		o        tengo.Object
		expected string
	}{
		{tengo.TrueValue, "true"},
		{tengo.FalseValue, "false"},
		{tengo.UndefinedValue, "null"},
		{&tengo.Int{Value: -42}, "-42"},
		{&tengo.Char{Value: 'A'}, "65"},
		{&tengo.Bytes{Value: []byte("foo")}, `"Zm9v"`},
		{&tengo.Bytes{}, `""`},
		{&tengo.Float{Value: 0}, "0"},
		{&tengo.Float{Value: -1.5}, "-1.5"},
		{&tengo.Float{Value: 123456789}, "123456789"},
		{&tengo.Float{Value: 1e20}, "100000000000000000000"},
		{&tengo.Float{Value: 1e21}, "1e+21"},
		{&tengo.Float{Value: 1e-6}, "0.000001"},
		{&tengo.Float{Value: 1e-7}, "1e-7"},
		{&tengo.Float{Value: -1.5e-10}, "-1.5e-10"},
		{&tengo.String{Value: "a\"b\\c\n\r\t\x01\x1f"},
			`"a\"b\\c\n\r\t\u0001\u001f"`},
		{&tengo.String{Value: "<>&\u2028é"}, "\"<>&\u2028é\""},
		{&tengo.Time{Value: time.Date(2009, 11, 10, 23, 0, 0, 5, time.UTC)},
			`"2009-11-10T23:00:00.000000005Z"`},
		{&tengo.Array{}, "[]"},
		{&tengo.ImmutableArray{}, "[]"},
		{&tengo.Array{Value: []tengo.Object{
			&tengo.Int{Value: 1}, tengo.UndefinedValue,
			&tengo.ImmutableArray{Value: []tengo.Object{tengo.TrueValue}},
		}}, "[1,null,[true]]"},
		{&tengo.Map{}, "{}"},
		{&tengo.ImmutableMap{}, "{}"},
		{&tengo.ImmutableMap{Value: map[string]tengo.Object{
			"a\"b": &tengo.Map{Value: map[string]tengo.Object{
				"c": &tengo.Int{Value: 1}}},
		}}, `{"a\"b":{"c":1}}`},
	} {
		b, err := json.Encode(c.o)
		require.NoError(t, err, c.expected)
		require.Equal(t, c.expected, string(b))
	}

	for _, o := range []tengo.Object{
		&tengo.Map{Value: map[string]tengo.Object{
			"a": &tengo.Int{Value: 1}, "b": &tengo.Int{Value: 2}}},
		&tengo.ImmutableMap{Value: map[string]tengo.Object{
			"a": &tengo.Int{Value: 1}, "b": &tengo.Int{Value: 2}}},
	} {
		b, err := json.Encode(o)
		require.NoError(t, err)
		s := string(b)
		require.True(t, s == `{"a":1,"b":2}` || s == `{"b":2,"a":1}`, s)
	}

	for _, o := range []tengo.Object{
		inf,
		&tengo.Float{Value: math.Inf(-1)},
		&tengo.Float{Value: math.NaN()},
		&tengo.Time{Value: time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)},
		&tengo.Array{Value: []tengo.Object{&tengo.Int{Value: 1}, inf}},
		&tengo.ImmutableArray{Value: []tengo.Object{inf}},
		&tengo.Map{Value: map[string]tengo.Object{"a": inf}},
		&tengo.ImmutableMap{Value: map[string]tengo.Object{"a": inf}},
	} {
		b, err := json.Encode(o)
		require.Error(t, err, o.String())
		require.Nil(t, b, o.String())
	}
}

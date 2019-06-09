package json_test

import (
	gojson "encoding/json"
	"testing"

	"github.com/d5/tengo"
	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/stdlib/json"
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

	testJSONEncodeDecode(t, true)
	testJSONEncodeDecode(t, false)

	testJSONEncodeDecode(t, ARR{})
	testJSONEncodeDecode(t, ARR{0})
	testJSONEncodeDecode(t, ARR{false})
	testJSONEncodeDecode(t, ARR{1, 2, 3, "four", false})
	testJSONEncodeDecode(t, ARR{1, 2, 3, "four", false, MAP{"a": 0, "b": "bee", "bool": true}})

	testJSONEncodeDecode(t, MAP{})
	testJSONEncodeDecode(t, MAP{"a": 0})
	testJSONEncodeDecode(t, MAP{"a": 0, "b": "bee"})
	testJSONEncodeDecode(t, MAP{"a": 0, "b": "bee", "bool": true})

	testJSONEncodeDecode(t, MAP{"a": 0, "b": "bee", "arr": ARR{1, 2, 3, "four"}})
	testJSONEncodeDecode(t, MAP{"a": 0, "b": "bee", "arr": ARR{1, 2, 3, MAP{"a": false, "b": 109.4}}})
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

func testDecodeError(t *testing.T, input string) {
	_, err := json.Decode([]byte(input))
	assert.Error(t, err)
}

func testJSONEncodeDecode(t *testing.T, v interface{}) bool {
	o, err := tengo.FromInterface(v)
	if !assert.NoError(t, err) {
		return false
	}

	b, err := json.Encode(o)
	if !assert.NoError(t, err) {
		return false
	}

	a, err := json.Decode(b)
	if !assert.NoError(t, err, string(b)) {
		return false
	}

	vj, err := gojson.Marshal(v)
	if !assert.NoError(t, err) {
		return false
	}

	aj, err := gojson.Marshal(tengo.ToInterface(a))
	if !assert.NoError(t, err) {
		return false
	}

	return assert.Equal(t, vj, aj)
}

package stdlib_test

import (
	"math"
	"strings"
	"testing"

	"github.com/d5/tengo/v2"
)

func TestJSON(t *testing.T) {
	module(t, "json").call("encode", 5).
		expect([]byte("5"))
	module(t, "json").call("encode", "foobar").
		expect([]byte(`"foobar"`))
	module(t, "json").call("encode", MAP{"foo": 5}).
		expect([]byte("{\"foo\":5}"))
	module(t, "json").call("encode", IMAP{"foo": 5}).
		expect([]byte("{\"foo\":5}"))
	module(t, "json").call("encode", ARR{1, 2, 3}).
		expect([]byte("[1,2,3]"))
	module(t, "json").call("encode", IARR{1, 2, 3}).
		expect([]byte("[1,2,3]"))
	module(t, "json").call("encode", MAP{"foo": "bar"}).
		expect([]byte("{\"foo\":\"bar\"}"))
	module(t, "json").call("encode", MAP{"foo": 1.8}).
		expect([]byte("{\"foo\":1.8}"))
	module(t, "json").call("encode", MAP{"foo": true}).
		expect([]byte("{\"foo\":true}"))
	module(t, "json").call("encode", MAP{"foo": '8'}).
		expect([]byte("{\"foo\":56}"))
	module(t, "json").call("encode", MAP{"foo": []byte("foo")}).
		expect([]byte("{\"foo\":\"Zm9v\"}")) // json encoding returns []byte as base64 encoded string
	module(t, "json").
		call("encode", MAP{"foo": ARR{"bar", 1, 1.8, '8', true}}).
		expect([]byte("{\"foo\":[\"bar\",1,1.8,56,true]}"))
	module(t, "json").
		call("encode", MAP{"foo": IARR{"bar", 1, 1.8, '8', true}}).
		expect([]byte("{\"foo\":[\"bar\",1,1.8,56,true]}"))
	module(t, "json").
		call("encode", MAP{"foo": ARR{ARR{"bar", 1}, ARR{"bar", 1}}}).
		expect([]byte("{\"foo\":[[\"bar\",1],[\"bar\",1]]}"))
	module(t, "json").
		call("encode", MAP{"foo": MAP{"string": "bar"}}).
		expect([]byte("{\"foo\":{\"string\":\"bar\"}}"))
	module(t, "json").
		call("encode", MAP{"foo": IMAP{"string": "bar"}}).
		expect([]byte("{\"foo\":{\"string\":\"bar\"}}"))
	module(t, "json").
		call("encode", MAP{"foo": MAP{"map1": MAP{"string": "bar"}}}).
		expect([]byte("{\"foo\":{\"map1\":{\"string\":\"bar\"}}}"))
	module(t, "json").
		call("encode", ARR{ARR{"bar", 1}, ARR{"bar", 1}}).
		expect([]byte("[[\"bar\",1],[\"bar\",1]]"))

	module(t, "json").call("decode", `5`).
		expect(5)
	module(t, "json").call("decode", `"foo"`).
		expect("foo")
	module(t, "json").call("decode", `[1,2,3,"bar"]`).
		expect(ARR{1, 2, 3, "bar"})
	module(t, "json").call("decode", `{"foo":5}`).
		expect(MAP{"foo": 5})
	module(t, "json").call("decode", `{"foo":2.5}`).
		expect(MAP{"foo": 2.5})
	module(t, "json").call("decode", `{"foo":true}`).
		expect(MAP{"foo": true})
	module(t, "json").call("decode", `{"foo":"bar"}`).
		expect(MAP{"foo": "bar"})
	module(t, "json").call("decode", `{"foo":[1,2,3,"bar"]}`).
		expect(MAP{"foo": ARR{1, 2, 3, "bar"}})

	module(t, "json").
		call("indent", []byte("{\"foo\":[\"bar\",1,1.8,56,true]}"), "", "  ").
		expect([]byte(`{
  "foo": [
    "bar",
    1,
    1.8,
    56,
    true
  ]
}`))

	module(t, "json").
		call("html_escape", []byte(
			`{"M":"<html>foo &`+"\xe2\x80\xa8 \xe2\x80\xa9"+`</html>"}`)).
		expect([]byte(
			`{"M":"\u003chtml\u003efoo \u0026\u2028 \u2029\u003c/html\u003e"}`))
}

func TestJSONInputTypes(t *testing.T) {
	module(t, "json").call("decode", []byte(`{"foo":[1,"bar"]}`)).
		expect(MAP{"foo": ARR{1, "bar"}})
	module(t, "json").call("indent", `{"foo":1}`, ">", "\t").
		expect([]byte("{\n>\t\"foo\": 1\n>}"))
	module(t, "json").call("html_escape", `"<a>&"`).
		expect([]byte(`"\u003ca\u003e\u0026"`))
}

func TestJSONErrors(t *testing.T) {
	wrongNum := tengo.ErrWrongNumArguments
	expectArgErrors(t, "json", []argErrorCase{
		{"decode", ARR{}, wrongNum},
		{"decode", ARR{"1", "1"}, wrongNum},
		{"decode", ARR{1}, invalidArg("first", "bytes/string", "int")},
		{"encode", ARR{}, wrongNum},
		{"encode", ARR{1, 1}, wrongNum},
		{"indent", ARR{"{}", ""}, wrongNum},
		{"indent", ARR{1, "", ""}, invalidArg("first", "bytes/string", "int")},
		{"indent", ARR{"{}", tengo.UndefinedValue, ""},
			invalidArg("prefix", "string(compatible)", "undefined")},
		{"indent", ARR{"{}", "", tengo.UndefinedValue},
			invalidArg("indent", "string(compatible)", "undefined")},
		{"html_escape", ARR{}, wrongNum},
		{"html_escape", ARR{"a", "a"}, wrongNum},
		{"html_escape", ARR{1}, invalidArg("first", "bytes/string", "int")},
	})

	module(t, "json").call("decode", `{"foo":`).
		expect(&tengo.Error{Value: &tengo.String{
			Value: "unexpected end of JSON input"}})
	module(t, "json").call("decode", []byte(`[1,]`)).
		expect(&tengo.Error{Value: &tengo.String{
			Value: "invalid character ']' looking for beginning of value"}})
	module(t, "json").call("indent", `{"foo"`, "", " ").expectErrorObject()
	module(t, "json").call("indent", []byte(`[1,]`), "", " ").
		expectErrorObject()
	module(t, "json").call("encode", math.Inf(1)).
		expect(&tengo.Error{Value: &tengo.String{
			Value: "unsupported float value"}})
}

func TestJSONLimits(t *testing.T) {
	curMaxBytesLen := tengo.MaxBytesLen
	defer func() { tengo.MaxBytesLen = curMaxBytesLen }()
	tengo.MaxBytesLen = 12

	module(t, "json").call("encode", "0123456789").
		expect([]byte(`"0123456789"`))
	module(t, "json").call("encode", "01234567890").expectError()

	module(t, "json").call("html_escape", "<>").
		expect([]byte(`\u003c\u003e`))
	module(t, "json").call("html_escape", "<>&").expectError()
	module(t, "json").call("html_escape", "\xe2\x80\xa8\xe2\x80\xa8").
		expect([]byte(`\u2028\u2028`))
	module(t, "json").call("html_escape", "\xe2\x80\xa8\xe2\x80\xa9a").
		expectError()

	module(t, "json").call("indent", "[1]", "", " ").
		expect([]byte("[\n 1\n]"))
	module(t, "json").call("indent", "[1,2,3]", "", " ").expectError()

	// a 2 KB input would indent to tens of gigabytes
	tengo.MaxBytesLen = 1 << 20
	deep := strings.Repeat("[", 1000) + strings.Repeat("]", 1000)
	module(t, "json").call("indent", deep, "", strings.Repeat(" ", 4096)).
		expectError()
}

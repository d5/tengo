package stdlib_test

import "testing"

func TestJSON(t *testing.T) {
	module(t, "json").call("encode", mockRuntime{}, 5).expect([]byte("5"))
	module(t, "json").call("encode", mockRuntime{}, "foobar").expect([]byte(`"foobar"`))
	module(t, "json").call("encode", mockRuntime{}, MAP{"foo": 5}).expect([]byte("{\"foo\":5}"))
	module(t, "json").call("encode", mockRuntime{}, IMAP{"foo": 5}).expect([]byte("{\"foo\":5}"))
	module(t, "json").call("encode", mockRuntime{}, ARR{1, 2, 3}).expect([]byte("[1,2,3]"))
	module(t, "json").call("encode", mockRuntime{}, IARR{1, 2, 3}).expect([]byte("[1,2,3]"))
	module(t, "json").call("encode", mockRuntime{}, MAP{"foo": "bar"}).expect([]byte("{\"foo\":\"bar\"}"))
	module(t, "json").call("encode", mockRuntime{}, MAP{"foo": 1.8}).expect([]byte("{\"foo\":1.8}"))
	module(t, "json").call("encode", mockRuntime{}, MAP{"foo": true}).expect([]byte("{\"foo\":true}"))
	module(t, "json").call("encode", mockRuntime{}, MAP{"foo": '8'}).expect([]byte("{\"foo\":56}"))
	module(t, "json").call("encode", mockRuntime{}, MAP{"foo": []byte("foo")}).expect([]byte("{\"foo\":\"Zm9v\"}")) // json encoding returns []byte as base64 encoded string
	module(t, "json").call("encode", mockRuntime{}, MAP{"foo": ARR{"bar", 1, 1.8, '8', true}}).expect([]byte("{\"foo\":[\"bar\",1,1.8,56,true]}"))
	module(t, "json").call("encode", mockRuntime{}, MAP{"foo": IARR{"bar", 1, 1.8, '8', true}}).expect([]byte("{\"foo\":[\"bar\",1,1.8,56,true]}"))
	module(t, "json").call("encode", mockRuntime{}, MAP{"foo": ARR{ARR{"bar", 1}, ARR{"bar", 1}}}).expect([]byte("{\"foo\":[[\"bar\",1],[\"bar\",1]]}"))
	module(t, "json").call("encode", mockRuntime{}, MAP{"foo": MAP{"string": "bar"}}).expect([]byte("{\"foo\":{\"string\":\"bar\"}}"))
	module(t, "json").call("encode", mockRuntime{}, MAP{"foo": IMAP{"string": "bar"}}).expect([]byte("{\"foo\":{\"string\":\"bar\"}}"))
	module(t, "json").call("encode", mockRuntime{}, MAP{"foo": MAP{"map1": MAP{"string": "bar"}}}).expect([]byte("{\"foo\":{\"map1\":{\"string\":\"bar\"}}}"))
	module(t, "json").call("encode", mockRuntime{}, ARR{ARR{"bar", 1}, ARR{"bar", 1}}).expect([]byte("[[\"bar\",1],[\"bar\",1]]"))

	module(t, "json").call("decode", mockRuntime{}, `5`).expect(5.0)
	module(t, "json").call("decode", mockRuntime{}, `"foo"`).expect("foo")
	module(t, "json").call("decode", mockRuntime{}, `[1,2,3,"bar"]`).expect(ARR{1.0, 2.0, 3.0, "bar"})
	module(t, "json").call("decode", mockRuntime{}, `{"foo":5}`).expect(MAP{"foo": 5.0})
	module(t, "json").call("decode", mockRuntime{}, `{"foo":2.5}`).expect(MAP{"foo": 2.5})
	module(t, "json").call("decode", mockRuntime{}, `{"foo":true}`).expect(MAP{"foo": true})
	module(t, "json").call("decode", mockRuntime{}, `{"foo":"bar"}`).expect(MAP{"foo": "bar"})
	module(t, "json").call("decode", mockRuntime{}, `{"foo":[1,2,3,"bar"]}`).expect(MAP{"foo": ARR{1.0, 2.0, 3.0, "bar"}})

	module(t, "json").call("indent", mockRuntime{}, []byte("{\"foo\":[\"bar\",1,1.8,56,true]}"), "", "  ").expect([]byte(`{
  "foo": [
    "bar",
    1,
    1.8,
    56,
    true
  ]
}`))

	module(t, "json").call("html_escape", mockRuntime{},
		[]byte(`{"M":"<html>foo &`+"\xe2\x80\xa8 \xe2\x80\xa9"+`</html>"}`)).expect([]byte(`{"M":"\u003chtml\u003efoo \u0026\u2028 \u2029\u003c/html\u003e"}`))
}

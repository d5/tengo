package stdlib_test

import "testing"

func TestJSON(t *testing.T) {
	module(t, "json").call("stringify", 5).expect("5")
	module(t, "json").call("stringify", "foobar").expect(`"foobar"`)
	module(t, "json").call("stringify", MAP{"foo": 5}).expect("{\"foo\":5}")
	module(t, "json").call("stringify", IMAP{"foo": 5}).expect("{\"foo\":5}")
	module(t, "json").call("stringify", ARR{1, 2, 3}).expect("[1,2,3]")
	module(t, "json").call("stringify", IARR{1, 2, 3}).expect("[1,2,3]")
	module(t, "json").call("stringify", MAP{"foo": "bar"}).expect("{\"foo\":\"bar\"}")
	module(t, "json").call("stringify", MAP{"foo": 1.8}).expect("{\"foo\":1.8}")
	module(t, "json").call("stringify", MAP{"foo": true}).expect("{\"foo\":true}")
	module(t, "json").call("stringify", MAP{"foo": '8'}).expect("{\"foo\":56}")
	module(t, "json").call("stringify", MAP{"foo": []byte("foo")}).expect("{\"foo\":\"Zm9v\"}") // json encoding returns []byte as base64 encoded string
	module(t, "json").call("stringify", MAP{"foo": ARR{"bar", 1, 1.8, '8', true}}).expect("{\"foo\":[\"bar\",1,1.8,56,true]}")
	module(t, "json").call("stringify", MAP{"foo": IARR{"bar", 1, 1.8, '8', true}}).expect("{\"foo\":[\"bar\",1,1.8,56,true]}")
	module(t, "json").call("stringify", MAP{"foo": ARR{ARR{"bar", 1}, ARR{"bar", 1}}}).expect("{\"foo\":[[\"bar\",1],[\"bar\",1]]}")
	module(t, "json").call("stringify", MAP{"foo": MAP{"string": "bar", "int": 1, "float": 1.8, "char": '8', "bool": true}}).expect("{\"foo\":{\"bool\":true,\"char\":56,\"float\":1.8,\"int\":1,\"string\":\"bar\"}}")
	module(t, "json").call("stringify", MAP{"foo": IMAP{"string": "bar", "int": 1, "float": 1.8, "char": '8', "bool": true}}).expect("{\"foo\":{\"bool\":true,\"char\":56,\"float\":1.8,\"int\":1,\"string\":\"bar\"}}")
	module(t, "json").call("stringify", MAP{"foo": MAP{"map1": MAP{"string": "bar"}, "map2": MAP{"int": "1"}}}).expect("{\"foo\":{\"map1\":{\"string\":\"bar\"},\"map2\":{\"int\":\"1\"}}}")
	module(t, "json").call("stringify", ARR{ARR{"bar", 1}, ARR{"bar", 1}}).expect("[[\"bar\",1],[\"bar\",1]]")

	module(t, "json").call("parse", `5`).expect(5.0)
	module(t, "json").call("parse", `"foo"`).expect("foo")
	module(t, "json").call("parse", `[1,2,3,"bar"]`).expect(ARR{1.0, 2.0, 3.0, "bar"})
	module(t, "json").call("parse", `{"foo":5}`).expect(MAP{"foo": 5.0})
	module(t, "json").call("parse", `{"foo":2.5}`).expect(MAP{"foo": 2.5})
	module(t, "json").call("parse", `{"foo":true}`).expect(MAP{"foo": true})
	module(t, "json").call("parse", `{"foo":"bar"}`).expect(MAP{"foo": "bar"})
	module(t, "json").call("parse", `{"foo":[1,2,3,"bar"]}`).expect(MAP{"foo": ARR{1.0, 2.0, 3.0, "bar"}})
}

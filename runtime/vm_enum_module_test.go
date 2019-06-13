package runtime_test

import (
	"testing"

	"github.com/d5/tengo"
)

func TestEnumModule(t *testing.T) {
	testEnumModule(t, `out = enum.key(0, 20)`, 0)
	testEnumModule(t, `out = enum.key(10, 20)`, 10)
	testEnumModule(t, `out = enum.value(0, 0)`, 0)
	testEnumModule(t, `out = enum.value(10, 20)`, 20)

	testEnumModule(t, `out = enum.all([], enum.value)`, true)
	testEnumModule(t, `out = enum.all([1], enum.value)`, true)
	testEnumModule(t, `out = enum.all([true, 1], enum.value)`, true)
	testEnumModule(t, `out = enum.all([true, 0], enum.value)`, false)
	testEnumModule(t, `out = enum.all([true, 0, 1], enum.value)`, false)
	testEnumModule(t, `out = enum.all(immutable([true, 0, 1]), enum.value)`, false) // immutable-array
	testEnumModule(t, `out = enum.all({}, enum.value)`, true)
	testEnumModule(t, `out = enum.all({a:1}, enum.value)`, true)
	testEnumModule(t, `out = enum.all({a:true, b:1}, enum.value)`, true)
	testEnumModule(t, `out = enum.all(immutable({a:true, b:1}), enum.value)`, true) // immutable-map
	testEnumModule(t, `out = enum.all({a:true, b:0}, enum.value)`, false)
	testEnumModule(t, `out = enum.all({a:true, b:0, c:1}, enum.value)`, false)
	testEnumModule(t, `out = enum.all("123", enum.value)`, true)
	testEnumModuleError(t, `enum.all(0, enum.value)`, "invalid type for argument 'x'")

	testEnumModule(t, `out = enum.any([], enum.value)`, false)
	testEnumModule(t, `out = enum.any([1], enum.value)`, true)
	testEnumModule(t, `out = enum.any([true, 1], enum.value)`, true)
	testEnumModule(t, `out = enum.any([true, 0], enum.value)`, true)
	testEnumModule(t, `out = enum.any([true, 0, 1], enum.value)`, true)
	testEnumModule(t, `out = enum.any(immutable([true, 0, 1]), enum.value)`, true) // immutable-array
	testEnumModule(t, `out = enum.any([false], enum.value)`, false)
	testEnumModule(t, `out = enum.any([false, 0], enum.value)`, false)
	testEnumModule(t, `out = enum.any({}, enum.value)`, false)
	testEnumModule(t, `out = enum.any({a:1}, enum.value)`, true)
	testEnumModule(t, `out = enum.any({a:true, b:1}, enum.value)`, true)
	testEnumModule(t, `out = enum.any({a:true, b:0}, enum.value)`, true)
	testEnumModule(t, `out = enum.any({a:true, b:0, c:1}, enum.value)`, true)
	testEnumModule(t, `out = enum.any(immutable({a:true, b:0, c:1}), enum.value)`, true) // immutable-map
	testEnumModule(t, `out = enum.any({a:false}, enum.value)`, false)
	testEnumModule(t, `out = enum.any({a:false, b:0}, enum.value)`, false)
	testEnumModule(t, `out = enum.any("0", enum.value)`, true)
	testEnumModule(t, `out = enum.any("123", enum.value)`, true)
	testEnumModuleError(t, `enum.any(0, enum.value)`, "invalid type for argument 'x'")

	testEnumModule(t, `out = enum.chunk([], 1)`, ARR{})
	testEnumModule(t, `out = enum.chunk([1], 1)`, ARR{ARR{1}})
	testEnumModule(t, `out = enum.chunk([1,2,3], 1)`, ARR{ARR{1}, ARR{2}, ARR{3}})
	testEnumModule(t, `out = enum.chunk([1,2,3], 2)`, ARR{ARR{1, 2}, ARR{3}})
	testEnumModule(t, `out = enum.chunk([1,2,3], 3)`, ARR{ARR{1, 2, 3}})
	testEnumModule(t, `out = enum.chunk([1,2,3], 4)`, ARR{ARR{1, 2, 3}})
	testEnumModule(t, `out = enum.chunk([1,2,3,4], 3)`, ARR{ARR{1, 2, 3}, ARR{4}})
	testEnumModule(t, `out = enum.chunk([], 0)`, tengo.UndefinedValue)      // size=0: undefined
	testEnumModule(t, `out = enum.chunk([1], 0)`, tengo.UndefinedValue)     // size=0: undefined
	testEnumModule(t, `out = enum.chunk([1,2,3], 0)`, tengo.UndefinedValue) // size=0: undefined
	testEnumModuleError(t, `enum.chunk({a:1,b:2,c:3}, 1)`, "invalid type for argument 'x'")
	testEnumModuleError(t, `enum.chunk(0, 1)`, "invalid type for argument 'x'")
	testEnumModuleError(t, `enum.chunk("123", 1)`, "invalid type for argument 'x'")

	testEnumModule(t, `out = enum.at([], 0)`, tengo.UndefinedValue)
	testEnumModule(t, `out = enum.at([], 1)`, tengo.UndefinedValue)
	testEnumModule(t, `out = enum.at([], -1)`, tengo.UndefinedValue)
	testEnumModule(t, `out = enum.at(["one"], 0)`, "one")
	testEnumModule(t, `out = enum.at(["one"], 1)`, tengo.UndefinedValue)
	testEnumModule(t, `out = enum.at(["one"], -1)`, tengo.UndefinedValue)
	testEnumModule(t, `out = enum.at(["one","two","three"], 0)`, "one")
	testEnumModule(t, `out = enum.at(["one","two","three"], 1)`, "two")
	testEnumModule(t, `out = enum.at(["one","two","three"], 2)`, "three")
	testEnumModule(t, `out = enum.at(["one","two","three"], -1)`, tengo.UndefinedValue)
	testEnumModule(t, `out = enum.at(["one","two","three"], 3)`, tengo.UndefinedValue)
	testEnumModuleError(t, `enum.at(["one","two","three"], "1")`, "invalid index type")
	testEnumModule(t, `out = enum.at({}, "a")`, tengo.UndefinedValue)
	testEnumModule(t, `out = enum.at({a:"one"}, "a")`, "one")
	testEnumModule(t, `out = enum.at({a:"one"}, "b")`, tengo.UndefinedValue)
	testEnumModule(t, `out = enum.at({a:"one",b:"two",c:"three"}, "a")`, "one")
	testEnumModule(t, `out = enum.at({a:"one",b:"two",c:"three"}, "b")`, "two")
	testEnumModule(t, `out = enum.at({a:"one",b:"two",c:"three"}, "c")`, "three")
	testEnumModule(t, `out = enum.at({a:"one",b:"two",c:"three"}, "d")`, tengo.UndefinedValue)
	testEnumModule(t, `out = enum.at({a:"one",b:"two",c:"three"}, 'a')`, "one")
	testEnumModule(t, `out = enum.at("abc", 1)`, 'b')
	testEnumModuleError(t, `enum.at(0, 1)`, "not indexable")

	testEnumModule(t, `out=0; enum.each([],func(k,v){out+=v})`, 0)
	testEnumModule(t, `out=0; enum.each([1,2,3],func(k,v){out+=v})`, 6)
	testEnumModule(t, `out=0; enum.each([1,2,3],func(k,v){out+=k})`, 3)
	testEnumModule(t, `out=0; enum.each({a:1,b:2,c:3},func(k,v){out+=v})`, 6)
	testEnumModule(t, `out=""; enum.each({a:1,b:2,c:3},func(k,v){out+=k}); out=len(out)`, 3)
	testEnumModuleError(t, `enum.each(5,func(k,v){})`, "invalid type for argument 'x'")
	testEnumModule(t, `out=0; enum.each("123",func(k,v){out+=v})`, byte(150)) // char additions

	testEnumModule(t, `out = enum.filter([], enum.value)`, ARR{})
	testEnumModule(t, `out = enum.filter([false,1,2], enum.value)`, ARR{1, 2})
	testEnumModule(t, `out = enum.filter([false,1,0,2], enum.value)`, ARR{1, 2})
	testEnumModule(t, `out = enum.filter({}, enum.value)`, ARR{})
	testEnumModuleError(t, `enum.filter(0, enum.value)`, "invalid type for argument 'x'")
	testEnumModule(t, `out = enum.filter("123", enum.value)`, ARR{'1', '2', '3'})

	testEnumModule(t, `out = enum.find([], enum.value)`, tengo.UndefinedValue)
	testEnumModule(t, `out = enum.find([0], enum.value)`, tengo.UndefinedValue)
	testEnumModule(t, `out = enum.find([1], enum.value)`, 1)
	testEnumModule(t, `out = enum.find([false,0,undefined,1], enum.value)`, 1)
	testEnumModule(t, `out = enum.find([1,2,3], enum.value)`, 1)
	testEnumModule(t, `out = enum.find({}, enum.value)`, tengo.UndefinedValue)
	testEnumModule(t, `out = enum.find({a:0}, enum.value)`, tengo.UndefinedValue)
	testEnumModule(t, `out = enum.find({a:1}, enum.value)`, 1)
	testEnumModule(t, `out = enum.find({a:false,b:0,c:undefined,d:1}, enum.value)`, 1)
	//testEnumModule(t, `out = enum.find({a:1,b:2,c:3}, enum.value)`, 1)
	testEnumModuleError(t, `enum.find(0, enum.value)`, "invalid type for argument 'x'")
	testEnumModule(t, `out = enum.find("123", enum.value)`, '1')

	testEnumModule(t, `out = enum.find_key([], enum.value)`, tengo.UndefinedValue)
	testEnumModule(t, `out = enum.find_key([0], enum.value)`, tengo.UndefinedValue)
	testEnumModule(t, `out = enum.find_key([1], enum.value)`, 0)
	testEnumModule(t, `out = enum.find_key([false,0,undefined,1], enum.value)`, 3)
	testEnumModule(t, `out = enum.find_key([1,2,3], enum.value)`, 0)
	testEnumModule(t, `out = enum.find_key({}, enum.value)`, tengo.UndefinedValue)
	testEnumModule(t, `out = enum.find_key({a:0}, enum.value)`, tengo.UndefinedValue)
	testEnumModule(t, `out = enum.find_key({a:1}, enum.value)`, "a")
	testEnumModule(t, `out = enum.find_key({a:false,b:0,c:undefined,d:1}, enum.value)`, "d")
	//testEnumModule(t, `out = enum.find_key({a:1,b:2,c:3}, enum.value)`, "a")
	testEnumModuleError(t, `enum.find_key(0, enum.value)`, "invalid type for argument 'x'")
	testEnumModule(t, `out = enum.find_key("123", enum.value)`, 0)

	testEnumModule(t, `out = enum.map([], enum.value)`, ARR{})
	testEnumModule(t, `out = enum.map([1,2,3], enum.value)`, ARR{1, 2, 3})
	testEnumModule(t, `out = enum.map([1,2,3], enum.key)`, ARR{0, 1, 2})
	testEnumModule(t, `out = enum.map([1,2,3], func(k,v) { return v*2 })`, ARR{2, 4, 6})
	testEnumModule(t, `out = enum.map({}, enum.value)`, ARR{})
	testEnumModule(t, `out = enum.map({a:1}, func(k,v) { return v*2 })`, ARR{2})
	testEnumModuleError(t, `enum.map(0, enum.value)`, "invalid type for argument 'x'")
	//testEnumModule(t, `out = enum.map("123", enum.value)`, ARR{'1', '2', '3'})
}

func testEnumModule(t *testing.T, input string, expected interface{}) {
	expect(t, `enum := import("enum"); `+input, Opts().Stdlib(), expected)
}

func testEnumModuleError(t *testing.T, input, expected string) {
	expectError(t, `enum := import("enum"); `+input, Opts().Stdlib(), expected)
}

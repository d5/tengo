package runtime_test

import (
	"testing"

	"github.com/d5/tengo/objects"
)

type StringArrayIterator struct {
	objects.ObjectImpl
	strArr *StringArray
	idx    int
}

func (i *StringArrayIterator) TypeName() string {
	return "string-array-iterator"
}

func (i *StringArrayIterator) String() string {
	return ""
}

func (i *StringArrayIterator) Next() bool {
	i.idx++
	return i.idx <= len(i.strArr.Value)
}

func (i *StringArrayIterator) Key() objects.Object {
	return &objects.Int{Value: int64(i.idx - 1)}
}

func (i *StringArrayIterator) Value() objects.Object {
	return &objects.String{Value: i.strArr.Value[i.idx-1]}
}

func (o *StringArray) Iterate() objects.Iterator {
	return &StringArrayIterator{
		strArr: o,
	}
}

func (o *StringArray) CanIterate() bool {
	return true
}

func TestIterable(t *testing.T) {
	strArr := func() *StringArray { return &StringArray{Value: []string{"one", "two", "three"}} }

	expect(t, `for i, s in arr { out += i }`, Opts().Symbol("arr", strArr()).Skip2ndPass(), 3)
	expect(t, `for i, s in arr { out += s }`, Opts().Symbol("arr", strArr()).Skip2ndPass(), "onetwothree")
	expect(t, `for i, s in arr { out += s + i }`, Opts().Symbol("arr", strArr()).Skip2ndPass(), "one0two1three2")
}

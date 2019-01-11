package objects

import "github.com/d5/tengo/compiler/token"

type Iterator interface {
	Object
	Next() bool
	Key() Object
	Value() Object
}

type ArrayIterator struct {
	v []Object
	i int
	l int
}

func NewArrayIterator(v *Array) Iterator {
	return &ArrayIterator{
		v: v.Value,
		l: len(v.Value),
	}
}

func (i *ArrayIterator) TypeName() string {
	return "array-iterator"
}

func (i *ArrayIterator) String() string {
	return "<array-iterator>"
}

func (i *ArrayIterator) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

func (i *ArrayIterator) IsFalsy() bool {
	return true
}

func (i *ArrayIterator) Equals(Object) bool {
	return false
}

func (i *ArrayIterator) Copy() Object {
	return &ArrayIterator{v: i.v, i: i.i, l: i.l}
}

func (i *ArrayIterator) Next() bool {
	i.i++
	return i.i <= i.l
}

func (i *ArrayIterator) Key() Object {
	return &Int{int64(i.i - 1)}
}

func (i *ArrayIterator) Value() Object {
	return i.v[i.i-1]
}

type MapIterator struct {
	v map[string]Object
	k []string
	i int
	l int
}

func NewMapIterator(v *Map) Iterator {
	var keys []string
	for k := range v.Value {
		keys = append(keys, k)
	}

	return &MapIterator{
		v: v.Value,
		k: keys,
		l: len(keys),
	}
}

func (i *MapIterator) TypeName() string {
	return "map-iterator"
}

func (i *MapIterator) String() string {
	return "<map-iterator>"
}

func (i *MapIterator) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

func (i *MapIterator) IsFalsy() bool {
	return true
}

func (i *MapIterator) Equals(Object) bool {
	return false
}

func (i *MapIterator) Copy() Object {
	return &MapIterator{v: i.v, k: i.k, i: i.i, l: i.l}
}

func (i *MapIterator) Next() bool {
	i.i++
	return i.i <= i.l
}

func (i *MapIterator) Key() Object {
	k := i.k[i.i-1]

	return &String{Value: k}
}

func (i *MapIterator) Value() Object {
	k := i.k[i.i-1]

	return i.v[k]
}

type StringIterator struct {
	v []rune
	i int
	l int
}

func NewStringIterator(v *String) Iterator {
	r := []rune(v.Value)

	return &StringIterator{
		v: r,
		l: len(r),
	}
}

func (i *StringIterator) TypeName() string {
	return "string-iterator"
}

func (i *StringIterator) String() string {
	return "<string-iterator>"
}

func (i *StringIterator) BinaryOp(op token.Token, rhs Object) (Object, error) {
	return nil, ErrInvalidOperator
}

func (i *StringIterator) IsFalsy() bool {
	return true
}

func (i *StringIterator) Equals(Object) bool {
	return false
}

func (i *StringIterator) Copy() Object {
	return &StringIterator{v: i.v, i: i.i, l: i.l}
}

func (i *StringIterator) Next() bool {
	i.i++
	return i.i <= i.l
}

func (i *StringIterator) Key() Object {
	return &Int{int64(i.i - 1)}
}

func (i *StringIterator) Value() Object {
	return &Char{Value: i.v[i.i-1]}
}

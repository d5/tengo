package tengo

// Iterator represents an iterator for underlying data type.
type Iterator interface {
	Object

	// Next returns true if there are more elements to iterate.
	Next() bool

	// Key returns the key or index value of the current element.
	Key() Object

	// Value returns the value of the current element.
	Value() Object
}

// ArrayIterator is an iterator for an array.
type ArrayIterator struct {
	ObjectImpl
	v []Object
	i int
	l int
}

// ArrayIteratorTN is the array iterator type name
const ArrayIteratorTN = "array-iterator"

// TypeName returns the name of the type.
func (i *ArrayIterator) TypeName() string {
	return ArrayIteratorTN
}

func (i *ArrayIterator) String() string {
	return TypeString(ArrayIteratorTN)
}

// IsFalsy returns true if the value of the type is falsy.
func (i *ArrayIterator) IsFalsy() bool {
	return true
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (i *ArrayIterator) Equals(Object) bool {
	return false
}

// Copy returns a copy of the type.
func (i *ArrayIterator) Copy() Object {
	return &ArrayIterator{v: i.v, i: i.i, l: i.l}
}

// Next returns true if there are more elements to iterate.
func (i *ArrayIterator) Next() bool {
	i.i++
	return i.i <= i.l
}

// Key returns the key or index value of the current element.
func (i *ArrayIterator) Key() Object {
	return &Int{Value: int64(i.i - 1)}
}

// Value returns the value of the current element.
func (i *ArrayIterator) Value() Object {
	return i.v[i.i-1]
}

// HasLen returns whether the Object has a length value.
func (i *ArrayIterator) HasLen() bool {
	return true
}

// Len returns the Objects length value.
func (i *ArrayIterator) Len() int {
	return i.l
}

// BytesIterator represents an iterator for a string.
type BytesIterator struct {
	ObjectImpl
	v []byte
	i int
	l int
}

// BytesIteratorTN is the bytes iterator type name
const BytesIteratorTN = "bytes-iterator"

// TypeName returns the name of the type.
func (i *BytesIterator) TypeName() string {
	return BytesIteratorTN
}

func (i *BytesIterator) String() string {
	return TypeString(BytesIteratorTN)
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (i *BytesIterator) Equals(Object) bool {
	return false
}

// Copy returns a copy of the type.
func (i *BytesIterator) Copy() Object {
	return &BytesIterator{v: i.v, i: i.i, l: i.l}
}

// Next returns true if there are more elements to iterate.
func (i *BytesIterator) Next() bool {
	i.i++
	return i.i <= i.l
}

// Key returns the key or index value of the current element.
func (i *BytesIterator) Key() Object {
	return &Int{Value: int64(i.i - 1)}
}

// Value returns the value of the current element.
func (i *BytesIterator) Value() Object {
	return &Int{Value: int64(i.v[i.i-1])}
}

// HasLen returns whether the Object has a length value.
func (i *BytesIterator) HasLen() bool {
	return true
}

// Len returns the Objects length value.
func (i *BytesIterator) Len() int {
	return i.l
}

// MapIterator represents an iterator for the map.
type MapIterator struct {
	ObjectImpl
	v map[string]Object
	k []string
	i int
	l int
}

// MapIteratorTN is the map iterator type name
const MapIteratorTN = "map-iterator"

// TypeName returns the name of the type.
func (i *MapIterator) TypeName() string {
	return MapIteratorTN
}

func (i *MapIterator) String() string {
	return TypeString(MapIteratorTN)
}

// IsFalsy returns true if the value of the type is falsy.
func (i *MapIterator) IsFalsy() bool {
	return true
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (i *MapIterator) Equals(Object) bool {
	return false
}

// Copy returns a copy of the type.
func (i *MapIterator) Copy() Object {
	return &MapIterator{v: i.v, k: i.k, i: i.i, l: i.l}
}

// Next returns true if there are more elements to iterate.
func (i *MapIterator) Next() bool {
	i.i++
	return i.i <= i.l
}

// Key returns the key or index value of the current element.
func (i *MapIterator) Key() Object {
	k := i.k[i.i-1]
	return &String{Value: k}
}

// Value returns the value of the current element.
func (i *MapIterator) Value() Object {
	k := i.k[i.i-1]
	return i.v[k]
}

// HasLen returns whether the Object has a length value.
func (i *MapIterator) HasLen() bool {
	return true
}

// Len returns the Objects length value.
func (i *MapIterator) Len() int {
	return i.l
}

// StringIterator represents an iterator for a string.
type StringIterator struct {
	ObjectImpl
	v []rune
	i int
	l int
}

// StringIteratorTN is the string iterator type name
const StringIteratorTN = "string-iterator"

// TypeName returns the name of the type.
func (i *StringIterator) TypeName() string {
	return StringIteratorTN
}

func (i *StringIterator) String() string {
	return TypeString(StringIteratorTN)
}

// IsFalsy returns true if the value of the type is falsy.
func (i *StringIterator) IsFalsy() bool {
	return true
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (i *StringIterator) Equals(Object) bool {
	return false
}

// Copy returns a copy of the type.
func (i *StringIterator) Copy() Object {
	return &StringIterator{v: i.v, i: i.i, l: i.l}
}

// Next returns true if there are more elements to iterate.
func (i *StringIterator) Next() bool {
	i.i++
	return i.i <= i.l
}

// Key returns the key or index value of the current element.
func (i *StringIterator) Key() Object {
	return &Int{Value: int64(i.i - 1)}
}

// Value returns the value of the current element.
func (i *StringIterator) Value() Object {
	return &Char{Value: i.v[i.i-1]}
}

// HasLen returns whether the Object has a length value.
func (i *StringIterator) HasLen() bool {
	return true
}

// Len returns the Objects length value.
func (i *StringIterator) Len() int {
	return i.l
}

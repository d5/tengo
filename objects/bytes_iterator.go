package objects

// BytesIterator represents an iterator for a string.
type BytesIterator struct {
	ObjectImpl
	v []byte
	i int
	l int
}

// TypeName returns the name of the type.
func (i *BytesIterator) TypeName() string {
	return "bytes-iterator"
}

func (i *BytesIterator) String() string {
	return "<bytes-iterator>"
}

// Equals returns true if the value of the type
// is equal to the value of another object.
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

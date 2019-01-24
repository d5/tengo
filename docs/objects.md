# Tengo Objects

## Table of Contents

- [Objects](#objects)
  - [Runtime Object Types](#runtime-object-types)
  - [User Object Types](#user-object-types)
- [Callable Objects](#callable-objects)
- [Indexable Objects](#indexable-objects)
- [Index-Assignable Objects](#index-assignable-objects)
- [Iterable Objects](#iterable-objects)
  - [Iterator Interface](#iterator-interface)

## Objects

All object types in Tengo implement [Object](https://godoc.org/github.com/d5/tengo/objects#Object) interface. 

```golang
TypeName() string
```

TypeName method should return the name of the type. Type names are not directly used by the runtime _(except when it reports a run-time error)_, but, it is generally a good idea to keep it short but distinguishable from other types.

```golang
String() string
```
String method should return a string representation of the underlying value. The value returned by String method will be used whenever string formatting for the value is required, most commonly when being converted into String value.

```golang
BinaryOp(op token.Token, rhs Object) (res Object, err error)
```

In Tengo, a type can overload binary operators (`+`, `-`, `*`, `/`, `%`, `&`, `|`, `^`, `&^`, `>>`, `<<`, `>`, `>=`; _note that `<` and `<=` operators are not overloadable as they're simply implemented by switching left-hand side and right-hand side of `>`/`>=` operator_) by implementing BinaryOp method. BinaryOp method takes the operator `op` and the right-hand side object `rhs`, and, should return a resulting value `res`. 

**Error value vs runtime error**

If BinaryOp method returns an error `err` (the second return value), it will be treated as a run-time error, which will halt the execution (`VM.Run() error`) and will return the error to the user. All runtime type implementations, for example, will return an `ErrInvalidOperator` error when the given operator is not supported by the type. 

Alternatively the method can return an `Error` value as its result `res` (the first return value), which will not halt the runtime and will be treated like any other values. As a dynamically typed language, the receiver (another expression or statement) can determine how to translate `Error` value returned from binary operator expression. 

```golang
IsFalsy() bool
```

IsFalsy method should return true if the underlying value is considered to be [falsy](https://github.com/d5/tengo/wiki/Variable-Types#objectisfalsy).

```golang
Equals(o Object) bool
```

Equals method should return true if the underlying value is considered to be equal to the underlying value of another object `o`. When comparing values of different types, the runtime does not guarantee or force anything, but, it's generally a good idea to make the result consistent. For example, a custom integer type may return true when comparing against String value, but, it should return the same result for the same inputs.

```golang
Copy() Object
```

Copy method should a _new_ copy of the same object. All primitive and composite value types implement this method to return a deep-copy of the value, which is recommended for other user types _(as `copy` builtin function uses this Copy method)_, but, it's not a strict requirement by the runtime.

### Runtime Object Types 

These are the Tengo runtime object types:

- Primitive value types: [Int](https://godoc.org/github.com/d5/tengo/objects#Int), [String](https://godoc.org/github.com/d5/tengo/objects#String), [Float](https://godoc.org/github.com/d5/tengo/objects#Float), [Bool](https://godoc.org/github.com/d5/tengo/objects#ArrayIterator), [Char](https://godoc.org/github.com/d5/tengo/objects#Char), [Bytes](https://godoc.org/github.com/d5/tengo/objects#Bytes)
- Composite value types: [Array](https://godoc.org/github.com/d5/tengo/objects#Array), [Map](https://godoc.org/github.com/d5/tengo/objects#Map), [ImmutableMap](https://godoc.org/github.com/d5/tengo/objects#ImmutableMap)
- Functions: [CompiledFunction](https://godoc.org/github.com/d5/tengo/objects#CompiledFunction), [BuiltinFunction](https://godoc.org/github.com/d5/tengo/objects#BuiltinFunction), [UserFunction](https://godoc.org/github.com/d5/tengo/objects#UserFunction)
- [Iterators](https://godoc.org/github.com/d5/tengo/objects#Iterator): [StringIterator](https://godoc.org/github.com/d5/tengo/objects#StringIterator), [ArrayIterator](https://godoc.org/github.com/d5/tengo/objects#ArrayIterator), [MapIterator](https://godoc.org/github.com/d5/tengo/objects#MapIterator), [ImmutableMapIterator](https://godoc.org/github.com/d5/tengo/objects#ImmutableMapIterator)
- [Error](https://godoc.org/github.com/d5/tengo/objects#Error)
- [Undefined](https://godoc.org/github.com/d5/tengo/objects#Undefined)
- Other internal objects: [Closure](https://godoc.org/github.com/d5/tengo/objects#Closure), [CompiledModule](https://godoc.org/github.com/d5/tengo/objects#CompiledModule), [Break](https://godoc.org/github.com/d5/tengo/objects#Break), [Continue](https://godoc.org/github.com/d5/tengo/objects#Continue), [ReturnValue](https://godoc.org/github.com/d5/tengo/objects#ReturnValue)

### User Object Types

Basically Tengo runtime treats and manages both the runtime types and user types exactly the same way as long as they implement Object interface. You can add values of the custom user types (via either [Script.Add](https://godoc.org/github.com/d5/tengo/script#Script.Add) method or by directly manipulating the symbol table and the global variables), and, use them directly in Tengo code.

Here's an example user type, `Time`:

```golang
type Time struct {
	Value time.Time
}

func (t *Time) TypeName() string {
	return "time"
}

func (t *Time) String() string {
	return t.Value.Format(time.RFC3339)
}

func (t *Time) BinaryOp(op token.Token, rhs objects.Object) (objects.Object, error) {
	switch rhs := rhs.(type) {
	case *Time:
		switch op {
		case token.Sub:
			return &objects.Int{
				Value: t.Value.Sub(rhs.Value).Nanoseconds(),
			}, nil
		}
	case *objects.Int:
		switch op {
		case token.Add:
			return &Time{
				Value: t.Value.Add(time.Duration(rhs.Value)),
			}, nil
		case token.Sub:
			return &Time{
				Value: t.Value.Add(-time.Duration(rhs.Value)),
			}, nil
		}
	}

	return nil, objects.ErrInvalidOperator
}

func (t *Time) IsFalsy() bool {
	return t.Value.IsZero()
}

func (t *Time) Equals(o objects.Object) bool {
	if o, ok := o.(*Time); ok {
		return t.Value.Equal(o.Value)
	}

	return false
}

func (t *Time) Copy() objects.Object {
	return &Time{Value: t.Value}
}
```

Now the Tengo runtime recognizes `Time` type, and, any `Time` values can be used directly in the Tengo code:

```golang
s := script.New([]byte(`
	a := currentTime + 10000  // Time + Int = Time
	b := a - currentTime      // Time - Time = Int
`))

// add Time value 'currentTime'
err := s.Add("currentTime", &Time{Value: time.Now()}) 
if err != nil {
	panic(err)
}

c, err := s.Run()
if err != nil {
	panic(err)
}

fmt.Println(c.Get("b")) // "10000"
```

## Callable Objects

Any types that implement [Callable](https://godoc.org/github.com/d5/tengo/objects#Callable) interface (in addition to Object interface), values of such types can be used as if they are functions. 

```golang
type Callable interface {
	Call(args ...Object) (ret Object, err error)
}
```

To make `Time` a callable value, add Call method to the previous implementation:

```golang
func (t *Time) Call(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}

	format, ok := objects.ToString(args[0])
	if !ok {
		return nil, objects.ErrInvalidTypeConversion
	}

	return &objects.String{Value: t.Value.Format(format)}, nil
}
```

Now `Time` values can be "called" like this:

```golang
s := script.New([]byte(`
	a := currentTime + 10000  // Time + Int = Time
	b := a("15:04:05")        // call 'a'
`))

// add Time value 'currentTime'
err := s.Add("currentTime", &Time{Value: time.Now()}) 
if err != nil {
	panic(err)
}

c, err := s.Run()
if err != nil {
	panic(err)
}

fmt.Println(c.Get("b")) // something like "21:15:27"
```

## Indexable Objects

If the type implements [Indexable](https://godoc.org/github.com/d5/tengo/objects#Indexable) interface, it enables dot selector (`value = object.index`) or indexer (`value = object[index]`) syntax for its values.

```golang
type Indexable interface {
	IndexGet(index Object) (value Object, err error)
}
```

If the implementation returns an error (`err`), the VM will treat it as a run-time error. Many runtime types such as Map and Array also implement the same interface:

```golang
func (o *Map) IndexGet(index Object) (res Object, err error) {
	strIdx, ok := index.(*String)
	if !ok {
		err = ErrInvalidIndexType
		return
	}

	val, ok := o.Value[strIdx.Value]
	if !ok {
		val = UndefinedValue
	}

	return val, nil
}
```

Array and Map implementation forces the type of index Object (Int and String respectively), but, it's not required behavior by the VM. It is completely okay to take various index types (or to do type coercion) as long as its result is consistent. 

By convention, Array or Array-like types return `ErrIndexOutOfBounds` error (as a runtime error) when the index is invalid (out of the bounds), and, Map or Map-like types return `Undefined` value when the key does not exist. But, again this is not a requirement, and, the type can implement the behavior however it fits.

## Index-Assignable Objects

If the type implements [IndexAssignable](https://godoc.org/github.com/d5/tengo/objects#IndexAssignable) interface, the values of that type allow assignment using dot selector (`object.index = value`) or indexer (`object[index] = value`) in the assignment statements.

```golang
type IndexAssignable interface {
	IndexSet(index, value Object) error
}
```

Map, Array, and a couple of other runtime types also implement the same interface:

```golang
func (o *Map) IndexSet(index, value Object) (err error) {
	strIdx, ok := ToString(index)
	if !ok {
		err = ErrInvalidTypeConversion
		return
	}

	o.Value[strIdx] = value

	return nil
}
```

Array and Map implementation forces the type of index Object (Int and String respectively), but, it's not required behavior by the VM. It is completely okay to take various index types (or to do type coercion) as long as its result is consistent. 

By convention, Array or Array-like types return `ErrIndexOutOfBounds` error (as a runtime error) when the index is invalid (out of the bounds). But, this is not a requirement, and, the type can implement the behavior however it fits.


## Iterable Objects

Values of the types that implement [Iterable](https://godoc.org/github.com/d5/tengo/objects#Iterable) interface can be used in `for-in` statements (`for key, value in object { ... }`).

```golang
type Iterable interface {
	Iterate() Iterator
}
```

This Iterate method should return another object that implements [Iterator](https://godoc.org/github.com/d5/tengo/objects#Iterator) interface.

### Iterator Interface

```golang
Next() bool
```

Next method should return true if there are more elements to iterate. When used with `for-in` statements, the compiler uses Key and Value methods to populate the current element's key (or index) and value from the object that this iterator represents. The runtime will stop iterating in `for-in` statement when this method returns false. 

```golang
Key() Object
```

Key method should return a key (or an index) Object for the current element of the underlying object. It should return the same value until Next method is called again. By convention, iterators for the map or map-like objects returns the String key, and, iterators for array or array-like objects returns the Int index. But, it's not a requirement by the VM.

```golang
Value() Object
```

Value method should return a value Object for the current element of the underlying object. It should return the same value until Next method is called again.
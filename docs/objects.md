# Object Types

## Table of Contents

- [Tengo Objects](#tengo-objects)
  - [Object Interface](#object-interface)
  - [Callable Interface](#callable-interface)
  - [Indexable Interface](#indexable-interface)
  - [Index-Assignable Interface](#index-assignable-interface)
  - [Iterable Interface](#iterable-interface)
    - [Iterator Interface](#iterator-interface)
- [Runtime Object Types](#runtime-object-types)
- [User Object Types](#user-object-types)

## Tengo Objects

In Tengo, all object types _(both [runtime types](#runtime-object-types) and [user types](#user-object-types))_ must implement [Object](https://godoc.org/github.com/d5/tengo/objects#Object) interface. And some types may implement other optional interfaces ([Callable](https://godoc.org/github.com/d5/tengo/objects#Callable), [Indexable](https://godoc.org/github.com/d5/tengo/objects#Indexable), [IndexAssignable](https://godoc.org/github.com/d5/tengo/objects#IndexAssignable), [Iterable](https://godoc.org/github.com/d5/tengo/objects#Iterable)) to support additional language features.  

### Object Interface

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

IsFalsy method should return true if the underlying value is considered to be [falsy](https://github.com/d5/tengo/blob/master/docs/runtime-types.md#objectisfalsy).

```golang
Equals(o Object) bool
```

Equals method should return true if the underlying value is considered to be equal to the underlying value of another object `o`. When comparing values of different types, the runtime does not guarantee or force anything, but, it's generally a good idea to make the result consistent. For example, a custom integer type may return true when comparing against String value, but, it should return the same result for the same inputs.

```golang
Copy() Object
```

Copy method should return a _new_ copy of the object. Builtin function `copy` uses this method to copy values. Default implementation of all runtime types return a deep-copy values, but, it's not a requirement by the runtime.


### Callable Interface

If the type implements [Callable](https://godoc.org/github.com/d5/tengo/objects#Callable) interface, its values can be invoked as if they were functions. 

```golang
type Callable interface {
	Call(args ...Object) (ret Object, err error)
}
```

### Indexable Interface

If the type implements [Indexable](https://godoc.org/github.com/d5/tengo/objects#Indexable) interface, its values support dot selector (`value = object.index`) and indexer (`value = object[index]`) syntax.

```golang
type Indexable interface {
	IndexGet(index Object) (value Object, err error)
}
```

If `IndexGet` returns an error (`err`), the VM will treat it as a run-time error. 

Array and Map implementation forces the type of index Object to be Int and String respectively, but, it's not a required behavior of the VM. It is completely okay to take various index types as long as it is consistent. 

By convention, Array or Array-like types and Map or Map-like types return `Undefined` value when the key does not exist. But, again, this is not a required behavior.

### Index-Assignable Interface

If the type implements [IndexAssignable](https://godoc.org/github.com/d5/tengo/objects#IndexAssignable) interface, its values support assignment using dot selector (`object.index = value`) and indexer (`object[index] = value`) in the assignment statements.

```golang
type IndexAssignable interface {
	IndexSet(index, value Object) error
}
```

Array and Map implementation forces the type of index Object to be Int and String respectively, but, it's not a required behavior of the VM. It is completely okay to take various index types as long as it is consistent. 

### Iterable Interface

If the type implements [Iterable](https://godoc.org/github.com/d5/tengo/objects#Iterable) interface, its values can be used in `for-in` statements (`for key, value in object { ... }`).

```golang
type Iterable interface {
	Iterate() Iterator
}
```

This Iterate method should return another object that implements [Iterator](https://godoc.org/github.com/d5/tengo/objects#Iterator) interface.

#### Iterator Interface

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

## Runtime Object Types

These are the basic types Tengo runtime supports out of the box:

- Primitive value types: [Int](https://godoc.org/github.com/d5/tengo/objects#Int), [String](https://godoc.org/github.com/d5/tengo/objects#String), [Float](https://godoc.org/github.com/d5/tengo/objects#Float), [Bool](https://godoc.org/github.com/d5/tengo/objects#ArrayIterator), [Char](https://godoc.org/github.com/d5/tengo/objects#Char), [Bytes](https://godoc.org/github.com/d5/tengo/objects#Bytes), [Time](https://godoc.org/github.com/d5/tengo/objects#Time)
- Composite value types: [Array](https://godoc.org/github.com/d5/tengo/objects#Array), [ImmutableArray](https://godoc.org/github.com/d5/tengo/objects#ImmutableArray), [Map](https://godoc.org/github.com/d5/tengo/objects#Map), [ImmutableMap](https://godoc.org/github.com/d5/tengo/objects#ImmutableMap)
- Functions: [CompiledFunction](https://godoc.org/github.com/d5/tengo/objects#CompiledFunction), [BuiltinFunction](https://godoc.org/github.com/d5/tengo/objects#BuiltinFunction), [UserFunction](https://godoc.org/github.com/d5/tengo/objects#UserFunction)
- [Iterators](https://godoc.org/github.com/d5/tengo/objects#Iterator): [StringIterator](https://godoc.org/github.com/d5/tengo/objects#StringIterator), [ArrayIterator](https://godoc.org/github.com/d5/tengo/objects#ArrayIterator), [MapIterator](https://godoc.org/github.com/d5/tengo/objects#MapIterator), [ImmutableMapIterator](https://godoc.org/github.com/d5/tengo/objects#ImmutableMapIterator)
- [Error](https://godoc.org/github.com/d5/tengo/objects#Error)
- [Undefined](https://godoc.org/github.com/d5/tengo/objects#Undefined)
- Other internal objects: [Closure](https://godoc.org/github.com/d5/tengo/objects#Closure), [Break](https://godoc.org/github.com/d5/tengo/objects#Break), [Continue](https://godoc.org/github.com/d5/tengo/objects#Continue), [ReturnValue](https://godoc.org/github.com/d5/tengo/objects#ReturnValue)

See [Runtime Types](https://github.com/d5/tengo/blob/master/docs/runtime-types.md) for more details on these runtime types.

## User Object Types

Users can easily extend and add their own types by implementing the same [Object](https://godoc.org/github.com/d5/tengo/objects#Object) interface, and, Tengo runtime will treat them in the same way as its runtime types with no performance overhead. 

Here's an example user type implementation, `StringArray`:

```golang
type StringArray struct {
	Value []string
}

func (o *StringArray) String() string {
	return strings.Join(o.Value, ", ")
}

func (o *StringArray) BinaryOp(op token.Token, rhs objects.Object) (objects.Object, error) {
	if rhs, ok := rhs.(*StringArray); ok {
		switch op {
		case token.Add:
			if len(rhs.Value) == 0 {
				return o, nil
			}
			return &StringArray{Value: append(o.Value, rhs.Value...)}, nil
		}
	}

	return nil, objects.ErrInvalidOperator
}

func (o *StringArray) IsFalsy() bool {
	return len(o.Value) == 0
}

func (o *StringArray) Equals(x objects.Object) bool {
	if x, ok := x.(*StringArray); ok {
		if len(o.Value) != len(x.Value) {
			return false
		}

		for i, v := range o.Value {
			if v != x.Value[i] {
				return false
			}
		}

		return true
	}

	return false
}

func (o *StringArray) Copy() objects.Object {
	return &StringArray{
		Value: append([]string{}, o.Value...),
	}
}

func (o *StringArray) TypeName() string {
	return "string-array"
}
```

You can use a user type via either [Script.Add](https://godoc.org/github.com/d5/tengo/script#Script.Add) or by directly manipulating the symbol table and the global variables. Here's an example code to add `StringArray` to the script:

```golang
// script that uses 'my_list'
s := script.New([]byte(`
	print(my_list + "three")
`))

myList := &StringArray{Value: []string{"one", "two"}}
s.Add("my_list", myList)  // add StringArray value 'my_list' 
s.Run()                   // prints "one, two, three" 
```

It can also implement `Indexable` and `IndexAssinable` interfaces:

```golang
func (o *StringArray) IndexGet(index objects.Object) (objects.Object, error) {
	intIdx, ok := index.(*objects.Int)
	if ok {
		if intIdx.Value >= 0 && intIdx.Value < int64(len(o.Value)) {
			return &objects.String{Value: o.Value[intIdx.Value]}, nil
		}

		return nil, objects.ErrIndexOutOfBounds
	}

	strIdx, ok := index.(*objects.String)
	if ok {
		for vidx, str := range o.Value {
			if strIdx.Value == str {
				return &objects.Int{Value: int64(vidx)}, nil
			}
		}

		return objects.UndefinedValue, nil
	}

	return nil, objects.ErrInvalidIndexType
}

func (o *StringArray) IndexSet(index, value objects.Object) error {
	strVal, ok := objects.ToString(value)
	if !ok {
		return objects.ErrInvalidIndexValueType
	}

	intIdx, ok := index.(*objects.Int)
	if ok {
		if intIdx.Value >= 0 && intIdx.Value < int64(len(o.Value)) {
			o.Value[intIdx.Value] = strVal
			return nil
		}

		return objects.ErrIndexOutOfBounds
	}

	return objects.ErrInvalidIndexType
}
```

If we implement `Callabale` interface:

```golang
func (o *StringArray) Call(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}

	s1, ok := objects.ToString(args[0])
	if !ok {
		return nil, objects.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}

	for i, v := range o.Value {
		if v == s1 {
			return &objects.Int{Value: int64(i)}, nil
		}
	}

	return objects.UndefinedValue, nil
}
```

Then it can be "invoked":

```golang
s := script.New([]byte(`
	print(my_list("two"))
`))

myList := &StringArray{Value: []string{"one", "two", "three"}}
s.Add("my_list", myList)  // add StringArray value 'my_list' 
s.Run()                   // prints "1" (index of "two")
```

We can also make `StringArray` iterable:

```golang
func (o *StringArray) Iterate() objects.Iterator {
	return &StringArrayIterator{
		strArr: o,
	}
}

type StringArrayIterator struct {
	objectImpl
	strArr *StringArray
	idx    int
}

func (i *StringArrayIterator) TypeName() string {
	return "string-array-iterator"
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
```


# Builtin Functions

## format

Returns a formatted string. The first argument must be a String object. See [this](https://github.com/d5/tengo/blob/master/docs/formatting.md) for more details on formatting.

```golang
a := [1, 2, 3]
s := format("Foo: %v", a) // s == "Foo: [1, 2, 3]"
```

## len

Returns the number of elements if the given variable is array, string, map, or module map.

```golang
v := [1, 2, 3]
l := len(v) // l == 3
```

## copy

Creates a copy of the given variable. `copy` function calls `Object.Copy` interface method, which is expected to return a deep-copy of the value it holds.

```golang
v1 := [1, 2, 3]
v2 := v1
v3 := copy(v1)
v1[1] = 0
print(v2[1]) // "0"; 'v1' and 'v2' referencing the same array
print(v3[1]) // "2"; 'v3' not affected by 'v1'
```

## append

Appends object(s) to an array (first argument) and returns a new array object. (Like Go's `append` builtin.) Currently, this function takes array type only.

```golang
v := [1]
v = append(v, 2, 3) // v == [1, 2, 3]
```

## type_name

Returns the type_name of an object.

```golang
type_name(1) // int
type_name("str") // string
type_name([1, 2, 3]) // array
```

## string

Tries to convert an object to string object. See [Runtime Types](https://github.com/d5/tengo/blob/master/docs/runtime-types.md) for more details on type conversion.

```golang
x := string(123) //  x == "123"
```


Optionally it can take the second argument, which will be returned if the first argument cannot be converted to string. Note that the second argument does not have to be string.

```golang
v = string(undefined, "foo")  // v == "foo"
v = string(undefined, false)  // v == false 
```

## int

Tries to convert an object to int object. See [this](https://github.com/d5/tengo/blob/master/docs/runtime-types.md) for more details on type conversion.

```golang
v := int("123") //  v == 123
```

Optionally it can take the second argument, which will be returned if the first argument cannot be converted to int. Note that the second argument does not have to be int.

```golang
v = int(undefined, 10)    // v == 10
v = int(undefined, false) // v == false 
```

## bool

Tries to convert an object to bool object. See [this](https://github.com/d5/tengo/blob/master/docs/runtime-types.md) for more details on type conversion.

```golang
v := bool(1) //  v == true
```

## float

Tries to convert an object to float object. See [this](https://github.com/d5/tengo/blob/master/docs/runtime-types.md) for more details on type conversion.

```golang
v := float("19.84") //  v == 19.84
```

Optionally it can take the second argument, which will be returned if the first argument cannot be converted to float. Note that the second argument does not have to be float.

```golang
v = float(undefined, 19.84)    // v == 19.84
v = float(undefined, false)    // v == false 
```

## char

Tries to convert an object to char object. See [this](https://github.com/d5/tengo/blob/master/docs/runtime-types.md) for more details on type conversion.

```golang
v := char(89) //  v == 'Y'
```

Optionally it can take the second argument, which will be returned if the first argument cannot be converted to float. Note that the second argument does not have to be float.

```golang
v = char(undefined, 'X')    // v == 'X'
v = char(undefined, false)  // v == false 
```

## bytes

Tries to convert an object to bytes object. See [this](https://github.com/d5/tengo/blob/master/docs/runtime-types.md) for more details on type conversion.

```golang
v := bytes("foo") //  v == [102 111 111]
```

Optionally it can take the second argument, which will be returned if the first argument cannot be converted to float. Note that the second argument does not have to be float.

```golang
v = bytes(undefined, bytes("foo"))    // v == bytes("foo")
v = bytes(undefined, false)           // v == false 
```

If you pass an int to `bytes()` function, it will create a new byte object with the given size.

```golang
v := bytes(100)
```

## time

Tries to convert an object to time value.

```golang
v := time(1257894000) // 2009-11-10 23:00:00 +0000 UTC
```

## is_string

Returns `true` if the object's type is string. Or it returns `false`.

## is_int

Returns `true` if the object's type is int. Or it returns `false`.

## is_bool

Returns `true` if the object's type is bool. Or it returns `false`.

## is_float

Returns `true` if the object's type is float. Or it returns `false`.

## is_char

Returns `true` if the object's type is char. Or it returns `false`.

## is_bytes

Returns `true` if the object's type is bytes. Or it returns `false`.

## is_error

Returns `true` if the object's type is error. Or it returns `false`.

## is_undefined

Returns `true` if the object's type is undefined. Or it returns `false`.

## is_function

Returns `true` if the object's type is function or closure. Or it returns `false`. Note that `is_function` returns `false` for builtin functions and user-provided callable objects. 

## is_callable

Returns `true` if the object is callable (e.g. function, closure, builtin function, or user-provided callable objects). Or it returns `false`.

## is_array

Returns `true` if the object's type is array. Or it returns `false`.

## is_immutable_array

Returns `true` if the object's type is immutable array. Or it returns `false`.

## is_map

Returns `true` if the object's type is map. Or it returns `false`.

## is_immutable_map

Returns `true` if the object's type is immutable map. Or it returns `false`.

## is_iterable

Returns `true` if the object's type is iterable: array, immutable array, map, immutable map, string, and bytes are iterable types in Tengo.

## is_time

Returns `true` if the object's type is time. Or it returns `false`.

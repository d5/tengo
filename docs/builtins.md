# Builtin Functions

## print

Prints a string representation of the given variable to the standard output.

```golang
v := [1, 2, 3]
print(v)  // "[1, 2, 3]"

print(1, 2, 3)
// "1"
// "2"
// "3"
```

## printf

Prints a formatted string to the standard output. It does not append the newline character at the end. The first argument must a String object. It's same as Go's `fmt.Printf`.

```golang
a := [1, 2, 3]
printf("foo %v", a)  // "foo [1, 2, 3]"
```

## sprintf

Returns a formatted string. The first argument must be a String object. It's the same as Go's `fmt.Sprintf`.

```golang
a := [1, 2, 3]
b := sprintp("foo %v", a)  // b == "foo [1, 2, 3]" 
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

## string

Tries to convert an object to string object. See [this](https://github.com/d5/tengo/wiki/Variable-Types) for more details on type conversion.

```golang
x := string(123) //  v == "123"
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

## is_string

Returns `true` if the object is string. Or it returns `false`.

## is_int

Returns `true` if the object is int. Or it returns `false`.

## is_bool

Returns `true` if the object is bool. Or it returns `false`.

## is_float

Returns `true` if the object is float. Or it returns `false`.

## is_char

Returns `true` if the object is char. Or it returns `false`.

## is_bytes

Returns `true` if the object is bytes. Or it returns `false`.

## is_error

Returns `true` if the object is error. Or it returns `false`.

## is_undefined

Returns `true` if the object is undefined. Or it returns `false`.
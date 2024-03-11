# Builtin Functions

## format

Returns a formatted string. The first argument must be a String object. See
[this](https://github.com/d5/tengo/blob/master/docs/formatting.md) for more
details on formatting.

```golang
a := [1, 2, 3]
s := format("Foo: %v", a) // s == "Foo: [1, 2, 3]"
```

## len

Returns the number of elements if the given variable is array, string, map, or
module map.

```golang
v := [1, 2, 3]
l := len(v) // l == 3
```

## copy

Creates a copy of the given variable. `copy` function calls `Object.Copy`
interface method, which is expected to return a deep-copy of the value it holds.

```golang
v1 := [1, 2, 3]
v2 := v1
v3 := copy(v1)
v1[1] = 0
print(v2[1]) // "0"; 'v1' and 'v2' referencing the same array
print(v3[1]) // "2"; 'v3' not affected by 'v1'
```

## append

Appends object(s) to an array (first argument) and returns a new array object.
(Like Go's `append` builtin.) Currently, this function takes array type only.

```golang
v := [1]
v = append(v, 2, 3) // v == [1, 2, 3]
```

## delete

Deletes the element with the specified key from the map type.
First argument must be a map type and second argument must be a string type.
(Like Go's `delete` builtin except keys are always string).
`delete` returns `undefined` value if successful and it mutates given map.

```golang
v := {key: "value"}
delete(v, "key") // v == {}
```

```golang
v := {key: "value"}
delete(v, "missing") // v == {"key": "value"}
```

```golang
delete({}) // runtime error, second argument is missing
delete({}, 1) // runtime error, second argument must be a string type
```

## splice

Deletes and/or changes the contents of a given array and returns
deleted items as a new array. `splice` is similar to
JS `Array.prototype.splice()` except splice is a builtin function and
first argument must an array. First argument must be an array, and
if second and third arguments are provided those must be integers
otherwise runtime error is returned.

Usage:

`deleted_items := splice(array[, start[, delete_count[, item1[, item2[, ...]]]])`

```golang
v := [1, 2, 3]
items := splice(v, 0) // items == [1, 2, 3], v == []
```

```golang
v := [1, 2, 3]
items := splice(v, 1) // items == [2, 3], v == [1]
```

```golang
v := [1, 2, 3]
items := splice(v, 0, 1) // items == [1], v == [2, 3]
```

```golang
// deleting
v := ["a", "b", "c"]
items := splice(v, 1, 2) // items == ["b", "c"], v == ["a"]
// splice(v, 1, 3) or splice(v, 1, 99) has same effect for this example
```

```golang
// appending
v := ["a", "b", "c"]
items := splice(v, 3, 0, "d", "e") // items == [], v == ["a", "b", "c", "d", "e"]
```

```golang
// replacing
v := ["a", "b", "c"]
items := splice(v, 2, 1, "d") // items == ["c"], v == ["a", "b", "d"]
```

```golang
// inserting
v := ["a", "b", "c"]
items := splice(v, 0, 0, "d", "e") // items == [], v == ["d", "e", "a", "b", "c"]
```

```golang
// deleting and inserting
v := ["a", "b", "c"]
items := splice(v, 1, 1, "d", "e") // items == ["b"], v == ["a", "d", "e", "c"]
```

## go

Starts an independent concurrent goroutine which runs fn(arg1, arg2, ...)

If fn is CompiledFunction, the current running VM will be cloned to create
a new VM in which the CompiledFunction will be running.
The fn can also be any object that has Call() method, such as BuiltinFunction,
in which case no cloned VM will be created.
Returns a goroutineVM object that has wait, result, abort methods.

The goroutineVM will not exit unless:
1. All its descendant goroutineVMs exit
2. It calls abort()
3. Its goroutineVM object abort() is called on behalf of its parent VM
The latter 2 cases will trigger aborting procedure of all the descendant
goroutineVMs, which will further result in #1 above.

```golang
var := 0

f1 := func(a,b) { var = 10; return a+b }
f2 := func(a,b,c) { var = 11; return a+b+c }

gvm1 := go(f1,1,2)
gvm2 := go(f2,1,2,5)

fmt.println(gvm1.result()) // 3
fmt.println(gvm2.result()) // 8
fmt.println(var) // 10 or 11
```

* wait() waits for the goroutineVM to complete up to timeout seconds and
returns true if the goroutineVM exited(successfully or not) within the
timeout. It waits forever if the optional timeout not specified,
or timeout < 0.
* abort() triggers the termination process of the goroutineVM and all
its descendant VMs.
* result() waits the goroutineVM to complete, returns Error object if
any runtime error occurred during the execution, otherwise returns the
result value of fn(arg1, arg2, ...)

### 1 client 1 server

Below is a simple client server example:

```golang
reqChan := makechan(8)
repChan := makechan(8)

client := func(interval) {
	reqChan.send("hello")
	for i := 0; true; i++ {
		fmt.println(repChan.recv())
		times.sleep(interval*times.second)
		reqChan.send(i)
	}
}

server := func() {
	for {
		req := reqChan.recv()
		if req == "hello" {
			fmt.println(req)
			repChan.send("world")
		} else {
			repChan.send(req+100)
		}
	}
}

gClient := go(client, 2)
gServer := go(server)

if ok := gClient.wait(5); !ok {
	gClient.abort()
}
gServer.abort()

//output:
//hello
//world
//100
//101
```

### n client n server, channel in channel

```golang
sharedReqChan := makechan(128)

client = func(name, interval, timeout) {
	print := func(s) {
		fmt.println(name, s)
	}
	print("started")

	repChan := makechan(1)
	msg := {chan:repChan}

	msg.data = "hello"
	sharedReqChan.send(msg)
	print(repChan.recv())

	for i := 0; i * interval < timeout; i++ {
		msg.data = i
		sharedReqChan.send(msg)
		print(repChan.recv())
		times.sleep(interval*times.second)
	}
}

server = func(name) {
	print := func(s) {
		fmt.println(name, s)
	}
	print("started")

	for {
		req := sharedReqChan.recv()
		if req.data == "hello" {
			req.chan.send("world")
		} else {
			req.chan.send(req.data+100)
		}
	}
}

clients := func() {
	for i :=0; i < 5; i++ {
		go(client, format("client %d: ", i), 1, 4)
	}
}

servers := func() {
	for i :=0; i < 2; i++ {
		go(server, format("server %d: ", i))
	}
}

// After 4 seconds, all clients should have exited normally
gclts := go(clients)
// If servers exit earlier than clients, then clients may be
// blocked forever waiting for the reply chan, because servers
// were aborted with the req fetched from sharedReqChan before
// sending back the reply.
// In such case, do below to abort() the clients manually
//go(func(){times.sleep(6*times.second); gclts.abort()})

// Servers are infinite loop, abort() them after 5 seconds
gsrvs := go(servers)
if ok := gsrvs.wait(5); !ok {
	gsrvs.abort()
}

// Main VM waits here until all the child "go" finish

// If somehow the main VM is stuck, that is because there is
// at least one child VM that has not exited as expected, we
// can do abort() to force exit.
abort()

//output:
//3
//8
//hello
//world
//100
//101

//unordered output:
//client 4: started
//server 0: started
//client 4: world
//client 4: 100
//client 3: started
//client 3: world
//client 3: 100
//client 2: started
//client 2: world
//client 2: 100
//client 0: started
//client 0: world
//client 0: 100
//client 1: started
//client 1: world
//client 1: 100
//server 1: started
//client 1: 101
//client 2: 101
//client 4: 101
//client 0: 101
//client 3: 101
//client 3: 102
//client 0: 102
//client 2: 102
//client 1: 102
//client 4: 102
//client 0: 103
//client 3: 103
//client 2: 103
//client 1: 103
//client 4: 103

```

## abort
Triggers the termination process of the current VM and all its descendant VMs.

## makechan

Makes a channel to send/receive object and returns a chan object that has
send, recv, close methods.

```golang
unbufferedChan := makechan()
bufferedChan := makechan(128)

// Send will block if the channel is full.
bufferedChan.send("hello") // send string
bufferedChan.send(55) // send int
bufferedChan.send([66, makechan(1)]) // channel in channel

// Receive will block if the channel is empty.
obj := bufferedChan.recv()

// Send to a closed channel causes panic.
// Receive from a closed channel returns undefined value.
unbufferedChan.close()
bufferedChan.close()
```

On the time the VM that the chan is running in is aborted, the sending
or receiving call returns immediately.

## type_name

Returns the type_name of an object.

```golang
type_name(1) // int
type_name("str") // string
type_name([1, 2, 3]) // array
```

## string

Tries to convert an object to string object. See
[Runtime Types](https://github.com/d5/tengo/blob/master/docs/runtime-types.md)
for more details on type conversion.

```golang
x := string(123) //  x == "123"
```

Optionally it can take the second argument, which will be returned if the first
argument cannot be converted to string. Note that the second argument does not
have to be string.

```golang
v = string(undefined, "foo")  // v == "foo"
v = string(undefined, false)  // v == false
```

## int

Tries to convert an object to int object. See
[this](https://github.com/d5/tengo/blob/master/docs/runtime-types.md)
for more details on type conversion.

```golang
v := int("123") //  v == 123
```

Optionally it can take the second argument, which will be returned if the first
argument cannot be converted to int. Note that the second argument does not have
to be int.

```golang
v = int(undefined, 10)    // v == 10
v = int(undefined, false) // v == false
```

## bool

Tries to convert an object to bool object. See
[this](https://github.com/d5/tengo/blob/master/docs/runtime-types.md) for more
details on type conversion.

```golang
v := bool(1) //  v == true
```

## float

Tries to convert an object to float object. See
[this](https://github.com/d5/tengo/blob/master/docs/runtime-types.md) for more
details on type conversion.

```golang
v := float("19.84") //  v == 19.84
```

Optionally it can take the second argument, which will be returned if the first
argument cannot be converted to float. Note that the second argument does not
have to be float.

```golang
v = float(undefined, 19.84)    // v == 19.84
v = float(undefined, false)    // v == false
```

## char

Tries to convert an object to char object. See
[this](https://github.com/d5/tengo/blob/master/docs/runtime-types.md) for more
details on type conversion.

```golang
v := char(89) //  v == 'Y'
```

Optionally it can take the second argument, which will be returned if the first
argument cannot be converted to float. Note that the second argument does not
have to be float.

```golang
v = char(undefined, 'X')    // v == 'X'
v = char(undefined, false)  // v == false
```

## bytes

Tries to convert an object to bytes object. See
[this](https://github.com/d5/tengo/blob/master/docs/runtime-types.md) for more
details on type conversion.

```golang
v := bytes("foo") //  v == [102 111 111]
```

Optionally it can take the second argument, which will be returned if the first
argument cannot be converted to float. Note that the second argument does not
have to be float.

```golang
v = bytes(undefined, bytes("foo"))    // v == bytes("foo")
v = bytes(undefined, false)           // v == false
```

If you pass an int to `bytes()` function, it will create a new byte object with
the given size.

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

Returns `true` if the object's type is function or closure. Or it returns
`false`. Note that `is_function` returns `false` for builtin functions and
user-provided callable objects.

## is_callable

Returns `true` if the object is callable (e.g. function, closure, builtin
function, or user-provided callable objects). Or it returns `false`.

## is_array

Returns `true` if the object's type is array. Or it returns `false`.

## is_immutable_array

Returns `true` if the object's type is immutable array. Or it returns `false`.

## is_map

Returns `true` if the object's type is map. Or it returns `false`.

## is_immutable_map

Returns `true` if the object's type is immutable map. Or it returns `false`.

## is_iterable

Returns `true` if the object's type is iterable: array, immutable array, map,
immutable map, string, and bytes are iterable types in Tengo.

## is_time

Returns `true` if the object's type is time. Or it returns `false`.

# Tengo Syntax

Tengo's syntax is designed to be familiar to Go developers while being a bit simpler and more streamlined.

## Comments

Tengo supports line comments (`//...`) and block comments (`/* ... */`).

```golang
/* 
  multi-line block comments 
*/

a := 5 // line comments
```
> [Run in Playground](https://tengolang.com/?s=02e384399a0397b0a752f08604ccb244d1a6cb37)

## Types and Assignment

Tengo is a dynamically typed language, and, you can initialize the variables using `:=` operator. 

```golang
a := 1984 		// int
b := "aomame"		// string
c := -9.22		// float
d := true		// bool
e := '九'		// char
f := [1, false, "foo"]	// array
g := {			// map
    h: 439,
    i: 12.34,
    j: [0, 9, false]
}
k := func(l, m) {	// function
    return l + m
}
```
> [Run in Playground](https://tengolang.com/?s=f8626a711769502ce20e4560ace65c0e9c1279f4)

After the variable is initialized, it can be re-assigned different value using `=` operator. 

```golang
a := 1928		// int
a = "foo"		// string
f := func() {
    a := false		// 'a' is defined in the function scope
    a = [1, 2, 3]	// and thus does not affect 'a' in global scope.
}
a == "foo" 		// still "foo"
```
> [Run in Playground](https://tengolang.com/?s=1d39bc2af5c51417df82b32db47a0e6a156d48ec)


Type is not directly specified, but, you can use type-coercion functions to convert between types.

```golang
s1 := string(1984)  // "1984"
i2 := int("-999")   // -999
f3 := float(-51)    // -51.0
b4 := bool(1)       // true
c5 := char("X")     // 'X'
```
> [Run in Playground](https://tengolang.com/?s=8d57905b82959eb244e9bbd2111e12ee04a33045)

_See [Variable Types](https://github.com/d5/tengo/wiki/Variable-Types) for more details on the variable types._

## Indexing

You can use the dot selector (`.`) and indexer (`[]`) operator to read or write elements of arrays, strings, or maps.

```golang
["one", "two", "three"][1]	// == "two"

m := {
    a: 1,
    b: [2, 3, 4],
    c: func() { return 10 }
}
m.a				// == 1
m["b"][1]			// == 3
m.c()				// == 10
m.x = 5				// add 'x' to map 'm'
//m.b[5] = 0			// but this is an error: index out of bounds
```
> [Run in Playground](https://tengolang.com/?s=d510c75ed8f06ef1e22c1aaf8a7d4565c793514c)

For sequence types (string, bytes, array), you can use slice operator (`[:]`) too.

```golang
a := [1, 2, 3, 4, 5][1:3]	// == [2, 3]
b := [1, 2, 3, 4, 5][3:]	// == [4, 5]
c := [1, 2, 3, 4, 5][:3]	// == [1, 2, 3]
d := "hello world"[2:10]	// == "llo worl"
```
> [Run in Playground](https://tengolang.com/?s=214ab490bb24549578770984985f6b161aed915d)

## Conditional Expression

Tengo supports the ternary conditional expression (`cond ? expr1 : expr2`):

```golang
a := true ? 1 : -1    // a == 1

min := func(a, b) {
    return a < b ? a : b
}
b := min(5, 10)       // b == 5
```

> [Run in Playground](https://tengolang.com/?s=24724cc03e8ef2e56b7851017f8c2a577897961e)

## Functions

In Tengo, functions are first-class citizen, and, it also supports closures, functions that captures variables in outer scopes. In the following example, the function returned from `adder` is capturing `base` variable.

```golang
adder := func(base) {
    return func(x) { return base + x }	// capturing 'base'
}
add5 := adder(5)
nine := add5(4)		// == 9
```
> [Run in Playground](https://tengolang.com/?s=fba79990473d5b38cc944dfa225d38580ddaf422)

## Flow Control

For flow control, Tengo currently supports **if-else**, **for**, **for-in** statements.

```golang
// IF-ELSE
if a < 0 {
    // ...
} else if a == 0 {
    // ...
} else {
    // ...
}

// IF with init statement
if a := 0; a < 10 {
    // ...
} else {
    // ...
}

// FOR
for a:=0; a<10; a++ {
    // ...
}

// FOR condition-only (like WHILE in other languages)
for a < 10 {
    // ...
}

// FOR-IN
for x in [1, 2, 3] {		// array: element
    // ...
}
for i, x in [1, 2, 3] {		// array: index and element
    // ...
} 
for k, v in {k1: 1, k2: 2} {	// map: key and value
    // ...
}
```

## Immutable Values

Basically, all values of the primitive types (Int, Float, String, Bytes, Char, Bool) are immutable.

```golang
s := "12345"
s[1] = 'b'   // error: String is immutable
s = "foo"    // ok: this is not mutating the value 
             //  but updating reference 's' with another String value
```

The composite types (Array, Map) are mutable by default, but, you can make them immutable using `immutable` expression.

```golang
a := [1, 2, 3]
a[1] = "foo"    // ok: array is mutable

b := immutable([1, 2, 3])
b[1] = "foo"    // error: 'b' references to an immutable array.
b = "foo"       // ok: this is not mutating the value of array
                //  but updating reference 'b' with different value
``` 

Not that, if you copy (using `copy` builtin function) an immutable value, it will return a "mutable" copy. Also, immutability is not applied to the individual elements of the array or map value, unless they are explicitly made immutable.

```golang
a := immutable({b: 4, c: [1, 2, 3]})
a.b = 5     // error
a.c[1] = 5  // ok: because 'a.c' is not immutable

a = immutable({b: 4, c: immutable([1, 2, 3])}) 
a.c[1] = 5  // error
```

## Errors

An error object is created using `error` expression. An error can contain value of any types, and, the underlying value can be read using `.value` selector.
 
```golang
err1 := error("oops")   // error with string value
err2 := error(1+2+3)    // error with int value
if is_error(err1) {     // 'is_error' builtin function
    err_val := err1.value   // get underlying value 
}  
``` 
> [Run in Playground](https://tengolang.com/?s=5eaba4289c9d284d97704dd09cb15f4f03ad05c1)

## Modules

You can load other scripts as modules using `import` expression.

Main script:
```golang
mod1 := import("./mod1") // assuming mod1.tengo file exists in the current directory 
                         // same as 'import("./mod1.tengo")' or 'import("mod1")'
mod1.func1(a)            // module function 
a += mod1.foo            // module variable
//mod1.foo = 5           // error: module variables are read-only
```

`mod1.tengo` file:

```golang
func1 := func(x) { print(x) }
foo := 2
```

Basically, `import` expression returns all the global variables defined in the module as a Map-like value. One can access the functions or variables defined in the module using `.` selector or `["key"]` indexer, but, module variables are immutable.

Also, you can use `import` to load the [Standard Library](https://github.com/d5/tengo/blob/master/docs/stdlib.md).

```golang
math := import("math")
a := math.abs(-19.84) // == 19.84
```

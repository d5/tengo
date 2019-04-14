# Tengo Syntax

Tengo's syntax is designed to be familiar to Go developers while being a bit simpler and more streamlined.

**You can test the Tengo code in online [Playground](https://tengolang.com).**

## Comments

Tengo supports line comments (`//...`) and block comments (`/* ... */`).

```golang
/* 
  multi-line block comments 
*/

a := 5 // line comments
```

## Values and Variables

In Tengo, everything is a value. 

```golang
19 + 84                // int values
"aomame" + `kawa`      // string values
-9.22 + 1e10           // float values
true || false          // bool values
'九' > '9'              // char values
[1, false, "foo"]      // array value
{a: 12.34, b: "bar"}   // map value
func() { /*...*/ }     // function value
```

_See [Runtime Types](https://github.com/d5/tengo/blob/master/docs/runtime-types.md) for the full list of Tengo runtime types._

And the values can be assigned to variables using `:=` and `=` operators.

```golang
a := 19.84              // 'a' has float value '19.84'
a = "foo bar"           // 'a' now has string value "foo bar"
f := func() { /*...*/ } // 'f' has a function value
```

In Tengo, all values have the underlying types, but, the variables are not directly associated with the types. Variables simply reference the values, and, they can even be re-assigned values with different types.

Symantic of `:=` and `=` operators are the same as Go. `:=` is used to define a new variable (symbol) in the current scope. `=` is used to re-assign value to an existing variable (symbol) defined in the current scope or its outer scopes.

```golang
a := 1234              // 'a' in global scope
b := "foo"             // 'b' in global scope
func() {
    a = -1984          // 'a' from global scope is re-assigned different int value
    b := "bar"         // a new 'b' variable is defined in the function scope
}()
```

Unlike Go, there's no declarations in Tengo. Everything is assignment.

```golang
var a                  // compile error
func b { /*....*/ }    // compile error
```

Also there's no pointers in Tengo.

## Type Coercion

Tengo as a dynamically typed language, the type is not directly specified, but, you can use builtin functions to convert a value into a differen type:

```golang
s1 := string(1984)  // "1984"
i2 := int("-999")   // -999
f3 := float(-51)    // -51.0
b4 := bool(1)       // true
c5 := char("X")     // 'X'
```

_See [Builtin Functions](https://github.com/d5/tengo/blob/master/docs/builtins.md) and [Operators](https://github.com/d5/tengo/blob/master/docs/operators.md) for more details on type coercions._

## Operators

### Unary Operators

| Operator | Usage | Types |
| :---: | :---: | :---: |
| `+`   | same as `0 + x` | int, float |
| `-`   | same as `0 - x` | int, float |
| `!`   | logical NOT | all types* |
| `^`   | bitwise complement | int |

_In Tengo, all values can be either [truthy or falsy](https://github.com/d5/tengo/blob/d5-patch-1/docs/runtime-types.md#objectisfalsy)._

### Binary Operators

| Operator | Usage | Types |
| :---: | :---: | :---: |
| `==` | equal | all types |
| `!=` | not equal | all types |
| `&&` | logical AND | all types |
| `\|\|` | logical OR | all types |
| `+`   | add/concat | int, float, string, char, time, array |
| `-`   | subtract | int, float, char, time |
| `*`   | multiply | int, float |
| `/`   | divide | int, float |
| `&`   | bitwise AND | int |
| `\|`   | bitwise OR | int |
| `^`   | bitwise XOR | int |
| `&^`   | bitclear (AND NOT) | int |
| `<<`   | shift left | int |
| `>>`   | shift right | int |
| `<`   | less than | int, float, char, time |
| `<=`   | less than or equal to | int, float, char, time |
| `>`   | greater than | int, float, char, time |
| `>=`   | greater than or equal to | int, float, char, time |

_See [Operators](https://github.com/d5/tengo/blob/d5-patch-1/docs/operators.md) for more details._

### Ternary Operators

Unlike Go, Tengo has a ternary conditional operator `?:`.

```golang
a := b > 4 ? "big" : false
```

### Assignment and Increment Operators

| Operator | Usage |
| :---: | :---: |
| `+=` | `(lhs) = (lhs) + (rhs)` |
| `-=` | `(lhs) = (lhs) - (rhs)` |
| `*=` | `(lhs) = (lhs) * (rhs)` |
| `/=` | `(lhs) = (lhs) / (rhs)` |
| `%=` | `(lhs) = (lhs) % (rhs)` |
| `&=` | `(lhs) = (lhs) & (rhs)` |
| `\|=` | `(lhs) = (lhs) \| (rhs)` |
| `&^=` | `(lhs) = (lhs) &^ (rhs)` |
| `^=` | `(lhs) = (lhs) ^ (rhs)` |
| `<<=` | `(lhs) = (lhs) << (rhs)` |
| `>>=` | `(lhs) = (lhs) >> (rhs)` |
| `++` | `(lhs) = (lhs) + 1` |
| `--` | `(lhs) = (lhs) - 1` |

### Operator Precedences

Unary operators have the highest precedence, and, ternary operator has the lowest precendece. There are five precedence levels for binary operators. Multiplication operators bind strongest, followed by addition operators, comparison operators, && (logical AND), and finally || (logical OR):

|Precedence|Operator|
| :---: | :---: |
| 5 | `*`  `/`  `%`  `<<`  `>>`  `&`  `&^` |
| 4 | `+`  `-`  `\|`  `^` |
| 3 | `==`  `!=`  `<`  `<=`  `>`  `>=` |
| 2 | `&&` |
| 1 | `\|\|` |

Like Go, `++` and `--` operators form statements, not expressions, they fall outside the operator hierarchy. 

### Selector and Indexer

You can use the dot selector (`.`) and indexer (`[]`) operator to read or write elements of arrays, strings, or maps.

Reading a nonexistent index returns `Undefined` value.

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
m["b"][5]                       // == undefined
m["b"][5].d                     // == undefined
//m.b[5] = 0			// but this is an error: index out of bounds
```

For sequence types (string, bytes, array), you can use slice operator (`[:]`) too.

```golang
a := [1, 2, 3, 4, 5][1:3]	// == [2, 3]
b := [1, 2, 3, 4, 5][3:]	// == [4, 5]
c := [1, 2, 3, 4, 5][:3]	// == [1, 2, 3]
d := "hello world"[2:10]	// == "llo worl"
c := [1, 2, 3, 4, 5][-1:10]    // == [1, 2, 3, 4, 5]
```

## Conditional Expression

Tengo supports the ternary conditional expression (`cond ? expr1 : expr2`):

```golang
a := true ? 1 : -1    // a == 1

min := func(a, b) {
    return a < b ? a : b
}
b := min(5, 10)       // b == 5
```

## Functions

In Tengo, functions are first-class citizen, and, it also supports closures, functions that captures variables in outer scopes. In the following example, the function returned from `adder` is capturing `base` variable.

```golang
adder := func(base) {
    return func(x) { return base + x }	// capturing 'base'
}
add5 := adder(5)
nine := add5(4)		// == 9
```

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

The compound types (Array, Map) are mutable by default, but, you can make them immutable using `immutable` expression.

```golang
a := [1, 2, 3]
a[1] = "foo"    // ok: array is mutable

b := immutable([1, 2, 3])
b[1] = "foo"    // error: 'b' references to an immutable array.
b = "foo"       // ok: this is not mutating the value of array
                //  but updating reference 'b' with different value
``` 

Note that, if you copy (using `copy` builtin function) an immutable value, it will return a "mutable" copy. Also, immutability is not applied to the individual elements of the array or map value, unless they are explicitly made immutable.

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

## Modules

You can load other scripts as modules using `import` expression.

Main script:

```golang
sum := import("./sum")   // assuming sum.tengo file exists in the current directory 
                         // same as 'import("./sum.tengo")' or 'import("sum")'
fmt.print(sum(10))       // module function 
```

`sum.tengo` file:

```golang
base := 5

export func(x) {
    return x + base
}
```

In Tengo, modules are very similar to functions.

- `import` expression loads the module and execute like a function. 
- Module should return a value using `export` statement.
    - Module can return `export` any Tengo objects: int, string, map, array, function, etc.
    - `export` in a module is like `return` in a function: it stops execution and return a value to the importing code.
    - `export`-ed values are always immutable.
    - If the module does not have any `export` statement, `import` expression simply returns `undefined`. _(Just like the function that has no `return`.)_  
    - Note that `export` statement is completely ignored and not evaluated if the code is executed as a regular script.  

Also, you can use `import` to load the [Standard Library](https://github.com/d5/tengo/blob/master/docs/stdlib.md).

```golang
math := import("math")
a := math.abs(-19.84) // == 19.84
```

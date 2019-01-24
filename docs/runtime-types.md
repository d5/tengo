# Tengo Runtime Types

- **Int**: signed 64bit integer
- **String**: string
- **Float**: 64bit floating point
- **Bool**: boolean
- **Char**: character (`rune` in Go)
- **Bytes**: byte array (`[]byte` in Go)
- **Array**: objects array (`[]Object` in Go)
- **Map**: objects map with string keys (`map[string]Object` in Go)
- **ImmutableMap**: immutable object map with string keys (`map[string]Object` in Go)
- **Error**: an error with underlying Object value of any type
- **Undefined**: undefined

## Type Conversion/Coercion Table
|src\dst  |Int      |String        |Float    |Bool      |Char   |Bytes  |Array  |Map    |IMap|Error  |Undefined|
| :---:   | :---:   | :---:        | :---:   | :---:    | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
|Int      |   -     |_strconv_     |float64(v)|!IsFalsy()| rune(v)|**X**|**X**|**X**|**X**|**X**|**X**|
|String   |_strconv_|   -          |_strconv_|!IsFalsy()|**X**|[]byte(s)|**X**|**X**|**X**|**X**|**X**|
|Float    |int64(f) |_strconv_     | -       |!IsFalsy()|**X**|**X**|**X**|**X**|**X**|**X**|**X**|
|Bool     |1 / 0    |"true" / "false"|**X**    |   -   |**X**|**X**|**X**|**X**|**X**|**X**|**X**|
|Char     |int64(c) |string(c)     |**X**    |!IsFalsy()|   -   |**X**|**X**|**X**|**X**|**X**|**X**|
|Bytes    |**X**    |string(y)|**X**    |!IsFalsy()|**X**|   -   |**X**|**X**|**X**|**X**|**X**|
|Array    |**X**    |"[...]"       |**X**    |!IsFalsy()|**X**|**X**|   -   |**X**|**X**|**X**|**X**|
|Map      |**X**    |"{...}"       |**X**    |!IsFalsy()|**X**|**X**|**X**|   -   |**X**|**X**|**X**|
|IMap  |**X**    |"{...}"       |**X**    |!IsFalsy()|**X**|**X**|**X**|**X**|   -   |**X**|**X**|
|Error    |**X**    |"error: ..."  |**X**    |false|**X**|**X**|**X**|**X**|**X**|   -   |**X**|
|Undefined|**X**    |**X**|**X**    |false|**X**|**X**|**X**|**X**|**X**|**X**|   -    |

_* **X**: No conversion; Typed value functions for `script.Variable` will return zero values._  
_* strconv: converted using Go's conversion functions from `strconv` package._  
_* IsFalsy(): use [Object.IsFalsy()](#objectisfalsy) function_

## Object.IsFalsy()

`Object.IsFalsy()` interface method is used to determine if a given value should evaluate to `false` (e.g. for condition expression of `if` statement).

- **Int**: `n == 0`
- **String**: `len(s) == 0`
- **Float**: `isNaN(f)`
- **Bool**: `!b`
- **Char**: `c == 0`
- **Bytes**: `len(bytes) == 0`
- **Array**: `len(arr) == 0`
- **Map**: `len(map) == 0`
- **ImmutableMap**: `len(map) == 0`
- **Error**: `true` _(Error is always falsy)_
- **Undefined**: `true` _(Undefined is always falsy)_

## Type Conversion Builtin Functions

- `string(x)`: tries to convert `x` into string; returns `undefined` if failed
- `int(x)`: tries to convert `x` into int; returns `undefined` if failed
- `bool(x)`: tries to convert `x` into bool; returns `undefined` if failed
- `float(x)`: tries to convert `x` into float; returns `undefined` if failed
- `char(x)`: tries to convert `x` into char; returns `undefined` if failed
- `bytes(x)`: tries to convert `x` into bytes; returns `undefined` if failed
  - `bytes(N)`: as a special case this will create a Bytes variable with the given size `N` (only if `N` is int)

## Type Checking Builtin Functions

- `is_string(x)`: returns `true` if `x` is string; `false` otherwise
- `is_int(x)`: returns `true` if `x` is int; `false` otherwise
- `is_bool(x)`: returns `true` if `x` is bool; `false` otherwise
- `is_float(x)`: returns `true` if `x` is float; `false` otherwise
- `is_char(x)`: returns `true` if `x` is char; `false` otherwise
- `is_bytes(x)`: returns `true` if `x` is bytes; `false` otherwise
- `is_error(x)`: returns `true` if `x` is error; `false` otherwise
- `is_undefined(x)`: returns `true` if `x` is undefined; `false` otherwise
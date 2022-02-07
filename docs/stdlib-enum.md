# Module - "enum"

```golang
enum := import("enum")
```

## Functions

- `all(x, fn) => bool`: returns true if the given function `fn` evaluates to a
  truthy value on all of the items in `x`. It returns undefined if `x` is not
  enumerable.
- `any(x, fn) => bool`: returns true if the given function `fn` evaluates to a
  truthy value on any of the items in `x`. It returns undefined if `x` is not
  enumerable.
- `chunk(x, size) => [object]`: returns an array of elements split into groups
  the length of size. If `x` can't be split evenly, the final chunk will be the
  remaining elements. It returns undefined if `x` is not array.
- `at(x, key) => object`: returns an element at the given index (if `x` is
  array) or key (if `x` is map). It returns undefined if `x` is not enumerable.
- `each(x, fn)`: iterates over elements of `x` and invokes `fn` for each
  element. `fn` is invoked with two arguments: `key` and `value`. `key` is an
  int index if `x` is array. `key` is a string key if `x` is map. It does not
  iterate and returns undefined if `x` is not enumerable.`
- `filter(x, fn) => [object]`: iterates over elements of `x`, returning an
  array of all elements `fn` returns truthy for. `fn` is invoked with two
  arguments: `key` and `value`. `key` is an int index if `x` is array. It returns
  undefined if `x` is not array.
- `find(x, fn) => object`: iterates over elements of `x`, returning value of
  the first element `fn` returns truthy for. `fn` is invoked with two
  arguments: `key` and `value`. `key` is an int index if `x` is array. `key` is
  a string key if `x` is map. It returns undefined if `x` is not enumerable.
- `find_key(x, fn) => int/string`: iterates over elements of `x`, returning key
  or index of the first element `fn` returns truthy for. `fn` is invoked with
  two arguments: `key` and `value`. `key` is an int index if `x` is array.
  `key` is a string key if `x` is map. It returns undefined if `x` is not
  enumerable.
- `map(x, fn) => [object]`: creates an array of values by running each element
  in `x` through `fn`. `fn` is invoked with two arguments: `key` and `value`.
  `key` is an int index if `x` is array. `key` is a string key if `x` is map.
  It returns undefined if `x` is not enumerable.
- `key(k, _) => object`: returns the first argument.
- `value(_, v) => object`: returns the second argument.

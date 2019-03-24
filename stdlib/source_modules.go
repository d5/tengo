// Code generated using gensrcmods.go; DO NOT EDIT.

package stdlib

// SourceModules are source type standard library modules.
var SourceModules = map[string]string{
	"enum": `is_enumerable := func(v) {
    return is_array(v) || is_map(v) || is_string(v) || is_bytes(v) || is_immutable_array(v) || is_immutable_map(v)
}

is_sequence := func(v) {
    return is_array(v) || is_string(v) || is_bytes(v) || is_immutable_array(v)
}

export {
  // all returns true if the given function fn evaluates to a truthy value on
  // all of the items in the enumerable.
  all: func(enumerable, fn) {
    for k, v in enumerable {
      if !fn(k, v) { return false }
    }
    return true
  },
  // any returns true if the given function fn evaluates to a truthy value on
  // any of the items in the enumerable.
  any: func(enumerable, fn) {
    for k, v in enumerable {
      if fn(k, v) { return true }
    }
    return false
  },
  // chunk returns an array of elements split into groups the length of size.
  // If the enumerable can't be split evenly, the final chunk will be the
  // remaining elements. It returns an empty array if the given argument is
  // not a sequence type.
  chunk: func(enumerable, size) {
    if !is_sequence(enumerable) { return [] }

    numElements := len(enumerable)

    if !numElements {
      return []
    }
    
    res := []
    idx := 0
    for idx < numElements {
      res = append(res, enumerable[idx:idx+size])
      idx += size
    }
    return res
  },
  // at returns an element at the given index.
  // It returns undefined if index is out of bounds.
  at: func(enumerable, index) {
    return enumerable[index]
  },
  each: func(enumerable, fn) {
    for k, v in enumerable {
      fn(k, v)
    }
  },
  filter: func(enumerable, fn) {
    if !is_sequence(enumerable) { return enumerable }

    dst := []
    for k, v in enumerable {
      if fn(k, v) { dst = append(dst, v) }
    }
    return dst
  }
}
`,
}

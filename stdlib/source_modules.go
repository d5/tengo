// Code generated using gensrcmods.go; DO NOT EDIT.

package stdlib

import "github.com/d5/tengo/objects"

// SourceModules are source type standard library modules.
var SourceModules = map[string]*objects.SourceModule{
	"enum": {Src: []byte(`export {
  // all returns true if the given function fn evaluates to a truthy value on
  // all of the items in the enumerable.
  all: func(enumerable, fn) {
    for k, v in enumerable {
      if !fn(k, v) {
        return false
      }
    }
    return true
  },
  // any returns true if the given function fn evaluates to a truthy value on
  // any of the items in the enumerable.
  any: func(enumerable, fn) {
    for k, v in enumerable {
      if fn(k, v) {
        return true
      }
    }
    return false
  },
  // chunk returns an array of elements split into groups the length of size.
  // If the enumerable can't be split evenly, the final chunk will be the
  // remaining elements.
  chunk: func(enumerable, size) {
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
  }
}
`)},
}

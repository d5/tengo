// Code generated using gensrcmods.go; DO NOT EDIT.

package stdlib

import "github.com/d5/tengo/objects"

// SourceModules are source type standard library modules.
var SourceModules = map[string]*objects.SourceModule{
	"enum": {Src: []byte(`export {
  each: func(x, f) {
    for k, v in x {
      f(k, v)
    }
  },
  all: func(x, f) {
    for k, v in x {
      if !f(k, v) {
        return false
      }
    }
    return true
  },
  any: func(x, f) {
    for k, v in x {
      if f(k, v) {
        return true
      }
    }
    return false
  }
}
`)},
}

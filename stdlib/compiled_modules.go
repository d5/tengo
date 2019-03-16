// Code generated using genmods.go; DO NOT EDIT.

package stdlib

// CompiledModules are Tengo-written standard libraries.
var CompiledModules = map[string]string{
	"enum": `export {
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
  },
}
`,
}

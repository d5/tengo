package tengo

// CallableFunc is a function signature for the callable functions.
type CallableFunc = func(rt Interop, args ...Object) (ret Object, err error)

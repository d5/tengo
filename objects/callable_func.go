package objects

// CallableFunc is a function signature for the callable functions.
type CallableFunc = func(rt Runtime, args ...Object) (ret Object, err error)

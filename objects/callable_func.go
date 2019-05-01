package objects

// CallableFunc is a function signature for the callable functions.
type CallableFunc = func(hooks RuntimeHooks, args ...Object) (ret Object, err error)

package tengo

// Interop an interface for Go-Tengo interoperability.
type Interop interface {
	// InteropCall can be used to invoke a Tengo function from Go functions.
	InteropCall(callable Object, args ...Object) (ret Object, err error)
}

package objects

// Runtime represents an interface to interact with the currently executing VM
type Runtime interface {
	// Call should call the Object passed in fn with the arguments passed in args.
	// If the function returns successfully, the return value is returned.
	// If the VM encounters an error during execution, the error is returned.
	Call(fn Object, args ...Object) (Object, error)
}

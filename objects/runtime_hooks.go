package objects

// RuntimeHooks provides access to certain functionality
// on the runtime for the currently executing script
type RuntimeHooks interface {
	// Call will call the provided object with the given arguments.
	// Providing an object that is not callable will result in an error.
	// If the vm encounters an error while executing the called function,
	// it will be returned. Otherwise, the return value from the function
	// is returned.
	Call(value Object, args ...Object) (Object, error)
}

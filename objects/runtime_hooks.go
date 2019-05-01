package objects

type RuntimeHooks interface {
	Call(value Object, args ...Object) (Object, error)
}

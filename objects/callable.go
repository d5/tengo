package objects

type Callable interface {
	Call(args ...Object) (ret Object, err error)
}

package objects

var (
	TrueValue      Object = &Bool{Value: true}
	FalseValue     Object = &Bool{Value: false}
	UndefinedValue Object = &Undefined{}
)

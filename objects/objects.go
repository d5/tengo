package objects

var (
	// TrueValue represents a true value.
	TrueValue Object = &Bool{Value: true}

	// FalseValue represents a false value.
	FalseValue Object = &Bool{Value: false}

	// UndefinedValue represents an undefined value.
	UndefinedValue Object = &Undefined{}
)

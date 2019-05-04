package objects

// Spread represents a value that will been spread onto the stack by the VM
type Spread struct {
	ObjectImpl
	Values []Object
}

// TypeName returns the name of the type.
func (s *Spread) TypeName() string { return "spread-values" }

// String returns a string representation of the type's value.
func (s *Spread) String() string { return "<" + s.TypeName() + ">" }

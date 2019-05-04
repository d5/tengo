package objects

type Spread struct {
	ObjectImpl
	Values []Object
}

func (s *Spread) TypeName() string { return "spread-values" }
func (s *Spread) String() string   { return "<" + s.TypeName() + ">" }

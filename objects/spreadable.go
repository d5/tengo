package objects

// Spreadable represents an object that can be spread
type Spreadable interface {
	// Spread should return a slice of Objects
	Spread() []Object
}

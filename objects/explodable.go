package objects

// Explodable represents an object that can be exploded
type Explodable interface {
	// Explode should return a slice of Objects
	Explode() []Object
}

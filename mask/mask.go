package mask

// Mask is the base interface for all collision-shapes.
// It's assigned to an Entity and can check overlaps.
type Mask interface {
	// Parent entity must expose position, origin, width/height, etc.
	SetParent(p Parent)
	// Collide againtst another Mask
	Collide(other Mask) bool
	// Update any internal state (e. g. recalc bounds)
	Update()
}

// Parent is what a Mask needs to know about its Entity.
// You can adapt this to your Entity interface.
type Parent interface {
	X() float32
	Y() float32
	OriginX() float32
	OriginY() float32
	Width() float32
	Height() float32
}

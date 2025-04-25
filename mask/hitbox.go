package mask

// Hitbox is a simple rectangular mask.
type Hitbox struct {
	parent           Parent
	XOff, YOff, W, H float32
}

// NewHitbox creates a rectangle mask with offset & size.
func NewHitbox(xoff, yoff, w, h float32) *Hitbox {
	return &Hitbox{XOff: xoff, YOff: yoff, W: w, H: h}
}

func (m *Hitbox) SetParent(p Parent) {
	m.parent = p
}

func (m *Hitbox) Collide(other Mask) bool {
	// dispatch based on type of other; for simplicity assume other is also *Hitbox
	if o, ok := other.(*Hitbox); ok {
		ax := m.parent.X() + m.XOff
		ay := m.parent.Y() + m.YOff
		bx := o.parent.X() + o.XOff
		by := o.parent.Y() + o.YOff
		return ax+m.W > bx &&
			ay+m.H > by &&
			ax < bx+o.W &&
			ay < by+o.H
	}
	// fallback
	return other.Collide(m)
}

func (m *Hitbox) Update() {
	// nothing to recalc for a simple box
}

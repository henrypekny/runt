package runt

import (
	"github.com/henrypekny/runt/graphics"
	"github.com/henrypekny/runt/mask"
)

// Interp is the current interpolation factor (0–1).
// The Engine sets this each frame before drawing.
var Interp float32

// BaseEntity provides position, layer, visibility, a Graphic,
// an (optional) Mask, and fixed‐timestep interpolation support.
// It also implements mask.Parent so Hitboxes can query its bounds.
type BaseEntity struct {
	// rawPosition is where we store the “true” X,Y each tick.
	// We rename these to avoid colliding with the X() and Y() methods.
	rawX, rawY float32

	// prevRawPosition is used for interpolation.
	prevRawX, prevRawY float32

	// LayerID is the rendering layer.
	LayerID int

	// Visible controls whether we Render.
	Visible bool

	// Graphic is drawn each frame (can be nil).
	Graphic graphics.Graphic

	// Mask, if non‐nil, is used for collision.
	Mask mask.Mask

	// Hitbox dimensions & offset, for mask.Parent methods.
	hitboxX, hitboxY          float32
	hitboxWidth, hitboxHeight float32
}

// NewBaseEntity creates one at (x,y) on the given layer.
// It initializes previous position so interpolation starts from the same point.
func NewBaseEntity(x, y float32, layer int) *BaseEntity {
	return &BaseEntity{
		rawX:         x,
		rawY:         y,
		prevRawX:     x,
		prevRawY:     y,
		LayerID:      layer,
		Visible:      true,
		hitboxWidth:  0, // by default no hitbox
		hitboxHeight: 0,
	}
}

// Snapshot stores the current rawX/rawY into prevRawX/prevRawY.
// Called by Engine once per physics tick before Update().
func (e *BaseEntity) Snapshot() {
	e.prevRawX = e.rawX
	e.prevRawY = e.rawY
}

// Update advances any graphic animations (but not movement).
func (e *BaseEntity) Update(dt float64) {
	if e.Graphic != nil {
		e.Graphic.Update(dt)
	}
}

// Render snaps to integer pixels and draws the Graphic.
// If Interp>0 we interpolate between prev and current.
func (e *BaseEntity) Render() {
	// nothing to draw?
	if !e.Visible || e.Graphic == nil || !e.Graphic.IsVisible() {
		return
	}

	// camera offset
	cx, cy := CurrentWorld.CameraX, CurrentWorld.CameraY

	// choose interpolated or direct position
	var drawX, drawY float32
	if Interp > 0 {
		drawX = e.prevRawX + (e.rawX-e.prevRawX)*Interp
		drawY = e.prevRawY + (e.rawY-e.prevRawY)*Interp
	} else {
		drawX = e.rawX
		drawY = e.rawY
	}

	// if it's an Image, push our computed drawX/drawY into it
	switch g := e.Graphic.(type) {
	case *graphics.Image:
		g.X = drawX
		g.Y = drawY
	case *graphics.Text:
		g.SetPosition(drawX, drawY)
	}

	// finally draw it
	e.Graphic.Render(cx, cy)
}

// Layer implements the runt.Entity interface.
func (e *BaseEntity) Layer() int {
	return e.LayerID
}

// Position accessors: these satisfy mask.Parent.

func (e *BaseEntity) X() float32 {
	return e.rawX
}
func (e *BaseEntity) Y() float32 {
	return e.rawY
}

// Origin for the mask is the hitbox offset.
func (e *BaseEntity) OriginX() float32 {
	return e.hitboxX
}
func (e *BaseEntity) OriginY() float32 {
	return e.hitboxY
}

// Width/Height for the mask.
func (e *BaseEntity) Width() float32 {
	return e.hitboxWidth
}
func (e *BaseEntity) Height() float32 {
	return e.hitboxHeight
}

// SetHitbox installs a rectangular hitbox of size (w,h) with local
// offset (ox,oy).  It creates a new Hitbox mask, sets its parent
// (so it can query X/Y/width/height), and records the bounds locally.
func (e *BaseEntity) SetHitbox(w, h, ox, oy float32) {
	hb := mask.NewHitbox(ox, oy, w, h)
	// Now that BaseEntity implements mask.Parent, this compiles:
	hb.SetParent(e)
	e.Mask = hb

	// store these values so our X(), Y(), Width(), Height() work
	e.hitboxX = ox
	e.hitboxY = oy
	e.hitboxWidth = w
	e.hitboxHeight = h
}

// MoveBy is a helper to adjust rawX/rawY; you can implement your own
// collision‐aware move logic here.
func (e *BaseEntity) MoveBy(dx, dy float32) {
	e.rawX += dx
	e.rawY += dy
}

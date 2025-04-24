// runt/entity.go
package runt

import "github.com/henrypekny/runt/graphics"

// Interp is the current interpolation factor (0–1).
// Engine sets this each frame before drawing.
var Interp float32

// BaseEntity provides X,Y, LayerID and an optional Graphic,
// plus support for fixed-timestep interpolation.
type BaseEntity struct {
	X, Y         float32
	prevX, prevY float32 // for interpolation
	LayerID      int
	Visible      bool
	Graphic      graphics.Graphic
}

// NewBaseEntity creates one at (x,y) on layer.
func NewBaseEntity(x, y float32, layer int) *BaseEntity {
	return &BaseEntity{
		X:       x,
		Y:       y,
		prevX:   x,
		prevY:   y,
		LayerID: layer,
		Visible: true,
	}
}

// Snapshot stores the current X/Y into prevX/prevY.
// Called by Engine once per physics tick before Update().
func (e *BaseEntity) Snapshot() {
	e.prevX = e.X
	e.prevY = e.Y
}

func (e *BaseEntity) Update(dt float64) {
	if e.Graphic != nil {
		e.Graphic.Update(dt)
	}
}

// Render snaps to integer pixels and draws.
// If Interp>0, we linearly interpolate between prev and curr.
func (e *BaseEntity) Render() {
	if !e.Visible || e.Graphic == nil || !e.Graphic.IsVisible() {
		return
	}
	cx, cy := CurrentWorld.CameraX, CurrentWorld.CameraY

	// pick the interpolated or raw position
	var drawX, drawY float32
	if Interp > 0 {
		drawX = e.prevX + (e.X-e.prevX)*Interp
		drawY = e.prevY + (e.Y-e.prevY)*Interp
	} else {
		drawX = e.X
		drawY = e.Y
	}

	// sync into the graphic (if Image)
	if img, ok := e.Graphic.(*graphics.Image); ok {
		img.X = drawX
		img.Y = drawY
	}

	e.Graphic.Render(cx, cy)
}

func (e *BaseEntity) Layer() int {
	return e.LayerID
}

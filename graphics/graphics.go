package graphics

// Graphic is the minimal interface that all renderable things must satisfy.
type Graphic interface {
	// Update is called once per‐frame before Render.
	Update(dt float64)
	// Render draws the graphic, given the world camera offset.
	Render(cameraX, cameraY float32)
	// Visible controls whether Render actually does anything.
	SetVisible(bool)
	IsVisible() bool
}

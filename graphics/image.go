package graphics

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Image is a non-animated texture with position, origin, scale, rotation,
// tint (Color) & parallax (ScrollX/Y).
type Image struct {
	Texture rl.Texture2D // GPU texture handle

	// World position
	X, Y float32

	// Per-axis scale factors (multiplied by Scale)
	ScaleX, ScaleY float32
	// Uniform extra scale
	Scale float32

	// Rotation in degrees around its center
	Rotation float32

	// Parallax scrolling factors (1 == follow camera exactly)
	ScrollX, ScrollY float32

	// Which sub-rectangle of the texture to draw
	SrcRec rl.Rectangle

	// Tint color & alpha override
	Color rl.Color

	// Visibility flag
	visible bool
}

// NewImage loads the image at `path`, uploads it to the GPU with point-filtering,
// and returns an Image whose pivot is automatically its center.
func NewImage(path string) *Image {
	// 1) load CPU-side image and upload
	imgCPU := rl.LoadImage(path)
	tex := rl.LoadTextureFromImage(imgCPU)
	rl.UnloadImage(imgCPU)

	// 2) enforce nearest-neighbour filtering
	rl.SetTextureFilter(tex, rl.FilterPoint)

	// 3) wrap in our Image struct
	w := float32(tex.Width)
	h := float32(tex.Height)
	return &Image{
		Texture: tex,
		SrcRec:  rl.NewRectangle(0, 0, w, h),
		ScaleX:  1, ScaleY: 1,
		Scale:    1,
		Rotation: 0,
		ScrollX:  1, ScrollY: 1,
		Color:   rl.White,
		visible: true,
	}
}

// NewImageFromTexture wraps an existing Texture2D in an Image,
// re-applying point-filter for consistency and centering pivot.
func NewImageFromTexture(tex rl.Texture2D) *Image {
	rl.SetTextureFilter(tex, rl.FilterPoint)
	w := float32(tex.Width)
	h := float32(tex.Height)
	return &Image{
		Texture: tex,
		SrcRec:  rl.NewRectangle(0, 0, w, h),
		ScaleX:  1, ScaleY: 1,
		Scale:    1,
		Rotation: 0,
		ScrollX:  1, ScrollY: 1,
		Color:   rl.White,
		visible: true,
	}
}

// Update is a no-op for static Images.
func (img *Image) Update(dt float64) {}

// IsVisible reports current visibility.
func (img *Image) IsVisible() bool { return img.visible }

// SetVisible toggles drawing.
func (img *Image) SetVisible(v bool) { img.visible = v }

// Render draws the Image at its world position, rotating & scaling around
// the center of the sprite.  camX,camY are the camera offsets.
func (img *Image) Render(camX, camY float32) {
	if !img.visible {
		return
	}

	// destination size
	w := img.SrcRec.Width * img.ScaleX * img.Scale
	h := img.SrcRec.Height * img.ScaleY * img.Scale

	// world-space draw position: pivot is drawn at (img.X,img.Y)
	dstX := img.X - camX*img.ScrollX
	dstY := img.Y - camY*img.ScrollY
	dst := rl.NewRectangle(dstX, dstY, w, h)

	// pivot inside that quad is its center
	origin := rl.NewVector2(w/2, h/2)

	// full-precision draw with rotation & scale
	rl.DrawTexturePro(
		img.Texture,
		img.SrcRec,
		dst,
		origin,
		img.Rotation,
		img.Color,
	)
}

// NewCircle creates a filled circle Image (transparent outside).
func NewCircle(radius int, col rl.Color) *Image {
	size := radius * 2
	img := rl.GenImageColor(size, size, rl.NewColor(0, 0, 0, 0))
	rl.ImageDrawCircle(img, int32(radius), int32(radius), int32(radius), col)
	tex := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	return NewImageFromTexture(tex)
}

// NewGradientLinear builds a CPU-side linear gradient and wraps it.
func NewGradientLinear(
	width, height int,
	direction int, // 0=vertical, 90=horizontal, etc.
	start, end rl.Color,
) *Image {
	img := rl.GenImageGradientLinear(width, height, direction, start, end)
	tex := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	return NewImageFromTexture(tex)
}

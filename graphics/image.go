package graphics

import rl "github.com/gen2brain/raylib-go/raylib"

// Image is a non-animated texture with position, origin, scale, rotation,
// tint (Color) & parallax (ScrollX/Y).
type Image struct {
	// GPU texture handle
	Texture rl.Texture2D

	// World position
	X, Y float32

	// Local transform origin (rotation/scale pivot)
	OriginX, OriginY float32

	// Per-axis scale factors (multiplied by Scale)
	ScaleX, ScaleY float32
	// Uniform extra scale
	Scale float32

	// Rotation in degrees around (OriginX,OriginY)
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

// NewImage loads a texture from disk, applies nearest-neighbour filtering,
// and initializes all transforms to their defaults.
func NewImage(path string) *Image {
	tex := rl.LoadTexture(path)
	rl.SetTextureFilter(tex, rl.FilterPoint) // pixel-perfect by default

	return &Image{
		Texture:  tex,
		SrcRec:   rl.NewRectangle(0, 0, float32(tex.Width), float32(tex.Height)),
		Color:    rl.White, // no tint, fully opaque
		ScaleX:   1,
		ScaleY:   1,
		Scale:    1,
		Rotation: 0,
		ScrollX:  1,
		ScrollY:  1,
		visible:  true,
	}
}

// Update is a no-op for static Images, but satisfies the common Graphic API.
func (img *Image) Update(dt float64) {}

// Render draws the Image at its world position, snapping to integer
// pixels for crispness, then applying scale & rotation if needed.
// camX/camY is the current camera offset in world space.
func (img *Image) Render(camX, camY float32) {
	if !img.visible {
		return
	}

	// 1) world-space position including parallax:
	rawX := img.X - img.OriginX - camX*img.ScrollX
	rawY := img.Y - img.OriginY - camY*img.ScrollY

	// 2) snap to integer for pixel perfection
	dstX := float32(int(rawX + 0.5))
	dstY := float32(int(rawY + 0.5))

	// 3) compute destination rectangle size
	w := img.SrcRec.Width * img.ScaleX * img.Scale
	h := img.SrcRec.Height * img.ScaleY * img.Scale
	dst := rl.NewRectangle(dstX, dstY, w, h)

	// 4) origin vector for DrawTexturePro
	origin := rl.NewVector2(img.OriginX, img.OriginY)

	// 5) if no rotation & no extra scale, use the faster DrawTextureRec
	effScaleX := img.ScaleX * img.Scale
	effScaleY := img.ScaleY * img.Scale
	if img.Rotation == 0 && effScaleX == 1 && effScaleY == 1 {
		rl.DrawTextureRec(
			img.Texture,
			img.SrcRec,
			rl.NewVector2(dstX, dstY),
			img.Color,
		)
	} else {
		// otherwise handle transform via DrawTexturePro
		rl.DrawTexturePro(
			img.Texture,
			img.SrcRec,
			dst,
			origin,
			img.Rotation,
			img.Color,
		)
	}
}

// SetVisible toggles whether this Image will be drawn.
func (img *Image) SetVisible(v bool) { img.visible = v }

// IsVisible reports the current visibility state.
func (img *Image) IsVisible() bool { return img.visible }

// NewImageFromTexture wraps an existing Texture2D in an Image,
// applying point filtering on it.
func NewImageFromTexture(tex rl.Texture2D) *Image {
	rl.SetTextureFilter(tex, rl.FilterPoint)
	return &Image{
		Texture:  tex,
		SrcRec:   rl.NewRectangle(0, 0, float32(tex.Width), float32(tex.Height)),
		Color:    rl.White,
		ScaleX:   1,
		ScaleY:   1,
		Scale:    1,
		Rotation: 0,
		ScrollX:  1,
		ScrollY:  1,
		visible:  true,
	}
}

// NewRect creates a solid-color rectangle Image of the given size.
func NewRect(width, height int, col rl.Color) *Image {
	// generate a CPU-side image filled with col
	img := rl.GenImageColor(width, height, col)
	tex := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	return NewImageFromTexture(tex)
}

// NewCircle creates a filled circle Image (transparent outside).
func NewCircle(radius int, col rl.Color) *Image {
	size := radius * 2
	// start with a fully transparent image
	img := rl.GenImageColor(size, size, rl.NewColor(0, 0, 0, 0))
	// draw a solid circle into it
	rl.ImageDrawCircle(img,
		int32(radius), int32(radius), int32(radius),
		col,
	)
	tex := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	return NewImageFromTexture(tex)
}

// NewGradientLinear builds a CPU-side linear gradient at any angle (deg),
// uploads it as a Texture2D, then wraps in our Image type.
func NewGradientLinear(
	width, height int,
	direction int, // 0 = vertical, 90 = left→right, etc.
	start, end rl.Color,
) *Image {
	// 1) generate CPU image
	gradImg := rl.GenImageGradientLinear(
		width,
		height,
		direction,
		start,
		end,
	)
	// 2) upload & free CPU image
	tex := rl.LoadTextureFromImage(gradImg)
	rl.UnloadImage(gradImg)

	// 3) wrap & return
	return NewImageFromTexture(tex)
}

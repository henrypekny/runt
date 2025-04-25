// runt/graphics/text.go
package graphics

import (
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/henrypekny/runt/loader"
)

// Align controls horizontal positioning of each line.
type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

// Text draws one or more lines of text with optional word-wrap and alignment.
type Text struct {
	font    rl.Font
	content string
	x, y    float32
	size    float32
	spacing float32
	color   rl.Color

	wordWrap bool
	maxWidth float32
	align    Align

	visible bool
}

// NewText creates a Text at (x,y).  It asks the loader for
// “VT323-Regular.ttf” at the requested size, falling back
// to the default font if needed.
func NewText(s string, x, y, size float32, c rl.Color) *Text {
	// loader.LoadFont will search all your dev/asset paths,
	// load+cache font, and set TextureFilter to POINT for you.
	fnt := loader.LoadFont("VT323-Regular.ttf", int32(size))

	return &Text{
		font:     fnt,
		content:  s,
		x:        x,
		y:        y,
		size:     size,
		spacing:  1,
		color:    c,
		wordWrap: false,
		maxWidth: 0,
		align:    AlignLeft,
		visible:  true,
	}
}

func (t *Text) SetWordWrap(on bool, maxWidth float32) { t.wordWrap, t.maxWidth = on, maxWidth }
func (t *Text) SetAlign(a Align)                      { t.align = a }
func (t *Text) SetVisible(v bool)                     { t.visible = v }
func (t *Text) Update(dt float64)                     {}
func (t *Text) IsVisible() bool                       { return t.visible }

// Render draws each line, applying camera offset, wrap and alignment.
func (t *Text) Render(camX, camY float32) {
	if !t.visible {
		return
	}
	x0, y0 := t.x-camX, t.y-camY
	lines := strings.Split(t.content, "\n")

	if t.wordWrap && t.maxWidth > 0 {
		var wrapped []string
		for _, line := range lines {
			wrapped = append(wrapped, wrapLine(line, t.font, t.size, t.spacing, t.maxWidth)...)
		}
		lines = wrapped
	}

	for i, line := range lines {
		meas := rl.MeasureTextEx(t.font, line, t.size, t.spacing)
		var dx float32
		switch t.align {
		case AlignCenter:
			dx = (t.maxWidth - meas.X) / 2
		case AlignRight:
			dx = t.maxWidth - meas.X
		}
		pos := rl.Vector2{X: x0 + dx, Y: y0 + float32(i)*(meas.Y+t.spacing)}
		rl.DrawTextEx(t.font, line, pos, t.size, t.spacing, t.color)
	}
}

func (t *Text) Width() float32 {
	var max float32
	for _, line := range strings.Split(t.content, "\n") {
		if w := rl.MeasureTextEx(t.font, line, t.size, t.spacing).X; w > max {
			max = w
		}
	}
	return max
}

func (t *Text) Height() float32 {
	h := rl.MeasureTextEx(t.font, "M", t.size, t.spacing).Y
	lines := float32(len(strings.Split(t.content, "\n")))
	return h*lines + t.spacing*(lines-1)
}

func (t *Text) SetText(s string) { t.content = s }

// wrapLine splits a single line into multiple so none exceed maxWidth.
func wrapLine(s string, font rl.Font, size, spacing, maxWidth float32) []string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}
	var out []string
	cur := ""
	for _, w := range words {
		next := w
		if cur != "" {
			next = cur + " " + w
		}
		if rl.MeasureTextEx(font, next, size, spacing).X > maxWidth {
			out = append(out, cur)
			cur = w
		} else {
			cur = next
		}
	}
	out = append(out, cur)
	return out
}

// SetPosition moves the Text to (x,y) in world space.
// You can call this from BaseEntity.Render if you wrap Text in an Entity.
func (t *Text) SetPosition(x, y float32) {
	t.x = x
	t.y = y
}

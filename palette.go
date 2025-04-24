// runt/palette.go
package runt

import rl "github.com/gen2brain/raylib-go/raylib"

// Color is our alias for rl.Color; users of runt only ever see runt.Color.
type Color = rl.Color

// NewColor constructs a new Color without ever importing rl in user code.
func NewColor(r, g, b, a uint8) Color {
	return rl.NewColor(r, g, b, a)
}

// Endesga-16 palette with easy names:
var (
	Sand     = NewColor(0xE4, 0xA6, 0x72, 0xFF) // light tan
	Rust     = NewColor(0xB8, 0x6F, 0x50, 0xFF) // muted brown
	Chestnut = NewColor(0x74, 0x3F, 0x39, 0xFF) // dark reddish-brown
	Charcoal = NewColor(0x3F, 0x28, 0x32, 0xFF) // near-black gray

	Crimson = NewColor(0x9E, 0x28, 0x35, 0xFF) // deep red
	Scarlet = NewColor(0xE5, 0x3B, 0x44, 0xFF) // vivid red
	Amber   = NewColor(0xFB, 0x92, 0x2B, 0xFF) // warm orange
	Mustard = NewColor(0xFF, 0xE7, 0x62, 0xFF) // golden yellow

	Lime   = NewColor(0x63, 0xC6, 0x4D, 0xFF) // bright green
	Forest = NewColor(0x32, 0x73, 0x45, 0xFF) // dark green
	Teal   = NewColor(0x19, 0x3D, 0x3F, 0xFF) // deep teal
	Slate  = NewColor(0x4F, 0x67, 0x81, 0xFF) // muted blue-gray

	Sky   = NewColor(0xAF, 0xBF, 0xD2, 0xFF) // pale sky blue
	White = NewColor(0xFF, 0xFF, 0xFF, 0xFF) // pure white
	Cyan  = NewColor(0x2C, 0xE8, 0xF4, 0xFF) // bright cyan
	Azure = NewColor(0x04, 0x84, 0xD1, 0xFF) // vivid blue
)

// Endesga16 is the ordered slice of the above colors.
var Endesga16 = []Color{
	Sand, Rust, Chestnut, Charcoal,
	Crimson, Scarlet, Amber, Mustard,
	Lime, Forest, Teal, Slate,
	Sky, White, Cyan, Azure,
}

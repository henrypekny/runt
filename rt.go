package runt

import rl "github.com/gen2brain/raylib-go/raylib"

// --------------------------------------------------------------------------------
// GLOBAL STATE (the “FP” bag of statics)
// --------------------------------------------------------------------------------

var (
	// logical viewport size
	Width, Height         int
	HalfWidth, HalfHeight float32

	// timing & framerate
	AssignedFPS  int
	FrameRate    float64
	Fixed        bool
	TimeInFrames bool
	Rate         float64 = 1.0
	Elapsed      float64

	// the active World
	CurrentWorld *World

	// default background clear color
	BackgroundColor = rl.Color{R: 30, G: 30, B: 30, A: 255}
)

func init() {
	// seed CurrentWorld so it’s never nil
	CurrentWorld = NewWorld()
}

// Resize updates the logical viewport globals.
func Resize(w, h int) {
	Width = w
	Height = h
	HalfWidth = float32(w) / 2
	HalfHeight = float32(h) / 2
}

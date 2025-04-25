// runt/rt.go
package runt

import (
	"math"
	"math/rand"
	"time"
)

// -----------------------------------------------------------------------------
// Global state (mirrors FP.as)
// -----------------------------------------------------------------------------

// Screen dimensions (set via Resize).
var (
	Width, Height         int     // full resolution
	HalfWidth, HalfHeight float32 // half resolution
)

// Timing
var (
	Fixed        bool        // fixed vs. variable timestep
	TimeInFrames bool        // if true, Elapsed is in frames rather than seconds
	FrameRate    float64     // measured frames per second
	AssignedFPS  int         // target FPS
	Elapsed      float64     // time since last frame (seconds)
	Rate         float64 = 1 // timescale multiplier for Elapsed
)

// BackgroundColor is the default clear‐screen color.  Engine uses this in ClearBackground.
var BackgroundColor Color = Charcoal

// CurrentWorld is the active World.  Engine and your Game code will call
// CurrentWorld.Update/Render/Entities() etc.
var CurrentWorld *World

// Camera offset
var CameraX, CameraY float32

// Random seed state
var (
	_seed        int64     = time.Now().UnixNano() & 0x7FFFFFFF
	LastTimeFlag time.Time = time.Now()
)

// -----------------------------------------------------------------------------
// Initialization & camera
// -----------------------------------------------------------------------------

// Resize sets the virtual screen size and recalculates half-width/height.
func Resize(w, h int) {
	Width, Height = w, h
	HalfWidth = float32(w) / 2
	HalfHeight = float32(h) / 2
}

// SetCamera moves the camera offset.
func SetCamera(x, y float32) {
	CameraX, CameraY = x, y
}

// ResetCamera zeroes out the camera.
func ResetCamera() {
	CameraX, CameraY = 0, 0
}

// -----------------------------------------------------------------------------
// Time utilities
// -----------------------------------------------------------------------------

// TimeFlag returns the milliseconds since the last TimeFlag() call.
func TimeFlag() time.Duration {
	now := time.Now()
	diff := now.Sub(LastTimeFlag)
	LastTimeFlag = now
	return diff
}

// -----------------------------------------------------------------------------
// Randomness (port of FP.rand, FP.random, etc.)
// -----------------------------------------------------------------------------

// RandSeed sets the PRNG seed.
func RandSeed(seed int64) {
	if seed <= 0 {
		seed = 1
	}
	_seed = seed % 2147483647
	rand.Seed(_seed)
}

// Random returns a float64 in [0,1).
func Random() float64 {
	// linear congruential for demonstration (or use math/rand)
	_seed = (_seed * 16807) % 2147483647
	return float64(_seed) / 2147483647
}

// Rand returns an int in [0, n).
func Rand(n int) int {
	return int(Random() * float64(n))
}

// RandomizeSeed picks a new seed from time.
func RandomizeSeed() {
	RandSeed(time.Now().UnixNano())
}

// -----------------------------------------------------------------------------
// Slice & math utilities
// -----------------------------------------------------------------------------

// RemoveElement removes the first occurrence of `el` from slice `s`. Returns true if removed.
func RemoveElement[T comparable](s *[]T, el T) bool {
	for i, v := range *s {
		if v == el {
			*s = append((*s)[:i], (*s)[i+1:]...)
			return true
		}
	}
	return false
}

// Choose picks one element at random.
func Choose[T any](options ...T) T {
	return options[Rand(len(options))]
}

// Sign returns 1, -1 or 0 depending on x’s sign.
func Sign(x float64) int {
	if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	}
	return 0
}

// Approach moves `value` towards `target` by at most `amt`.
func Approach(value, target, amt float64) float64 {
	if value < target-amt {
		return value + amt
	} else if value > target+amt {
		return value - amt
	}
	return target
}

// Lerp linearly interpolates between a and b by t in [0,1].
func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

// ColorLerp blends two Colors by t in [0,1].
func ColorLerp(c1, c2 Color, t float32) Color {
	if t <= 0 {
		return c1
	}
	if t >= 1 {
		return c2
	}
	r := uint8(float32(c1.R) + (float32(c2.R)-float32(c1.R))*t)
	g := uint8(float32(c1.G) + (float32(c2.G)-float32(c1.G))*t)
	b := uint8(float32(c1.B) + (float32(c2.B)-float32(c1.B))*t)
	a := uint8(float32(c1.A) + (float32(c2.A)-float32(c1.A))*t)
	return NewColor(r, g, b, a)
}

// Distance between two points.
func Distance(x1, y1, x2, y2 float64) float64 {
	dx, dy := x2-x1, y2-y1
	return math.Hypot(dx, dy)
}

// DistanceRects returns separation between two rectangles (0 if overlapping).
func DistanceRects(x1, y1, w1, h1, x2, y2, w2, h2 float64) float64 {
	// axis‐aligned rectangle distance (exact port of FP.distanceRects)
	if x1 < x2+w2 && x2 < x1+w1 {
		if y1 < y2+h2 && y2 < y1+h1 {
			return 0
		}
		if y1 > y2 {
			return y1 - (y2 + h2)
		}
		return y2 - (y1 + h1)
	}
	if y1 < y2+h2 && y2 < y1+h1 {
		if x1 > x2 {
			return x1 - (x2 + w2)
		}
		return x2 - (x1 + w1)
	}
	// diagonal case
	cx1 := x1 + w1/2
	cy1 := y1 + h1/2
	cx2 := x2 + w2/2
	cy2 := y2 + h2/2
	return Distance(cx1, cy1, cx2, cy2)
}

// DistanceRectPoint returns separation point→rect (0 if inside).
func DistanceRectPoint(px, py, rx, ry, rw, rh float64) float64 {
	// port of FP.distanceRectPoint
	if px >= rx && px <= rx+rw {
		if py >= ry && py <= ry+rh {
			return 0
		}
		if py > ry {
			return py - (ry + rh)
		}
		return ry - py
	}
	if py >= ry && py <= ry+rh {
		if px > rx {
			return px - (rx + rw)
		}
		return rx - px
	}
	// corner → use Euclidean
	cx := rx + rw/2
	cy := ry + rh/2
	return Distance(px, py, cx, cy)
}

// Clamp clamps x to [min,max].
func Clamp(x, min, max float64) float64 {
	if min > max {
		min, max = max, min
	}
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

// ClampInRect clamps (x,y) into the given box with optional padding.
func ClampInRect(x, y, rx, ry, rw, rh, pad float64) (float64, float64) {
	return Clamp(x, rx+pad, rx+rw-pad), Clamp(y, ry+pad, ry+rh-pad)
}

// Scale linearly remaps value from [min,max] → [min2,max2].
func Scale(value, min, max, min2, max2 float64) float64 {
	return min2 + ((value-min)/(max-min))*(max2-min2)
}

// ScaleClamp remaps and clamps into [min2,max2].
func ScaleClamp(value, min, max, min2, max2 float64) float64 {
	v := Scale(value, min, max, min2, max2)
	return Clamp(v, min2, max2)
}

// Next/Prev cycle through a slice of options.
func Next[T comparable](current T, options []T, loop bool) T {
	for i, v := range options {
		if v == current {
			if loop {
				return options[(i+1)%len(options)]
			}
			if i+1 < len(options) {
				return options[i+1]
			}
		}
	}
	return current
}
func Prev[T comparable](current T, options []T, loop bool) T {
	for i, v := range options {
		if v == current {
			if loop {
				return options[(i-1+len(options))%len(options)]
			}
			if i-1 >= 0 {
				return options[i-1]
			}
		}
	}
	return current
}

// Swap switches between a and b.
func Swap[T comparable](current, a, b T) T {
	if current == a {
		return b
	}
	return a
}

// Shuffle randomly permutes a slice.
func Shuffle[T any](s []T) {
	for i := range s {
		j := Rand(i + 1)
		s[i], s[j] = s[j], s[i]
	}
}

// Frames returns a slice of indices from ‘from’ to ‘to’, stepping by skip+1.
func Frames(from, to, skip int) []int {
	step := skip + 1
	var seq []int
	if from <= to {
		for x := from; x <= to; x += step {
			seq = append(seq, x)
		}
	} else {
		for x := from; x >= to; x -= step {
			seq = append(seq, x)
		}
	}
	return seq
}

// Sort and SortBy would just wrap sort.Slice or sort.SliceStable as needed.
// …

// Color construction helpers (port of FP.getColorRGB/HSV etc.)

// GetColorRGB returns a 24-bit color.
func GetColorRGB(r, g, b uint8) Color {
	return NewColor(r, g, b, 0xFF)
}

// (HSV and channel extractors can be added similarly…)
// -----------------------------------------------------------------------------

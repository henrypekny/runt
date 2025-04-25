package runt

import (
	"fmt"
	"math"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Game is your application’s entrypoint interface.
//
//	– Create()   is called once at startup.
//	– Update(dt) is called one or more times per frame (dt = seconds).
//	– Draw(interp) is called every frame; interp is 0–1 in fixed mode, always 0 in variable mode.
type Game interface {
	Create()
	Update(dt float64)
	Draw(interp float32)
}

// Engine drives the window, main loop, timing and background.
// It supports both fixed‐timestep (with interpolation) and variable‐timestep modes.
type Engine struct {
	game         Game          // the user’s Game implementation
	title        string        // window title
	fps          int           // desired frame rate (and tick rate in fixed mode)
	bg           Color         // clear color for the backbuffer (alias for rl.Color)
	fixed        bool          // true → fixed‐timestep + interpolation
	tickRate     time.Duration // time per physics tick (1/fps)
	maxElapsed   float64       // clamp on dt to avoid spiral-of-death
	maxFrameSkip int           // max physics steps per frame
	paused       bool          // when true, Update(dt) is skipped
}

// NewEngine constructs an Engine but does not open the window.
//
//	w, h   – window dimensions
//	title  – window title
//	fps    – target framerate (also physics tick rate in fixed mode)
//	game   – your Game instance
//	fixed  – whether to use fixed‐timestep + interpolation
func NewEngine(
	w, h int,
	title string,
	fps int,
	game Game,
	fixed bool,
) *Engine {
	// Configure our package-level state (screen size, FPS)
	Resize(w, h)
	AssignedFPS = fps

	return &Engine{
		game:         game,
		title:        title,
		fps:          fps,
		bg:           BackgroundColor, // default from palette.go
		fixed:        fixed,           // fixed vs variable
		tickRate:     time.Second / time.Duration(fps),
		maxElapsed:   1.0 / 10.0, // clamp dt at 100ms
		maxFrameSkip: 5,          // avoid too many physics steps
		paused:       false,      // start unpaused
	}
}

// Run opens the window, initializes audio, and enters the main loop.
// It handles timing, update, interpolation, and drawing.
func (e *Engine) Run() {
	// --- Initialize Raylib ---
	rl.InitWindow(int32(Width), int32(Height), e.title)
	defer rl.CloseWindow()
	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()
	rl.SetTargetFPS(int32(e.fps))

	// Let the Game set itself up.
	e.game.Create()

	// Buffer for dt statistics.
	const sampleCount = 120
	dts := make([]float64, 0, sampleCount)

	previous := rl.GetTime()
	var lag float64

	// Main loop.
	for !rl.WindowShouldClose() {
		// ---- 1) Measure Δt ----
		now := rl.GetTime()
		dt := now - previous
		previous = now

		// Clamp dt to avoid spiral-of-death.
		if dt > e.maxElapsed {
			dt = e.maxElapsed
		}

		// Collect stats.
		dts = append(dts, dt)
		if len(dts) > sampleCount {
			dts = dts[1:]
		}
		if len(dts) == sampleCount {
			// Print min/max/mean/stddev every sampleCount frames.
			min, max, sum := dts[0], dts[0], 0.0
			for _, v := range dts {
				if v < min {
					min = v
				}
				if v > max {
					max = v
				}
				sum += v
			}
			mean := sum / sampleCount
			varVariance := 0.0
			for _, v := range dts {
				d := v - mean
				varVariance += d * d
			}
			stddev := math.Sqrt(varVariance / sampleCount)
			fmt.Printf("[%s] dt over last %d frames: min=%.5f max=%.5f mean=%.5f σ=%.5f\n",
				time.Now().Format("15:04:05"), sampleCount, min, max, mean, stddev)
			dts = dts[:0]
		}

		// Apply any global time‐scale.
		Elapsed = dt * Rate

		// ---- 2) Update ----
		if !e.paused {
			if e.fixed {
				// Fixed‐timestep mode.
				lag += dt
				step := e.tickRate.Seconds()

				// Cap physics steps.
				if lag > step*float64(e.maxFrameSkip) {
					lag = step * float64(e.maxFrameSkip)
				}

				// Snapshot all entities for interpolation.
				for _, ent := range CurrentWorld.Entities() {
					if s, ok := ent.(interface{ Snapshot() }); ok {
						s.Snapshot()
					}
				}

				// Run fixed‐size physics steps.
				for lag >= step {
					e.game.Update(step)
					lag -= step
				}
			} else {
				// Variable‐timestep mode.
				e.game.Update(dt)
			}
		}

		// ---- 3) Render ----
		rl.BeginDrawing()
		rl.ClearBackground(e.bg) // Color is our alias for rl.Color

		// Camera transform.
		camX, camY := CurrentWorld.CameraX, CurrentWorld.CameraY
		rl.BeginMode2D(rl.NewCamera2D(
			rl.Vector2{X: camX, Y: camY},
			rl.Vector2{X: 0, Y: 0},
			0, 1,
		))

		// Draw with interpolation factor.
		if e.fixed {
			alpha := float32(lag / e.tickRate.Seconds())
			Interp = alpha
			e.game.Draw(alpha)
		} else {
			Interp = 0
			e.game.Draw(0)
		}

		rl.EndMode2D()
		rl.EndDrawing()

		// ---- 4) FPS ----
		FrameRate = float64(rl.GetFPS())
	}
}

// SetBackground updates the clear color at runtime.
func (e *Engine) SetBackground(c Color) {
	e.bg = c
}

// Pause suspends further Update(dt) calls until Resume() is called.
func (e *Engine) Pause() {
	e.paused = true
}

// Resume re‐enables Update(dt) calls.
func (e *Engine) Resume() {
	e.paused = false
}

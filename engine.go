// runt/engine.go
package runt

import (
	"fmt"
	"math"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Game is your application’s entrypoint interface.
//   - Create() is called once at startup.
//   - Update(dt) is called one or more times per frame (dt = seconds).
//   - Draw(interp) is called every frame; interp is 0–1 in fixed mode, always 0 in variable mode.
type Game interface {
	Create()
	Update(dt float64)
	Draw(interp float32)
}

// Engine drives the window, main loop, timing and background.
type Engine struct {
	game         Game
	title        string
	fps          int
	bg           rl.Color
	fixed        bool
	tickRate     time.Duration
	maxElapsed   float64 // clamp on dt to avoid spiral-of-death
	maxFrameSkip int
	paused       bool
}

// NewEngine builds (but does not open) a windowed engine.
// Pass fixed=true to enable fixed-timestep + interpolation.
func NewEngine(
	w, h int,
	title string,
	fps int,
	game Game,
	fixed bool,
) *Engine {
	Resize(w, h)
	AssignedFPS = fps

	return &Engine{
		game:         game,
		title:        title,
		fps:          fps,
		bg:           BackgroundColor,
		fixed:        fixed,
		tickRate:     time.Second / time.Duration(fps),
		maxElapsed:   1.0 / 10.0, // clamp dt to 100ms
		maxFrameSkip: 5,
		paused:       false,
	}
}

// Run opens the window, initializes audio, and enters the main loop.
func (e *Engine) Run() {
	rl.InitWindow(int32(Width), int32(Height), e.title)
	defer rl.CloseWindow()

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	rl.SetTargetFPS(int32(e.fps))

	// allow the game to set up
	e.game.Create()

	// rolling buffer for Δt statistics
	const sampleCount = 120
	dts := make([]float64, 0, sampleCount)
	frame := 0

	previous := rl.GetTime()
	var lag float64

	for !rl.WindowShouldClose() {
		// --- 1) measure Δt ---
		now := rl.GetTime()
		dt := now - previous
		previous = now

		if dt > e.maxElapsed {
			dt = e.maxElapsed
		}
		// record for stats
		dts = append(dts, dt)
		if len(dts) > sampleCount {
			dts = dts[1:]
		}
		frame++
		if frame%sampleCount == 0 && len(dts) == sampleCount {
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
				diff := v - mean
				varVariance += diff * diff
			}
			stddev := math.Sqrt(varVariance / sampleCount)
			fmt.Printf("[%s] dt over last %d frames: min=%.5f  max=%.5f  mean=%.5f  σ=%.5f\n",
				time.Now().Format("15:04:05"), sampleCount, min, max, mean, stddev)
		}

		// --- 2) update ---
		Elapsed = dt * Rate
		if !e.paused {
			if e.fixed {
				// accumulate unprocessed time
				lag += dt
				step := e.tickRate.Seconds()

				// avoid spiral
				if lag > step*float64(e.maxFrameSkip) {
					lag = step * float64(e.maxFrameSkip)
				}

				// take a snapshot for interpolation
				for _, ent := range CurrentWorld.Entities() {
					if s, ok := ent.(interface{ Snapshot() }); ok {
						s.Snapshot()
					}
				}

				// run fixed-size physics steps
				for lag >= step {
					e.game.Update(step)
					lag -= step
				}
			} else {
				// variable-timestep
				e.game.Update(dt)
			}
		}

		// --- 3) render ---
		rl.BeginDrawing()
		rl.ClearBackground(e.bg)

		camX, camY := CurrentWorld.CameraX, CurrentWorld.CameraY
		rl.BeginMode2D(rl.NewCamera2D(
			rl.Vector2{X: camX, Y: camY},
			rl.Vector2{X: 0, Y: 0},
			0, 1,
		))

		if e.fixed {
			// pass interpolation α = lag/step to Draw
			alpha := float32(lag / e.tickRate.Seconds())
			Interp = alpha
			e.game.Draw(alpha)
		} else {
			// no interpolation
			Interp = 0
			e.game.Draw(0)
		}

		rl.EndMode2D()
		rl.EndDrawing()

		// --- 4) FPS ---
		FrameRate = float64(rl.GetFPS())
	}
}

// SetBackground changes the clear color at runtime.
func (e *Engine) SetBackground(c rl.Color) {
	e.bg = c
}

// Pause suspends calls to Update(dt).
func (e *Engine) Pause() {
	e.paused = true
}

// Resume re-enables calls to Update(dt).
func (e *Engine) Resume() {
	e.paused = false
}

package runt

import (
	"sort"
)

// Entity must implement Update, Render and Layer.
type Entity interface {
	Update(dt float64)
	Render()
	Layer() int
}

// World holds and updates/draws a set of Entities.
type World struct {
	entities         []Entity
	CameraX, CameraY float32
	UseCamera        bool
}

// NewWorld makes an empty World.
func NewWorld() *World {
	return &World{
		entities:  make([]Entity, 0, 16),
		UseCamera: false,
	}
}

// Add a new Entity.
func (w *World) Add(e Entity) {
	w.entities = append(w.entities, e)
}

// Update calls each Entity.Update.
func (w *World) Update(dt float64) {
	for _, e := range w.entities {
		e.Update(dt)
	}
}

// Render sorts by layer then draws each Entity.Render.
func (w *World) Render() {
	sort.SliceStable(w.entities, func(i, j int) bool {
		return w.entities[i].Layer() < w.entities[j].Layer()
	})
	for _, e := range w.entities {
		e.Render()
	}
}

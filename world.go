package runt

import (
	"reflect"
	"sort"
)

// Entity must implement Update, Render and Layer.
type Entity interface {
	Update(dt float64)
	Render()
	Layer() int
}

// World holds and updates/draws a set of Entities in layers, with
// proper add/remove queues, camera support, and z‐order operations.
type World struct {
	// active entities, grouped by layer
	layers     map[int][]Entity
	layerOrder []int // distinct layers, kept sorted

	// staging queues for safe add/remove during update
	addQueue    []Entity
	removeQueue []Entity

	// optional camera
	CameraX, CameraY float32
	UseCamera        bool

	// simple recycling pool: type name -> []Entity
	pool map[string][]Entity

	// fast‐lookup counts by type string
	typeCounts map[string]int
}

// Entities returns a flat slice of all Entities in this World,
// in ascending layer → insertion order.
func (w *World) Entities() []Entity {
	// Make sure any pending adds/removes are applied.
	w.FlushQueues()

	// First compute how many total entities we have, so we can pre-alloc.
	total := 0
	for _, layer := range w.layerOrder {
		total += len(w.layers[layer])
	}

	// Build a single slice, in layer-order.
	all := make([]Entity, 0, total)
	for _, layer := range w.layerOrder {
		all = append(all, w.layers[layer]...)
	}
	return all
}

// NewWorld makes an empty World.
func NewWorld() *World {
	return &World{
		layers:      make(map[int][]Entity),
		layerOrder:  make([]int, 0, 8),
		addQueue:    make([]Entity, 0, 16),
		removeQueue: make([]Entity, 0, 16),
		pool:        make(map[string][]Entity),
		typeCounts:  make(map[string]int),
	}
}

// Add queues an Entity for addition at the end of this frame.
func (w *World) Add(e Entity) {
	w.addQueue = append(w.addQueue, e)
}

// Remove queues an Entity for removal at the end of this frame.
func (w *World) Remove(e Entity) {
	w.removeQueue = append(w.removeQueue, e)
}

// FlushQueues integrates all queued add/removes.
// Call this once per frame (e.g. at end of Update or start of Render).
func (w *World) FlushQueues() {
	// --- Removals ---
	for _, e := range w.removeQueue {
		layer := e.Layer()
		list := w.layers[layer]
		for i, ent := range list {
			if ent == e {
				// decrement the type-count
				typeName := reflect.TypeOf(e).Elem().Name()
				w.typeCounts[typeName]--

				// remove this entity from its layer slice
				w.layers[layer] = append(list[:i], list[i+1:]...)
				break
			}
		}
	}
	// clear removal queue
	w.removeQueue = w.removeQueue[:0]

	// --- Additions ---
	for _, e := range w.addQueue {
		layer := e.Layer()
		if _, ok := w.layers[layer]; !ok {
			// first time we’ve seen this layer: initialize and record it
			w.layers[layer] = make([]Entity, 0, 8)
			w.layerOrder = append(w.layerOrder, layer)
			sort.Ints(w.layerOrder)
		}
		// append the new entity
		w.layers[layer] = append(w.layers[layer], e)

		// increment the type-count
		typeName := reflect.TypeOf(e).Elem().Name()
		w.typeCounts[typeName]++
	}
	// clear addition queue
	w.addQueue = w.addQueue[:0]
}

// Count returns the total number of active Entities in the world.
func (w *World) Count() int {
	w.FlushQueues()
	total := 0
	for _, layer := range w.layerOrder {
		total += len(w.layers[layer])
	}
	return total
}

// LayerCount returns how many Entities live on the given layer.
func (w *World) LayerCount(layer int) int {
	w.FlushQueues()
	return len(w.layers[layer])
}

// TypeCount returns how many Entities of a given Go type name are in the world.
// This is now a constant‐time map lookup.
func (w *World) TypeCount(typeName string) int {
	w.FlushQueues()
	return w.typeCounts[typeName]
}

// ForEach calls fn on every Entity in the world, in layer‐order.
func (w *World) ForEach(fn func(Entity)) {
	w.FlushQueues()
	for _, layer := range w.layerOrder {
		for _, e := range w.layers[layer] {
			fn(e)
		}
	}
}

// Update all active Entities
func (w *World) Update(dt float64) {
	w.FlushQueues()
	for _, layer := range w.layerOrder {
		for _, e := range w.layers[layer] {
			e.Update(dt)
		}
	}
}

// Render all Entities in front→back order
func (w *World) Render() {
	w.FlushQueues()
	for _, layer := range w.layerOrder {
		for _, e := range w.layers[layer] {
			e.Render()
		}
	}
}

// BringToFront, SendToBack, BringForward, SendBackward omitted for brevity...

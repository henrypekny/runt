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
}

// NewWorld makes an empty World.
func NewWorld() *World {
	return &World{
		layers:     make(map[int][]Entity),
		layerOrder: make([]int, 0, 8),

		addQueue:    make([]Entity, 0, 16),
		removeQueue: make([]Entity, 0, 16),

		UseCamera: false,
		pool:      make(map[string][]Entity),
	}
}

// Entities returns a flat slice of all Entities in this World,
// in ascending layer→insertion order.
func (w *World) Entities() []Entity {
	// make sure we see the up-to-date set
	w.FlushQueues()

	// pre-allocate for speed
	var total int
	for _, layer := range w.layerOrder {
		total += len(w.layers[layer])
	}
	all := make([]Entity, 0, total)

	// append each layer's entities in layer-order
	for _, layer := range w.layerOrder {
		all = append(all, w.layers[layer]...)
	}
	return all
}

// Add queues an Entity for addition at end of this frame.
func (w *World) Add(e Entity) {
	w.addQueue = append(w.addQueue, e)
}

// Remove queues an Entity for removal at end of this frame.
func (w *World) Remove(e Entity) {
	w.removeQueue = append(w.removeQueue, e)
}

// FlushQueues integrates all queued add/removes.
// Call this once per frame (e.g. at end of Update or start of Render).
func (w *World) FlushQueues() {
	// process removals
	for _, e := range w.removeQueue {
		layer := e.Layer()
		list := w.layers[layer]
		for i, ent := range list {
			if ent == e {
				w.layers[layer] = append(list[:i], list[i+1:]...)
				break
			}
		}
	}
	w.removeQueue = w.removeQueue[:0]

	// process additions
	for _, e := range w.addQueue {
		layer := e.Layer()
		if _, ok := w.layers[layer]; !ok {
			w.layers[layer] = make([]Entity, 0, 8)
			w.layerOrder = append(w.layerOrder, layer)
			sort.Ints(w.layerOrder)
		}
		w.layers[layer] = append(w.layers[layer], e)
	}
	w.addQueue = w.addQueue[:0]
}

// Update all active Entities
func (w *World) Update(dt float64) {
	// first, flush any pending adds/removes from last frame
	w.FlushQueues()

	// then update
	for _, layer := range w.layerOrder {
		for _, e := range w.layers[layer] {
			e.Update(dt)
		}
	}
}

// Render all visible Entities, front->back.
func (w *World) Render() {
	// if you want to defer flush until just before render:
	w.FlushQueues()

	for _, layer := range w.layerOrder {
		for _, e := range w.layers[layer] {
			e.Render()
		}
	}
}

// BringToFront moves e to the highest index in its layer slice.
func (w *World) BringToFront(e Entity) {
	layer := e.Layer()
	list := w.layers[layer]
	// find and remove
	for i, ent := range list {
		if ent == e {
			list = append(list[:i], list[i+1:]...)
			list = append(list, e)
			break
		}
	}
	w.layers[layer] = list
}

// SendToBack moves e to index 0 in its layer slice.
func (w *World) SendToBack(e Entity) {
	layer := e.Layer()
	list := w.layers[layer]
	for i, ent := range list {
		if ent == e {
			list = append([]Entity{e}, append(list[:i], list[i+1:]...)...)
			break
		}
	}
	w.layers[layer] = list
}

// BringForward swaps e with the one in front of it.
func (w *World) BringForward(e Entity) {
	layer := e.Layer()
	list := w.layers[layer]
	for i := range list {
		if list[i] == e && i < len(list)-1 {
			list[i], list[i+1] = list[i+1], list[i]
			break
		}
	}
}

// SendBackward swaps e with the one behind it.
func (w *World) SendBackward(e Entity) {
	layer := e.Layer()
	list := w.layers[layer]
	for i := range list {
		if list[i] == e && i > 0 {
			list[i], list[i-1] = list[i-1], list[i]
			break
		}
	}
}

// (Optional) Recycle pushes e into a type‐based pool for reuse.
func (w *World) Recycle(typeName string, e Entity) {
	w.pool[typeName] = append(w.pool[typeName], e)
}

// (Optional) Create tries to pull from pool before allocating new.
func (w *World) Create(typeName string, ctor func() Entity) Entity {
	if lst := w.pool[typeName]; len(lst) > 0 {
		e := lst[len(lst)-1]
		w.pool[typeName] = lst[:len(lst)-1]
		return e
	}
	return ctor()
}

// ForEach calls fn on every Entity in the world, in layer‐order.
func (w *World) ForEach(fn func(Entity)) {
	// make sure queues are flushed so we see the up‐to‐date list
	w.FlushQueues()
	for _, layer := range w.layerOrder {
		for _, e := range w.layers[layer] {
			fn(e)
		}
	}
}

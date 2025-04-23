package engine

// Version of the engine API:
const Version = "v0.1.0"

// Init is called by the game on startup.
func Init() {
    println("runt engine initialized, version", Version)
}

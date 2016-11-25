package graph

// #include "konig.h"
import "C"

// Init _must_ be called on the main thread. It initializes the graphics
// and compute systems. The graph window dimensions will be width x height
// unless fullscreen is set, in which case, the current default monitor
// resolution is used.
func Init(width, height int, fullscreen bool) {
	var fullScreenInt int
	if fullscreen {
		fullScreenInt = 1
	}
	C.init(C.int(width), C.int(height), C.int(fullScreenInt))
}

// Update _must_ be called on the main thread. It processes user events
// and refreshes the view. Returns true if the user wants to quit.
func Update() bool {
	var wantsQuitInt = C.update()
	return wantsQuitInt == 1
}

// Shutdown _must_ be called ont he main thread. It destroys all
// GPU objects, detaches the internal event handlers and closes the
// graph view.
func Shutdown() {
	C.shutdown()
}

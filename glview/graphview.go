package glview

// #include "graphview.h"
import "C"

// GraphView implements the GraphView interface by providing
// an OpenGL graph renderer.
type GraphView struct {
	width, height int
	fullscreen    bool
}

// New returns a new OpenGL graphview. This can be called on any thread.
func New(width, height int, fullscreen bool) *GraphView {
	return &GraphView{width, height, fullscreen}
}

// Startup _must_ be called on the main thread. It initializes the graphics
// and compute systems. The graph window dimensions will be width x height
// unless fullscreen is set, in which case, the current default monitor
// resolution is used.
func (g *GraphView) Startup() {
	var fullScreenInt int
	if g.fullscreen {
		fullScreenInt = 1
	}
	C.graphview_init(C.int(g.width), C.int(g.height), C.int(fullScreenInt))
}

// Update _must_ be called on the main thread. It processes user events
// and refreshes the view. Returns true if the user wants to quit.
func (g *GraphView) Update() (dt float64, shouldQuit bool) {
	var cDT C.double
	var wantsQuitInt = C.graphview_update(&cDT)
	return float64(cDT), wantsQuitInt == 1
}

// Shutdown _must_ be called ont he main thread. It destroys all
// GPU objects, detaches the internal event handlers and closes the
// graph view.
func (g *GraphView) Shutdown() {
	C.graphview_shutdown()
}

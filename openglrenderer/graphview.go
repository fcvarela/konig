package openglrenderer

// #include "graphview.h"
import "C"

// Renderer implements the konig.Renderer interface by providing
// an OpenGL graph renderer.
type Renderer struct {
	width, height int
	fullscreen    bool
}

// New returns a new OpenGL renderer. This can be called on any thread.
func New(width, height int, fullscreen bool) *Renderer {
	return &Renderer{width, height, fullscreen}
}

// Startup _must_ be called on the main thread. It initializes the graphics
// and compute systems. The graph window dimensions will be width x height
// unless fullscreen is set, in which case, the current default monitor
// resolution is used.
func (g *Renderer) Startup() {
	var fullScreenInt int
	if g.fullscreen {
		fullScreenInt = 1
	}
	C.graphview_init(C.int(g.width), C.int(g.height), C.int(fullScreenInt))
}

// Render _must_ be called on the main thread. It processes user events
// and refreshes the view. Returns true if the user wants to quit.
func (g *Renderer) Render() (dt float32, shouldQuit bool) {
	var cDT C.float
	var wantsQuitInt = C.graphview_update(&cDT)
	return float32(cDT), wantsQuitInt == 1
}

// Shutdown _must_ be called ont he main thread. It destroys all
// GPU objects, detaches the internal event handlers and closes the
// graph view.
func (g *Renderer) Shutdown() {
	C.graphview_shutdown()
}

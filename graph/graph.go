package graph

/*
// Windows Build Tags
// ----------------
// GLFW Options:
#cgo windows CFLAGS: -D_GLFW_WIN32

// Linker Options:
#cgo windows LDFLAGS: -lopengl32 -lgdi32


// Darwin Build Tags
// ----------------
// GLFW Options:
#cgo darwin CFLAGS: -D_GLFW_COCOA -D_GLFW_USE_CHDIR -D_GLFW_USE_MENUBAR -Wno-deprecated-declarations

// Linker Options:
#cgo darwin LDFLAGS: -framework Cocoa -framework OpenGL -framework IOKit -framework CoreVideo -lglfw3


// Linux Build Tags
// ----------------
// GLFW Options:
#cgo linux,!wayland CFLAGS: -D_GLFW_X11
#cgo linux,wayland CFLAGS: -D_GLFW_WAYLAND

// Linker Options:
#cgo linux,!wayland LDFLAGS: -lGL -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lm -lXinerama -ldl -lrt
#cgo linux,wayland LDFLAGS: -lGL -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lm -lXinerama -ldl -lrt


// FreeBSD Build Tags
// ----------------
// GLFW Options:
#cgo freebsd,!wayland CFLAGS: -D_GLFW_X11 -D_GLFW_HAS_GLXGETPROCADDRESSARB -D_GLFW_HAS_DLOPEN
#cgo freebsd,wayland CFLAGS: -D_GLFW_WAYLAND -D_GLFW_HAS_DLOPEN

// Linker Options:
#cgo freebsd,!wayland LDFLAGS: -lGL -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lm -lXinerama
#cgo freebsd,wayland LDFLAGS: -lGL -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lm -lXinerama

#include "konig.h"
*/
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
	if wantsQuitInt == 1 {
		return true
	}
	return false
}

// Shutdown _must_ be called ont he main thread. It destroys all
// GPU objects, detaches the internal event handlers and closes the
// graph view.
func Shutdown() {
	C.shutdown()
}

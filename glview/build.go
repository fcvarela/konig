package glview

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
#cgo darwin LDFLAGS: -framework Cocoa -framework OpenGL -framework IOKit -framework CoreVideo -lglfw -lGLEW


// Linux Build Tags
// ----------------
// GLFW Options:
#cgo linux,!wayland CFLAGS: -D_GLFW_X11
#cgo linux,wayland CFLAGS: -D_GLFW_WAYLAND

// Linker Options:
#cgo linux,!wayland LDFLAGS: -lGL -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lm -lXinerama -ldl -lrt -lglfw -lGLEW
#cgo linux,wayland LDFLAGS: -lGL -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lm -lXinerama -ldl -lrt -lglfw -lGLEW


// FreeBSD Build Tags
// ----------------
// GLFW Options:
#cgo freebsd,!wayland CFLAGS: -D_GLFW_X11 -D_GLFW_HAS_GLXGETPROCADDRESSARB -D_GLFW_HAS_DLOPEN
#cgo freebsd,wayland CFLAGS: -D_GLFW_WAYLAND -D_GLFW_HAS_DLOPEN

// Linker Options:
#cgo freebsd,!wayland LDFLAGS: -lGL -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lm -lXinerama -lglfw -lGLEW
#cgo freebsd,wayland LDFLAGS: -lGL -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lm -lXinerama -lglfw -lGLEW
*/
import "C"

package konig

// Renderer defines an interface used to render a graph
type Renderer interface {
	// Startup provides a post-init hook for the view
	// to perform any initialization it may require (create windows, etc)
	Startup()

	// Render tells the view to render and returns elapsed time since last
	// call and a bool signaling the app to quit. Implementations should
	// keep track of when they are called.
	// This is useful for visual renderers (ie: OpenGL) to return a frame duration
	// which can then be used by other packages to solve the layout.
	// User events which tell the app the user wants to quit (close window, etc)
	// should trigger shouldQuit to be true.
	Render() (dt float32, shouldQuit bool)

	// Shutdown provides a pre-quit hook for the view
	// to perform any cleanup it may require (close windows, etc)
	Shutdown()
}

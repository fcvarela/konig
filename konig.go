package konig

// Solver defines the solve interface
type Solver interface {
	// Step iterates the solver by recomputing the graph layout
	// based on a time step and slices containing nodes and edges.
	Step(dt float64, nodes []Node, edges []Edge) error
}

// GraphView defines an interface used to render a graph
type View interface {
	// Startup provides a post-init hook for the view
	// to perform any initialization it may require (create windows, etc)
	Startup()

	// Update tells the view to render and returns elapsed time since last
	// call and a bool signaling the app to quit. Implementations should
	// keep track of when they are called.
	// This is useful for visual renderers (ie: OpenGL) to return a frame duration
	// which can then be used by other packages to solve the layout.
	// User events which tell the app the user wants to quit (close window, etc)
	// should trigger shouldQuit to be true.
	Update() (dt float64, shouldQuit bool)

	// Shutdown provides a pre-quit hook for the view
	// to perform any cleanup it may require (close windows, etc)
	Shutdown()
}

var (
	solver Solver
	view   View
)

// Startup sets the solver and view and initializes them
func Startup(s Solver, v View) {
	solver = s
	view = v

	view.Startup()
}

// Shutdown provides a shutdown hook
func Shutdown() {
	view.Shutdown()
}

func Step() bool {
	dt, quit := view.Update()
	graphsLock.Lock()
	for _, g := range graphs {
		if err := solver.Step(dt, g.Nodes, g.Edges); err != nil {
			panic(err)
		}
	}
	defer graphsLock.Unlock()

	return quit
}

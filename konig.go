package konig

var (
	solver   Solver
	renderer Renderer
	dt       float32
)

// Startup sets the solver and renderer and initializes them.
func Startup(s Solver, r Renderer) {
	solver = s
	renderer = r

	renderer.Startup()
}

// Shutdown provides a shutdown hook
func Shutdown() {
	renderer.Shutdown()
}

// Step is called by the main loop and provides konig
// with a hook to step the selected solver and renderer.
func Step() bool {
	graphsLock.Lock()

	for _, g := range graphs {
		if err := solver.Solve(&g, dt); err != nil {
			panic(err)
		}
	}

	ldt, quit := renderer.Render()
	defer graphsLock.Unlock()

	dt = ldt
	return quit
}

package konig

// SolverLayout defines the type of layout a Solver
// will use to solve the graph.
type SolverLayout int

// The following is a list of supported layout types
// Not all solvers will necessarily implement all types,
// so errors should be returned in those cases.
const (
	SolverLayoutForceDirected SolverLayout = iota
	SolverLayoutRadial        SolverLayout = iota
)

// Solver defines an interface for solving graph layouts
type Solver interface {
	// Solve iterates the solver by recomputing the graph layout.
	Solve(g *Graph, dt float32) error
}

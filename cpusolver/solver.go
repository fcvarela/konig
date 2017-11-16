package cpusolver

import (
	"github.com/fcvarela/konig"
	"github.com/golang/glog"
)

// Solver implements the solver interface by laying out the graph
// in the cpu
type Solver struct{}

// New returns a new Solver
func New() *Solver {
	return &Solver{}
}

// Step implements the solver interface
func (s *Solver) Step(dt float64, nodes []konig.Node, edges []konig.Edge) error {
	glog.Infof("Stepping: %f\n", dt)
	return nil
}

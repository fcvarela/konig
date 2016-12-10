package cpusolver

import (
	"github.com/fcvarela/konig"
	"github.com/golang/glog"
)

type Solver struct{}

func New() *Solver {
	return &Solver{}
}

func (s *Solver) Step(dt float64, nodes []konig.Node, edges []konig.Edge) error {
	glog.Infof("Stepping: %f\n", dt)
	return nil
}

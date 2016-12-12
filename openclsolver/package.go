// Package openclsolver implements the konig.Solver interface
// by providing an OpenCL based graph solver. This solver is
// only compatible with OpenGL renderers as it requires
// object sharing to avoid gl/cl DMA transfers.
package openclsolver

import "C"

import (
	"github.com/fcvarela/konig"
	"github.com/golang/glog"
)

// Solver implements the konig.Solver interface by providing
// an OpenCL based graph solver.
type Solver struct {
	InputNodeBufferMap  map[konig.NodeBufferHandle]C.cl_mem
	OutputNodeBufferMap map[konig.NodeBufferHandle]C.cl_mem
	EdgeBufferMap       map[konig.EdgeBufferHandle]C.cl_mem
}

// New returns a new solver initialized with the passed
// configuration options.
func New() *Solver {
	return &Solver{
		InputNodeBufferMap:  make(map[konig.NodeBufferHandle]C.cl_mem),
		OutputNodeBufferMap: make(map[konig.NodeBufferHandle]C.cl_mem),
		EdgeBufferMap:       make(map[konig.EdgeBufferHandle]C.cl_mem),
	}
}

// Solve implements the konig.Solver interface
func (s *Solver) Solve(g *konig.Graph, dt float32) error {
	// if the graph's VRAM handles are 0, their memory hasn't
	// been mapped yet.
	if g.InputNodeBuffer.(uint32) == 0 {
		glog.Info("Aborting solver iteration, no VRAM yet...")
	}

	// if our handler maps don't know about them, then
	// we need to queue in a capture.
	/*
	   cl_vbo_in = clCreateFromGLBuffer(context, CL_MEM_READ_ONLY, vbo_in, NULL);
	   cl_vbo_out = clCreateFromGLBuffer(context, CL_MEM_READ_WRITE, vbo_out, NULL);
	   cl_edge_vbo = clCreateFromGLBuffer(context, CL_MEM_READ_ONLY, edge_vbo, NULL);

	   clEnqueueAcquireGLObjects(queue, 1, &cl_vbo_in, 0, 0, 0);
	   clEnqueueAcquireGLObjects(queue, 1, &cl_vbo_out, 0, 0, 0);
	   clEnqueueAcquireGLObjects(queue, 1, &cl_edge_vbo, 0, 0, 0);
	*/

	glog.Infof("YAY: %f", dt)

	return nil
}

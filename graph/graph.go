// Package graph provides an implementation of a graph whose nodes and edges are represented
// in sequential memory so that they can be used directly by compute and graphics kernels.
// the convention for thread safery is to use a mutex on outer operations (public api)
// and _not_ on inner operations (lowercase).
package graph

import "sync"

// Handle is an opaque type which represents a graph
type Handle uint32

// NodeHandle is an opaque type which represents a graph node
type NodeHandle uint32

// EdgeHandle is an opaque type which represents a graph edge
type EdgeHandle uint32

// std140 alignment for opengl/opencl
type node struct {
	position     [4]float32
	velocity     [4]float32
	acceleration [4]float32
	active       [4]float32
}

type edge struct {
	node1ID NodeHandle
	node2ID NodeHandle
	active  uint32
	padding uint32
}

// graph contains all data necessary to manage a graph where nodes and edges
// are continuous memory regions. this makes it really easy to use as a GPU
// buffer we can use directly with both compute kernels and draw operations.
type graph struct {
	nodes       []node
	edges       []edge
	freeNodeSet map[NodeHandle]struct{}
	freeEdgeSet map[EdgeHandle]struct{}

	// contains a list of edges pointing to and from each node
	// we use this to delete those edges when a node is deleted
	nodeEdgeIndex map[NodeHandle][]EdgeHandle
}

var (
	graphsLock sync.Mutex
	graphs     []graph
)

func init() {
	graphs = make([]graph, 0)
}

// New returns a handle for a newly created graph
func New() Handle {
	graphsLock.Lock()
	defer graphsLock.Unlock()

	var g = graph{
		nodes:         make([]node, 0),
		edges:         make([]edge, 0),
		freeNodeSet:   make(map[NodeHandle]struct{}),
		freeEdgeSet:   make(map[EdgeHandle]struct{}),
		nodeEdgeIndex: make(map[NodeHandle][]EdgeHandle),
	}

	graphs = append(graphs, g)
	return Handle(len(graphs) - 1)
}

// NewNode returns a handle for a newly created graph node
func (g Handle) NewNode() NodeHandle {
	graphsLock.Lock()
	defer graphsLock.Unlock()

	// do we have any items in our freeset? return the first one after deleting it
	for k := range graphs[g].freeNodeSet {
		// delete it from the freeNodeSet
		delete(graphs[g].freeNodeSet, k)

		// make it active
		graphs[g].nodes[k].active = [4]float32{1.0, 1.0, 1.0, 1.0}

		// return it
		return k
	}

	// got here? we don't have anything in our free set, add a new node
	graphs[g].nodes = append(graphs[g].nodes, node{})
	return NodeHandle(len(graphs[g].nodes) - 1)
}

// DeleteNode deletes a graph node and all edges connected to or from it
func (g Handle) DeleteNode(n NodeHandle) {
	graphsLock.Lock()
	defer graphsLock.Unlock()

	// make it inactive so the cl kernel and shaders won't use it
	graphs[g].nodes[n].active = [4]float32{0.0, 0.0, 0.0, 0.0}

	// add the node to the graphs freeNode list so that it can be used
	// the next time a node is created.
	graphs[g].freeNodeSet[n] = struct{}{}

	// delete any edge connected to or from this node, they are indexed in nodeEdgeIndex
	for _, e := range graphs[g].nodeEdgeIndex[n] {
		g.deleteEdge(e)
	}

	// reset the index for this node
	graphs[g].nodeEdgeIndex[n] = graphs[g].nodeEdgeIndex[n][:0]
}

// NewEdge returns a handle for a newly created graph edge
func (g Handle) NewEdge(n1, n2 NodeHandle) EdgeHandle {
	graphsLock.Lock()
	defer graphsLock.Unlock()

	// create the edge
	var newEdge = edge{active: 1, node1ID: n1, node2ID: n2}

	// do we have any items in our freeset? return the first one after deleting it
	for k := range graphs[g].freeEdgeSet {
		// remove it from freeEdgeSet
		delete(graphs[g].freeEdgeSet, k)

		// initialize it
		graphs[g].edges[k] = newEdge

		// index it
		g.indexEdge(k, []NodeHandle{n1, n2})

		// done
		return k
	}

	// got here? we don't have anything in our free set, add a new edge
	graphs[g].edges = append(graphs[g].edges, newEdge)

	// get a handle
	var newEdgeHandle = EdgeHandle(len(graphs[g].edges) - 1)

	// index it
	g.indexEdge(newEdgeHandle, []NodeHandle{n1, n2})

	return newEdgeHandle
}

func (g Handle) deleteEdge(e EdgeHandle) {
	// add the edge to the graphs freeNode list so that it can be used
	// the next time an edge is created.
	graphs[g].freeEdgeSet[e] = struct{}{}
}

// DeleteEdge deletes a graph edge
func (g Handle) DeleteEdge(e EdgeHandle) {
	graphsLock.Lock()
	defer graphsLock.Unlock()

	g.deleteEdge(e)
}

// indexEdge adds this edge to both node's list of referencing edges
// this will be checked whenever one of those nodes is deleted to
// determine what edges need to be removed
func (g Handle) indexEdge(e EdgeHandle, nodes []NodeHandle) {
	for _, n := range nodes {
		if graphs[g].nodeEdgeIndex[n] == nil {
			graphs[g].nodeEdgeIndex[n] = make([]EdgeHandle, 0)
		}
		graphs[g].nodeEdgeIndex[n] = append(graphs[g].nodeEdgeIndex[n], e)
	}
}

// Package graph provides an implementation of a graph whose nodes and edges are represented
// in sequential memory so that they can be used directly by compute and graphics kernels.
// the convention for thread safery is to use a mutex on outer operations (public api)
// and _not_ on inner operations (lowercase).
package graph

import (
	"fmt"
	"sync"

	"github.com/satori/go.uuid"
)

// Handle is an opaque type which represents a graph
type Handle string

// NodeHandle is an opaque type which represents a graph node
type NodeHandle string

// EdgeHandle is an opaque type which represents a graph edge
type EdgeHandle string

// std140 alignment for opengl/opencl
type node struct {
	position     [4]float32
	velocity     [4]float32
	acceleration [4]float32
	active       [4]float32
}

type edge struct {
	node1Index uint32
	node2Index uint32
}

// graph contains all data necessary to manage a graph where nodes and edges
// are continuous memory regions. this makes it really easy to use as a GPU
// buffer we can use directly with both compute kernels and draw operations.
type graph struct {
	nodes       []node
	edges       []edge
	freeNodeSet map[uint32]struct{}
	freeEdgeSet map[uint32]struct{}
	nodeHandles map[NodeHandle]uint32
	edgeHandles map[EdgeHandle]uint32

	// contains a list of edges pointing to and from each node
	// we use this to delete those edges when a node is deleted
	nodeEdgeIndex map[NodeHandle][]EdgeHandle
}

var (
	graphsLock   sync.Mutex
	graphs       []graph
	graphHandles map[Handle]uint32
)

var (
	errGraphNotFound = "Graph %s does not exist"
	errNodeNotFound  = "Node %s does not exist"
	errEdgeNotFound  = "Edge %s does not exist"
)

func init() {
	graphs = make([]graph, 0)
	graphHandles = make(map[Handle]uint32)
}

// New returns a handle for a newly created graph
func New() Handle {
	graphsLock.Lock()
	defer graphsLock.Unlock()

	// initialize it
	var g = graph{
		nodes:         make([]node, 0),
		edges:         make([]edge, 0),
		nodeHandles:   make(map[NodeHandle]uint32),
		edgeHandles:   make(map[EdgeHandle]uint32),
		freeNodeSet:   make(map[uint32]struct{}),
		freeEdgeSet:   make(map[uint32]struct{}),
		nodeEdgeIndex: make(map[NodeHandle][]EdgeHandle),
	}

	// add it to the graph slice
	graphs = append(graphs, g)

	// create a handle
	var handle = Handle(uuid.NewV4().String())

	// map it
	graphHandles[handle] = uint32(len(graphs) - 1)

	return Handle(handle)
}

// NewNode returns a handle for a newly created graph node
func NewNode(g Handle) (NodeHandle, error) {
	graphsLock.Lock()
	defer graphsLock.Unlock()

	gidx, ok := graphHandles[g]
	if !ok {
		return NodeHandle(""), fmt.Errorf(errGraphNotFound, g)
	}

	var nodeIndex uint32
	var found bool

	// do we have any items in our free index
	for k := range graphs[gidx].freeNodeSet {
		delete(graphs[gidx].freeNodeSet, k)
		found = true
		nodeIndex = k
	}

	if !found {
		nodeIndex = uint32(len(graphs[gidx].nodes))
		graphs[gidx].nodes = append(graphs[gidx].nodes, node{})
	}

	// make active
	graphs[gidx].nodes[nodeIndex].active = [4]float32{1.0, 1.0, 1.0, 1.0}

	// create a handle
	var handle = NodeHandle(uuid.NewV4().String())

	// map it
	graphs[gidx].nodeHandles[handle] = nodeIndex

	return handle, nil
}

// DeleteNode deletes a graph node and all edges connected to or from it
func DeleteNode(g Handle, n NodeHandle) error {
	graphsLock.Lock()
	defer graphsLock.Unlock()

	// resolve the graph handle
	gidx, ok := graphHandles[g]
	if !ok {
		return fmt.Errorf(errGraphNotFound, g)
	}

	// resolve the node handle
	index, ok := graphs[gidx].nodeHandles[n]
	if !ok {
		return fmt.Errorf(errNodeNotFound, n)
	}

	// delete the handle
	delete(graphs[gidx].nodeHandles, n)

	// make it inactive so the cl kernel and shaders won't use it
	graphs[gidx].nodes[index].active = [4]float32{0.0, 0.0, 0.0, 0.0}

	// add the node to the graphs freeNode list so that it can be used
	// the next time a node is created.
	graphs[gidx].freeNodeSet[index] = struct{}{}

	// delete any edge connected to or from this node, they are indexed in nodeEdgeIndex
	for _, e := range graphs[gidx].nodeEdgeIndex[n] {
		deleteEdge(g, e)
	}

	// reset the index for this node
	graphs[gidx].nodeEdgeIndex[n] = graphs[gidx].nodeEdgeIndex[n][:0]

	return nil
}

// NewEdge returns a handle for a newly created graph edge
func NewEdge(g Handle, nh1, nh2 NodeHandle) (EdgeHandle, error) {
	graphsLock.Lock()
	defer graphsLock.Unlock()

	gidx, ok := graphHandles[g]
	if !ok {
		return EdgeHandle(""), fmt.Errorf(errGraphNotFound, g)
	}

	var edgeIndex uint32
	var found bool

	// do we have any items in our free index
	for k := range graphs[gidx].freeEdgeSet {
		delete(graphs[gidx].freeEdgeSet, k)
		found = true
		edgeIndex = k
	}

	if !found {
		edgeIndex = uint32(len(graphs[gidx].edges))
		graphs[gidx].edges = append(graphs[gidx].edges, edge{})
	}

	// resolve n1 and n2
	n1, ok := graphs[gidx].nodeHandles[nh1]
	if !ok {
		return EdgeHandle(""), fmt.Errorf(errNodeNotFound, nh1)
	}

	n2, ok := graphs[gidx].nodeHandles[nh2]
	if !ok {
		return EdgeHandle(""), fmt.Errorf(errNodeNotFound, nh2)
	}

	// set it
	graphs[gidx].edges[edgeIndex] = edge{node1Index: n1, node2Index: n2}

	// get a handle
	var handle = EdgeHandle(uuid.NewV4().String())

	// map it
	graphs[gidx].edgeHandles[handle] = edgeIndex

	// index it
	indexEdge(gidx, handle, []NodeHandle{nh1, nh2})

	return handle, nil
}

func deleteEdge(g Handle, e EdgeHandle) error {
	// resolve graph
	gidx, ok := graphHandles[g]
	if !ok {
		return fmt.Errorf(errGraphNotFound, g)
	}

	// resolve edge
	eidx, ok := graphs[gidx].edgeHandles[e]
	if !ok {
		return fmt.Errorf(errEdgeNotFound, e)
	}

	// delete the edge
	delete(graphs[gidx].edgeHandles, e)

	// add the edge to the graphs freeNode list so that it can be used
	// the next time an edge is created.
	graphs[gidx].freeEdgeSet[eidx] = struct{}{}

	// it will still end up in the gpu so make it a zero-length line
	graphs[gidx].edges[eidx].node1Index = 0
	graphs[gidx].edges[eidx].node2Index = 0

	return nil
}

// DeleteEdge deletes a graph edge
func DeleteEdge(g Handle, e EdgeHandle) error {
	graphsLock.Lock()
	defer graphsLock.Unlock()

	return deleteEdge(g, e)
}

// indexEdge adds this edge to both node's list of referencing edges
// this will be checked whenever one of those nodes is deleted to
// determine what edges need to be removed
func indexEdge(gidx uint32, e EdgeHandle, nodes []NodeHandle) {
	for _, n := range nodes {
		if graphs[gidx].nodeEdgeIndex[n] == nil {
			graphs[gidx].nodeEdgeIndex[n] = make([]EdgeHandle, 0)
		}
		// append the edge to it
		graphs[gidx].nodeEdgeIndex[n] = append(graphs[gidx].nodeEdgeIndex[n], e)
	}
}

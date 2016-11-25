package graph

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
	custom       [4]float32
}

type edge struct {
	node1ID  NodeHandle
	node2ID  NodeHandle
	padding0 uint32
	padding1 uint32
}

// graph contains all data necessary to manage a graph where nodes and edges
// are continuous memory regions. this makes it really easy to use as a GPU
// buffer we can draw directly from as well as passing them to compute
// kernels
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
	graphs []graph
)

func init() {
	graphs = make([]graph, 0)
}

// New returns a handle for a newly created graph
func New() Handle {
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
	// do we have any items in our freeset? return the first one after deleting it
	for k := range graphs[g].freeNodeSet {
		// delete it from the freeNodeSet
		delete(graphs[g].freeNodeSet, k)

		// return it
		return k
	}

	// got here? we don't have anything in our free set, add a new node
	graphs[g].nodes = append(graphs[g].nodes, node{})
	return NodeHandle(len(graphs[g].nodes) - 1)
}

// DeleteNode deletes a graph node and all edges connected to or from it
func (g Handle) DeleteNode(n NodeHandle) {
	// add the node to the graphs freeNode list so that it can be used
	// the next time a node is created.
	graphs[g].freeNodeSet[n] = struct{}{}
}

// NewEdge returns a handle for a newly created graph edge
func (g Handle) NewEdge(n1, n2 NodeHandle) EdgeHandle {
	// do we have any items in our freeset? return the first one after deleting it
	for k := range graphs[g].freeEdgeSet {
		// remove it from freeEdgeSet
		delete(graphs[g].freeEdgeSet, k)

		// set it to the passed values
		graphs[g].edges[k].node1ID = n1
		graphs[g].edges[k].node2ID = n2

		// return it
		return k
	}

	// got here? we don't have anything in our free set, add a new edge
	graphs[g].edges = append(graphs[g].edges, edge{
		node1ID: n1,
		node2ID: n2,
	})
	return EdgeHandle(len(graphs[g].edges) - 1)
}

// DeleteEdge deletes a graph edge
func (g Handle) DeleteEdge(e EdgeHandle) {
	// add the edge to the graphs freeNode list so that it can be used
	// the next time an edge is created.
	graphs[g].freeEdgeSet[e] = struct{}{}
}

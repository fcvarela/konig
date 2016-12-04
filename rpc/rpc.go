//Package rpc implements a gRPC service
//go:generate ./codegen.sh
package rpc

import (
	"errors"
	"fmt"
	"net"

	"github.com/fcvarela/konig/graph"
	"github.com/golang/glog"
	"google.golang.org/grpc"

	"github.com/mwitkow/go-grpc-middleware"
	"golang.org/x/net/context"
)

var (
	errGraphNotFound      = errors.New("graph not found")
	errFailedToAddNode    = errors.New("failed to add a node")
	errFailedToAddEdge    = errors.New("failed to add a edge")
	errFailedToRemoveNode = errors.New("failed to remove a node")
	errFailedToRemoveEdge = errors.New("failed to remove a edge")
)

//KonigService struct implements the gRPC handlers
type KonigService struct {
}

//CreateGraph creates a graph
func (s *KonigService) CreateGraph(ctx context.Context, in *CreateGraphRequest) (*CreateGraphResponse, error) {
	handle := graph.New()
	return &CreateGraphResponse{Graph: &Graph{Handle: string(handle)}}, nil
}

//AddEdges adds edges to a graph
func (s *KonigService) AddEdges(ctx context.Context, in *AddEdgesRequest) (*AddEdgesResponse, error) {
	if in.Graph == nil || in.Graph.Handle == "" {
		glog.Errorf("%s", errGraphNotFound)
		return nil, errGraphNotFound
	}
	g := graph.Handle(in.Graph.Handle)
	var edges []*Edge
	for _, e := range in.Edges {
		n1 := graph.NodeHandle(e.Node1.Handle)
		n2 := graph.NodeHandle(e.Node2.Handle)
		h, err := graph.NewEdge(g, n1, n2)
		if err != nil {
			glog.Errorf("%s graph: %s node1: %s node2: %s err: %s", errFailedToAddEdge, in.Graph.Handle, e.Node1.Handle, e.Node2.Handle, err)
			continue
		}
		e.Handle = string(h)
		edges = append(edges, e)
	}
	return &AddEdgesResponse{Graph: in.Graph, Edges: edges}, nil
}

//AddNodes adds nodes to a graph
func (s *KonigService) AddNodes(ctx context.Context, in *AddNodesRequest) (*AddNodesResponse, error) {
	if in.Graph == nil || in.Graph.Handle == "" {
		glog.Errorf("%s", errGraphNotFound)
		return nil, errGraphNotFound
	}
	g := graph.Handle(in.Graph.Handle)
	var nodes []*Node
	for range in.Nodes {
		h, err := graph.NewNode(g)
		if err != nil {
			glog.Errorf("%s %s %s", errFailedToAddNode, in.Graph.Handle, err)
			return nil, err
		}
		n := &Node{Handle: string(h)}
		nodes = append(nodes, n)
	}
	return &AddNodesResponse{Graph: in.Graph, Nodes: nodes}, nil
}

//RemoveNodes removes nodes from a graph
func (s *KonigService) RemoveNodes(ctx context.Context, in *RemoveNodesRequest) (*RemoveNodesResponse, error) {
	if in.Graph == nil || in.Graph.Handle == "" {
		glog.Errorf("%s", errGraphNotFound)
		return nil, errGraphNotFound
	}
	g := graph.Handle(in.Graph.Handle)
	for _, node := range in.Nodes {
		glog.Infof("remove node")
		n := graph.NodeHandle(node.Handle)
		if err := graph.DeleteNode(g, n); err != nil {
			glog.Errorf("%s %s %s: %s", errFailedToRemoveNode, in.Graph.Handle, node.Handle, err)
		}
	}
	return &RemoveNodesResponse{}, nil
}

//RemoveEdges removes edges from a graph
func (s *KonigService) RemoveEdges(ctx context.Context, in *RemoveEdgesRequest) (*RemoveEdgesResponse, error) {
	if in.Graph == nil || in.Graph.Handle == "" {
		glog.Errorf("%s", errGraphNotFound)
		return nil, errGraphNotFound
	}
	g := graph.Handle(in.Graph.Handle)
	for _, edge := range in.Edges {
		e := graph.EdgeHandle(edge.Handle)
		if err := graph.DeleteEdge(g, e); err != nil {
			glog.Errorf("%s %s %s: %s", errFailedToRemoveEdge, in.Graph.Handle, edge.Handle, err)
		}
	}
	return &RemoveEdgesResponse{}, nil
}

//Start starts an RPC server
func Start(host string, port uint) {
	srv := &KonigService{}
	opts := []grpc.ServerOption{}
	opts = append(opts, grpc_middleware.WithUnaryServerChain(RequestLogger))
	s := grpc.NewServer(opts...)
	RegisterKonigRPCServer(s, srv)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		glog.Fatalf("Could not start listener: %s", err)
	}
	glog.Infof("Starting gRPC server at: %s", listener.Addr().String())
	glog.Fatal(s.Serve(listener))

	// read stop handler, close rpc server/clean shutdown
	select {}
}

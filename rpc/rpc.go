package rpc

import (
	"errors"
	"fmt"
	"net"

	"github.com/golang/glog"
	"google.golang.org/grpc"

	"github.com/mwitkow/go-grpc-middleware"
	"golang.org/x/net/context"
)

type KonigService struct {
	GraphName   string
	NumVertices uint32
	NumEdges    uint32
	Vertices    []*Vertex
	Edges       []*Edge
}

func (s *KonigService) grapStats() *GraphStats {
	return &GraphStats{
		Vertices: uint32(len(s.Vertices)),
		Edges:    uint32(len(s.Edges)),
	}
}

func (s *KonigService) CreateGraph(ctx context.Context, in *CreateGraphRequest) (*CreateGraphResponse, error) {
	if len(in.Graph.Name) == 0 {
		return nil, errors.New("invalid graph name")
	}
	s.GraphName = in.Graph.Name
	s.NumEdges = in.Graph.NumEdges
	s.NumVertices = in.Graph.NumVertices
	s.Vertices = make([]*Vertex, s.NumVertices)
	for _, v := range in.Graph.Vertices {
		s.Vertices = append(s.Vertices, v)
	}
	s.Edges = make([]*Edge, s.NumEdges)
	for _, e := range in.Graph.Edges {
		s.Edges = append(s.Edges, e)
	}
	return &CreateGraphResponse{Stats: s.grapStats()}, nil
}

func (s *KonigService) AddEdges(ctx context.Context, in *AddEdgesRequest) (*AddEdgesResponse, error) {
	if s.GraphName == "" {
		return nil, errors.New("graph not initialized")
	}
	for _, e := range in.Edges {
		s.Edges = append(s.Edges, e)
	}
	return &AddEdgesResponse{Stats: s.grapStats()}, nil
}

func (s *KonigService) AddVertices(ctx context.Context, in *AddVerticesRequest) (*AddVerticesResponse, error) {
	if s.GraphName == "" {
		return nil, errors.New("graph not initialized")
	}
	for _, v := range in.Vertices {
		glog.Infof("add vertex: %s", v.Label)
		s.Vertices = append(s.Vertices, v)
	}
	return &AddVerticesResponse{Stats: s.grapStats()}, nil
}

func (s *KonigService) RemoveVertices(ctx context.Context, in *RemoveVerticesRequest) (*RemoveVerticesResponse, error) {
	//FIXME: think of a strategy
	return nil, errors.New("not implemented")
}

func (s *KonigService) RemoveEdges(ctx context.Context, in *RemoveEdgesRequest) (*RemoveEdgesResponse, error) {
	//FIXME: lock lists
	el := make([]*Edge, s.NumEdges)
	removed := false
	for _, e := range in.Edges {
		for _, ge := range s.Edges {
			if e.X.Label == ge.X.Label && e.Y.Label == ge.Y.Label {
				removed = true
				continue
			}
			el = append(el, ge)
		}
	}
	if !removed {
		return nil, errors.New("no edges were remove")
	}
	s.Edges = el
	return &RemoveEdgesResponse{Stats: s.grapStats()}, nil
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

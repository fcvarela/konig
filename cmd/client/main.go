package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"gopkg.in/urfave/cli.v2"

	"github.com/fcvarela/konig/rpc"
)

var (
	addNodeCommand = &cli.Command{
		Name:    "node",
		Aliases: []string{"n"},
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("invalid number of parameters")
			}
			g := c.Args().First()
			client := c.App.Metadata["client"].(rpc.KonigRPCClient)
			resp, err := client.AddNodes(
				context.Background(),
				&rpc.AddNodesRequest{
					Graph: &rpc.Graph{
						Handle: g,
					},
					Nodes: []*rpc.Node{
						&rpc.Node{},
					},
				})
			if err != nil {
				return err
			}
			data, err := json.Marshal(resp)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", data)
			return nil
		},
	}

	addEdgeCommand = &cli.Command{
		Name:    "edge",
		Aliases: []string{"e"},
		Action: func(c *cli.Context) error {
			if c.NArg() != 3 {
				return errors.New("invalid number of parameters")
			}
			g := c.Args().Get(0)
			n1 := c.Args().Get(1)
			n2 := c.Args().Get(2)
			client := c.App.Metadata["client"].(rpc.KonigRPCClient)
			resp, err := client.AddEdges(
				context.Background(),
				&rpc.AddEdgesRequest{
					Graph: &rpc.Graph{
						Handle: g,
					},
					Edges: []*rpc.Edge{
						&rpc.Edge{
							Node1: &rpc.Node{
								Handle: n1,
							},
							Node2: &rpc.Node{
								Handle: n2,
							},
						},
					},
				})
			if err != nil {
				return err
			}
			data, err := json.Marshal(resp)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", data)
			return nil
		},
	}

	removeNodeCommand = &cli.Command{
		Name:    "node",
		Aliases: []string{"n"},
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
				return errors.New("invalid number of parameters")
			}
			g := c.Args().Get(0)
			n := c.Args().Get(1)
			client := c.App.Metadata["client"].(rpc.KonigRPCClient)
			resp, err := client.RemoveNodes(
				context.Background(),
				&rpc.RemoveNodesRequest{
					Graph: &rpc.Graph{
						Handle: g,
					},
					Nodes: []*rpc.Node{
						&rpc.Node{
							Handle: n,
						},
					},
				})
			if err != nil {
				return err
			}
			data, err := json.Marshal(resp)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", data)
			return nil
		},
	}

	removeEdgeCommand = &cli.Command{
		Name:    "edge",
		Aliases: []string{"e"},
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
				return errors.New("invalid number of parameters")
			}
			g := c.Args().Get(0)
			e := c.Args().Get(1)
			client := c.App.Metadata["client"].(rpc.KonigRPCClient)
			resp, err := client.RemoveEdges(
				context.Background(),
				&rpc.RemoveEdgesRequest{
					Graph: &rpc.Graph{
						Handle: g,
					},
					Edges: []*rpc.Edge{
						&rpc.Edge{
							Handle: e,
						},
					},
				})

			if err != nil {
				return err
			}
			data, err := json.Marshal(resp)
			if err != nil {
				return nil
			}
			fmt.Printf("%s\n", data)
			return nil
		},
	}
)

func initClient() rpc.KonigRPCClient {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "localhost", 1234), grpc.WithInsecure(), grpc.WithTimeout(10*time.Second))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := rpc.NewKonigRPCClient(conn)
	return client
}

func main() {
	app := cli.App{}
	app.Name = "konig"
	app.Commands = []*cli.Command{
		&cli.Command{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "creates a graph",
			Action: func(c *cli.Context) error {
				client := c.App.Metadata["client"].(rpc.KonigRPCClient)
				resp, err := client.CreateGraph(
					context.Background(),
					&rpc.CreateGraphRequest{
						Graph: &rpc.Graph{},
					})
				if err != nil {
					return err
				}
				data, err := json.Marshal(resp)
				if err != nil {
					return err
				}
				fmt.Printf("%s\n", data)
				return nil
			},
		},
		&cli.Command{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "adds a vertex or an edge",
			Subcommands: []*cli.Command{
				addNodeCommand,
				addEdgeCommand,
			},
		},
		&cli.Command{
			Name:    "remove",
			Aliases: []string{"r"},
			Usage:   "removes a node or an edge",
			Subcommands: []*cli.Command{
				removeNodeCommand,
				removeEdgeCommand,
			},
		},
	}
	client := initClient()
	app.Metadata = make(map[string]interface{}, 1)
	app.Metadata["client"] = client
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

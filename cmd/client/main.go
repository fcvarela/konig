package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fcvarela/konig/rpc"
	"google.golang.org/grpc"
	"gopkg.in/urfave/cli.v2"
)

var (
	addVertexCommand = &cli.Command{
		Name:    "vertex",
		Aliases: []string{"v"},
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("invalid number of parameters")
			}
			label := c.Args().First()
			client := c.App.Metadata["client"].(rpc.KonigRPCClient)
			resp, err := client.AddVertices(
				context.Background(),
				&rpc.AddVerticesRequest{
					Vertices: []*rpc.Vertex{
						&rpc.Vertex{
							Label: label,
						},
					},
				})
			if err == nil {
				data, err := json.Marshal(resp)
				if err != nil {
					return err
				}
				fmt.Printf("%s\n", data)
			}
			return err
		},
	}

	addEdgeCommand = &cli.Command{
		Name:    "edge",
		Aliases: []string{"e"},
		Action: func(c *cli.Context) error {

			client := c.App.Metadata["client"].(rpc.KonigRPCClient)
			resp, err := client.AddEdges(context.Background(), &rpc.AddEdgesRequest{})
			if err == nil {
				data, err := json.Marshal(resp)
				if err != nil {
					return err
				}
				fmt.Printf("%s\n", data)
			}
			return err
		},
	}

	removeVertexCommand = &cli.Command{
		Name:    "vertex",
		Aliases: []string{"v"},
		Action: func(c *cli.Context) error {
			client := c.App.Metadata["client"].(rpc.KonigRPCClient)
			resp, err := client.RemoveVertices(context.Background(), &rpc.RemoveVerticesRequest{})
			if err == nil {
				data, err := json.Marshal(resp)
				if err != nil {
					return err
				}
				fmt.Printf("%s\n", data)
			}
			return err
		},
	}

	removeEdgeCommand = &cli.Command{
		Name:    "edge",
		Aliases: []string{"e"},
		Action: func(c *cli.Context) error {
			client := c.App.Metadata["client"].(rpc.KonigRPCClient)
			resp, err := client.RemoveEdges(context.Background(), &rpc.RemoveEdgesRequest{})
			if err == nil {
				data, err := json.Marshal(resp)
				if err != nil {
					return err
				}
				fmt.Printf("%s\n", data)
			}
			return err
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
				if c.NArg() != 1 {
					return errors.New("invalid number of parameters")
				}
				name := c.Args().First()
				client := c.App.Metadata["client"].(rpc.KonigRPCClient)
				resp, err := client.CreateGraph(context.Background(), &rpc.CreateGraphRequest{Graph: &rpc.Graph{Name: name}})
				if err == nil {
					data, err := json.Marshal(resp)
					if err != nil {
						return err
					}
					fmt.Printf("%s\n", data)
				}
				return err
			},
		},
		&cli.Command{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "adds a vertex or an edge",
			Subcommands: []*cli.Command{
				addVertexCommand,
				addEdgeCommand,
			},
		},
		&cli.Command{
			Name:    "remove",
			Aliases: []string{"r"},
			Usage:   "removes a node or an edge",
			Subcommands: []*cli.Command{
				removeVertexCommand,
				removeEdgeCommand,
			},
		},
	}
	client := initClient()
	app.Metadata = make(map[string]interface{}, 1)
	app.Metadata["client"] = client
	app.Run(os.Args)
}

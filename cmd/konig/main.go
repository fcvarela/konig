package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"

	"github.com/fcvarela/konig/graph"
	"github.com/fcvarela/konig/rpc"
	"github.com/golang/glog"
)

var (
	signalhandlerChannel = make(chan os.Signal, 1)
)

func init() {
	// we need to process input and draw on the main thread
	runtime.LockOSThread()

	// parse flags and enable stderr out
	flag.Parse()
	//flag.Lookup("alsologtostderr").Value.Set("true")

	// install clean shutdown signal handler
	signal.Notify(signalhandlerChannel, os.Interrupt)
}

func main() {
	go rpc.StartRPC()

	// 0, 0, fullscreen. graph ignores width and height when full screen is set
	graph.Init(0, 0, true)

	// wait on sigint
	for {
		var stop = false
		select {
		case <-signalhandlerChannel:
			stop = true
		default:
			stop = graph.Update()
		}
		if stop {
			break
		}
	}

	glog.Info("Got abort signal, stopping...")
	graph.Shutdown()
}

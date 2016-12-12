package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"

	"github.com/fcvarela/konig"
	"github.com/fcvarela/konig/openclsolver"
	"github.com/fcvarela/konig/openglrenderer"
	"github.com/fcvarela/konig/rpc"
	"github.com/golang/glog"
)

var (
	grpcPort             uint
	grpcHost             string
	signalhandlerChannel = make(chan os.Signal, 1)
)

func setupGlogStderr() {
	if f := flag.Lookup("alsologtostderr"); f == nil {
		panic("Cannot find alsologtostderr flag")
	} else {
		if err := f.Value.Set("true"); err != nil {
			panic("Error setting alsologtostderr to true")
		}
	}
}

func init() {
	// setup command flag
	flag.UintVar(&grpcPort, "port", 1234, "port to be used by the grpc server")
	flag.StringVar(&grpcHost, "host", "0.0.0.0", "rpc host")
	flag.Parse()

	// we need to process input and draw on the main thread
	// which is the one we always start on
	runtime.LockOSThread()

	// parse flags
	flag.Parse()

	// enable glog stderr
	setupGlogStderr()

	// install clean shutdown signal handler
	signal.Notify(signalhandlerChannel, os.Interrupt)
}

func main() {
	// placeholder
	go rpc.Start(grpcHost, grpcPort)

	// init the renderer
	renderer := openglrenderer.New(1280, 720, false)

	// init the solver
	solver := openclsolver.New()

	// call defered startup
	konig.Startup(solver, renderer)

	// wait on sigint
	for {
		var stop = false
		select {
		case <-signalhandlerChannel:
			stop = true
		default:
			stop = konig.Step()
		}
		if stop {
			break
		}
	}

	glog.Info("Got abort signal, stopping...")
	konig.Shutdown()
}

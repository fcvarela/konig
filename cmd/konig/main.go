package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"

	"github.com/fcvarela/konig/graphview"
	"github.com/fcvarela/konig/rpc"
	"github.com/golang/glog"
)

var (
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
	go rpc.StartRPC()

	// 0, 0, fullscreen. graph ignores width and height when full screen is set
	graphview.Init(0, 0, true)

	// wait on sigint
	for {
		var stop = false
		select {
		case <-signalhandlerChannel:
			stop = true
		default:
			stop = graphview.Update()
		}
		if stop {
			break
		}
	}

	glog.Info("Got abort signal, stopping...")
	graphview.Shutdown()
}

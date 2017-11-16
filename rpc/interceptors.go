package rpc

import (
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// LoggingInterceptor logs and times all grpc calls
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	t := time.Now()
	resp, err := handler(ctx, req)
	glog.Infof("call %s took %dns", info.FullMethod, time.Since(t).Nanoseconds())
	return resp, err
}

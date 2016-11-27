package rpc

import (
	"strings"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//RequestLogger logs and times all grpc calls
func RequestLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	t := time.Now()
	resp, err := handler(ctx, req)
	method := strings.Replace(info.FullMethod, "/rpc.KonigRPC/", "", -1)
	glog.Infof("call %s took %dns", method, time.Since(t).Nanoseconds())
	return resp, err
}

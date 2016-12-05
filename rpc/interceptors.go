package rpc

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//Interceptor struct holds the needed configuration to run interceptors
type Interceptor struct {
	Debug bool
}

//NewInterceptor creates and configures a new interceptor
func NewInterceptor(debug bool) *Interceptor {
	return &Interceptor{Debug: debug}
}

//RequestLogger logs and times all grpc calls
func (i *Interceptor) RequestLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	t := time.Now()
	method := strings.Replace(info.FullMethod, "/rpc.KonigRPC/", "", -1)
	if i.Debug {
		data, err := json.Marshal(req)
		if err != nil {
			glog.Infof("unable to log request payload for method %s: %s", method, err)
		} else {
			glog.Infof("%s << %s", method, data)
		}
	}
	resp, err := handler(ctx, req)
	if i.Debug {
		data, mErr := json.Marshal(resp)
		if err != nil {
			glog.Infof("unable to log response payload for method %s: %s", method, mErr)
		} else {
			glog.Infof("%s >> %s", method, data)
		}
	}
	glog.Infof("call %s took %dns", method, time.Since(t).Nanoseconds())
	return resp, err
}

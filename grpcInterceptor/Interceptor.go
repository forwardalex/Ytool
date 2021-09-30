package grpcInterceptor

import (
	"Ytool/codes"
	"Ytool/log"
	"context"
	"github.com/tal-tech/go-zero/core/breaker"
	"google.golang.org/grpc"
	"path"
)

// 请求日志打印
func ReqLogInterceptor() grpc.UnaryServerInterceptor {
	interceptor := grpc.UnaryServerInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log.Infof("grpc method:", info.FullMethod, ",requestBody is:", req)
		resp, err = handler(ctx, req)
		if err != nil {
			log.Errorf("grpc method:", info.FullMethod, "Error,", err)
		}
		log.Infof("grpc method:", info.FullMethod, ",responseBody is:", resp)
		return resp, err
	})
	return interceptor
}

func BreakerInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 基于请求方法进行熔断
	breakerName := path.Join(cc.Target(), method)
	return breaker.DoWithAcceptable(breakerName, func() error {
		// 真正发起调用
		return invoker(ctx, method, req, reply, cc, opts...)
		// codes.Acceptable判断哪种错误需要加入熔断错误计数
	}, codes.Acceptable)
}

// 填坑测试
func BreakerInterceptor2() grpc.UnaryServerInterceptor {
	interceptor := grpc.UnaryServerInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log.Infof("grpc method:", info.FullMethod, ",requestBody is:", req)
		breaker.DoWithAcceptable(info.FullMethod, func() (err error) {
			resp, err = handler(ctx, req)
			if err != nil {
				log.Errorf("grpc method:", info.FullMethod, "Error,", err)
			}
			log.Infof("grpc method:", info.FullMethod, ",responseBody is:", resp)
			return err
		}, codes.Acceptable)
		return resp, err
	})
	return interceptor
}

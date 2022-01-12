package grpcInterceptor

import (
	"context"
	"fmt"
	"github.com/forwardalex/Ytool/codes"
	"github.com/forwardalex/Ytool/log"
	"github.com/forwardalex/Ytool/proto"
	"github.com/forwardalex/Ytool/svc"
	"github.com/tal-tech/go-zero/core/breaker"
	"google.golang.org/grpc"
	"net"
	"path"
)

// 请求日志打印
func ReqLogInterceptor() grpc.UnaryServerInterceptor {
	interceptor := grpc.UnaryServerInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log.Infof(context.Background(), "grpc method:", info.FullMethod, ",requestBody is:", req)
		resp, err = handler(ctx, req)
		if err != nil {
			log.Errorf(context.Background(), "grpc method:", info.FullMethod, "Error,", err)
		}
		log.Infof(context.Background(), "grpc method:", info.FullMethod, ",responseBody is:", resp)
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

// 填坑测试  拦截器相关于一个aop
func BreakerInterceptor2() grpc.UnaryServerInterceptor {
	interceptor := grpc.UnaryServerInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log.Infof(context.Background(), "grpc method:", info.FullMethod, ",requestBody is:", req)
		breaker.DoWithAcceptable(info.FullMethod, func() (err error) {
			resp, err = handler(ctx, req)
			if err != nil {
				log.Errorf(context.Background(), "grpc method:", info.FullMethod, "Error,", err)
			}
			log.Infof(context.Background(), "grpc method:", info.FullMethod, ",responseBody is:", resp)
			return err
		}, codes.Acceptable)
		return resp, err
	})
	return interceptor
}

func Test() {
	lis, err := net.Listen("tcp", ":50056")
	if err != nil {
		log.Fatal(context.Background(), "failed to listen: %v", err)
	}
	size := 100 * 1024 * 1024
	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(size),
		grpc.MaxSendMsgSize(size),
		grpc.ChainUnaryInterceptor(BreakerInterceptor2()),
	)
	if err := s.Serve(lis); err != nil {
		log.Fatal(context.Background(), "failed: %v", err)
	}
	s.RegisterService(&proto.AopTest_ServiceDesc, &svc.Svc{})
	fmt.Println("Proxy started")
}

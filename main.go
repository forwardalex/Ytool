package main

import (
	"context"
	"fmt"
	grpchttpproxy "github.com/forwardalex/Ytool/grpc/grpcproxy"
	"github.com/forwardalex/Ytool/grpc/grpcproxy/testuse"
	"github.com/forwardalex/Ytool/log"
	pb "github.com/forwardalex/Ytool/proto"
	"github.com/forwardalex/Ytool/store/db"
	"github.com/forwardalex/Ytool/tool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"runtime"
)

var (
	port       = ":50051"
	httpPort   = 80
	httpPorStr = ":80"
)

func main() {
	fmt.Println(runtime.Version())
	tool.Init()
	fmt.Println("welcome")
	//test.TestBlame()
	//grpcInterceptor.Test()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(context.TODO(), err)
	}
	size := 100 * 1024 * 1024

	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(size),
		grpc.MaxSendMsgSize(size),
	)
	pb.RegisterHelloServiceServer(s, &testuse.HelloService{})
	fmt.Println("server registered: 0.0.0.0", port)
	reflection.Register(s)
	err = db.GetRedisConn().Ping(context.Background()).Err()
	if err != nil {
		fmt.Println("====", err)
		return
	}
	fmt.Println("redis ok")
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(context.TODO(), err)
		}
		log.Info(context.TODO(), err, "info ")
	}()
	l, err := net.Listen("tcp", httpPorStr)
	if err != nil {
		log.Fatal(context.TODO(), err)
	}
	grpchttpproxy.StartWithInterceptor(l, "api/access/pb/cmd", "localhost"+port, nil)
}

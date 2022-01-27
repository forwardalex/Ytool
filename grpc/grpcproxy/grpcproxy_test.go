package grpchttpproxy

import (
	"context"
	"fmt"
	"github.com/forwardalex/Ytool/grpc/grpcproxy/testuse"
	"github.com/forwardalex/Ytool/log"
	pb "github.com/forwardalex/Ytool/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"testing"
)

var (
	port     = ":50051"
	httpPort = 80
)

//先开启一个grpc服务测试
func TestReflect(t *testing.T) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(context.TODO(), "failed to listen: ", err)
	}
	size := 100 * 1024 * 1024

	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(size),
		grpc.MaxSendMsgSize(size),
	)
	pb.RegisterHelloServiceServer(s, &testuse.HelloService{})
	fmt.Println("server registered: 0.0.0.0", port)
	reflection.Register(s)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(context.TODO(), "failed:", err)
		}
		log.Info(context.TODO(), "Proxy started", "")
	}()
	StartWithInterceptor(httpPort, "api/access/pb/cmd", "localhost"+port, nil)
}

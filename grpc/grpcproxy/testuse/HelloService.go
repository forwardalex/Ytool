package testuse

import (
	"context"
	"github.com/forwardalex/Ytool/log"
	pb "github.com/forwardalex/Ytool/proto"
)

// HelloService server
type HelloService struct {
	pb.UnimplementedHelloServiceServer
}

//hello
func (server *HelloService) Hello(ctx context.Context,
	req *pb.CheckResultReq) (resp *pb.CheckResultResp, err error) {

	resp = &pb.CheckResultResp{ResponseCode: "code1"}
	return resp, nil
}

//hello2
func (server *HelloService) Hello2(ctx context.Context,
	req *pb.CheckResultReq) (resp *pb.CheckResultResp, err error) {
	log.Info(ctx, "server ", " ok")
	resp = &pb.CheckResultResp{ResponseCode: "code2"}
	return resp, nil
}

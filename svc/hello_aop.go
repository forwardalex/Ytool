package svc

import (
	"context"
	"github.com/forwardalex/Ytool/log"
	"github.com/forwardalex/Ytool/proto"
)

type Svc struct {
	proto.UnimplementedAopTestServer
}

func (s *Svc) Hello(ctx context.Context, req *proto.CheckResultReq) (resp *proto.CheckResultResp, err error) {
	log.Info(context.Background(), "mock test", "")
	resp = &proto.CheckResultResp{
		ResponseCode: "1",
	}
	return resp, nil
}

package svc

import (
	"Ytool/log"
	"Ytool/proto"
	"context"
)

type Svc struct {
	proto.UnimplementedAopTestServer
}

func (s *Svc) Hello(ctx context.Context, req *proto.CheckResultReq) (resp *proto.CheckResultResp, err error) {
	log.Info("mock test", "")
	resp = &proto.CheckResultResp{
		ResponseCode: "1",
	}
	return resp, nil
}

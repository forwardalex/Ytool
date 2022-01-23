package svc

var str = `
type %s struct {
	proto.UnimplementedAopTestServer
}

func (s *%s) Hello(ctx context.Context, req *proto.%sReq) (resp *proto.%sResp, err error) {
	return resp, nil
}
`
